package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointcfg "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	"github.com/golang/protobuf/ptypes/any"
)

type EndpointServer struct {
}

func (es *EndpointServer) StreamEndpoints(req endpoint.EndpointDiscoveryService_StreamEndpointsServer) error {
	// This should be dynamic
	echoEndpoints := &endpointcfg.ClusterLoadAssignment{
		ClusterName: "echo",
		Endpoints: []*endpointcfg.LocalityLbEndpoints{{
			LbEndpoints: []*endpointcfg.LbEndpoint{{
				HostIdentifier: &endpointcfg.LbEndpoint_Endpoint{
					Endpoint: &endpointcfg.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Address: "192.168.65.2", // Update this line with the ip address from docker (get it with "host host.docker.internal" ran inside the container)
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: 8081,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}

	endpointsMsg, _ := anypb.New(echoEndpoints)

	for {
		_, err := req.Recv()
		if err == io.EOF {
			fmt.Println("Done")
			return nil
		}

		req.Send(&discovery.DiscoveryResponse{
			TypeUrl: endpointsMsg.TypeUrl,
			Resources: []*any.Any{
				endpointsMsg,
			},
		})

		time.Sleep(5 * time.Second)
	}
}

func (es *EndpointServer) DeltaEndpoints(endpoint.EndpointDiscoveryService_DeltaEndpointsServer) error {
	fmt.Println("DeltaEndpoints")
	return errors.New("not implemented")
}

func (es *EndpointServer) FetchEndpoints(ctx context.Context, req *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	fmt.Println("FetchEndpoints")
	return nil, errors.New("not implemented")
}

func main() {
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

	s := grpc.NewServer(opts...)
	endpoint.RegisterEndpointDiscoveryServiceServer(s, &EndpointServer{})
	fmt.Println("Endpoint discovery running on port 9092")
	s.Serve(lis)
}
