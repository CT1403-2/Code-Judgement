package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  `Starts the server and begins listening for requests`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Server started")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add flags specific to this command
	serveCmd.Flags().IntP("port", "p", 8080, "Port to run the server on")
}
