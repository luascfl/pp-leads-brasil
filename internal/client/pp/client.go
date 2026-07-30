package pp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"pp-leads-brasil/internal/client/casadados"
	"pp-leads-brasil/internal/usecase"
)

// Client defines the interface for enrichment workflows.
type Client interface {
	RunCompanyGoat(cnpj string) (interface{}, error)
	RunContactGoat(target string) (interface{}, error)
	RunEnrich(cnpj string) (interface{}, error)
	RunApifyActor(actorID string, input interface{}) (interface{}, error)
}

type PPClient struct {
	APIKey            string
	CasaDadosAPIKey   string
	LeadTablePath     string
	CasaDadosBaseURL  string
	HTTPClient        *http.Client
	CompanyGoatBin    string
	ContactGoatBin    string
	ScrapeCreatorsBin string
	CommandTimeout    time.Duration
}

func (c *PPClient) RunCompanyGoat(cnpj string) (interface{}, error) {
	cd := c.casaDadosClient()
	lead, ok, err := cd.FindLeadByCNPJ(cnpj)
	if err != nil {
		return nil, err
	}

	company, companyErr := cd.SearchCompanyByCNPJ(cnpj)
	if !ok && companyErr != nil {
		return nil, companyErr
	}

	var goat any
	if ok {
		goat, _ = c.callCompanyGoat(lead, company)
	}
	if !ok {
		return map[string]any{
			"status":       "company-goat-external-only",
			"source":       "Casa dos Dados + company-goat quando disponível",
			"cnpj":         cnpj,
			"company":      company,
			"company_goat": goatStatus(goat, nil),
		}, nil
	}

	payload := map[string]any{
		"status":       "company-goat-local",
		"source":       "lead table + Casa dos Dados + company-goat quando disponível",
		"cnpj":         lead["CNPJ"],
		"company":      company,
		"company_goat": goatStatus(goat, nil),
	}
	if companyErr != nil {
		payload["casadosdados_status"] = map[string]any{"status": "unavailable", "error": companyErr.Error()}
	}
	if goat == nil {
		payload["company_goat"] = map[string]any{"status": "unavailable", "error": "company-goat-pp-cli não encontrado ou sem resposta útil"}
	}
	return payload, nil
}

func (c *PPClient) RunContactGoat(target string) (interface{}, error) {
	needle := strings.TrimSpace(target)
	if needle == "" {
		return nil, fmt.Errorf("contact target is required")
	}

	cd := c.casaDadosClient()
	leads, err := cd.LoadLeads()
	if err != nil {
		return nil, err
	}

	localMatches := make([]map[string]any, 0, 2)
	var matchedLead casadados.LeadRecord
	for _, lead := range leads {
		if matchesLead(target, lead) {
			localMatches = append(localMatches, casadados.CompanyPayload(lead))
			if matchedLead == nil {
				matchedLead = lead
			}
		}
	}
	payload := map[string]any{
		"status":        "contact-goat-local",
		"input":         target,
		"local_matches": localMatches,
	}

	if matchedLead != nil {
		contactGoat, err := c.callContactGoat(matchedLead)
		scrapeCreators, scrapeErr := c.callScrapeCreators(matchedLead)
		payload["contact_goat"] = goatStatus(contactGoat, err)
		payload["scrape_creators"] = goatStatus(scrapeCreators, scrapeErr)
		payload["prospecting_links"] = ProspectingLinks(matchedLead, localMatches[0], contactGoat, scrapeCreators)
		return payload, nil
	}

	contactGoat, err := c.callRawContactGoat(target)
	payload["contact_goat"] = goatStatus(contactGoat, err)
	if len(localMatches) == 0 && contactGoat == nil {
		return nil, casadados.ErrNotFound
	}
	return payload, nil
}

func (c *PPClient) RunEnrich(cnpj string) (interface{}, error) {
	cd := c.casaDadosClient()
	lead, ok, err := cd.FindLeadByCNPJ(cnpj)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, casadados.ErrNotFound
	}

	company, companyErr := cd.SearchCompanyByCNPJ(cnpj)
	companyGoat, companyGoatErr := c.callCompanyGoat(lead, company)
	contactGoat, contactGoatErr := c.callContactGoat(lead)
	scrapeCreators, scrapeCreatorsErr := c.callScrapeCreators(lead)

	payload := EnrichmentPayload(lead, company, companyErr)
	payload["company_goat"] = goatStatus(companyGoat, companyGoatErr)
	payload["contact_goat"] = goatStatus(contactGoat, contactGoatErr)
	payload["scrape_creators"] = goatStatus(scrapeCreators, scrapeCreatorsErr)
	payload["prospecting_links"] = ProspectingLinks(lead, company, contactGoat, scrapeCreators)
	payload["sources_used"] = []string{"local_lead_table", "casa_dos_dados", "company-goat", "contact-goat", "scrape-creators"}
	return payload, nil
}

