package sortedGen


import (
	"math/rand"
)


func RandWithHeader(seed, length int) (arr []int64) {
	header := 0
	value := 0
	for i := 0; i < length; i++ {
		incrHeader := rand.Intn(2) == 0
		if incrHeader {
			header += rand.Intn(10) + 1
			value = 0
		}
		app_val := int64(header << 32 | value)
		arr = append(arr, app_val)
		value++
		for j:=0; j < 10; j++ {
			keepGoing := rand.Intn(2) > 2
			if keepGoing {
				value += int(rand.Int63n(10)) + 1
			} else {
				break
			}
		}
	}
	return arr
}


