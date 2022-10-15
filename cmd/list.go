package cmd

import (
	"fmt"
	"github.com/bakyazi/gopass/config"
	"github.com/bakyazi/gopass/sheetsapi"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all passwords from DB",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		db, err := sheetsapi.NewPasswordDB(cfg.CredentialFile, cfg.SheetId)
		if err != nil {
			panic(err)
		}
		passw, err := db.GetPasswords(1, 5)
		if err != nil {
			fmt.Println("Failed to list all passwords", err.Error())
			return
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 8, 8, 2, '\t', 0)

		defer w.Flush()
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "Row", "Site", "Username", "Password")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "----", "----", "----", "----")
		for _, p := range passw {
			fmt.Fprintf(w, "\n %d\t%s\t%s\t%s\t", p.Row, p.Site, p.Username, p.Password)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
