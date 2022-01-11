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
	length int
	head   *ListItem
	tail   *ListItem
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l list) Len() int {
	return l.length
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem, first := l.isFirstItem(v)
	if !first {
		temp := l.head
		l.head = newItem
		l.head.Next = temp
		temp.Prev = l.head
	}
	l.length++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem, first := l.isFirstItem(v)
	if !first {
		temp := l.tail
		l.tail = newItem
		l.tail.Prev = temp
		temp.Next = l.tail
	}
	l.length++
	return newItem
}

func (l *list) isFirstItem(v interface{}) (*ListItem, bool) {
	newItem := &ListItem{Value: v}
	if l.head == nil {
		l.head = newItem
		l.tail = newItem
		return newItem, true
	}
	return newItem, false
}

func (l *list) Remove(i *ListItem) {
	temp := l.head
	if temp != nil && temp.Value == i.Value {
		l.head = temp.Next
		l.head.Prev = nil
		return
	}
	var prev *ListItem
	for temp != nil && temp.Value != i.Value {
		prev = temp
		temp = temp.Next
	}
	if temp == nil {
		return
	}
	prev.Next = temp.Next
	if temp.Next != nil {
		temp.Next.Prev = prev
	} else {
		l.tail = prev
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head.Value == i.Value {
		return
	}
	cur := l.tail
	for {
		if cur != nil && cur.Value != i.Value {
			cur = cur.Prev
		} else {
			if cur.Next == nil {
				cur.Prev.Next = nil
			} else {
				cur.Next.Prev = cur.Prev
				cur.Prev.Next = cur.Next
			}
			break
		}
	}
	if i.Next == nil {
		l.tail = i.Prev
	}
	i.Next = l.head
	i.Prev = nil

	l.head.Prev = i
	l.head = i
}

func NewList() List {
	return new(list)
}
