package algos


func IntersectWithHeaderMarks(lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64) (lhs_indices []int, rhs_indices []int) {

    lhsHeaderIdx := 0
    rhsHeaderIdx := 0

    for lhsHeaderIdx < len(lhsHeaders) && rhsHeaderIdx < len(rhsHeaders) {
        if (lhsHeaders[lhsHeaderIdx] >> 32) < (rhsHeaders[rhsHeaderIdx] >> 32) {
            lhsHeaderIdx++
        } else if (lhsHeaders[lhsHeaderIdx] >> 32) > (rhsHeaders[rhsHeaderIdx] >> 32) {
            rhsHeaderIdx++
        } else {
            // get the index into lhs / rhs
            lhsSliceStart := int(lhsHeaders[lhsHeaderIdx] & 0xFFFFFFFF)
            rhsSliceStart := int(rhsHeaders[rhsHeaderIdx] & 0xFFFFFFFF)
            lhsSliceNext := len(lhs)
            rhsSliceNext := len(rhs)
            if lhsHeaderIdx + 1 < len(lhsHeaders) {
                lhsSliceNext = int(lhsHeaders[lhsHeaderIdx+1] & 0xFFFFFFFF)
            } 
            if rhsHeaderIdx + 1 < len(rhsHeaders) {
                rhsSliceNext = int(rhsHeaders[rhsHeaderIdx+1] & 0xFFFFFFFF)
            }

            lhsSlice := lhs[lhsSliceStart:lhsSliceNext]
            rhsSlice := rhs[rhsSliceStart:rhsSliceNext]
            // two pointer intesrect
            i, j := 0, 0
            for i < len(lhsSlice) && j < len(rhsSlice) {
                if lhsSlice[i] < rhsSlice[j] {
                    i++
                } else if lhsSlice[i] > rhsSlice[j] {
                    j++
                } else {
                    lhs_indices = append(lhs_indices, i + lhsSliceStart)
                    rhs_indices = append(rhs_indices, j + rhsSliceStart)
                    i++
                    j++
                }
            }
            rhsHeaderIdx++
            lhsHeaderIdx++
        }
    }
    return lhs_indices, rhs_indices
}

