package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/bakyazi/gopass/config"
	"github.com/bakyazi/gopass/sheetsapi"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {

		site, err := cmd.Flags().GetString("site")
		if err != nil {
			fmt.Println("Cannot parse 'site' flag", err.Error())
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("Cannot parse 'username' flag", err.Error())
		}

		isPrint, err := cmd.Flags().GetBool("print")
		if err != nil {
			fmt.Println("Cannot parse 'username' flag", err.Error())
		}
		cfg := config.GetConfig()

		db, err := sheetsapi.NewPasswordDB(cfg.CredentialFile, cfg.SheetId)
		if err != nil {
			panic(err)
		}
		passw, err := db.GetPassword(site, username)
		if err != nil {
			fmt.Println("Failed to list all passwords", err.Error())
			return
		}

		err = clipboard.WriteAll(passw)
		if err != nil {
			fmt.Println("Failed to write password into clipboard")
			return
		}
		if isPrint {
			fmt.Println(passw)
		}
		fmt.Println("Success!")

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().String("site", "", "Enter website")
	_ = getCmd.MarkFlagRequired("site")

	getCmd.Flags().String("username", "", "Enter username")
	_ = getCmd.MarkFlagRequired("username")

	getCmd.Flags().Bool("print", false, "Enter print flag to print password to console")

}
