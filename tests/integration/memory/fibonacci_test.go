//go:build integration
// +build integration

package memory_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/seregproj/fibonacci_slice/api/proto"
)

type GetSliceServiceSuite struct {
	suite.Suite
	fibonacciClient pb.FibonacciServiceClient
	ctx             context.Context
}

func (s *GetSliceServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	conn, err := grpc.Dial(net.JoinHostPort(os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT")), grpc.WithInsecure())
	s.Require().NoError(err)

	s.fibonacciClient = pb.NewFibonacciServiceClient(conn)
}

func (s *GetSliceServiceSuite) TestValidGetSlice() {
	tests := []struct {
		x   uint64
		y   uint64
		exp []uint64
	}{
		{
			x:   0,
			y:   0,
			exp: []uint64{0},
		},
		{
			x:   1,
			y:   1,
			exp: []uint64{1},
		},
		{
			x:   2,
			y:   2,
			exp: []uint64{1},
		},
		{
			x:   0,
			y:   1,
			exp: []uint64{0, 1},
		},
		{
			x:   0,
			y:   2,
			exp: []uint64{0, 1, 1},
		},
		{
			x:   0,
			y:   3,
			exp: []uint64{0, 1, 1, 2},
		},
		{
			x:   0,
			y:   4,
			exp: []uint64{0, 1, 1, 2, 3},
		},
		{
			x:   0,
			y:   5,
			exp: []uint64{0, 1, 1, 2, 3, 5},
		},
		{
			x:   1,
			y:   2,
			exp: []uint64{1, 1},
		},
		{
			x:   1,
			y:   3,
			exp: []uint64{1, 1, 2},
		},
		{
			x:   1,
			y:   4,
			exp: []uint64{1, 1, 2, 3},
		},
		{
			x:   2,
			y:   3,
			exp: []uint64{1, 2},
		},
		{
			x:   2,
			y:   4,
			exp: []uint64{1, 2, 3},
		},
		{
			x:   3,
			y:   4,
			exp: []uint64{2, 3},
		},
		{
			x:   3,
			y:   5,
			exp: []uint64{2, 3, 5},
		},
		{
			x:   3,
			y:   5,
			exp: []uint64{2, 3, 5},
		},
	}

	for i, tt := range tests {
		s.Run(fmt.Sprintf("case %d", i), func() {
			resp, err := s.fibonacciClient.GetSlice(s.ctx, &pb.GetSliceRequest{X: tt.x, Y: tt.y})
			s.Require().NoError(err)
			s.Require().Equal(tt.exp, resp.GetItems())
		})
	}
}

func (s *GetSliceServiceSuite) TestInvalidArgsGetSlice() {
	var exp []uint64

	tests := []struct {
		x uint64
		y uint64
	}{
		{
			x: 1,
			y: 0,
		},
		{
			x: 2,
			y: 1,
		},
		{
			x: 6,
			y: 5,
		},
		{
			x: 5,
			y: 3,
		},
	}

	for i, tt := range tests {
		s.Run(fmt.Sprintf("case %d", i), func() {
			resp, err := s.fibonacciClient.GetSlice(s.ctx, &pb.GetSliceRequest{X: tt.x, Y: tt.y})
			s.Require().Error(err)
			s.Require().Equal(exp, resp.GetItems())
			st, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(codes.InvalidArgument, st.Code())
			s.Require().Equal("invalid index numbers", st.Message())
		})
	}
}

func TestGetSliceSuite(t *testing.T) {
	suite.Run(t, new(GetSliceServiceSuite))
}
