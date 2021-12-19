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

	t.Run("delete items due to queue size", func(t *testing.T) {
		c := NewCache(2)

		c.Set("frst", 10)
		c.Set("sec", 20)

		c.Set("third", 30)
		val, ok := c.Get("frst")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("element that was used the most long ago will be pushed out", func(t *testing.T) {
		c := NewCache(3)

		c.Set("frst", 10)
		c.Set("sec", 20)
		c.Set("third", 30)

		c.Set("frst", 1)
		c.Get("third")

		c.Set("fourth", 40)

		val, ok := c.Get("sec")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("complex scenario", func(t *testing.T) {
		c := NewCache(3)

		c.Set("frst", 10)
		c.Set("sec", 20)
		c.Set("third", 30)

		c.Set("fourth", 40) // "fourth" should throw out last one (which is "frst")
		val, ok := c.Get("frst")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("sec") // "sec" has been moved to first position, now last should be "third"
		require.True(t, ok)
		require.Equal(t, 20, val)

		c.Set("five", 50) // "five" should throw out last one (which is "third")
		val, ok = c.Get("third")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
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
