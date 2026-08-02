package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newOperationCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operation",
		Short: "Create and explicitly apply bounded external operation plans",
	}
	var inputPath string
	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "Create a side-effect-free operation plan from JSON input",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if inputPath == "" {
				return usageErr(fmt.Errorf("--input is required"))
			}
			data, err := os.ReadFile(inputPath)
			if err != nil {
				return err
			}
			var body map[string]any
			if err := json.Unmarshal(data, &body); err != nil {
				return fmt.Errorf("parsing operation input: %w", err)
			}
			client, err := flags.newClient()
			if err != nil {
				return err
			}
			response, _, err := client.PostWithParams(cmd.Context(), "/operations/plan", map[string]string{}, body)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			return printOutputWithFlags(cmd.OutOrStdout(), response, flags)
		},
	}
	planCmd.Flags().StringVar(&inputPath, "input", "", "JSON file containing profile, kind, operator and targets")

	applyCmd := &cobra.Command{
		Use:   "apply <plan-id>",
		Short: "Apply a stored operation plan after explicit approval",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("yes") && !cmd.InheritedFlags().Changed("yes") {
				return fmt.Errorf("confirmation required: pass --yes explicitly; --agent is not approval")
			}
			client, err := flags.newClient()
			if err != nil {
				return err
			}
			response, _, err := client.PostWithParams(cmd.Context(), "/operations/"+args[0]+"/apply", map[string]string{}, map[string]any{"approved": true})
			if err != nil {
				return classifyAPIError(err, flags)
			}
			return printOutputWithFlags(cmd.OutOrStdout(), response, flags)
		},
	}
	cmd.AddCommand(planCmd, applyCmd)
	return cmd
}
