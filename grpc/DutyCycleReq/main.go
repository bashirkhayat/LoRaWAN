package main

import (
	"fmt"
	"log"
	"net"

	"github.com/brocaar/loraserver/api/nc"
	"github.com/brocaar/loraserver/api/ns"
	"github.com/brocaar/lorawan"
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
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.GetDevEui())
	// log.Printf("HandleRXInfo method called (DevEUI: %s)", devEUI)
	p := &lorawan.DutyCycleReqPayload{255}
	// construct the mac-command
	mac := lorawan.MACCommand{
		CID:     lorawan.DutyCycleReq,
		Payload: p,
	}
	macBytes, err := mac.MarshalBinary()
	if err != nil {
		log.Printf("marshal mac-command error: %s")
		return nil, err
	}
	// send the mac-command to the network-server
	_, err = n.nsClient.CreateMACCommandQueueItem(ctx, &ns.CreateMACCommandQueueItemRequest{
		DevEui: req.GetDevEui(),
		Cid:    uint32(lorawan.DutyCycleReq),
		Commands: [][]byte{
			macBytes,
		},
	})
	if err != nil {
		return nil, err
	}
	out := new(empty.Empty)
	log.Printf("Sending DutyCycleReq to the end-device %s", devEUI)
	// log.Printf("End-device %s is turned off, in order to turn it back again please push the black button on the device itself", devEUI)
	return out, nil
}

// HandleDataUpMACCommand handles the uplink mac-commands, which is called
// by the network-server each time an acknowledgement was received for a
// mac-command scheduled through the API (as is done in the above method).
func (n *NetworkControllerAPI) HandleUplinkMACCommand(ctx context.Context, req *nc.HandleUplinkMACCommandRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.GetDevEui())

	log.Printf("HandleDataUpMACCommand method called (DevEUI: %s)", devEUI)

	// make sure the CID (command identifier) is of type DevStatusAns.
	if req.Cid != uint32(lorawan.DutyCycleAns) {
		return nil, fmt.Errorf("unexpected CID: %s", lorawan.CID(req.Cid))
	}

	for _, macBytes := range req.Commands {
		var mac lorawan.MACCommand
		if err := mac.UnmarshalBinary(true, macBytes); err != nil { // true since this is an uplink mac-command
			log.Printf("unmarshal mac-command error: %s", err)
			return nil, err
		}
	}
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
	log.Fatal(ncServer.Serve(ln))

}