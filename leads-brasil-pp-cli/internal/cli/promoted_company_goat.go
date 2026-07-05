package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCompanyGoatPromotedCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "company-goat <cnpj>",
		Short:       "Run company-goat enrichment for a CNPJ",
		Long:        "Run company-goat enrichment for a CNPJ",
		Example:     "  leads-brasil-pp-cli company-goat 11.370.755/0001-02 --agent",
		Annotations: map[string]string{"pp:method": "POST", "pp:path": "/company-goat/{cnpj}"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			path := "/company-goat/{cnpj}"
			if len(args) < 1 {
				if flags.asJSON {
					if printErr := printJSONFiltered(cmd.OutOrStdout(), map[string]any{"error": "cnpj is required", "usage": fmt.Sprintf("%s <%s>", cmd.CommandPath(), "cnpj")}, flags); printErr != nil {
						return printErr
					}
				}
				return usageErr(fmt.Errorf("cnpj is required\nUsage: %s <%s>", cmd.CommandPath(), "cnpj"))
			}
			path = replacePathParam(path, "cnpj", args[0])
			params := map[string]string{}
			body := map[string]any{}
			data, _, err := c.PostWithParams(cmd.Context(), path, params, body)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			if wantsHumanTable(cmd.OutOrStdout(), flags) {
				var items []map[string]any
				if json.Unmarshal(data, &items) == nil && len(items) > 0 {
					if err := printAutoTable(cmd.OutOrStdout(), items); err != nil {
						return err
					}
					if len(items) >= 25 {
						fmt.Fprintf(os.Stderr, "\nShowing %d results. To narrow: add --limit, --json --select, or filter flags.\n", len(items))
					}
					return nil
				}
			}
			return printOutputWithFlags(cmd.OutOrStdout(), data, flags)
		},
	}
	return cmd
}
