package app

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidIndexNumbers = errors.New("invalid index numbers")

type FibonacciCalc struct {
	cache  Cache
	logger Logger
}

type Logger interface {
	Info(text string)
	Error(text string)
}

type Cache interface {
	Get(ctx context.Context, k uint64) (uint64, error)
	Set(ctx context.Context, k uint64, v uint64) error
}

func NewFibonacciCalc(cache Cache, logger Logger) *FibonacciCalc {
	return &FibonacciCalc{
		cache:  cache,
		logger: logger,
	}
}

func (fc *FibonacciCalc) GetSlice(ctx context.Context, x, y uint64) ([]uint64, error) {
	var res []uint64
	if x > y {
		return res, ErrInvalidIndexNumbers
	}

	res = make([]uint64, 0, y-x+1)
	for i := x; i <= y; i++ {
		val, err := fc.calc(ctx, i)
		if err != nil {
			fc.logger.Error(fmt.Sprintf("cant calc for: %v", i))

			return res, err
		}

		res = append(res, val)
	}

	fc.logger.Info(fmt.Sprintf("result for %v-%v: %v", x, y, res))

	return res, nil
}

func (fc *FibonacciCalc) calc(ctx context.Context, n uint64) (uint64, error) {
	if n == 0 {
		return 0, nil
	}

	if n == 1 {
		return 1, nil
	}

	val, err := fc.cache.Get(ctx, n)
	if err == nil {
		return val, nil
	}

	fc.logger.Info(fmt.Sprintf("%v item is not present in cache, need to calc", n))

	prev, err := fc.calc(ctx, n-1)
	if err != nil {
		return 0, err
	}

	beforePrev, err := fc.calc(ctx, n-2)
	if err != nil {
		return 0, err
	}

	val = prev + beforePrev
	err = fc.cache.Set(ctx, n, val)
	if err != nil {
		return 0, err
	}

	return val, nil
}
