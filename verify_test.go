package verify_test

import (
	"strings"
	"testing"

	"github.com/codyoss/verify"
)

func TestItInputTypes(t *testing.T) {
	var va Aer = valueAer{}
	var pa Aer = &pointerAer{}

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"bool input", true, true},
		{"string input", "", true},
		{"int input", 0, true},
		{"float64 input", 1.4, true},
		{"*string input", new(string), true},
		{"struct input", valueAer{}, false},
		{"*struct input", &pointerAer{}, false},
		{"interface struct input", va, false},
		{"interface *struct input", pa, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verify.It(tt.input)
			if (tt.wantErr && got == nil) || (!tt.wantErr && got != nil) {
				t.Errorf("wantErr is %v, while got is %v", tt.wantErr, got)
			}
		})
	}
}

func TestItMinSize(t *testing.T) {
	type A struct {
		A bool `verify:"minSize"`
	}
	type B struct {
		A bool `verify:"minSize=abc"`
	}
	type C struct {
		A bool `verify:"minSize=5"`
	}
	type D struct {
		A string `verify:"minSize=5"`
	}

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"missing value", A{}, true},
		{"can't parse value", B{}, true},
		{"field wrong type", C{}, true},
		{"field too short", D{}, true},
		{"works", D{"12345"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verify.It(tt.input)
			if (tt.wantErr && got == nil) || (!tt.wantErr && got != nil) {
				t.Errorf("wantErr is %v, while got is %v", tt.wantErr, got)
			}
		})
	}
}

func TestItMaxSize(t *testing.T) {
	type A struct {
		A bool `verify:"maxSize"`
	}
	type B struct {
		A bool `verify:"maxSize=abc"`
	}
	type C struct {
		A bool `verify:"maxSize=5"`
	}
	type D struct {
		A string `verify:"maxSize=5"`
	}

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"missing value", A{}, true},
		{"can't parse value", B{}, true},
		{"field wrong type", C{}, true},
		{"field too long", D{"123456"}, true},
		{"works", D{"12345"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verify.It(tt.input)
			if (tt.wantErr && got == nil) || (!tt.wantErr && got != nil) {
				t.Errorf("wantErr is %v, while got is %v", tt.wantErr, got)
			}
		})
	}
}

func TestItMin(t *testing.T) {
	type A struct {
		A bool `verify:"min"`
	}
	type B struct {
		A bool `verify:"min=abc"`
	}
	type C struct {
		A bool `verify:"min=3"`
	}
	type D struct {
		A float64 `verify:"min=2"`
	}
	type E struct {
		A int64 `verify:"min=2.1"`
	}
	type F struct {
		A float64 `verify:"min=2.1"`
	}
	type G struct {
		A int64 `verify:"min=2"`
	}

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"missing value", A{}, true},
		{"can't parse value", B{}, true},
		{"field wrong type", C{}, true},
		{"field type does not match tag type", D{3.1}, true},
		{"tag type does not match field type", E{3}, true},
		{"too small float", F{1.1}, true},
		{"too small int", G{1}, true},
		{"works float", F{3.1}, false},
		{"works int", G{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verify.It(tt.input)
			if (tt.wantErr && got == nil) || (!tt.wantErr && got != nil) {
				t.Errorf("wantErr is %v, while got is %v", tt.wantErr, got)
			}
		})
	}
}

func TestItMax(t *testing.T) {
	type A struct {
		A bool `verify:"max"`
	}
	type B struct {
		A bool `verify:"max=abc"`
	}
	type C struct {
		A bool `verify:"max=3"`
	}
	type D struct {
		A float64 `verify:"max=2"`
	}
	type E struct {
		A int64 `verify:"max=2.1"`
	}
	type F struct {
		A float64 `verify:"max=2.1"`
	}
	type G struct {
		A int64 `verify:"max=2"`
	}

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"missing value", A{}, true},
		{"can't parse value", B{}, true},
		{"field wrong type", C{}, true},
		{"field type does not match tag type", D{1.1}, true},
		{"tag type does not match field type", E{1}, true},
		{"too large float", F{3.1}, true},
		{"too large int", G{3}, true},
		{"works float", F{1.1}, false},
		{"works int", G{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verify.It(tt.input)
			if (tt.wantErr && got == nil) || (!tt.wantErr && got != nil) {
				t.Errorf("wantErr is %v, while got is %v", tt.wantErr, got)
			}
		})
	}
}

func TestItRequired(t *testing.T) {

	type Zero struct {
		A string
	}
	type A struct {
		A string `verify:"required"`
	}
	type B struct {
		A bool `verify:"required"`
	}
	type C struct {
		A int `verify:"required"`
	}
	type D struct {
		A float64 `verify:"required"`
	}
	type E struct {
		A Zero `verify:"required"`
	}
	type F struct {
		A *Zero `verify:"required"`
	}
	type G struct {
		A []int `verify:"required"`
	}

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"deafult string", A{}, true},
		{"deafult bool", B{}, true},
		{"deafult int", C{}, true},
		{"deafult float64", D{}, true},
		{"deafult struct", E{}, false},
		{"deafult *struct", F{}, true},
		{"deafult slice", G{}, true},
		{"non-deafult string", A{"a"}, false},
		{"non-deafult bool", B{true}, false},
		{"non-deafult int", C{1}, false},
		{"non-deafult float64", D{1.1}, false},
		{"non-deafult struct", E{Zero{"a"}}, false},
		{"non-deafult *struct", F{&Zero{"a"}}, false},
		{"non-deafult slice", G{[]int{1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verify.It(tt.input)
			if (tt.wantErr && got == nil) || (!tt.wantErr && got != nil) {
				t.Errorf("wantErr is %v, while got is %v", tt.wantErr, got)
			}
		})
	}
}

func TestItMultipleValidationsFail(t *testing.T) {
	type A struct {
		A int `verify:"required,max=-1"`
	}

	err := verify.It(A{})
	if err == nil || !strings.Contains(err.Error(), "zero value") || !strings.Contains(err.Error(), "greater than max") {
		t.Error("expected err two contain two messages")
	}

}

type Aer interface {
	A()
}

type valueAer struct{}

func (a valueAer) A() {}

type pointerAer struct{}

func (a *pointerAer) A() {}
