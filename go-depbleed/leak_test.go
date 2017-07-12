package depbleed

import (
	"errors"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"testing"
)

func TestLeakError(t *testing.T) {
	pkg := types.NewPackage("foo/bar", "bar")
	typename := types.NewTypeName(token.NoPos, pkg, "MyType", types.NewStruct(nil, nil))
	leak := Leak{
		Object: typename,
		err:    errors.New("fail"),
	}

	expected := "MyType: fail"
	err := leak.Error()

	if err != expected {
		t.Errorf("expected \"%s\" but got \"%s\"", expected, err)
	}
}

func TestLeaksSort(t *testing.T) {
	a1 := Leak{
		Position: token.Position{Filename: "a"},
	}
	b1 := Leak{
		Position: token.Position{Filename: "b"},
	}
	b2 := Leak{
		Position: token.Position{Filename: "b", Line: 2},
	}
	b3 := Leak{
		Position: token.Position{Filename: "b", Line: 2, Column: 3},
	}
	leaks := Leaks{b3, a1, b2, b1}
	expected := Leaks{a1, b1, b2, b3}

	if leaks.Len() != len(leaks) {
		t.Errorf("expected a length of %d, but got: %d", len(leaks), leaks.Len())
	}

	sort.Sort(leaks)

	if !reflect.DeepEqual(leaks, expected) {
		t.Errorf("expected:\n%v\ngot:\n%v", expected, leaks)
	}
}
