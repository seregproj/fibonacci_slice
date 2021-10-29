//go:build bench
// +build bench

package redis_test

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	pb "github.com/seregproj/fibonacci_slice/api/proto"
)

func TestGetSliceTime(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.Dial(net.JoinHostPort(os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT")), grpc.WithInsecure())
	require.NoError(t, err)

	rdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	require.NoError(t, err)
	redisClient := redis.NewClient(&redis.Options{DB: rdb, Addr: net.JoinHostPort(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))})
	redisClient.FlushDB(ctx)

	fibonacciClient := pb.NewFibonacciServiceClient(conn)

	bench := func(b *testing.B) {
		b.ResetTimer()
		_, err := fibonacciClient.GetSlice(ctx, &pb.GetSliceRequest{X: 1, Y: 100})
		require.NoError(t, err)

		b.StopTimer()

		redisClient.FlushDB(ctx)
	}

	result := testing.Benchmark(bench)
	t.Logf("time used: %s", result.T)

	require.Less(t, int64(result.T), int64(300*time.Millisecond), "calc is too slow")
}
