package di_test

import (
	"github.com/marczahn/simple-go-di/pkg/di"
	"testing"
)

func TestGetOrSetTypes(t *testing.T) {
	di.ResetReg()
	number := di.GetOrSet("test", func() int { return 1 }, false)
	if number != 1 {
		t.Errorf("number failed - want 1, got %v", number)
	}

	di.ResetReg()
	str := di.GetOrSet("test", func() string { return "some string" }, false)
	if str != "some string" {
		t.Errorf("string failed - want 1, got %v", str)
	}

	di.ResetReg()
	aIn := [3]int{1, 2, 3}
	a := di.GetOrSet("test", func() [3]int { return aIn }, false)
	if a != aIn {
		t.Errorf("array failed - want %v, got %v", aIn, a)
	}

	di.ResetReg()
	type T struct {
		field string
	}
	sIn := T{field: "value"}
	s := di.GetOrSet("test", func() T { return sIn }, false)
	if s != sIn {
		t.Errorf("struct failed - want %v, got %v", sIn, s)
	}

	di.ResetReg()
	ptrIn := &T{field: "value"}
	ptr := di.GetOrSet("test", func() *T { return ptrIn }, false)
	if *ptr != *ptrIn {
		t.Errorf("pointer failed - want %v, got %v", ptrIn, ptr)
	}
}

func TestGetOrSetOverwrite(t *testing.T) {
	di.ResetReg()
	_ = di.GetOrSet("test", func() int { return 1 }, false)
	old := di.GetOrSet[int]("test", func() int { return 2 }, false)
	if old != 1 {
		t.Errorf("number overwrite: false failed - want 1, got %v", old)
	}
	n := di.GetOrSet("test", func() int { return 2 }, true)
	if n != 2 {
		t.Errorf("overwrite: true failed - want 2, got %v", n)
	}
}

func TestGetOrSetNestedDependency(t *testing.T) {
	di.ResetReg()
	const testValue = "some value"
	type nestedDep struct {
		someValue string
	}
	type outerDep struct {
		dep nestedDep
	}
	getNestedDep := func() nestedDep {
		return di.GetOrSet(
			"nested dep",
			func() nestedDep { return nestedDep{someValue: testValue} },
			true,
		)
	}
	dep := di.GetOrSet("dep", func() outerDep { return outerDep{dep: getNestedDep()} }, true)

	if dep.dep.someValue != testValue {
		t.Error("nested di failed")
	}
}

func TestGetOrSetConcurrency(t *testing.T) {
	di.ResetReg()
	for i := 0; i < 100; i++ {
		go func(i int) {
			actual := di.GetOrSet("test", func() int { return i }, true)
			if actual != i {
				t.Errorf("invalid concurrency test value - want %d, got %d", i, actual)
			}
		}(i)
	}
}
