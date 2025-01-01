package sortedGen

import (
	"math/rand"
	"fmt"
	"slices"
)


// Return a random array of length length with values between 0 and maxVal
// sorted and with no duplicates
func RandUnique(length int, maxVal int64) (arr []int64, err error) {
    seen := make(map[int64]bool)
	if maxVal < int64(length) {
		return nil, fmt.Errorf("Cannot create unique array where maxVal must be greater than length")
	}
    for i := 0; i < length; i++ {
		for {
			val := rand.Int63n(maxVal)
			if _, ok := seen[val]; !ok {
				arr = append(arr, val)
				seen[val] = true
				break
			}
		}
		if len(arr) == length {
			break
		}
    }
    slices.Sort(arr)
	return arr, nil
}
