package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	//nolint:depguard
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
		c := NewCache(3)

		wasInCache := c.Set("key1", 1)
		require.False(t, wasInCache)

		wasInCache = c.Set("key2", 2)
		require.False(t, wasInCache)

		wasInCache = c.Set("key3", 3)
		require.False(t, wasInCache)

		// first element should be removed
		wasInCache = c.Set("key4", 4)
		require.False(t, wasInCache)

		val, ok := c.Get("key1")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("key4")
		require.True(t, ok)
		require.Equal(t, 4, val)

		wasInCache = c.Set("key2", 22)
		require.True(t, wasInCache)

		wasInCache = c.Set("key3", 33)
		require.True(t, wasInCache)

		// least recently used element should be removed
		wasInCache = c.Set("key5", 5)
		require.False(t, wasInCache)

		val, ok = c.Get("key4")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("key1", 1)
		require.False(t, wasInCache)

		wasInCache = c.Set("key2", 2)
		require.False(t, wasInCache)

		c.Clear()

		val, ok := c.Get("key1")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("key2")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(*testing.T) {
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
