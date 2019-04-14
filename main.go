package main

import (
	"log"
	"net"
	"os"

	"github.com/vovanushka/port-service/api"
	"github.com/vovanushka/port-service/repo"
	"google.golang.org/grpc"
)

// main start a gRPC server and waits for connection
func main() {
	// init new mgo session
	ms, err := repo.NewSession(getEnv("MONGO_URL", "127.0.0.1:27017"))
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	// construct repo for port entity
	portRepo := repo.NewPortRepo(ms.Copy(), getEnv("MONGO_DB_NAME", "port_service"), "port")
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", getEnv("PORT_SERVICE_ADDR", "127.0.0.1:7777"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a port server instance
	s := api.NewServer(portRepo)
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach port service to the server
	api.RegisterPortServer(grpcServer, s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// helper function to make easier setting default env parameters
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
