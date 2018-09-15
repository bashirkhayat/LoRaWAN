package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"encoding/hex"
	"github.com/brocaar/loraserver/api/nc"
	"github.com/brocaar/loraserver/api/ns"
	"github.com/brocaar/lorawan"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

import empty "github.com/golang/protobuf/ptypes/empty"

var recievedDevEUI lorawan.EUI64

// NetworkControllerAPI implements the NetworkController service.
// https://github.com/brocaar/loraserver/blob/master/api/nc/nc.proto
type NetworkControllerAPI struct {
	nsClient ns.NetworkServerServiceClient
}


// UnmarshalText implements encoding.TextUnmarshaler.
func UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	copy(recievedDevEUI[:], b)
	log.Printf("%s",recievedDevEUI)	
	return nil
}



// HandleRXInfo handles the RX info calls from the network-server, which is
// called on each uplink.
// Each time this method is called, the code below will schedule a mac-command
// for this node asking for the battery status and link margin
// (DevStatusReq mac-command).
func (n *NetworkControllerAPI) HandleUplinkMetaData(ctx context.Context, req *nc.HandleUplinkMetaDataRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.GetDevEui())
	//log.Printf("HandleRXInfo method called (DevEUI: %s)", devEUI)
	log.Printf("%s %s", devEUI, recievedDevEUI)
	if( devEUI != recievedDevEUI ) {
		out := new(empty.Empty)
		return out, nil	
	}
	// construct the mac-command
	mac := lorawan.MACCommand{
		CID:     lorawan.DevStatusReq,
		Payload: nil, // this mac-command does not have a payload
	}
	macBytes, err := mac.MarshalBinary()
	if err != nil {
		log.Printf("marshal mac-command error: %s", err)
		return nil, err
	}
	// send the mac-command to the network-server
	_, err = n.nsClient.CreateMACCommandQueueItem(ctx, &ns.CreateMACCommandQueueItemRequest{
		DevEui: req.GetDevEui(),
		Cid:    uint32(lorawan.DevStatusReq),
		Commands: [][]byte{
			macBytes,
		},
	})
	if err != nil {
		return nil, err
	}
	log.Printf("requested the status of the end-device %s", devEUI)	
	out := new(empty.Empty)
	return out, nil
}

// HandleDataUpMACCommand handles the uplink mac-commands, which is called
// by the network-server each time an acknowledgement was received for a
// mac-command scheduled through the API (as is done in the above method).
// The code below will simply print the battery and margin values.
func (n *NetworkControllerAPI) HandleUplinkMACCommand(ctx context.Context, req *nc.HandleUplinkMACCommandRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.GetDevEui())
	// log.Printf("HandleDataUpMACCommand method called (DevEUI: %s)", devEUI)
	// make sure the CID (command identifier) is of type DevStatusAns.
	if req.Cid != uint32(lorawan.DevStatusAns) {
		return nil, fmt.Errorf("unexpected CID: %s", lorawan.CID(req.Cid))
	}

	for _, macBytes := range req.Commands {
		var mac lorawan.MACCommand
		if err := mac.UnmarshalBinary(true, macBytes); err != nil { // true since this is an uplink mac-command
			log.Printf("unmarshal mac-command error: %s", err)
			return nil, err
		}

		pl, ok := mac.Payload.(*lorawan.DevStatusAnsPayload)
		if !ok {
			return nil, fmt.Errorf("expected *lorawan.DevStatusAnsPayload, got %T", mac.Payload)
		}

		log.Printf("battery: %d, margin: %d (DevEUI: %s)", 254 - pl.Battery, pl.Margin, devEUI)
		os.Exit(0)
	}
	out := new(empty.Empty)
	return out, nil
}

func main() {
	UnmarshalText([]byte(os.Args[1]))
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