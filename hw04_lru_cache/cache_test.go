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
		c, err := NewCache(10)
		require.Nil(t, err)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c, err := NewCache(5)
		require.Nil(t, err)

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
		c, err := NewCache(3)
		require.Nil(t, err)

		c.Set("1", 100)
		c.Set("2", 200)
		c.Set("3", 300)

		c.Set("4", 400)

		val, ok := c.Get("1")
		require.False(t, ok)
		require.Nil(t, val)

		c.Set("5", 500)

		val, ok = c.Get("2")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic last usage", func(t *testing.T) {
		c, err := NewCache(3)
		require.Nil(t, err)

		c.Set("1", 100)
		c.Set("2", 200)
		c.Set("3", 300)

		c.Set("1", 111)
		c.Set("4", 400)

		val, ok := c.Get("2")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("3")
		require.True(t, ok)
		require.Equal(t, val, 300)

		c.Set("5", 500)

		val, ok = c.Get("1")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear func", func(t *testing.T) {
		c, err := NewCache(3)
		require.Nil(t, err)

		c.Set("1", 100)
		c.Set("2", 200)
		c.Set("3", 300)

		c.Clear()

		val, ok := c.Get("1")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("2")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("3")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("negative capacity", func(t *testing.T) {
		c, err := NewCache(-10)
		require.Nil(t, c)
		require.NotNil(t, err)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c, err := NewCache(10)
	require.Nil(t, err)

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
