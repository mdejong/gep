package gene

import (
	"math"
	"reflect"
	"testing"

	bn "github.com/gmlewis/gep/functions/bool_nodes"
)

var nandTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false}, true},
	{[]bool{false, true}, true},
	{[]bool{true, false}, true},
	{[]bool{true, true}, false},
}

func validateNand(t *testing.T, g *Gene) {
	for i, n := range nandTests {
		r := g.EvalBool(n.in, bn.BoolAllGates)
		if r != n.out {
			t.Errorf("%v: nand.EvalBool(%#v, BoolAllGates) => %v, want %v", i, n.in, r, n.out)
		}
	}
}

func TestNand(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	validateNand(t, nand)
	nand = New("Or.And.Not.d0.Not.And.Or.d0.d0.d1.d1.d0.d1.d1.d1")
	validateNand(t, nand)
}

type mathTest struct {
	in  []float64
	out float64
}

var mathTests = []struct {
	gene  string
	tests []mathTest
}{
	{
		gene: "+.d0.d1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0",
		tests: []mathTest{
			mathTest{in: []float64{1.0, 2.0}, out: 3.0},
		},
	},
	{
		gene: "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0",
		tests: []mathTest{
			mathTest{in: []float64{0}, out: 0},
			mathTest{in: []float64{2.81}, out: -10.7061},
			mathTest{in: []float64{6}, out: -42},
			mathTest{in: []float64{7.043}, out: -56.646849},
			mathTest{in: []float64{8}, out: -72},
			mathTest{in: []float64{10}, out: -110},
			mathTest{in: []float64{11.38}, out: -140.8844},
			mathTest{in: []float64{12}, out: -156},
			mathTest{in: []float64{14}, out: -210},
			mathTest{in: []float64{15}, out: -240},
			mathTest{in: []float64{20}, out: -420},
		},
	},
	{
		gene: "-.*.*.*.d0./.d0.d0.d0.d0.d0.d0.d0",
		tests: []mathTest{
			mathTest{in: []float64{20.0}, out: 7980.0},
		},
	},
}

func validateMath(t *testing.T, g *Gene, in []float64, out float64) {
	r := g.EvalMath(in)
	if math.Abs(r-out) > 1e-10 {
		t.Errorf("%v: math.Eval(%#v) => %v, want %v", g, in, r, out)
	}
}

func TestMath(t *testing.T) {
	for _, v := range mathTests {
		g := New(v.gene)
		for _, n := range v.tests {
			validateMath(t, g, n.in, n.out)
		}
	}
}

func TestGetBoolArgOrder(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	got := nand.getBoolArgOrder(bn.BoolAllGates)
	want := [][]int{
		{1, 2}, {3, 4}, {5}, {6}, {7, 8}, {9, 10}, {11, 12}, nil, nil, nil, nil, nil, nil, nil, nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nand.GetBoolArgOrder() got %#v, want %#v", got, want)
	}
}

func TestDup(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	validateNand(t, nand) // Force evaluation
	g1 := nand.Dup()
	if err := CheckEqual(g1, nand); err != nil {
		t.Errorf("TestDup after Dup failed: g1 != nand: %v\n", err)
	}
	validateNand(t, nand) // Force evaluation
	validateNand(t, g1)

	g1 = New(mathTests[0].gene)
	test := mathTests[0].tests[0]
	validateMath(t, g1, test.in, test.out) // Force evaluation
	nand = g1.Dup()
	if err := CheckEqual(g1, nand); err != nil {
		t.Errorf("TestDup after Dup failed: g1 != nand: %v\n", err)
	}
	validateMath(t, g1, test.in, test.out) // Force evaluation
	validateMath(t, nand, test.in, test.out)
}

func TestMutate(t *testing.T) {
	headSize := 7
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 5
	funcs := []FuncWeight{
		{"Not", 1},
		{"And", 5},
		{"Or", 5},
	}
	g1 := RandomNew(headSize, tailSize, numTerminals, funcs)
	gn := g1.Dup()
	g1.Mutate()
	if err := CheckEqual(gn, g1); err == nil {
		t.Errorf("TestMutate failed: g1 == mux\n")
	}
}
