package algos

type IntersectFunc func(lhs []int64, rhs []int64) ([]int, []int)


func HeaderIndices(arr []int64) (headers []int64) {
    lastHeader := -1
    for i := 0; i < len(arr); i++ {
        headerVal := int(arr[i] >> 32)
        if headerVal != lastHeader {
            value := int64(headerVal << 32) | int64(i)
            headers = append(headers, value)
            lastHeader = int(arr[i] >> 32)
        }
    }
    return headers
}
