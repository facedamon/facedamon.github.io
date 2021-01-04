package doublylinkedlist

import (
	"fmt"
	"strings"
	"github.com/facedamon/gods/utils"
)

type element struct {
	value interface{}
	prev  *element
	next  *element
}

// List holds the elements, where each element points to
// the next and previous element
type List struct {
	first *element
	last  *element
	size  int
}

// New instantiates a new list and adds the passed values, if any, to the list
func New(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

// Add appends a value at the end of the list
func (list *List) Add(values ...interface{}) {
	for _, value := range values {
		newElement := &element{value: value, prev: list.last}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.last.next = newElement
			list.last = newElement
		}
		list.size++
	}
}

// Prepend prepends a values
func (list *List) Prepend(values ...interface{}) {
	for v := len(values) - 1; v >= 0; v-- {
		newElement := &element{value: values[v], next: list.first}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.first.prev = newElement
			list.first = newElement
		}
		list.size++
	}
}

func (list *List) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

// Clear removes all elements from the list
func (list *List) Clear() {
	list.size = 0
	list.last = nil
	list.first = nil
}

// Get returns the element at index.
// Second return parameter is true if index is within bounds of the array and array is not empty
func (list *List) Get(index int) (interface{}, bool) {
	if !list.withinRange(index) {
		return nil, false
	}
	if list.size-index < index {
		element := list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
		return element.value, true
	}
	element := list.first
	for e := 0; e != index; e, element = e+1, element.next {
	}
	return element.value, true
}

// Remove removes the element at the given index from the list.
func (list *List) Remove(index int) {
	if !list.withinRange(index) {
		return
	}
	if list.size == 1 {
		list.Clear()
		return
	}
	var element *element
	//二分法
	if list.size-index < index {
		element = list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
	} else {
		element = list.first
		for e := 0; e != index; e, element = e+1, element.next {
		}
	}
	if element == list.first {
		list.first = element.next
	}
	if element == list.last {
		list.last = element.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}

	element = nil
	list.size--
}

// Contains check if values (one or more) are present int the set.
// All values have to be present in the set for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all.
func (list *List) Contains(values ...interface{}) bool {
	if len(values) == 0 {
		return true
	}
	if list.size == 0 {
		return false
	}
	for _, value := range values {
		found := false
		for element := list.first; element != nil; element = element.next {
			if element.value = value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Values returns all elements in the list.
func (list *List) Values() []interface{} {
	values := make([]interface{}, list.size, list.size)
	for e, element := 0, list.first; element != nil; e, element = e+1, element.next {
		values[e] = element.value
	}
	return values
}

// IndexOf returns index of provided element
func (list *List) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.Values() {
		if element == value {
			return index
		}
	}
	return -1
}

// Empty returns true if list does not contain any elements.
func (list *List) Empty() bool {
	return list.size == 0
}

// Sort sortts values using
func (list *List) Sort(comparator utils.comparator) {
	if list.size < 2 {
		return
	}
	values := list.Values()
	utils.Sort(values, comparator)
	list.Clear()
	list.Add(values...)
} 

// Insert inserts values at specified index position shifting the value at that position and any subsequent elements to the right.
// Does not do anything if position is negative or bigger than list`s size
// Note: position equal to list`s size is valid
func (list *List) Insert(index int, values ...interface{}) {
	if !list.withinRange((index)) {
		if index == list.size {
			list.Add(values...)
		}
		return
	}
	list.size += len(values)

	var beforeElement *element
	var foundElement *element

	// determine traversal direction, last to first or first to last
	if list.size - index < index {
		foundElement = list.last
		for e := list.size - 1; e != index; e, foundElement = e-1, foundElement.prev {
			beforeElement = foundElement.prev
		}
	} else {
		foundElement = list.first
		for e := 0; e != index; e, foundElement = e+1, foundElement.next {
			beforeElement = foundElement
		}
	}

	if foundElement == list.first {
		oldNextElement := list.first
		for i, v := range values {
			newElement := &element{value: v}
			if i == 0 {
				list.first = newElement
			} else {
				newElement.prev = beforeElement
				beforeElement.next = newElement
			}
			beforeElement = newElement
		}
		oldNextElement.prev = beforeElement
		beforeElement.next = oldNextElement
	} else {
		oldNextElement := beforeElement.next
		for _, v := range values {
			newElement := &element{value: v}
			newElement.prev = beforeElement
			beforeElement.next = newElement
			beforeElement = newElement
		}
		oldNextElement.prev = beforeElement
		beforeElement.next = oldNextElement
	}
}

// String returns a string representation of container
func (list *List) String() string {
	str := "DoubleLinkedList\n"
	values := []string{}
	for element := list.first; element != nil; element = element.next {
		values = append(values, fmt.Sprintf("%f", element.value))
	}
	str += strings.Join(values, ",")
	return str
}