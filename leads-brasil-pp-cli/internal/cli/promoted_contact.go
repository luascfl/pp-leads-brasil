package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newContactPromotedCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contact <query>",
		Short: "Deep search for people/decision-makers inside a company",
		Long:  "Finds decision makers, titles, and links for a given company name, email, or URL.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := flags.newClient()
			if err != nil {
				return err
			}

			if len(args) < 1 {
				if flags.asJSON {
					if printErr := printJSONFiltered(cmd.OutOrStdout(), map[string]any{"error": "query is required", "usage": fmt.Sprintf("%s <%s>", cmd.CommandPath(), "query")}, flags); printErr != nil {
						return printErr
					}
				}
				return usageErr(fmt.Errorf("query is required\nUsage: %s <%s>", cmd.CommandPath(), "query"))
			}

			query := args[0]
			path := "/contact-goat"
			body := map[string]any{"company": query}

			data, _, err := c.PostWithParams(cmd.Context(), path, map[string]string{}, body)
			if err != nil {
				return classifyAPIError(err, flags)
			}

			if wantsHumanTable(cmd.OutOrStdout(), flags) {
				var items []map[string]any
				if json.Unmarshal(data, &items) == nil && len(items) > 0 {
					if err := printAutoTable(cmd.OutOrStdout(), items); err != nil {
						return err
					}
					return nil
				}
			}
			return printOutputWithFlags(cmd.OutOrStdout(), data, flags)
		},
	}
	return cmd
}
