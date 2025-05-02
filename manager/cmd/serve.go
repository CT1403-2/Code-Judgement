package main

import (
	"fmt"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"manger/internal/manager"
	"net"
	"net/http"
	"os"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  `Starts the server and begins listening for requests`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		fmt.Println("Server started")
		err := serve(port)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add flags specific to this command
	serveCmd.Flags().StringP("port", "p", os.Getenv("SERVER_PORT"), "Port to run the server on")
}

func serve(port string) error {
	fmt.Println()
	go serveStaticFiles(port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	m, err := manager.NewManager()
	if err != nil {
		return err
	}

	proto.RegisterManagerServer(grpcServer, m)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return err
}

func serveStaticFiles(port string) {
	fs := http.FileServer(http.Dir("../../front/")) // Your static directory
	http.Handle("/", fs)
	addr := ":" + port
	log.Printf("Serving static files on %v\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Static server failed: %v", err)
	}
}