func (c *PPClient) RunApifyActor(actorID string, input interface{}) (interface{}, error) {
	switch strings.TrimSpace(actorID) {
	case "company-goat":
		cnpj, err := cnpjFromInput(input)
		if err != nil {
			return nil, err
		}
		return c.RunCompanyGoat(cnpj)
	case "contact-goat":
		contact, err := contactFromInput(input)
		if err != nil {
			return nil, err
		}
		return c.RunContactGoat(contact)
	case "scrape-creators":
		contact, err := contactFromInput(input)
		if err != nil {
			return nil, err
		}
		cd := c.casaDadosClient()
		lead, ok, err := cd.FindLeadByCNPJ(contact)
		if err == nil && ok {
			return c.callScrapeCreators(lead)
		}
		return nil, fmt.Errorf("scrape-creators expects a lead resolvable by CNPJ in the current local table")
	default:
		return nil, fmt.Errorf("unsupported local actor %q", actorID)
	}
}

func (c *PPClient) casaDadosClient() *casadados.CasaDadosClient {
	if c == nil {
		return &casadados.CasaDadosClient{}
	}
	apiKey := c.CasaDadosAPIKey
	if apiKey == "" {
		apiKey = c.APIKey
	}
	return &casadados.CasaDadosClient{
		APIKey:        apiKey,
		LeadTablePath: c.LeadTablePath,
		BaseURL:       c.CasaDadosBaseURL,
		HTTPClient:    c.HTTPClient,
	}
}

func (c *PPClient) callCompanyGoat(lead casadados.LeadRecord, company any) (any, error) {
	bin := c.companyGoatBin()
	if bin == "" {
		return nil, fmt.Errorf("company-goat-pp-cli não encontrado")
	}
	args := []string{"snapshot", lead["lead"], "--agent"}
	if domain := companyDomain(lead, company); domain != "" {
		args = append(args, "--domain", domain)
	}
	return runCLIJSON(c.timeout(), bin, args...)
}

func (c *PPClient) callContactGoat(lead casadados.LeadRecord) (any, error) {
	bin := c.contactGoatBin()
	if bin == "" {
		return nil, fmt.Errorf("contact-goat-pp-cli não encontrado")
	}
	args := []string{"coverage", lead["lead"], "--agent", "--limit", "10"}
	return runCLIJSON(c.timeout(), bin, args...)
}

func (c *PPClient) callRawContactGoat(target string) (any, error) {
	bin := c.contactGoatBin()
	if bin == "" {
		return nil, fmt.Errorf("contact-goat-pp-cli não encontrado")
	}
	return runCLIJSON(c.timeout(), bin, "search", target, "--agent")
}

func (c *PPClient) callScrapeCreators(lead casadados.LeadRecord) (any, error) {
	bin := c.scrapeCreatorsBin()
	if bin == "" {
		return nil, fmt.Errorf("scrape-creators-pp-cli não encontrado")
	}

	result := map[string]any{"lead": lead["lead"]}
	var successes int
	var errs []string

	linkedinSearch, err := c.runScrapeCreatorsJSON("google", "list-search", "--query", fmt.Sprintf("\"%s\" site:linkedin.com/company", lead["lead"]), "--agent")
	if err == nil {
		result["linkedin_search"] = linkedinSearch
		if linkedinURL := firstURLFromResults(linkedinSearch); linkedinURL != "" {
			linkedinCompany, companyErr := c.runScrapeCreatorsJSON("linkedin", "list-company", "--url", linkedinURL, "--agent")
			if companyErr == nil {
				result["linkedin_company"] = linkedinCompany
				successes++
			} else {
				errs = append(errs, "linkedin company: "+companyErr.Error())
			}
		}
	} else {
		errs = append(errs, "linkedin search: "+err.Error())
	}

	instagramSearch, err := c.runScrapeCreatorsJSON("google", "list-search", "--query", fmt.Sprintf("\"%s\" site:instagram.com", lead["lead"]), "--agent")
	if err == nil {
		result["instagram_search"] = instagramSearch
		if instagramURL := firstURLFromResults(instagramSearch); instagramURL != "" {
			if handle := instagramHandleFromURL(instagramURL); handle != "" {
				instaProfile, profileErr := c.runScrapeCreatorsJSON("instagram", "list-profile", "--handle", handle, "--trim", "--agent")
				if profileErr == nil {
					result["instagram_profile"] = instaProfile
					successes++
				} else {
					errs = append(errs, "instagram profile: "+profileErr.Error())
				}
			}
		}
	} else {
		errs = append(errs, "instagram search: "+err.Error())
	}

	if site := strings.TrimSpace(lead["site"]); isBioServiceURL(site) {
		bioResult, bioErr := c.runScrapeCreatorsJSON("bio", "resolve", site, "--agent")
		if bioErr == nil {
			result["bio_resolve"] = bioResult
			successes++
		} else {
			errs = append(errs, "bio resolve: "+bioErr.Error())
		}
	}

	if len(errs) > 0 {
		result["notes"] = errs
	}
	if successes == 0 {
		return result, fmt.Errorf("scrape-creators sem dados utilizáveis; %s", strings.Join(errs, " | "))
	}
	return result, nil
}

