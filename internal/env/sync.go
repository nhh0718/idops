package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SyncInteractive prompts the user for missing variables and appends them to the target.
func SyncInteractive(source, target *EnvFile) error {
	diff := Compare(source, target)
	if len(diff.Missing) == 0 {
		fmt.Println("All variables are in sync!")
		return nil
	}

	fmt.Printf("Found %d missing variable(s). Enter values:\n\n", len(diff.Missing))
	reader := bufio.NewReader(os.Stdin)

	for _, key := range diff.Missing {
		defaultVal := source.Vars[key]
		if defaultVal != "" {
			fmt.Printf("  %s [default: %s]: ", key, defaultVal)
		} else {
			fmt.Printf("  %s: ", key)
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			input = defaultVal
		}

		target.Vars[key] = input
		target.Order = append(target.Order, key)

		// Carry over comment from source if available
		if comment, ok := source.Comments[key]; ok {
			target.Comments[key] = comment
		}
	}

	if err := target.Write(target.Path); err != nil {
		return fmt.Errorf("failed to write %s: %w", target.Path, err)
	}

	fmt.Printf("\nSynced %d variable(s) to %s\n", len(diff.Missing), target.Path)
	return nil
}

// InitFromExample creates a new .env from .env.example with all values prompted.
func InitFromExample(source *EnvFile, outputPath string) error {
	fmt.Printf("Generating %s from %s\n\n", outputPath, source.Path)
	reader := bufio.NewReader(os.Stdin)

	target := &EnvFile{
		Path:     outputPath,
		Vars:     make(map[string]string),
		Comments: source.Comments,
		Order:    source.Order,
	}

	for _, key := range source.Order {
		defaultVal := source.Vars[key]
		if defaultVal != "" {
			fmt.Printf("  %s [default: %s]: ", key, defaultVal)
		} else {
			fmt.Printf("  %s: ", key)
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			input = defaultVal
		}
		target.Vars[key] = input
	}

	if err := target.Write(outputPath); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	fmt.Printf("\nCreated %s with %d variable(s)\n", outputPath, len(target.Order))
	return nil
}
