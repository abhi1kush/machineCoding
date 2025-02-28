package tests

import (
	"fmt"
	"slices" // Go 1.18+ package for slice utilities (import "golang.org/x/exp/slices" if not in stdlib)
	"sort"
)

func SliceUse() {
	// Declaration method 1: Create a slice with a fixed size.
	size := 10
	slice := make([]int, size)
	fmt.Printf("Slice (method 1): %v\n", slice)
	// Declaration method 2: Short-hand with make.
	slice2 := make([]int, size)
	fmt.Printf("Slice (method 2): %v\n", slice2)
	// Declaration method 3: Literal empty slice.
	slice3 := []int{}
	fmt.Printf("Slice (method 3): %v\n", slice3)
	// Declaration method 4: Slice literal with initial values.
	slice4 := []int{4, 7, 1, 8, 3, 0, 2, 67, 34, 4, 1}
	// Adding an element using append.
	slice4 = append(slice4, 11)
	// Get length.
	length := len(slice4)
	fmt.Printf("Length: %v\n", length)
	// Remove last element.
	slice4 = slice4[:len(slice4)-1]
	// Remove first element.
	slice4 = slice4[1:]
	// Sort in ascending order.
	slices.Sort(slice4)
	fmt.Printf("Sorted Ascending: %v\n", slice4)
	// Sort in descending order using slices.SortFunc.
	slices.SortFunc(slice4, func(left, right int) int {
		return right - left // Swap order for descending sort.
	})
	fmt.Printf("Sorted Descending: %v\n", slice4)
	// Alternatively, using sort.Slice for descending order.
	sort.Slice(slice4, func(i, j int) bool {
		return slice4[i] > slice4[j]
	})
	fmt.Printf("Sorted Descending (sort.Slice): %v\n", slice4)
	// Reverse the slice.
	slices.Reverse(slice4)
	fmt.Printf("Reversed: %v\n", slice4)
}
func main() {
	SliceUse()
}
