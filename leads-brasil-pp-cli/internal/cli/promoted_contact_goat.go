package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newContactGoatPromotedCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "contact-goat <target>",
		Short:       "Run contact-goat enrichment for a company, URL, or contact",
		Long:        "Run contact-goat enrichment for a company, URL, or contact",
		Example:     "  leads-brasil-pp-cli contact-goat \"Facto Agência de Comunicação\" --agent",
		Annotations: map[string]string{"pp:method": "POST", "pp:path": "/contact-goat"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			if len(args) < 1 {
				if flags.asJSON {
					if printErr := printJSONFiltered(cmd.OutOrStdout(), map[string]any{"error": "target is required", "usage": fmt.Sprintf("%s <%s>", cmd.CommandPath(), "target")}, flags); printErr != nil {
						return printErr
					}
				}
				return usageErr(fmt.Errorf("target is required\nUsage: %s <%s>", cmd.CommandPath(), "target"))
			}
			path := "/contact-goat"
			params := map[string]string{}
			body := map[string]any{"contact": args[0], "company": args[0]}
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