func (c *PPClient) companyGoatBin() string {
	if c != nil && c.CompanyGoatBin != "" {
		return c.CompanyGoatBin
	}
	if path, err := exec.LookPath("company-goat-pp-cli"); err == nil {
		return path
	}
	return ""
}

func (c *PPClient) contactGoatBin() string {
	if c != nil && c.ContactGoatBin != "" {
		return c.ContactGoatBin
	}
	if path, err := exec.LookPath("contact-goat-pp-cli"); err == nil {
		return path
	}
	return ""
}

func (c *PPClient) scrapeCreatorsBin() string {
	if c != nil && c.ScrapeCreatorsBin != "" {
		return c.ScrapeCreatorsBin
	}
	if path, err := exec.LookPath("scrape-creators-pp-cli"); err == nil {
		return path
	}
	return ""
}

func (c *PPClient) runScrapeCreatorsJSON(args ...string) (any, error) {
	bin := c.scrapeCreatorsBin()
	if bin == "" {
		return nil, fmt.Errorf("scrape-creators-pp-cli não encontrado")
	}
	return runCLIJSON(c.timeout(), bin, args...)
}

func (c *PPClient) timeout() time.Duration {
	if c != nil && c.CommandTimeout > 0 {
		return c.CommandTimeout
	}
	return 45 * time.Second
}

func EnrichmentPayload(lead casadados.LeadRecord, company interface{}, companyErr error) map[string]any {
	cfg, _ := usecase.LoadFromEnv()
	if company == nil {
		company = casadados.CompanyPayload(lead)
	}
	payload := map[string]any{
		"status":       "enriched-local",
		"source":       "local lead table + Casa dos Dados quando configurado",
		"cnpj":         usecase.CNPJ(lead, cfg),
		"company":      company,
		"lead_context": usecase.LeadContext(lead),
	}
	if meta := cfg.Metadata(); meta != nil {
		payload["use_case"] = meta
	}
	if companyErr != nil {
		payload["casadosdados_status"] = map[string]any{"status": "unavailable", "error": companyErr.Error()}
	}
	return payload
}

func ProspectingLinks(lead casadados.LeadRecord, company interface{}, contactGoat interface{}, scrapeCreators interface{}) map[string]any {
	cfg, _ := usecase.LoadFromEnv()
	message := usecase.Message(lead, cfg)
	links := map[string]any{
		"whatsapp":             []map[string]string{},
		"link_whatsapp_pronto": []string{},
	}
	if email := cleanEmail(usecase.Contact(lead, cfg)); email != "" {
		mailto := mailtoLink(email, "Troca rápida sobre a operação do lead", message)
		links["email"] = email
		links["mailto"] = mailto
		links["link_email_pronto"] = mailto
	}
	phones := collectPhones(company, contactGoat, scrapeCreators)
	if len(phones) == 0 {
		phones = publicWebsitePhones(lead, company)
	}
	if len(phones) > 0 {
		waLinks := make([]map[string]string, 0, len(phones))
		readyLinks := make([]string, 0, len(phones))
		for _, phone := range phones {
			link := waLink(phone, message)
			waLinks = append(waLinks, map[string]string{"phone": phone, "wa_me": link, "link_whatsapp_pronto": link})
			readyLinks = append(readyLinks, link)
		}
		links["whatsapp"] = waLinks
		links["link_whatsapp_pronto"] = readyLinks
	}
	return links
}

func goatStatus(result any, err error) map[string]any {
	if err != nil {
		return map[string]any{"status": "unavailable", "error": err.Error()}
	}
	return map[string]any{"status": "ok", "result": result}
}

