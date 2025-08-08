package cmd

import (
	"context"

	"github.com/lost-melody/openra-tr-tools/pkg"
	"github.com/spf13/cobra"
)

var (
	extractOutputFile string
	extractRegexp     string
	patchFile         string
	patchOutputDir    string
)

// extractCmd represents the extract command.
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract strings from yaml files.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.ExtractStringsFromFile(context.Background(), args, extractOutputFile, extractRegexp)
		if err != nil {
			cmd.PrintErrf("failed to extract strings: %s.\n", err)
			return
		}
	},
}

// patchCmd represents the extract command.
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Patch strings from yaml files.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.PatchStringsInFile(context.Background(), args, patchFile, patchOutputDir)
		if err != nil {
			cmd.PrintErrf("failed to extract strings: %s.\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.Flags().StringVarP(&extractOutputFile, "output", "o", "-", "Output file to write strings into.")
	extractCmd.Flags().StringVarP(&extractRegexp, "regexp", "r", "", "Regexp for keys to extract.")

	rootCmd.AddCommand(patchCmd)
	patchCmd.Flags().StringVarP(&patchFile, "patch", "p", "-", "Patch file to read strings from.")
	patchCmd.Flags().StringVarP(&patchOutputDir, "output", "o", "output", "Output dir to write files into.")
}
