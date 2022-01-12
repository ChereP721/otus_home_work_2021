package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) MoveToFront(i *ListItem) {
	v := i.Value
	l.Remove(i)
	l.PushFront(v)
}

func (l *list) Remove(i *ListItem) {
	l.len--
	if l.len == 0 {
		l.front = nil
		l.back = nil

		return
	}

	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := ListItem{Value: v, Next: nil, Prev: l.back}
	if l.back != nil {
		l.back.Next = &li
	} else {
		l.front = &li
	}
	l.back = &li
	l.len++

	return &li
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := ListItem{Value: v, Next: l.front, Prev: nil}
	if l.front != nil {
		l.front.Prev = &li
	} else {
		l.back = &li
	}
	l.front = &li
	l.len++

	return &li
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func NewList() List {
	return new(list)
}
