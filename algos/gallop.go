package algos

// Galloping works to move quickly through one array when the other array is much smaller.
// 4882381383139334
func SortedIntersectGalloping(lhs []int64, rhs []int64) (lhs_indices []int, rhs_indices []int) {
    lhs_idx, rhs_idx := 0, 0
    gallop := 1
    for lhs_idx < len(lhs) && rhs_idx < len(rhs) {
        // Advance LHS past rhs
        for lhs_idx < len(lhs) && lhs[lhs_idx] < rhs[rhs_idx] {
            lhs_idx += gallop
            gallop <<= 1
        }
        lhs_idx -= (gallop >> 1)
        gallop = 1
        // Advance RHS past lhs
        for rhs_idx < len(rhs) && rhs[rhs_idx] < lhs[lhs_idx] {
            rhs_idx += gallop
            gallop <<= 1
        }
        rhs_idx -= (gallop >> 1)
        gallop = 1

        // Standard two pointer check
        if lhs[lhs_idx] < rhs[rhs_idx] {
            lhs_idx++
        } else if lhs[lhs_idx] > rhs[rhs_idx] {
            rhs_idx++
        } else {
            for lhs_idx < len(lhs) && rhs_idx < len(rhs) && lhs[lhs_idx] == rhs[rhs_idx] {
                lhs_indices = append(lhs_indices, lhs_idx)
                rhs_indices = append(rhs_indices, rhs_idx)
                lhs_idx++
                rhs_idx++
            }
        }
    }
    return lhs_indices, rhs_indices
}
