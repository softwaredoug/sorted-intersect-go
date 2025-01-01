package algos

func IntersectGallopToGallopWithHeaderMarks(lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64) (lhs_indices []int, rhs_indices []int) {

    lhsHeaderIdx := 0
    rhsHeaderIdx := 0
    gallop := 1

    for lhsHeaderIdx < len(lhsHeaders) && rhsHeaderIdx < len(rhsHeaders) {
        // Advance LHS past rhs
        for lhsHeaderIdx < len(lhsHeaders) && lhsHeaders[lhsHeaderIdx] < rhsHeaders[rhsHeaderIdx] {
            lhsHeaderIdx += gallop
            gallop <<= 1
        }
        lhsHeaderIdx -= gallop >> 1
        gallop = 1
        // Advance RHS past lhs
        for rhsHeaderIdx < len(rhsHeaders) && rhsHeaders[rhsHeaderIdx] < lhsHeaders[lhsHeaderIdx] {
            rhsHeaderIdx += gallop
            gallop <<= 1
        }
        rhsHeaderIdx -= gallop >> 1
        gallop = 1

        // two pointer check
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

			// gallop intersect in this slice
            i, j := 0, 0
			gallop = 1
            for i < len(lhsSlice) && j < len(rhsSlice) {
				for i < len(lhsSlice) && (lhsSlice[i] < rhsSlice[j]) {
                    i += gallop
					gallop <<= 1
				}
				i -= gallop >> 1
				gallop = 1
				for j < len(rhsSlice) && rhsSlice[j] < lhsSlice[i] {
					j += gallop
					gallop <<= 1
				}
				j -= gallop >> 1
				gallop = 1

				// two pointer check
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


func MakeHeaderIntesectGallopToGallopFn(
	lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64,
) IntersectFunc {
	return func(lhs []int64, rhs []int64) ([]int, []int) {
		return IntersectGallopToGallopWithHeaderMarks(lhs, rhs, lhsHeaders, rhsHeaders)
	}
}

