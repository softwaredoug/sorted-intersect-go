package algos


func SortedIntersectNaive(lhs []int64, rhs []int64) (lhs_indices []int, rhs_indices []int) {
    // Standard two pointer check for intersections
    i, j := 0, 0
    for i < len(lhs) && j < len(rhs) {
        if lhs[i] < rhs[j] {
            i++
        } else if lhs[i] > rhs[j] {
            j++
        } else {
            lhs_indices = append(lhs_indices, i)
            rhs_indices = append(rhs_indices, j)
            i++
            j++
        }
    }
    return lhs_indices, rhs_indices
}


func MakeHeaderIntesectFn(
    lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64,
) IntersectFunc {
    return func(lhs []int64, rhs []int64) ([]int, []int) {
        return IntersectWithHeaderMarks(lhs, rhs, lhsHeaders, rhsHeaders)
    }
}

