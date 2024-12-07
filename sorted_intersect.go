package main

import (
    "math/rand"
    "time"
    "fmt"
    "slices"

    "reflect"
    "runtime"
)


func timer(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}


func sortedIntersectNaive(lhs []int64, rhs []int64) (lhs_indices []int, rhs_indices []int) {
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


// Galloping works to move quickly through one array when the other array is much smaller.
func sortedIntersectGalloping(lhs []int64, rhs []int64) (lhs_indices []int, rhs_indices []int) {
    lhs_idx, rhs_idx := 0, 0
    lhs_gallop, rhs_gallop := 1, 1
    for lhs_idx < len(lhs) && rhs_idx < len(rhs) {
        // Advance LHS past rhs
        for lhs_idx < len(lhs) && lhs[lhs_idx] < rhs[rhs_idx] {
            lhs_idx += lhs_gallop
            lhs_gallop <<= 1
        }
        lhs_idx -= lhs_gallop >> 1
        lhs_gallop = 1
        // Advance RHS past lhs
        for rhs_idx < len(rhs) && rhs[rhs_idx] < lhs[lhs_idx] {
            rhs_idx += rhs_gallop
            rhs_gallop <<= 1
        }
        rhs_idx -= rhs_gallop >> 1
        rhs_gallop = 1

        // Standard two pointer check
        if lhs[lhs_idx] < rhs[rhs_idx] {
            lhs_idx++
        } else if lhs[lhs_idx] > rhs[rhs_idx] {
            rhs_idx++
        } else {
            lhs_indices = append(lhs_indices, lhs_idx)
            rhs_indices = append(rhs_indices, rhs_idx)
            lhs_idx++
            rhs_idx++
        }
    }
    return lhs_indices, rhs_indices
}

type intersectFunc func(lhs []int64, rhs []int64) ([]int, []int)

func profileCall(f intersectFunc, lhs []int64, rhs []int64, times int) (lhs_indices []int, rhs_indices []int) {
    funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
    defer timer(funcName)()
    for i := 0; i < times; i++ {
        lhs_indices, rhs_indices = f(lhs, rhs)
    }
    return lhs_indices, rhs_indices
}


func randArrs(seed, lengthLhs, lengthRhs int, maxVal int64) (lhs []int64, rhs []int64) {
    rand.Seed(int64(seed))
    for i := 0; i < lengthLhs; i++ {
        lhs = append(lhs, rand.Int63n(maxVal))
    }
    for i := 0; i < lengthRhs; i++ {
        rhs = append(rhs, rand.Int63n(maxVal))
    }
    slices.Sort(lhs)
    slices.Sort(rhs)
    return lhs, rhs
}


func randArrsWithHeaders(seed, lengthLhs, lengthRhs int) (lhs []int64, rhs []int64) {
    rand.Seed(int64(seed))
    header := 0
    value := 0
    for i := 0; i < lengthLhs; i++ {
        incrHeader := rand.Intn(10) == 0
        if incrHeader {
            header++
            value = 0
        }
        lhs = append(lhs, int64(header << 32 | value))
        value++
        for j:=0; j < 10; j++ {
            incrValue := rand.Intn(2) == 0
            if incrValue {
                value++
            } else {
                break
            }
        }
    }
    header = 0
    value = 0
    for i := 0; i < lengthRhs; i++ {
        incrHeader := rand.Intn(10) == 0
        if incrHeader {
            header++
            value = 0
        }
        rhs = append(rhs, int64(header << 32 | value))
        for j:=0; j < 10; j++ {
            incrValue := rand.Intn(2) == 0
            if incrValue {
                value++
            } else {
                break
            }
        }
    }
    return lhs, rhs
}


func headerIndices(arr []int64) (headers []int64) {
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


func intersectWithHeaderMarks(lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64) (lhs_indices []int, rhs_indices []int) {

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


func makeHeaderIntesectFn(
    lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64,
) intersectFunc {
    return func(lhs []int64, rhs []int64) ([]int, []int) {
        return intersectWithHeaderMarks(lhs, rhs, lhsHeaders, rhsHeaders)
    }
}

func main() {
    fmt.Println("******************")
    fmt.Println("Lobsided -- dense header")
    // lhs, rhs := randArrs(42, 100, 1000000, maxVal)
    lhs, rhs := randArrsWithHeaders(42, 100, 1000000)
    lhsHeaders := headerIndices(lhs)
    rhsHeaders := headerIndices(rhs)
    headerIntersectFn := makeHeaderIntesectFn(lhs, rhs, lhsHeaders, rhsHeaders)

    profileCall(sortedIntersectNaive, lhs, rhs, 100)
    profileCall(sortedIntersectGalloping, lhs, rhs, 100)
    profileCall(headerIntersectFn, lhs, rhs, 100)
   
    fmt.Println("******************")
    fmt.Println("Even -- dense header")
    lhs, rhs = randArrsWithHeaders(42, 1000000, 1000000)
    lhsHeaders = headerIndices(lhs)
    rhsHeaders = headerIndices(rhs)
    headerIntersectFn = makeHeaderIntesectFn(lhs, rhs, lhsHeaders, rhsHeaders)

    profileCall(sortedIntersectNaive, lhs, rhs, 100)
    profileCall(sortedIntersectGalloping, lhs, rhs, 100)
    profileCall(headerIntersectFn, lhs, rhs, 100)
    
    fmt.Println("******************")
    fmt.Println("Lobsided -- truly random")
    maxVal := int64(10000)
    lhs, rhs = randArrs(42, 100, 1000000, maxVal)
    lhsHeaders = headerIndices(lhs)
    rhsHeaders = headerIndices(rhs)
    headerIntersectFn = makeHeaderIntesectFn(lhs, rhs, lhsHeaders, rhsHeaders)

    profileCall(sortedIntersectNaive, lhs, rhs, 100)
    profileCall(sortedIntersectGalloping, lhs, rhs, 100)
    profileCall(headerIntersectFn, lhs, rhs, 100)
    
    fmt.Println("******************")
    fmt.Println("Even -- truly random")
    lhs, rhs = randArrs(42, 100000, 1000000, maxVal)
    lhsHeaders = headerIndices(lhs)
    rhsHeaders = headerIndices(rhs)
    headerIntersectFn = makeHeaderIntesectFn(lhs, rhs, lhsHeaders, rhsHeaders)

    profileCall(sortedIntersectNaive, lhs, rhs, 100)
    profileCall(sortedIntersectGalloping, lhs, rhs, 100)
    profileCall(headerIntersectFn, lhs, rhs, 100)
   
}
