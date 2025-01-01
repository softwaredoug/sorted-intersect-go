package main

import (
    "math/rand"
    "time"
    "fmt"
    "slices"

    "runtime"

    "reflect"

	"github.com/softwaredoug/sorted-intersect-go/algos"
	"github.com/softwaredoug/sorted-intersect-go/sortedGen"
)


func timer(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}


func randArrs(seed, lengthLhs, lengthRhs int, maxVal int64) (lhs []int64, rhs []int64) {
    rand.Seed(int64(seed))
	lhs, err := sortedGen.RandUnique(lengthLhs, maxVal)
	if err != nil {
		panic(err)
	}
	rhs, err = sortedGen.RandUnique(lengthRhs, maxVal)
	if err != nil {
		panic(err)
	}
    return lhs, rhs
}


func randArrsWithHeaders(seed, lengthLhs, lengthRhs int) (lhs []int64, rhs []int64) {
    rand.Seed(int64(seed))
	lhs = sortedGen.RandWithHeader(seed, lengthLhs)
	rhs = sortedGen.RandWithHeader(seed, lengthRhs)
    return lhs, rhs
}

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


func profileUnchecked(f algos.IntersectFunc, lhs []int64, rhs []int64, times int) (lhs_indices []int, rhs_indices []int) {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	defer timer(funcName)()
	for i := 0; i < times; i++ {
		lhs_indices, rhs_indices = f(lhs, rhs)
	}
	return lhs_indices, rhs_indices
}

func profileCall(f algos.IntersectFunc, lhs []int64, rhs []int64, times int,
                 lhs_check, rhs_check []int) (lhs_indices []int, rhs_indices []int) {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lhs_indices, rhs_indices = profileUnchecked(f, lhs, rhs, times)
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
			fmt.Printf("Duplicate value: %d\n", val)
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
        panic("Warning: lhs not unique")
    }
    if !isArrUnique(rhs) {
        panic("Warning: rhs not unique")
    }
    lhsHeaders := algos.HeaderIndices(lhs)
    rhsHeaders := algos.HeaderIndices(rhs)
    fmt.Printf("LHS Header len: %d / %d\n", len(lhsHeaders), len(lhs))
    fmt.Printf("RHS Header len: %d / %d\n", len(rhsHeaders), len(rhs))
    headerIntersectFn := algos.MakeHeaderIntesectFn(lhs, rhs, lhsHeaders, rhsHeaders)
    headerIntersectGallopToNaiveFn := algos.MakeHeaderIntesectGallopToNaiveFn(lhs, rhs, lhsHeaders, rhsHeaders)
    headerIntersectGallopToGallopFn := algos.MakeHeaderIntesectGallopToGallopFn(lhs, rhs, lhsHeaders, rhsHeaders)

    lhsResult, rhsResult := profileCall(algos.SortedIntersectNaive, lhs, rhs, 1, nil, nil)
    profileCall(algos.SortedIntersectGalloping, lhs, rhs, 1, lhsResult, rhsResult)
    profileCall(headerIntersectFn, lhs, rhs, 1, lhsResult, rhsResult)
    profileCall(headerIntersectGallopToNaiveFn, lhs, rhs, 1, lhsResult, rhsResult)
    profileCall(headerIntersectGallopToGallopFn, lhs, rhs, 1, lhsResult, rhsResult)
}





func main() {
    maxVal := int64(100000000)
    
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
    fmt.Println("******************")
    fmt.Println("******************")
    fmt.Println("Lobsided -- truly random")
    lhs, rhs = randArrs(42, 100, 1000000, maxVal)
    profileAll(lhs, rhs)
    
    fmt.Println("******************")
    fmt.Println("Even -- truly random")
    lhs, rhs = randArrs(42, 100000, 1000000, maxVal)
    profileAll(lhs, rhs)
   
}
