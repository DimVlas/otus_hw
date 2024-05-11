package hw04lrucache

type List interface {
	Len() int                          // Кол-во элементов в списке.
	Front() *ListItem                  // Первый элемент списка.
	Back() *ListItem                   // Последний элемент списка.
	PushFront(v interface{}) *ListItem // Добавление элемента в начало списка.
	PushBack(v interface{}) *ListItem  // Добавление элемента в конец списка.
	Remove(i *ListItem)                // Удаление элемента из списка.
	MoveToFront(i *ListItem)           // Переместить элемент вперед.
}

type ListItem struct {
	Value interface{}
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

// Создать новый двусвязный список.
func NewList() List {
	return new(list)
}

// Кол-во элементов в списке.
func (lst *list) Len() int {
	return lst.len
}

// Первый элемент списка.
func (lst *list) Front() *ListItem {
	return lst.front
}

// Последний элемент списка.
func (lst *list) Back() *ListItem {
	return lst.back
}

// Добавление элемента в начало списка.
func (lst *list) PushFront(v interface{}) *ListItem {
	itm := &ListItem{
		Value: v,
		Next:  lst.front,
	}
	lst.front = itm
	lst.len++

	if itm.Next == nil {
		lst.back = itm
		return itm
	}

	itm.Next.Prev = itm
	return itm
}

// Добавление элемента в конец списка.
func (lst *list) PushBack(v interface{}) *ListItem {
	itm := &ListItem{
		Value: v,
		Prev:  lst.back,
	}
	lst.back = itm
	lst.len++

	if itm.Prev == nil {
		lst.front = itm
		return itm
	}

	itm.Prev.Next = itm
	return itm
}

// Удаление элемента из списка.
func (lst *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev == nil && i.Next == nil { // единственный элемент
		lst.front = nil
		lst.back = nil
		lst.len--
		return
	}

	if i.Prev == nil { // первый элемент
		i.Next.Prev = nil
		lst.front = i.Next
		lst.len--
		return
	}

	if i.Next == nil { // последний элемент
		i.Prev.Next = nil
		lst.back = i.Prev
		lst.len--
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	lst.len--
}

// Переместить элемент вперед.
func (lst *list) MoveToFront(i *ListItem) {
	if i.Prev == nil { // первый элемент
		return
	}

	if i.Next != nil { // элемент из середины
		i.Prev.Next = i.Next // предыдущий ссылается на следующий
		i.Next.Prev = i.Prev // следующий ссылается на предыдущий
	} else { // последний элемент
		i.Prev.Next = nil
		lst.back = i.Prev
	}

	i.Prev = nil
	i.Next = lst.front
	lst.front = i
	i.Next.Prev = i
}
