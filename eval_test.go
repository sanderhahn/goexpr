package main

import (
	"fmt"
	"math"
	"testing"
)

func TestExpr(t *testing.T) {

	testEval := func(in string, expected float64) {
		val, err := Eval(in)
		if err != nil {
			t.Error(err)
		}
		if val != expected {
			fmt.Printf("stack %v\nops %v\n", e.machine.stack, e.machine.ops)
			fmt.Printf("%s != %g\n", in, expected)
			t.Errorf("%s != %g", in, expected)
		}
	}

	test(t, e.expr, "10", 2)
	test(t, e.expr, "1+2*3", 5)
	test(t, e.expr, "(1)", 3)
	test(t, e.expr, "(1))", -1)
	test(t, e.expr, ")", -1)
	test(t, e.expr, "+", -1)
	test(t, e.expr, "(1+1)", 5)
	test(t, e.expr, "1+)1", -1)
	test(t, e.expr, "(1)+1)", -1)
	test(t, e.expr, "1==1==1", -1)

	testEval("1-2", 1-2)
	testEval("1+2", 1+2)
	testEval("1/2", 1./2)
	testEval("1*2", 1*2)
	testEval("3%2", 3%2)
	testEval("1/2/4", 1./2./4.)
	testEval("1+2*3", 1+2*3)
	testEval("(1+2)*(1+2)", (1+2)*(1+2))
	testEval("4**3**2", math.Pow(4, math.Pow(3, 2)))

	testEval("1-2", 1-2)

	testEval("1==1", 1)
	testEval("1!=0", 1)
	testEval("1>=1", 1)
	testEval("1>0", 1)
	testEval("1<=1", 1)
	testEval("1<0", 0)

	testEval("1&&0", 0)
	testEval("0||1", 1)

	testEval("1<<1", 2)
	testEval("2>>1", 1)

	testEval("4&3", 4&3)
	testEval("4^3", 4^3)
	testEval("1|2", 1|2)
	testEval("3&2", 3&2)

	testEval(" 1 + 2 + 3 ", 6)
	testEval(" ( 1 ) ", 1)
	testEval(" ( 1 + 2 ) ", 3)
	testEval(" 1 + ( 2 + 3 ) ", 6)

	e.env["a"] = 5
	testEval("a", 5)

	testEvalVar := func(in string, name string, expected float64) {
		_, err := Eval(in)
		if err != nil {
			t.Error(err)
		}
		val, ok := e.env[name]
		if !ok || val != expected {
			t.Errorf("%s != %g", name, expected)
		}
	}

	testEvalVar("a^a", "a", 5)
	testEvalVar("b=3", "b", 3)
	testEvalVar(" a += 1", "a", 6)
	testEvalVar("a=a==6", "a", 1)
	testEval(" a += 1", 2)
	testEvalVar("c=1", "c", 1)
	testEvalVar("c+=2*3", "c", 7)

	testEval(".5", .5)
	testEval("2.", 2.)
	testEval("1.5", 1.5)
}