func runCLIJSON(timeout time.Duration, bin string, args ...string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	var out any
	stdoutText := strings.TrimSpace(stdout.String())
	parseOK := stdoutText != "" && json.Unmarshal(stdout.Bytes(), &out) == nil
	stderrText := strings.TrimSpace(stderr.String())
	stderrLower := strings.ToLower(stderrText)
	warningOnly := strings.Contains(stderrLower, "warning:") || strings.Contains(stderrLower, "happenstance unavailable")

	if err != nil {
		if parseOK && warningOnly {
			return out, nil
		}
		msg := stderrText
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("%s failed: %s", bin, msg)
	}
	if !parseOK {
		return nil, fmt.Errorf("%s returned non-json output", bin)
	}
	return out, nil
}

func companyDomain(lead casadados.LeadRecord, company any) string {
	cfg, _ := usecase.LoadFromEnv()
	if domain := extractDomain(usecase.Site(lead, cfg)); domain != "" {
		return domain
	}
	if m, ok := company.(map[string]any); ok {
		if contacts, ok := m["contacts"].(map[string]any); ok {
			if site, ok := contacts["site"].(string); ok {
				return extractDomain(site)
			}
		}
		if cdd, ok := m["casadosdados"].(map[string]any); ok {
			for _, key := range []string{"dominio", "site", "website"} {
				if value, ok := cdd[key].(string); ok {
					if domain := extractDomain(value); domain != "" {
						return domain
					}
				}
			}
		}
	}
	return ""
}

func collectPhones(values ...interface{}) []string {
	seen := map[string]bool{}
	phones := make([]string, 0, 2)
	var walk func(interface{})
	walk = func(v interface{}) {
		switch x := v.(type) {
		case map[string]any:
			for key, val := range x {
				lk := strings.ToLower(key)
				if strings.Contains(lk, "telefone") || strings.Contains(lk, "phone") || strings.Contains(lk, "whatsapp") || strings.Contains(lk, "ddd") {
					maybeAddPhone(val, seen, &phones)
				}
				walk(val)
			}
		case []any:
			for _, item := range x {
				walk(item)
			}
		}
	}
	for _, value := range values {
		walk(value)
	}
	return phones
}

