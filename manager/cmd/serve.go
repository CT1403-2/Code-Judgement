package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"manger/internal/manager"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  `Starts the server and begins listening for requests`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Server started")
		err := serve()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add flags specific to this command
	serveCmd.Flags().IntP("port", "p", 8080, "Port to run the server on")
}

func serve() error {
	m, err := manager.NewManager()
	if err != nil {
		return err
	}
	err = m.Start()
	return err
}
