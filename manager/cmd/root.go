package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manager",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Your code here
		fmt.Println("Hello from Cobra!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize flags here
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
