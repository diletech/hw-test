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

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

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

func TestListOperations(t *testing.T) {
	// Создаем новый список
	myList := NewList()

	// Проверяем начальную длину списка
	if myList.Len() != 0 {
		t.Errorf("Expected initial length: 0, got: %d", myList.Len())
	}

	// Добавляем элементы в список
	item1 := myList.PushFront(1)
	item2 := myList.PushBack(2)
	item3 := myList.PushBack(3)

	// Проверяем длину после добавления
	if myList.Len() != 3 {
		t.Errorf("Expected length after adding elements: 3, got: %d", myList.Len())
	}

	// Проверяем Front и Back
	if myList.Front() != item1 || myList.Back() != item3 {
		t.Errorf("Front or Back not as expected")
	}

	// Перемещаем элемент в начало списка
	myList.MoveToFront(item2)

	// Проверяем, что элемент переместился в начало 2 -> 1 -> 3
	if myList.Front().Value != item2.Value || myList.Back().Value != item3.Value {
		t.Errorf("MoveToFront failed")
	}

	// Удаляем элемент первый элемент, который стал 2
	myList.Remove(myList.Front())

	// Проверяем, что элемент успешно удален 1 -> 3
	if myList.Len() != 2 || myList.Front().Value != item1.Value || myList.Back().Value != item3.Value {
		t.Errorf("Remove failed")
	}
}

func TestListTraversalOrder(t *testing.T) {
	myList := NewList()

	// Добавляем элементы
	item1 := myList.PushBack(1)
	item2 := myList.PushBack(2)
	myList.PushBack(3)
	myList.PushFront(4)

	// Используем методы для List
	myList.Remove(item1)
	myList.MoveToFront(item2)

	// Ожидаемый порядок при обходе списка: 2 -> 4 -> 3
	expectedOrder := []int{2, 4, 3}

	// Создаем слайс для хранения значений при обходе списка
	actualOrder := make([]int, 0, len(expectedOrder))

	// Обходим список и сохраняем значения в actualOrder
	current := myList.Front()
	for current != nil {
		actualOrder = append(actualOrder, current.Value.(int))
		current = current.Next
	}

	// Проверяем, что порядок значений совпадает с ожидаемым
	for i, v := range expectedOrder {
		if actualOrder[i] != v {
			t.Errorf("Traversal order mismatch at index %d. Expected: %d, got: %d", i, v, actualOrder[i])
		}
	}
}
