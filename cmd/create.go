package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/bakyazi/gopass/config"
	"github.com/bakyazi/gopass/passgen"
	"github.com/bakyazi/gopass/sheetsapi"
	"github.com/spf13/cobra"
	"os"
)

var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {

		site, err := cmd.Flags().GetString("site")
		if err != nil {
			fmt.Println("Cannot parse 'site' flag", err.Error())
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("Cannot parse 'username' flag", err.Error())
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println("Cannot parse 'password' flag", err.Error())
		}

		auto, err := cmd.Flags().GetString("auto")
		if err != nil {
			fmt.Println("Cannot parse 'auto' flag", err.Error())
		}

		if auto != "" {
			password, err = passgen.Generate(auto)
			if err != nil {
				fmt.Println("Cannot generate password", err.Error())
				os.Exit(1)
			}
		}

		if password == "" {
			fmt.Println("You should set 'password' or 'auto' flag")
			return
		}

		cfg := config.GetConfig()

		db, err := sheetsapi.NewPasswordDB(cfg.CredentialFile, cfg.SheetId)
		if err != nil {
			panic(err)
		}

		err = db.CreatePassword(site, username, password)
		if err != nil {
			fmt.Println("Failed to create and save passwords", err.Error())
			return
		}

		err = clipboard.WriteAll(password)
		if err != nil {
			fmt.Println("Failed to write password into clipboard")
			return
		}
		fmt.Println("Success!")

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("site", "", "Enter website")
	_ = createCmd.MarkFlagRequired("site")

	createCmd.Flags().String("username", "", "Enter username")
	_ = createCmd.MarkFlagRequired("username")

	createCmd.Flags().String("password", "", "Enter password")

	createCmd.Flags().String("auto", "", "Enter auto flag to generate password")

}
