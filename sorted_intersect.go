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
    // fmt.Printf("lhs_indices: %v\n", lhs_indices[len(lhs_indices)-10:])
    // fmt.Printf("rhs_indices: %v\n", lhs_indices[len(rhs_indices)-10:])
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


func buildIndexMask(sorted_arr []int64, maxValue int64) int64 {
    var bitmask int64
    var currMin, currMax int64
    var i int64

    for i = 0; i < 64; i++ {
        currMin = i * (maxValue / 64)
        currMax = (i + 1) * (maxValue / 64)
        for _, val := range sorted_arr {
            if val >= currMin && val < currMax {
                bitmask |= 1 << i
                break
            }
        }
    }
    return bitmask
}


func main() {
    fmt.Println("******************")
    fmt.Println("Lobsided")
    lhs, rhs := randArrs(42, 100, 1000000, 10000)
    bitmaskLhs := buildIndexMask(lhs, 10000)
    bitmaskRhs := buildIndexMask(rhs, 10000)
    fmt.Printf("lhs bitmask: %0x\n", uint64(bitmaskLhs))
    fmt.Printf("rhs bitmask: %0x\n", uint64(bitmaskRhs))

    profileCall(sortedIntersectNaive, lhs, rhs, 100)
    profileCall(sortedIntersectGalloping, lhs, rhs, 100)
   
    fmt.Println("******************")
    fmt.Println("Even")
    lhs, rhs = randArrs(42, 1000000, 1000000, 10000)

    profileCall(sortedIntersectNaive, lhs, rhs, 100)
    profileCall(sortedIntersectGalloping, lhs, rhs, 100)
}
