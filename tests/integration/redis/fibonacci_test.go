//go:build integration
// +build integration

package redis_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	pb "github.com/seregproj/fibonacci_slice/api/proto"
)

type GetSliceServiceSuite struct {
	suite.Suite
	fibonacciClient pb.FibonacciServiceClient
	redisClient     *redis.Client
	ctx             context.Context
}

func (s *GetSliceServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	conn, err := grpc.Dial(net.JoinHostPort(os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT")), grpc.WithInsecure())
	s.Require().NoError(err)

	rdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(&redis.Options{DB: rdb, Addr: net.JoinHostPort(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))})
	s.redisClient.FlushDB(s.ctx)

	s.fibonacciClient = pb.NewFibonacciServiceClient(conn)
}

func (s *GetSliceServiceSuite) TestValidGetSlice() {
	tests := []struct {
		x         uint64
		y         uint64
		exp       []uint64
		redisData map[string]string
	}{
		{
			x:         0,
			y:         1,
			exp:       []uint64{0, 1},
			redisData: map[string]string{},
		},
		{
			x:   0,
			y:   2,
			exp: []uint64{0, 1, 1},
			redisData: map[string]string{
				"2": "1",
			},
		},
		{
			x:   0,
			y:   3,
			exp: []uint64{0, 1, 1, 2},
			redisData: map[string]string{
				"2": "1",
			},
		},
		{
			x:   0,
			y:   4,
			exp: []uint64{0, 1, 1, 2, 3},
			redisData: map[string]string{
				"2": "1",
				"3": "2",
			},
		},
		{
			x:         1,
			y:         2,
			exp:       []uint64{1, 1},
			redisData: map[string]string{},
		},
		{
			x:   1,
			y:   3,
			exp: []uint64{1, 1, 2},
			redisData: map[string]string{
				"2": "1",
				"3": "2",
			},
		},
		{
			x:   1,
			y:   4,
			exp: []uint64{1, 1, 2, 3},
			redisData: map[string]string{
				"2": "1",
				"3": "2",
				"4": "3",
			},
		},
		{
			x:   2,
			y:   3,
			exp: []uint64{1, 2},
			redisData: map[string]string{
				"2": "1",
				"3": "2",
			},
		},
		{
			x:   2,
			y:   4,
			exp: []uint64{1, 2, 3},
			redisData: map[string]string{
				"2": "1",
				"3": "2",
				"4": "3",
			},
		},
		{
			x:   3,
			y:   4,
			exp: []uint64{2, 3},
			redisData: map[string]string{
				"3": "2",
				"4": "3",
			},
		},
		{
			x:   3,
			y:   5,
			exp: []uint64{2, 3, 5},
			redisData: map[string]string{
				"3": "2",
				"4": "3",
				"5": "5",
			},
		},
	}

	for i, tt := range tests {
		s.Run(fmt.Sprintf("case %d", i), func() {
			resp, err := s.fibonacciClient.GetSlice(s.ctx, &pb.GetSliceRequest{X: tt.x, Y: tt.y})
			s.Require().NoError(err)
			s.Require().Equal(tt.exp, resp.GetItems())

			for k, v := range tt.redisData {
				res, err := s.redisClient.Get(s.ctx, k).Result()
				s.Require().NoError(err)
				s.Require().Equal(v, res)
			}

			s.redisClient.FlushDB(s.ctx)
		})
	}
}

func TestGetSliceSuite(t *testing.T) {
	suite.Run(t, new(GetSliceServiceSuite))
}
