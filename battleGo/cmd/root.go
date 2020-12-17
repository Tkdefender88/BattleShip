package cmd

import (
	"fmt"
	"os"

	"github.com/Tkdefender88/BattleShip/server"
	"github.com/spf13/cobra"
)

var (
	Address string
	RootCmd = &cobra.Command{
		Use: "battleShip",
		Short: "battleShip runs a simple battle ship server for webscience class",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Root Command ran")
			err := server.Start(server.StartServer, Address)
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Address, "port", "p", "30124" ,"Port to open dev server on")
}

