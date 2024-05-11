package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("single front", func(t *testing.T) {
		l := NewList()

		const testVal string = "TestValue"
		l.PushFront(testVal)
		first := l.Front()
		last := l.Back()

		require.Equal(t, 1, l.Len())
		require.Equal(t, first, last)
		require.Equal(t, testVal, last.Value)

		l.Remove(l.Front())
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

	})

	t.Run("single back", func(t *testing.T) {
		l := NewList()

		const testVal string = "TestValue"
		l.PushBack(testVal)
		first := l.Front()
		last := l.Back()

		require.Equal(t, 1, l.Len())
		require.Equal(t, first, last)
		require.Equal(t, testVal, last.Value)

		l.Remove(l.Back())
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("move to front", func(t *testing.T) {
		l := NewList()

		l.PushBack(11)
		l.PushBack(22)
		l.PushBack(33)
		require.Equal(t, 3, l.Len())

		require.Equal(t, 11, l.Front().Value)
		require.Equal(t, 22, l.Back().Prev.Value)
		require.Equal(t, 33, l.Back().Value)

		l.MoveToFront(l.Front().Next)
		l.MoveToFront(l.Back())

		require.Equal(t, 33, l.Front().Value)
		require.Equal(t, 22, l.Front().Next.Value)
		require.Equal(t, 11, l.Back().Value)

	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]

		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		require.Equal(t, 20, middle.Value)

		l.Remove(middle) // [10, 30]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 30, l.Back().Value)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
