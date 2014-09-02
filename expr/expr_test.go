package expr

import "testing"

func TestParse(t *testing.T) {

	e := Parser()

	test := func(input string, expected bool) {
		node := e.Parse(input)
		ok := node != nil
		if ok != expected {
			t.Errorf("input %s should match %v\ngrammar %s\n", input, expected, e.(*expressionParser).grammar)
		}
	}

	test("10", true)
	test("1+2*3", true)
	test("(1)", true)
	test("(1))", false)
	test(")", false)
	test("+", false)
	test("(1+1)", true)
	test("1+)1", false)
	test("(1)+1)", false)
	test("1==1==1", true)
	test("a+=1=1", false)
}
