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
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	if l.len == 0 {
		l.front = item
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}

	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	if l.len == 0 {
		l.front = item
		l.back = item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}

	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if l.len == 1 {
		l.front = nil
		l.back = nil
		l.len = 0
		return
	}

	var isFront, isBack bool
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		isBack = true
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		isFront = true
	}

	if isFront {
		l.front = i.Next
	}

	if isBack {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len == 1 {
		return
	}

	if i.Prev == nil {
		return
	}

	i.Prev.Next = i.Next

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}
