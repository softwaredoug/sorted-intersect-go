package main

import (
    "math/rand"
    "time"
    "fmt"
    "slices"

    "runtime"

    "reflect"
)


func timer(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}


func sortedIntersectNaive(lhs []int64, rhs []int64) (lhs_indices []int, rhs_indices []int) {
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


// Galloping works to move quickly through one array when the other array is much smaller.
// 4882381383139334
func sortedIntersectGalloping(lhs []int64, rhs []int64) (lhs_indices []int, rhs_indices []int) {
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
        incrHeader := rand.Intn(2) == 0
        if incrHeader {
            header += rand.Intn(10) + 1
            value = 0
        }
        lhs = append(lhs, int64(header << 32 | value))
        value++
        for j:=0; j < 10; j++ {
            keepGoing := rand.Intn(2) > 2
            if keepGoing {
                value += int(rand.Int63n(10))
            } else {
                break
            }
        }
    }
    // start in the middle
    maxLhs := int(lhs[len(lhs)-1] >> 32)
    header = maxLhs / 2
    value = 0
    for i := 0; i < lengthRhs; i++ {
        incrHeader := rand.Intn(10) == 0
        if incrHeader {
            header += rand.Intn(10) + 1
            value = 0
        }
        rhs = append(rhs, int64(header << 32 | value))
        for j:=0; j < 10; j++ {
            keepGoing := rand.Intn(2) == 0
            if keepGoing {
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


func intersectGallopWithHeaderMarks(lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64) (lhs_indices []int, rhs_indices []int) {

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


func makeHeaderIntesectGallopFn(
    lhs []int64, rhs []int64, lhsHeaders []int64, rhsHeaders []int64,
) intersectFunc {
    return func(lhs []int64, rhs []int64) ([]int, []int) {
        return intersectGallopWithHeaderMarks(lhs, rhs, lhsHeaders, rhsHeaders)
    }
}


type intersectFunc func(lhs []int64, rhs []int64) ([]int, []int)

func diffArrays(expected, actual []int, self []int64, other []int64) {
    for i := 0; i < len(expected); i++ {
        if expected[i] != actual[i] {
            valueAt := self[expected[i]]
            found := false
            fmt.Printf("At %d -- Expected: %d, Actual: %d\n", i, expected[i], actual[i])
            fmt.Printf("Expected Value Here %d\n", valueAt)
            // find valueAt in other
            j := int64(0)
            for j = 0; j < int64(len(other)); j++ {
                if other[j] == valueAt {
                    fmt.Printf("Value in other at %d: %d\n", j, other[j])
                    fmt.Printf("Length of other %d\n", len(other))
                    fmt.Printf("Length of self %d\n", len(self))
                    found = true
                    break
                }
            }
            if !found {
                fmt.Println("Value not found in other")
            }
            return
        }
    }
}

func profileCall(f intersectFunc, lhs []int64, rhs []int64, times int,
                 lhs_check, rhs_check []int) (lhs_indices []int, rhs_indices []int) {
    funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
    defer timer(funcName)()
    for i := 0; i < times; i++ {
        lhs_indices, rhs_indices = f(lhs, rhs)
    }
    if lhs_check != nil && rhs_check != nil {
        if !slices.Equal(lhs_check, lhs_indices) {
            msg := fmt.Sprintf("lhs check failed: %s", funcName)
            diffArrays(lhs_check, lhs_indices, lhs, rhs)
            // runtime.Breakpoint()
            panic(msg)
        }
        if !slices.Equal(rhs_check, rhs_indices) {
            msg := fmt.Sprintf("rhs check failed: %s", funcName)
            diffArrays(rhs_check, rhs_indices, rhs, lhs)
            panic(msg)
        }
    }
    return lhs_indices, rhs_indices
}




func isArrUnique(arr []int64) bool {
    seen := make(map[int64]bool)
    for _, val := range arr {
        if _, ok := seen[val]; ok {
            return false
        }
        seen[val] = true
    }
    return true
}


func isArrSorted(arr []int64) bool {
    for i := 1; i < len(arr); i++ {
        if arr[i] < arr[i-1] {
            fmt.Printf("arr[%d]: %d, arr[%d]: %d\n", i-1, arr[i-1], i, arr[i])
            fmt.Printf("header: %d, value: %d\n", arr[i] >> 32, arr[i] & 0xFFFFFFFF)
            fmt.Printf("header: %d, value: %d\n", arr[i-1] >> 32, arr[i-1] & 0xFFFFFFFF)
            return false
        }
    }
    return true
}


func profileAll(lhs, rhs []int64) {
    if !isArrSorted(lhs) {
        panic("lhs not sorted")
    }
    if !isArrSorted(rhs) {
        panic("rhs not sorted")
    }
    if !isArrUnique(lhs) {
        fmt.Println("Warning: lhs not unique")
    }
    if !isArrUnique(rhs) {
        fmt.Println("Warning: rhs not unique")
    }
    lhsHeaders := headerIndices(lhs)
    rhsHeaders := headerIndices(rhs)
    fmt.Printf("LHS Header len: %d / %d\n", len(lhsHeaders), len(lhs))
    fmt.Printf("RHS Header len: %d / %d\n", len(rhsHeaders), len(rhs))
    headerIntersectFn := makeHeaderIntesectFn(lhs, rhs, lhsHeaders, rhsHeaders)
    headerIntersectGallopFn := makeHeaderIntesectGallopFn(lhs, rhs, lhsHeaders, rhsHeaders)

    lhsResult, rhsResult := profileCall(sortedIntersectNaive, lhs, rhs, 1, nil, nil)
    profileCall(sortedIntersectGalloping, lhs, rhs, 1, lhsResult, rhsResult)
    profileCall(headerIntersectFn, lhs, rhs, 1, lhsResult, rhsResult)
    profileCall(headerIntersectGallopFn, lhs, rhs, 1, lhsResult, rhsResult)
}





func main() {
    for i := 0; i < 1000; i++ {
        fmt.Println("******************")
        fmt.Printf("Even -- data -- seed: %d\n", i)
        lhs, rhs := randArrsWithHeaders(i, 100, 100)
        profileAll(lhs, rhs)
    }

    fmt.Println("******************")
    fmt.Println("Even -- data with headers")
    lhs, rhs := randArrsWithHeaders(42, 1000000, 1000000)
    profileAll(lhs, rhs)
    
    fmt.Println("******************")
    fmt.Println("Lobsided -- lhs largest - data with headers")
    lhs, rhs = randArrsWithHeaders(42, 100000000, 100)
    profileAll(lhs, rhs)
    

    fmt.Println("******************")
    fmt.Println("Lobsided -- rhs largest - data with headers")
    lhs, rhs = randArrsWithHeaders(42, 100, 100000000)
    profileAll(lhs, rhs)

    fmt.Println("******************")
    fmt.Println("Lobsided -- truly random")
    maxVal := int64(10000)
    lhs, rhs = randArrs(42, 100, 1000000, maxVal)
    profileAll(lhs, rhs)
    
    fmt.Println("******************")
    fmt.Println("Even -- truly random")
    lhs, rhs = randArrs(42, 100000, 1000000, maxVal)
    profileAll(lhs, rhs)
   
}
