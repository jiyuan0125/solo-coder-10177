package treeset

import (
	"cmp"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/emirpasic/gods/v2/utils"
)

func TestDiagPointers(t *testing.T) {
	f1 := cmp.Compare[int]
	f2 := cmp.Compare[int]

	fmt.Printf("f1 code ptr: %x\n", reflect.ValueOf(f1).Pointer())
	fmt.Printf("f2 code ptr: %x\n", reflect.ValueOf(f2).Pointer())
	fmt.Printf("f1 funcval ptr: %x\n", *(*uintptr)(unsafe.Pointer(&f1)))
	fmt.Printf("f2 funcval ptr: %x\n", *(*uintptr)(unsafe.Pointer(&f2)))
	fmt.Printf("f1 == f2 (code): %v\n", reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer())
	fmt.Printf("f1 == f2 (funcval): %v\n", *(*uintptr)(unsafe.Pointer(&f1)) == *(*uintptr)(unsafe.Pointer(&f2)))

	factory := func(m int) utils.Comparator[int] {
		return func(a, b int) int {
			return (a - b) * m
		}
	}

	c1 := factory(1)
	c2 := factory(2)
	c3 := factory(1)

	fmt.Printf("\nc1 code ptr: %x\n", reflect.ValueOf(c1).Pointer())
	fmt.Printf("c2 code ptr: %x\n", reflect.ValueOf(c2).Pointer())
	fmt.Printf("c3 code ptr: %x\n", reflect.ValueOf(c3).Pointer())
	fmt.Printf("c1 funcval ptr: %x\n", *(*uintptr)(unsafe.Pointer(&c1)))
	fmt.Printf("c2 funcval ptr: %x\n", *(*uintptr)(unsafe.Pointer(&c2)))
	fmt.Printf("c3 funcval ptr: %x\n", *(*uintptr)(unsafe.Pointer(&c3)))
	fmt.Printf("c1 vs c2 (same code): %v\n", reflect.ValueOf(c1).Pointer() == reflect.ValueOf(c2).Pointer())
	fmt.Printf("c1 vs c2 (same funcval): %v\n", *(*uintptr)(unsafe.Pointer(&c1)) == *(*uintptr)(unsafe.Pointer(&c2)))
	fmt.Printf("c1 vs c3 (same funcval): %v\n", *(*uintptr)(unsafe.Pointer(&c1)) == *(*uintptr)(unsafe.Pointer(&c3)))
}
