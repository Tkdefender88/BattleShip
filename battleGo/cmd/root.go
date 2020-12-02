package cmd

import (
	"fmt"
	"os"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/server"
	"github.com/spf13/cobra"
)

var (
	Address string
	RootCmd = &cobra.Command{
		Use: "battleShip",
		Short: "battleShip runs a simple battle ship server for webscience class",
		Run: func(cmd *cobra.Command, args []string){
			fmt.Println("Root Command ran")
			server.Start(server.StartServer, Address)
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
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&Address, "port", "p", "30124" ,"Port to open dev server on")
}

func initConfig() {
    
}
