package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nhh0718/idops/internal/env"
	"github.com/nhh0718/idops/internal/ui"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment file manager",
}

var envCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare .env.example with .env",
	RunE:  runEnvCompare,
}

var envSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Interactive sync missing variables",
	RunE:  runEnvSync,
}

var envValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate .env format",
	RunE:  runEnvValidate,
}

var envInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate .env from .env.example",
	RunE:  runEnvInit,
}

var envShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display .env with masked secrets",
	RunE:  runEnvShow,
}

func init() {
	envCompareCmd.Flags().String("source", ".env.example", "source env file")
	envCompareCmd.Flags().String("target", ".env", "target env file")
	envSyncCmd.Flags().String("source", ".env.example", "source env file")
	envSyncCmd.Flags().String("target", ".env", "target env file")
	envInitCmd.Flags().String("source", ".env.example", "source env file")
	envInitCmd.Flags().String("output", ".env", "output file")
	envInitCmd.Flags().Bool("force", false, "overwrite existing file")
	envValidateCmd.Flags().String("file", ".env", "file to validate")
	envShowCmd.Flags().String("file", ".env", "file to show")
	envShowCmd.Flags().Bool("json", false, "output as JSON")

	envCmd.AddCommand(envCompareCmd, envSyncCmd, envValidateCmd, envInitCmd, envShowCmd)
	rootCmd.AddCommand(envCmd)
}

func runEnvCompare(cmd *cobra.Command, args []string) error {
	srcPath, _ := cmd.Flags().GetString("source")
	tgtPath, _ := cmd.Flags().GetString("target")

	source, err := env.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", srcPath, err)
	}
	target, err := env.Parse(tgtPath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", tgtPath, err)
	}

	diff := env.Compare(source, target)
	if !diff.HasDifferences() {
		fmt.Println(ui.RenderSuccess("All variables are in sync!"))
		return nil
	}

	if len(diff.Missing) > 0 {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("MISSING in %s (%d):", tgtPath, len(diff.Missing))))
		for _, k := range diff.Missing {
			fmt.Printf("  + %s\n", k)
		}
	}
	if len(diff.Extra) > 0 {
		fmt.Println(ui.RenderInfo(fmt.Sprintf("EXTRA in %s, not in %s (%d):", tgtPath, srcPath, len(diff.Extra))))
		for _, k := range diff.Extra {
			fmt.Printf("  ~ %s\n", k)
		}
	}
	return nil
}

func runEnvSync(cmd *cobra.Command, args []string) error {
	srcPath, _ := cmd.Flags().GetString("source")
	tgtPath, _ := cmd.Flags().GetString("target")

	source, err := env.Parse(srcPath)
	if err != nil {
		return err
	}
	target, err := env.Parse(tgtPath)
	if err != nil {
		return err
	}
	return env.SyncInteractive(source, target)
}

func runEnvValidate(cmd *cobra.Command, args []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	issues, err := env.Validate(filePath)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		fmt.Println(ui.RenderSuccess("No issues found!"))
		return nil
	}
	fmt.Printf("Found %d issue(s):\n", len(issues))
	for _, issue := range issues {
		fmt.Printf("  line %d: %s - %s\n", issue.Line, issue.Key, issue.Message)
	}
	return nil
}

func runEnvInit(cmd *cobra.Command, args []string) error {
	srcPath, _ := cmd.Flags().GetString("source")
	outPath, _ := cmd.Flags().GetString("output")
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("%s already exists (use --force to overwrite)", outPath)
		}
	}

	source, err := env.Parse(srcPath)
	if err != nil {
		return err
	}
	return env.InitFromExample(source, outPath)
}

func runEnvShow(cmd *cobra.Command, args []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	envFile, err := env.Parse(filePath)
	if err != nil {
		return err
	}

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		masked := make(map[string]string)
		for k, v := range envFile.Vars {
			masked[k] = env.MaskValue(k, v)
		}
		data, _ := json.MarshalIndent(masked, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	for _, key := range envFile.Order {
		value := env.MaskValue(key, envFile.Vars[key])
		fmt.Printf("  %s=%s\n", key, value)
	}
	return nil
}
