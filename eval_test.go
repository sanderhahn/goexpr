package main

import (
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
			t.Errorf("%s != %g", in, expected)
		}
	}

	test(t, statement, "10", 2)
	test(t, statement, "1+2*3", 5)
	test(t, statement, "(1)", 3)
	test(t, statement, "(1))", -1)
	test(t, statement, ")", -1)
	test(t, statement, "+", -1)
	test(t, statement, "(1+1)", 5)
	test(t, statement, "1+)1", -1)
	test(t, statement, "(1)+1)", -1)
	test(t, statement, "1==1==1", -1)

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

	calc.env["a"] = 5
	testEval("a", 5)

	testEvalVar := func(in string, name string, expected float64) {
		_, err := Eval(in)
		if err != nil {
			t.Error(err)
		}
		val, ok := calc.env[name]
		if !ok || val != expected {
			t.Errorf("%s != 5", name, expected)
		}
	}

	testEvalVar("a^a", "a", 5)
	testEvalVar("b=3", "b", 3)
	testEvalVar(" a += 1", "a", 6)
	testEvalVar("a=a==6", "a", 1)

	testEval(".5", .5)
	testEval("2.", 2.)
	testEval("1.5", 1.5)
}
