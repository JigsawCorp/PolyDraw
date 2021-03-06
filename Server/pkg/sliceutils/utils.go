package sliceutils

//PopInt remove an element from the array. The order will change
func PopInt(a *[]int, i int) int {
	// Remove the element at index i from a.
	lastElement := len(*a) - 1
	value := (*a)[i]
	(*a)[i] = (*a)[lastElement]
	(*a)[lastElement] = -1
	*a = (*a)[:lastElement]

	return value
}

//PopAtLocationInt remove as an element in the array and determines the end
func PopAtLocationInt(a *[]int, i int, lastElement int) int {
	// Remove the element at index i from a.
	value := (*a)[i]
	(*a)[i] = (*a)[lastElement]
	(*a)[lastElement] = -1

	return value
}
