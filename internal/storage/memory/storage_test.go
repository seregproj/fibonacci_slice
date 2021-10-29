package memorystorage_test

import (
	"context"
	"fmt"
	"testing"

	memorystorage "github.com/seregproj/fibonacci_slice/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("get item from empty map", func(t *testing.T) {
		ctx := context.Background()

		c := memorystorage.NewMemoryCache()
		v, err := c.Get(ctx, 5)
		var exp uint64
		require.Equal(t, exp, v)
		require.Equal(t, fmt.Errorf("not found"), err)
	})

	t.Run("set/get some items", func(t *testing.T) {
		ctx := context.Background()

		c := memorystorage.NewMemoryCache()
		c.Set(ctx, 1, 3)
		c.Set(ctx, 5, 7)

		// get a non-existent item if not empty map
		v, err := c.Get(ctx, 2)
		var exp uint64
		require.Equal(t, exp, v)
		require.Equal(t, fmt.Errorf("not found"), err)

		// get first valid item
		v, err = c.Get(ctx, 1)
		exp = 3
		require.Equal(t, exp, v)
		require.Nil(t, err)

		// get second valid item
		v, err = c.Get(ctx, 5)
		exp = 7
		require.Equal(t, exp, v)
		require.Nil(t, err)
	})
}

func TestSet(t *testing.T) {
	t.Run("set existing item", func(t *testing.T) {
		ctx := context.Background()
		c := memorystorage.NewMemoryCache()
		c.Set(ctx, 1, 7)
		c.Set(ctx, 1, 7)

		v, err := c.Get(ctx, 1)
		var exp uint64 = 7
		require.Equal(t, exp, v)
		require.Nil(t, err)
	})

	t.Run("set another val to an item, check no rewriting in mem", func(t *testing.T) {
		ctx := context.Background()

		c := memorystorage.NewMemoryCache()
		c.Set(ctx, 1, 7)
		c.Set(ctx, 1, 8)

		v, err := c.Get(ctx, 1)
		var exp uint64 = 7
		require.Equal(t, exp, v)
		require.Nil(t, err)
	})
}
