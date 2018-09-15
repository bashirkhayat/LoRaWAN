package main

import (
	// "fmt"
	"log"
	"net"
	// "os"

	"github.com/brocaar/loraserver/api/nc"
	"github.com/brocaar/loraserver/api/ns"
	// "github.com/brocaar/lorawan"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

import empty "github.com/golang/protobuf/ptypes/empty"

// NetworkControllerAPI implements the NetworkController service.
// https://github.com/brocaar/loraserver/blob/master/api/nc/nc.proto
type NetworkControllerAPI struct {
	nsClient ns.NetworkServerServiceClient
}

// HandleRXInfo handles the RX info calls from the network-server, which is
// called on each uplink.
// Each time this method is called, the code below will schedule a mac-command
func (n *NetworkControllerAPI) HandleUplinkMetaData(ctx context.Context, req *nc.HandleUplinkMetaDataRequest) (*empty.Empty, error) {
	out := new(empty.Empty)
	return out, nil
}

// HandleDataUpMACCommand handles the uplink mac-commands, which is called
// by the network-server each time an acknowledgement was received for a
// mac-command scheduled through the API (as is done in the above method).
// The code below will simply print the battery and margin values.
func (n *NetworkControllerAPI) HandleUplinkMACCommand(ctx context.Context, req *nc.HandleUplinkMACCommandRequest) (*empty.Empty, error) {

	out := new(empty.Empty)
	return out, nil
}

func main() {
	// allow insecure / non-tls connections
	grpcDialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	// connect to loraserver on localhost:8000
	nsConn, err := grpc.Dial("localhost:8000", grpcDialOpts...)
	if err != nil {
		log.Fatal(err)
	}

	// create the NetworkServer client
	nsClient := ns.NewNetworkServerServiceClient(nsConn)

	// create the NetWorkController service
	ncService := NetworkControllerAPI{
		nsClient: nsClient,
	}

	// setup the gRPC server
	grpcServerOpts := []grpc.ServerOption{}
	ncServer := grpc.NewServer(grpcServerOpts...)
	nc.RegisterNetworkControllerServiceServer(ncServer, &ncService)
	ln, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on 0.0.0.0:50051")
	// start the gRPC api server
	//nsConn.Close()
	log.Fatal(ncServer.Serve(ln))
}