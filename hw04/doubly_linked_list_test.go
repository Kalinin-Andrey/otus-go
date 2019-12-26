package hw04

import (
	"sort"
	"testing"
)

var TestItemsInts []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var TestItemsStrs []string = []string{"A", "B", "C", "D"}

var list List = *New()

func TestPushFront(t *testing.T) {
	items1 := make([]int, len(TestItemsInts))
	copy(items1, TestItemsInts)
	sort.Sort(sort.Reverse(sort.IntSlice(items1)))

	for _, i := range items1 {
		list.PushFront(i)
	}

	if int(list.Len()) != len(TestItemsInts) {
		t.Errorf("Длина списка не совпадает с ожидаемым. Ожидали %v ; Получили %v", len(TestItemsInts), list.Len())
	}
	i := 0

	for item, isEnd := list.NextItem(); !isEnd; item, isEnd = list.NextItem(){

		if item.Value() != TestItemsInts[i] {
			t.Errorf("Значение элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsInts[i], item.Value())
		}

		if i-1 >= 0 && item.Prev().Value() != TestItemsInts[i-1] {
			t.Errorf("Значение предыдущего элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsInts[i-1], item.Prev().Value())
		}

		if i+1 < len(TestItemsInts) && item.Next().Value() != TestItemsInts[i+1] {
			t.Errorf("Значение следующего элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsInts[i+1], item.Next().Value())
		}
		i++
	}
}

func TestRemove(t *testing.T) {
	TestPushFront(t)
	rmCount := 5
	i := 0

	for currentItem := list.First(); i < rmCount; i++  {
		list.Remove(*currentItem)
		currentItem = currentItem.Next()
	}
	j := i

	if int(list.Len()) != len(TestItemsInts) - rmCount {
		t.Errorf("Длина списка не совпадает с ожидаемым. Ожидали %v ; Получили %v", len(TestItemsInts) - rmCount, list.Len())
	}

	for item, isEnd := list.NextItem(); !isEnd; item, isEnd = list.NextItem(){

		if item.Value() != TestItemsInts[i] {
			t.Errorf("Значение элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsInts[i], item.Value())
		}

		if i-1 >= 0 && item.Prev().Value() != TestItemsInts[i-1] {
			t.Errorf("Значение предыдущего элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsInts[i-1], item.Prev().Value())
		}

		if i+1 < len(TestItemsInts) && item.Next().Value() != TestItemsInts[i+1] {
			t.Errorf("Значение следующего элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsInts[i+1], item.Next().Value())
		}
		i++
	}

	for currentItem := list.First(); j < len(TestItemsInts); j++  {
		list.Remove(*currentItem)
		currentItem = currentItem.Next()
	}

	if int(list.Len()) != 0 {
		t.Errorf("Длина списка не совпадает с ожидаемым. Ожидали %v ; Получили %v", len(TestItemsInts) - rmCount, list.Len())
	}

	if list.First() != nil {
		t.Errorf("После удаления всех элементов список не пуст: первый элемент не nil")
	}

	if list.Last() != nil {
		t.Errorf("После удаления всех элементов список не пуст: последний элемент не nil")
	}

}

func TestPushBack(t *testing.T) {

	for _, i := range TestItemsStrs {
		list.PushBack(i)
	}

	if int(list.Len()) != len(TestItemsStrs) {
		t.Errorf("Длина списка не совпадает с ожидаемым. Ожидали %v ; Получили %v", len(TestItemsStrs), list.Len())
	}
	i := 0

	for item, isEnd := list.NextItem(); !isEnd; item, isEnd = list.NextItem(){

		if item.Value() != TestItemsStrs[i] {
			t.Errorf("Значение элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsStrs[i], item.Value())
		}

		if i-1 >= 0 && item.Prev().Value() != TestItemsStrs[i-1] {
			t.Errorf("Значение предыдущего элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsStrs[i-1], item.Prev().Value())
		}

		if i+1 < len(TestItemsStrs) && item.Next().Value() != TestItemsStrs[i+1] {
			t.Errorf("Значение следующего элемента не совпадает с ожидаемым. Ожидали %v ; Получили %v", TestItemsStrs[i+1], item.Next().Value())
		}
		i++
	}
}
