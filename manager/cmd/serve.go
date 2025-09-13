package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/CT1403-2/Code-Judgement/manager/internal/manager"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  `Starts the server and begins listening for requests`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		fmt.Println("Server started")
		err := server(port)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	// Add flags specific to this command
	serveCmd.Flags().StringP("port", "p", "", "Port to run the server on")
}

func server(port string) error {
	addr := ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	m := cmux.New(lis)

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	grpcServer := grpc.NewServer()
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithAllowedRequestHeaders([]string{"x-grpc-web", "content-type"}),
	)
	man, err := manager.NewManager()
	if err != nil {
		return err
	}
	proto.RegisterManagerServer(grpcServer, man)

	httpServer := &http.Server{
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(req) || wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}
			filePath := filepath.Join("build/browser", req.URL.Path)
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(resp, req, filePath)
				return
			}
			http.ServeFile(resp, req, "build/browser/index.html")
		}),
	}

	go func() {
		log.Println("Starting grpc on " + port)
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	go func() {
		log.Println("Serving HTTP on", port)
		if err := httpServer.Serve(httpL); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start cmux
	log.Println("Starting multiplexer on", port)
	if err := m.Serve(); err != nil {
		log.Fatalf("cmux server error: %v", err)
	}
	return nil
}
