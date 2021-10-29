package grpcservice

import (
	"context"
	"errors"

	pb "github.com/seregproj/fibonacci_slice/api/proto"
	"github.com/seregproj/fibonacci_slice/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application interface {
	GetSlice(ctx context.Context, x, y uint64) ([]uint64, error)
}

type FibonacciService struct {
	app Application
	pb.UnimplementedFibonacciServiceServer
}

func NewFibonacciService(app Application) *FibonacciService {
	return &FibonacciService{app: app}
}

func (f *FibonacciService) GetSlice(ctx context.Context, req *pb.GetSliceRequest) (*pb.GetSliceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid index numbers")
	}

	l, err := f.app.GetSlice(ctx, req.GetX(), req.GetY())
	if err != nil {
		if errors.Is(err, app.ErrInvalidIndexNumbers) {
			return nil, status.Error(codes.InvalidArgument, "invalid index numbers")
		}

		return nil, status.Error(codes.Internal, "something goes wrong")
	}

	return &pb.GetSliceResponse{Items: l}, nil
}
