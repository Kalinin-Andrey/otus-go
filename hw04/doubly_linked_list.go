package hw04

import (

)

type List struct {
	len     uint
	first   *Item
	last    *Item
	current *Item
}

type Item struct {
	val         interface{}
	next        *Item
	prev        *Item
}

func New() *List{
	return &List{}
}

func (i *Item) Value() interface{} {
	return i.val
}

func (i *Item) Next() *Item {
	return i.next
}

func (i *Item) Prev() *Item {
	return i.prev
}

func (l *List) Len() uint {
	return l.len
}

func (l *List) First() *Item {
	return l.first
}

func (l *List) Last() *Item {
	return l.last
}

func (l *List) PushFront(v interface{}) {
	item := &Item{
		val:    v,
	}

	if l.first != nil {
		l.first.prev = item
		item.next = l.first
	}
	l.first = item

	if l.last == nil {
		l.last = item
	}
	l.len++
}

func (l *List) PushBack(v interface{}) {
	item := &Item{
		val:    v,
	}

	if l.last != nil {
		l.last.next = item
		item.prev = l.last
	}
	l.last = item

	if l.first == nil {
		l.first = item
	}
	l.len++
}

func (l *List) Remove(i Item) {

	if i.prev != nil {
		i.prev.next = i.next
	} else {
		l.first = i.next
	}

	if i.next != nil {
		i.next.prev = i.prev
	} else {
		l.last = i.prev
	}

	l.len--
}

func (l *List) NextItem() (*Item, bool) {
	isEnd := false
	if l.current == nil {
		l.current = l.first
	} else {
		l.current = l.current.next
	}

	if l.current == nil {
		isEnd = true
	}
	return l.current, isEnd
}
