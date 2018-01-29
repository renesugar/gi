package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test086SlicesPointToArrays(t *testing.T) {

	cv.Convey(`two slices of the same array should share the same memory`, t, func() {

		code := `
   a := [2]int64{1,3}
   b := a[:]
   c := a[1:]
   b[1]++
   c0 := c[0]
   a1 := a[1]
`
		// c0 should be 4, a1 should be 4
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
	a = _gi_NewArray({[0]=1LL, 3LL}, "int64", 2, 0LL);
  	b = _gi_NewSlice("int64", a);
  	c = __subslice(_gi_NewSlice("int64", a), 1);
  	_gi_SetRangeCheck(b, 1, (_gi_GetRangeCheck(b, 1) + (1LL)));
  	c0 = _gi_GetRangeCheck(c, 0);
  	a1 = a[1];
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "c0", 4)
		LuaMustInt64(vm, "a1", 4)

	})
}

func Test088SlicesFromArrays(t *testing.T) {

	cv.Convey(`a slices from an array should work against an array proxy`, t, func() {

		code := `
   a := [2]int64{88,99}
   b := a[:]
   b0 := b[0]
   b1 := b[1]
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
	a = _gi_NewArray({[0]=88LL, 99LL}, "int64", 2, 0LL);
  	b = _gi_NewSlice("int64", a);
 	b0 = _gi_GetRangeCheck(b, 0);
   	b1 = _gi_GetRangeCheck(b, 1);
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "b0", 88)
		LuaMustInt64(vm, "b1", 99)

	})
}