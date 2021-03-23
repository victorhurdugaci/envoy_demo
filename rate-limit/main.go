package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	ratelimit "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
)

type RateLimitServer struct {
	lastCall time.Time
}

func (rls *RateLimitServer) ShouldRateLimit(ctx context.Context, req *ratelimit.RateLimitRequest) (*ratelimit.RateLimitResponse, error) {
	log.Println("Handling Rate limit check request")

	now := time.Now()

	for _, descriptor := range req.Descriptors {
		for _, entry := range descriptor.Entries {
			if entry.Key != "request.method" {
				continue
			}

			rlKey := entry.Value + strconv.Itoa(int(now.Second()/10.0))

			rlVal := redisClient.Get(context.Background(), rlKey)
			rlIntVal := 0
			if rlVal != nil {
				rlIntVal, _ = rlVal.Int()
			}
			log.Printf("%s: %d\n", rlKey, rlIntVal)
			if rlIntVal > 5 {
				break // Reached limit
			}

			incrRes := redisClient.Incr(context.Background(), rlKey)
			if incrRes.Err() != nil {
				log.Println(incrRes.Err())
				break
			}
			expireRes := redisClient.Expire(context.Background(), rlKey, 30*time.Second)
			if expireRes.Err() != nil {
				log.Println(expireRes.Err())
				break
			}

			return &ratelimit.RateLimitResponse{
				OverallCode: ratelimit.RateLimitResponse_OK,
			}, nil
		}
	}

	return &ratelimit.RateLimitResponse{
		OverallCode: ratelimit.RateLimitResponse_OVER_LIMIT,
	}, nil
}

var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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
