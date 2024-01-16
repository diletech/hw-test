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
	newFront := &ListItem{Value: v, Next: l.front}
	if l.front != nil {
		l.front.Prev = newFront
	}
	l.front = newFront
	if l.back == nil {
		l.back = newFront
	}
	l.len++
	return newFront
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{Value: v, Prev: l.back}
	if l.back != nil {
		l.back.Next = newBack
	}
	l.back = newBack
	if l.front == nil {
		l.front = newBack
	}
	l.len++
	return newBack
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
