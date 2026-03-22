package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	internalssh "github.com/nhh0718/idops/internal/ssh"
)

func init() {
	sshCmd.AddCommand(sshKeysCmd)
	sshKeysCmd.AddCommand(sshKeysDeleteCmd)

	sshKeysCmd.Flags().Bool("json", false, "Xuất danh sách keys dạng JSON")
}

var sshKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Liệt kê SSH keys trong ~/.ssh/",
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := internalssh.ListKeys()
		if err != nil {
			return err
		}

		jsonOut, _ := cmd.Flags().GetBool("json")
		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(keys)
		}

		if len(keys) == 0 {
			fmt.Println("  Không tìm thấy SSH key nào trong ~/.ssh/")
			return nil
		}

		fmt.Printf("  Tìm thấy %d SSH key(s):\n\n", len(keys))
		for _, k := range keys {
			fmt.Printf("  📎 %s (%s)\n", k.Name, k.Type)
			if k.Comment != "" {
				fmt.Printf("     Comment: %s\n", k.Comment)
			}
			if k.Fingerprint != "" {
				fmt.Printf("     Fingerprint: %s\n", k.Fingerprint)
			}
			fmt.Println()
		}
		return nil
	},
}

var sshKeysDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Xóa SSH key pair",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := internalssh.DeleteKey(name); err != nil {
			return err
		}
		fmt.Printf("  Đã xóa key: %s\n", name)
		return nil
	},
}
