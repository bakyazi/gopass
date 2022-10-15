package cmd

import (
	"fmt"
	"github.com/bakyazi/gopass/config"
	"github.com/bakyazi/gopass/sheetsapi"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a password from DB",
	Run: func(cmd *cobra.Command, args []string) {

		site, err := cmd.Flags().GetString("site")
		if err != nil {
			fmt.Println("Cannot parse 'site' flag", err.Error())
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("Cannot parse 'username' flag", err.Error())
		}

		cfg := config.GetConfig()

		db, err := sheetsapi.NewPasswordDB(cfg.CredentialFile, cfg.SheetId)
		if err != nil {
			panic(err)
		}

		err = db.DeletePassword(site, username)
		if err != nil {
			fmt.Println("Failed to create and save passwords", err.Error())
			return
		}

		fmt.Println("Success!")

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().String("site", "", "Enter website")
	_ = deleteCmd.MarkFlagRequired("site")

	deleteCmd.Flags().String("username", "", "Enter username")
	_ = deleteCmd.MarkFlagRequired("username")

}
