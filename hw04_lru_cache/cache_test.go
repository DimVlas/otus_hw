package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)
		c.Set("100", 100)
		c.Set("200", 200)
		c.Clear()

		val, ok := c.Get("100")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("extrude", func(t *testing.T) {
		c := NewCache(5)
		c.Set("100", 100)
		c.Set("200", 200)
		c.Set("300", 300)
		c.Set("400", 400)
		c.Set("500", 500)
		// 500, 400, 300, 200, 100

		val, ok := c.Get("100") // 100, 500, 400, 300, 200
		require.True(t, ok)
		require.Equal(t, 100, val)

		ok = c.Set("1000", 1000) // 1000, 100, 500, 400, 300
		require.False(t, ok)

		val, ok = c.Get("200")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("300") // 100, 500, 400, 300, 200
		require.True(t, ok)
		require.Equal(t, 300, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
