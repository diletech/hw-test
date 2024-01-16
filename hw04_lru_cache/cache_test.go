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
		c := NewCache(3)

		// Добавляем 4 элемента
		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3)
		c.Set("four", 4)

		// Первый элемент должен быть вытолкнут из-за размера очереди
		_, ok := c.Get("one")
		require.False(t, ok, "Expected element 'one' to be evicted")

		// очищаем для переиспользования
		c.Clear()

		// Добавляем 3 элемента
		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3)

		// Дергаем элементы, чтобы обновить время их использования
		c.Get("one")
		c.Get("three")

		// Добавляем 4й элемент
		c.Set("four", 4)

		// Самый давно использовавшийся элемент должен быть вытолкнут
		_, exists := c.Get("two")
		require.False(t, exists, "Expected element 'two' to be evicted")
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Run("multi thread", func(t *testing.T) {
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
	})
}
