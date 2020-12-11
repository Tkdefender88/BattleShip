package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/server"
)

func init() {
	RootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs the server in dev mode",
	Long:  `Runs the server in dev mode where tls certs aren't in use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("dev mode called")
		if err := server.Start(server.StartDevServer, Address); err != nil {
			return err
		}
		return nil
	},
}