func extractDomain(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.Contains(strings.ToLower(raw), "sem site") {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	return host
}

func publicWebsitePhones(lead casadados.LeadRecord, company interface{}) []string {
	site := websiteURL(lead, company)
	if site == "" || looksLikeNonWebsite(site) {
		return nil
	}
	if !strings.Contains(site, "://") {
		site = "https://" + site
	}
	base, err := url.Parse(site)
	if err != nil {
		return nil
	}
	candidates := []string{base.String()}
	for _, path := range []string{"/contato", "/contato/", "/contact", "/contact/", "/fale-conosco", "/fale-conosco/"} {
		u := *base
		u.Path = path
		u.RawQuery = ""
		candidates = append(candidates, u.String())
	}
	client := &http.Client{Timeout: 8 * time.Second}
	seen := map[string]bool{}
	phones := make([]string, 0, 2)
	for _, candidate := range candidates {
		req, err := http.NewRequest(http.MethodGet, candidate, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "pp-leads-brasil/1.0")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
		_ = resp.Body.Close()
		if readErr != nil {
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			continue
		}
		for _, match := range phoneMatches(string(body), candidate, resp.Request.URL.String()) {
			maybeAddPhone(match, seen, &phones)
		}
	}
	return phones
}

func websiteURL(lead casadados.LeadRecord, company interface{}) string {
	cfg, _ := usecase.LoadFromEnv()
	if value := strings.TrimSpace(usecase.Site(lead, cfg)); value != "" {
		return value
	}
	if m, ok := company.(map[string]any); ok {
		if contacts, ok := m["contacts"].(map[string]any); ok {
			if site, ok := contacts["site"].(string); ok && strings.TrimSpace(site) != "" {
				return site
			}
		}
	}
	return ""
}

func looksLikeNonWebsite(raw string) bool {
	lower := strings.ToLower(strings.TrimSpace(raw))
	for _, marker := range []string{"sem site", "google.com/maps", "google.com.br/maps", "instagram.com", "linkedin.com"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func phoneMatches(values ...string) []string {
	seen := map[string]bool{}
	phones := make([]string, 0, 2)
	for _, value := range values {
		for _, match := range telHrefPattern.FindAllStringSubmatch(value, -1) {
			if len(match) > 1 {
				maybeAddPhone(match[1], seen, &phones)
			}
		}
		for _, match := range whatsappHrefPattern.FindAllStringSubmatch(value, -1) {
			if len(match) > 1 {
				maybeAddPhone(match[1], seen, &phones)
			}
		}
	}
	return phones
}

func isBioServiceURL(raw string) bool {
	host := extractDomain(raw)
	for _, candidate := range []string{"linktr.ee", "lnk.bio", "komi.io", "pillar.io", "linkme.bio"} {
		if host == candidate || strings.HasSuffix(host, "."+candidate) {
			return true
		}
	}
	return false
}

func firstURLFromResults(value interface{}) string {
	switch x := value.(type) {
	case map[string]any:
		for _, key := range []string{"url", "link", "profile_url", "company_url"} {
			if v, ok := x[key].(string); ok && strings.TrimSpace(v) != "" {
				return v
			}
		}
		if results, ok := x["results"]; ok {
			if url := firstURLFromResults(results); url != "" {
				return url
			}
		}
		if items, ok := x["items"]; ok {
			if url := firstURLFromResults(items); url != "" {
				return url
			}
		}
	case []any:
		for _, item := range x {
			if url := firstURLFromResults(item); url != "" {
				return url
			}
		}
	}
	return ""
}

func instagramHandleFromURL(raw string) string {
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	switch strings.ToLower(parts[0]) {
	case "p", "reel", "tv", "stories":
		return ""
	default:
		return parts[0]
	}
}

func matchesLead(target string, lead casadados.LeadRecord) bool {
	needle := strings.ToLower(strings.TrimSpace(target))
	if needle == "" {
		return false
	}
	candidates := []string{lead["lead"], lead["contato"], lead["site"], lead["cidade/UF"], lead["CNPJ"]}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		value := strings.ToLower(strings.TrimSpace(candidate))
		if value == needle || strings.Contains(value, needle) || strings.Contains(needle, value) {
			return true
		}
	}
	return false
}

func mailtoLink(email, subject, body string) string {
	q := url.Values{}
	q.Set("subject", subject)
	q.Set("body", body)
	return "mailto:" + email + "?" + q.Encode()
}

func waLink(phone, message string) string {
	digits := onlyDigits(phone)
	if strings.HasPrefix(digits, "55") {
		return "https://wa.me/" + digits + "?text=" + url.QueryEscape(message)
	}
	return "https://wa.me/55" + digits + "?text=" + url.QueryEscape(message)
}

func cleanEmail(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(strings.ToLower(value), "sem e-mail") || !strings.Contains(value, "@") {
		return ""
	}
	return value
}

func maybeAddPhone(value interface{}, seen map[string]bool, phones *[]string) {
	switch x := value.(type) {
	case string:
		digits := onlyDigits(x)
		if len(digits) >= 10 && len(digits) <= 13 && !seen[digits] {
			seen[digits] = true
			*phones = append(*phones, digits)
		}
	case float64:
		digits := onlyDigits(fmt.Sprintf("%.0f", x))
		if len(digits) >= 10 && len(digits) <= 13 && !seen[digits] {
			seen[digits] = true
			*phones = append(*phones, digits)
		}
	}
}

func onlyDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

var cnpjDigits = regexp.MustCompile(`\D+`)

var telHrefPattern = regexp.MustCompile(`(?i)href=["']tel:([^"'?#]+)`)
var whatsappHrefPattern = regexp.MustCompile(`(?i)(?:wa\.me/|api\.whatsapp\.com/send\?phone=)(\d{10,13})`)

func cnpjFromInput(input interface{}) (string, error) {
	switch v := input.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return "", fmt.Errorf("cnpj is required")
		}
		return cnpjDigits.ReplaceAllString(v, ""), nil
	case map[string]string:
		if strings.TrimSpace(v["cnpj"]) == "" {
			return "", fmt.Errorf("cnpj is required")
		}
		return cnpjDigits.ReplaceAllString(v["cnpj"], ""), nil
	case map[string]any:
		if cnpj, ok := v["cnpj"].(string); ok && strings.TrimSpace(cnpj) != "" {
			return cnpjDigits.ReplaceAllString(cnpj, ""), nil
		}
	}
	return "", fmt.Errorf("cnpj is required")
}

func contactFromInput(input interface{}) (string, error) {
	switch v := input.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return "", fmt.Errorf("contact is required")
		}
		return v, nil
	case map[string]string:
		for _, key := range []string{"url", "email", "contact", "company"} {
			if strings.TrimSpace(v[key]) != "" {
				return v[key], nil
			}
		}
	case map[string]any:
		for _, key := range []string{"url", "email", "contact", "company"} {
			if value, ok := v[key].(string); ok && strings.TrimSpace(value) != "" {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("contact is required")
}
