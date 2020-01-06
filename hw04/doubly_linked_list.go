package hw04

import (

)

// List of Items
type List struct {
	len     uint
	first   *Item
	last    *Item
	current *Item
}

// Item of the List
type Item struct {
	val         interface{}
	next        *Item
	prev        *Item
}

// New object of the List (constructor)
func New() *List{
	return &List{}
}

// Value of the Item object
func (i *Item) Value() interface{} {
	return i.val
}

// Next Item object
func (i *Item) Next() *Item {
	return i.next
}

// Prev (previous) Item object
func (i *Item) Prev() *Item {
	return i.prev
}

// Len (a length) of a list object
func (l *List) Len() uint {
	return l.len
}

// First Item object in a List object
func (l *List) First() *Item {
	return l.first
}

// Last Item object in a List object
func (l *List) Last() *Item {
	return l.last
}

// PushFront is the method for a pushing into front of a List object
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

// PushBack is the method for a pushing into back of a List object
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

// Remove a Item object ia a List object
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

// NextItem is the method for getting a next Item object in a List object and a list completion flag. For usage in a loop.
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
