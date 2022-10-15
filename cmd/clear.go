package cmd

import (
	"fmt"
	"github.com/bakyazi/gopass/config"
	"github.com/bakyazi/gopass/sheetsapi"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear password database",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.GetConfig()

		db, err := sheetsapi.NewPasswordDB(cfg.CredentialFile, cfg.SheetId)
		if err != nil {
			panic(err)
		}
		err = db.Clear()
		if err != nil {
			fmt.Println("Failed to create and save passwords", err.Error())
			return
		}

		fmt.Println("Success!")

	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
