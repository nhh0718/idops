package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	internalssh "github.com/nhh0718/idops/internal/ssh"
)

func init() {
	sshCmd.AddCommand(sshKeygenCmd)

	sshKeygenCmd.Flags().String("name", "id_ed25519", "Tên file key (trong ~/.ssh/)")
	sshKeygenCmd.Flags().String("type", "ed25519", "Loại key: ed25519 hoặc rsa")
	sshKeygenCmd.Flags().Int("bits", 4096, "Số bit RSA (chỉ áp dụng với --type rsa)")
	sshKeygenCmd.Flags().String("comment", "", "Comment (thường là email)")
	sshKeygenCmd.Flags().Bool("json", false, "Xuất kết quả dạng JSON")
	sshKeygenCmd.Flags().Bool("force", false, "Ghi đè nếu key đã tồn tại")
}

var sshKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Tạo SSH key pair mới",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		keyType, _ := cmd.Flags().GetString("type")
		bits, _ := cmd.Flags().GetInt("bits")
		comment, _ := cmd.Flags().GetString("comment")
		jsonOut, _ := cmd.Flags().GetBool("json")
		force, _ := cmd.Flags().GetBool("force")

		if keyType != "ed25519" && keyType != "rsa" {
			return fmt.Errorf("loại key không hợp lệ %q, chỉ hỗ trợ ed25519 hoặc rsa", keyType)
		}

		result, err := internalssh.GenerateKey(internalssh.KeygenOptions{
			Name:    name,
			Type:    keyType,
			Bits:    bits,
			Comment: comment,
			Force:   force,
		})
		if err != nil {
			return err
		}

		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if result.Output != "" {
			fmt.Print(result.Output)
		}
		fmt.Printf("\nĐã tạo SSH key: %s\n", result.PrivateKey)
		fmt.Printf("Public key:     %s\n", result.PublicKey)
		return nil
	},
}
