package main

import (
	"fmt"
	"log"
	"net"
	"os"

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
	r := lorawan.Redundancy{ChMaskCntl: 1, NbRep: 1}
	chMask := [16]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}
	p := &lorawan.LinkADRReqPayload{DataRate: 0, TXPower: 14, ChMask: chMask, Redundancy: r}
	// construct the mac-command
	mac := lorawan.MACCommand{
		CID:     lorawan.LinkADRReq,
		Payload: p, // this mac-command does not have a payload
	}
	macBytes, err := mac.MarshalBinary()
	if err != nil {
		log.Printf("marshal mac-command error: %s", err)
		return nil, err
	}
	// send the mac-command to the network-server
	_, err = n.nsClient.CreateMACCommandQueueItem(ctx, &ns.CreateMACCommandQueueItemRequest{
		DevEui: req.GetDevEui(),
		Cid:    uint32(lorawan.LinkADRReq),
		Commands: [][]byte{
			macBytes,
		},
	})
	if err != nil {
		return nil, err
	}
	out := new(empty.Empty)
	log.Printf("Checking the parameters of the end-device %s", devEUI)
	return out, nil
}

// HandleDataUpMACCommand handles the uplink mac-commands, which is called
// by the network-server each time an acknowledgement was received for a
// mac-command scheduled through the API (as is done in the above method).
// The code below will simply print the battery and margin values.
func (n *NetworkControllerAPI) HandleUplinkMACCommand(ctx context.Context, req *nc.HandleUplinkMACCommandRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.GetDevEui())

	log.Printf("HandleDataUpMACCommand method called (DevEUI: %s)", devEUI)

	// make sure the CID (command identifier) is of type DevStatusAns.
	if req.Cid != uint32(lorawan.LinkADRAns) {
		return nil, fmt.Errorf("unexpected CID: %s", lorawan.CID(req.Cid))
	}

	for _, macBytes := range req.Commands {
		var mac lorawan.MACCommand
		if err := mac.UnmarshalBinary(true, macBytes); err != nil { // true since this is an uplink mac-command
			log.Printf("unmarshal mac-command error: %s", err)
			return nil, err
		}

		pl, ok := mac.Payload.(*lorawan.LinkADRAnsPayload)
		if !ok {
			return nil, fmt.Errorf("expected *lorawan.LinkADRAnsPayload, got %T", mac.Payload)
		}

		log.Printf("ChannelMaskACK: %t, DataRateACK: %t, PowerACK:%t,  (DevEUI: %s)", pl.ChannelMaskACK, pl.DataRateACK, pl.PowerACK, devEUI)
		os.Exit(0)
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
	//nsConn.Close()
	log.Fatal(ncServer.Serve(ln))
}