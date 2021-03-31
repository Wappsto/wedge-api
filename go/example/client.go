package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Wappsto/wedge-api/go/wedge"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

var (
	NodeIdentity = &pb.NodeIdentity{
		Id: "node-1234567890",
	}

	myDevices = []*pb.Device{
		{
			Id:   1,
			Name: "PorcDevice",
			Value: []*pb.Value{
				{
					Id:         1,
					Name:       "living_room_temperature",
					Permission: "rw",
					Number: &pb.Number{
						Min:  -20,
						Max:  60,
						Step: 1,
					},
					State: []*pb.State{
						{
							Id:   1,
							Data: "18",
							Type: pb.Type_Report,
						},
						{
							Id:   2,
							Data: "2",
							Type: pb.Type_Control,
						},
					},
				},
			},
		},
		{
			Id:   6,
			Name: "ToBeDeleted",
			Value: []*pb.Value{
				{
					Id:         1,
					Name:       "do_not_exist",
					Permission: "r",
					Status:     pb.Status_ok,
					String_: &pb.String{
						Max: 1024,
					},
					State: []*pb.State{
						{
							Id:   1,
							Data: "2",
							Type: pb.Type_Report,
						},
					},
				},
			},
		},
	}
)

func NodeModel() *pb.SetModelRequest {

	req := &pb.SetModelRequest{
		Model: &pb.Model{
			Node:   NodeIdentity,
			Device: myDevices,
		},
	}
	return req
}

func CurrentState() *pb.SetStateRequest {

	req := &pb.SetStateRequest{
		Node: NodeIdentity,
		State: &pb.State{
			Id:   1,
			Data: "7",
		},
		ValueId:  1,
		DeviceId: 1,
	}

	return req
}

func AskForControl(client pb.WedgeClient, in *pb.GetControlRequest) {

	for {
		fmt.Println("Sending request for control.")
		resp, err := client.GetControl(context.Background(), in)
		if err != nil {
			fmt.Println("Error, maybe GW is down?")
			time.Sleep(time.Second * 3)
			continue
		}

		b, err := json.Marshal(resp)
		if err != nil {
			// TODO handling errors.
			fmt.Println("Marshal Error")
		}
		fmt.Printf("Control received: %s\n", string(b))

		if resp.Update != nil {
			for _, dev := range myDevices {
				if resp.Update.DeviceId == dev.Id {
					for _, val := range dev.Value {
						if resp.Update.ValueId == val.Id {
							for _, state := range val.State {
								if resp.Update.State.Id == state.Id {
									fmt.Printf("Control state: updating data from: '%s' to '%s'\n", state.Data, resp.Update.State.Data)
									state.Data = resp.Update.State.Data
								}
							} // state
						}
					} // value
				}
			} // device
		} // search for state to Update.
	} // for loop
}

func DriverLoop(client pb.WedgeClient) {

	var count int

	for {
		time.Sleep(time.Second * 15)

		newState := CurrentState()
		newState.State.Data = strconv.Itoa(count)
		count++
		resp, err := client.SetState(context.Background(), newState)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		b, _ := json.Marshal(resp)
		fmt.Printf("Set State response %s\n", string(b))
	}
}

func main() {
	// dial server
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error: can not connect with server %v", err)
		return
	}

	// create stream
	client := pb.NewWedgeClient(conn)

	in := &pb.GetControlRequest{
		Node: NodeIdentity,
	}

	model := NodeModel()
	resp, err := client.SetModel(context.Background(), model)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	b, _ := json.Marshal(resp)
	fmt.Printf("Set Model response %s\n", string(b))

	go AskForControl(client, in)
	go DriverLoop(client)

	//we will wait indefinatelly.
	select {}
}
