package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	ratelimit "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
)

type RateLimitServer struct {
	lastCall time.Time
}

func (rls *RateLimitServer) ShouldRateLimit(ctx context.Context, req *ratelimit.RateLimitRequest) (*ratelimit.RateLimitResponse, error) {
	log.Printf("Handling Rate limit check request")

	now := time.Now()
	if rls.lastCall.Add(5 * time.Second).After(now) {
		return &ratelimit.RateLimitResponse{
			OverallCode: ratelimit.RateLimitResponse_OVER_LIMIT,
		}, nil
	}

	rls.lastCall = now
	return &ratelimit.RateLimitResponse{
		OverallCode: ratelimit.RateLimitResponse_OK,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

	s := grpc.NewServer(opts...)

	ratelimit.RegisterRateLimitServiceServer(s, &RateLimitServer{})

	fmt.Printf("Rate limit running on port 9091")
	s.Serve(lis)
}
