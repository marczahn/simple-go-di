package di_test

import (
	"github.com/marczahn/simple-go-di/pkg/di"
	"testing"
)

func TestGetOrSetTypes(t *testing.T) {
	subjectI := di.NewInstance[int]()
	number := subjectI.GetOrSet(func() int { return 1 }, false)
	if number != 1 {
		t.Errorf("number failed - want 1, got %v", number)
	}
	subjectS := di.NewInstance[string]()
	str := subjectS.GetOrSet(func() string { return "string" }, false)
	if str != "string" {
		t.Errorf("string failed - want string, got %v", str)
	}

	subjectSl := di.NewInstance[[3]int]()
	aIn := [3]int{1, 2, 3}
	a := subjectSl.GetOrSet(func() [3]int { return aIn }, false)
	if a != aIn {
		t.Errorf("array failed - want %v, got %v", aIn, a)
	}

	type T struct {
		field string
	}
	subjectStruct := di.NewInstance[T]()
	sIn := T{field: "value"}
	s := subjectStruct.GetOrSet(func() T { return sIn }, false)
	if s != sIn {
		t.Errorf("struct failed - want %v, got %v", sIn, s)
	}

	ptrIn := &T{field: "value"}
	subjectPtr := di.NewInstance[*T]()
	ptr := subjectPtr.GetOrSet(func() *T { return ptrIn }, false)
	if *ptr != *ptrIn {
		t.Errorf("pointer failed - want %v, got %v", ptrIn, ptr)
	}
}

func TestGetOrSetOverwrite(t *testing.T) {
	subject := di.NewInstance[int]()
	_ = subject.GetOrSet(func() int { return 1 }, false)
	old := subject.GetOrSet(func() int { return 2 }, false)
	if old != 1 {
		t.Errorf("number overwrite: false failed - want 1, got %v", old)
	}
	n := subject.GetOrSet(func() int { return 2 }, true)
	if n != 2 {
		t.Errorf("overwrite: true failed - want 2, got %v", n)
	}
}

func TestGetOrSetNestedDependency(t *testing.T) {
	const testValue = "some value"
	type nestedDep struct {
		someValue string
	}
	type outerDep struct {
		dep nestedDep
	}
	getNestedDep := func() nestedDep {
		subject := di.NewInstance[nestedDep]()
		return subject.GetOrSet(
			func() nestedDep { return nestedDep{someValue: testValue} },
			true,
		)
	}
	subjectOuter := di.NewInstance[outerDep]()
	dep := subjectOuter.GetOrSet(func() outerDep { return outerDep{dep: getNestedDep()} }, true)

	if dep.dep.someValue != testValue {
		t.Error("nested di failed")
	}
}

func TestGetOrSetConcurrency(t *testing.T) {
	subject := di.NewInstance[int]()
	for i := 0; i < 100; i++ {
		go func(i int) {
			actual := subject.GetOrSet(func() int { return i }, true)
			if actual != i {
				t.Errorf("invalid concurrency test value - want %d, got %d", i, actual)
			}
		}(i)
	}
}
