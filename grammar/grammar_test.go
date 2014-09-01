package grammar

import (
	"fmt"
	"testing"
)

func test(t *testing.T, grammar Grammar, input string, ok int) {
	if n := grammar.Parse(input); n != ok {
		t.Errorf("input %s (%d != %d)\ngrammar %s\n", input, n, ok, grammar)
	}
}

func Test(t *testing.T) {
	test(t, Epsilon(), "1", 0)
	test(t, Eof(), "", 0)
	test(t, eof{}, "1", -1)
	test(t, Str("1"), "1", 1)
	test(t, str{"1"}, "2", -1)
	test(t, str{"12"}, "12", 2)
	test(t, And(str{"1"}, str{"2"}, eof{}), "12", 2)
	test(t, And(str{"1"}, str{"2"}, eof{}), "11", -1)
	test(t, Or(str{"1"}, str{"2"}, eof{}), "1", 1)
	test(t, Or(str{"1"}, str{"2"}, eof{}), "3", -1)
	test(t, Or(And(str{"1"}, str{"1"}), And(str{"1"}, str{"2"})), "12", 2)
	test(t, Loop(str{"1"}), "", 0)
	test(t, loop{str{"1"}}, "111", 3)
	test(t, And(loop{str{"1"}}, Err("woeps")), "111", -1)
	test(t, Rang('a', 'z'), "a", 1)
	test(t, rang{'a', 'z'}, "A", -1)
	test(t, Group("0123456789"), "9", 1)
	test(t, group{"0123456789"}, "a", -1)
	test(t, str{"\u2318"}, "\u2318", 3)
	test(t, group{"\u2318"}, "\u2318", 3)
	test(t, Loop1(str{"1"}), "", -1)
	test(t, Loop1(str{"1"}), "1", 1)
	test(t, Loop1(str{"1"}), "11", 2)
	test(t, And(str{"lang: "}, Loop1(rang{'a', 'z'}), eof{}), "lang: go", 8)
	test(t, Opt(str{"1"}), "", 0)
	test(t, opt{str{"1"}}, "1", 1)
	test(t, Sep1(str{"1"}, str{","}), "1,1,1", 5)
	test(t, And(Sep1(str{"1"}, str{","}), eof{}), "1,,1", -1)

	// the first succesful or is chosen and no backtracking is done
	test(t, And(Or(str{"*"}, str{"**"}), eof{}), "**", -1)
	// loop consumes greedy and this will always fail
	test(t, And(loop{str{"*"}}, str{"*"}, eof{}), "***", -1)

	test(t, Not(str{"1"}), "1", -1)
	test(t, not{str{"1"}}, "2", 0)
	test(t, Ahead(str{"1"}), "1", 0)

	test(t, And(str{"1"}, Notahead(str{"1"})), "12", 1)
	test(t, And(str{"1"}, Notahead(str{"1"})), "11", -1)

	test(t, Alts("1 2"), "2", 1)

	s := Str("")
	fmt.Sprintf("%v", Loop(Or(And(Str(""), Action(nil, Not(Rang('a', 'z'))), Ahead(Opt(Group("1"))), Ref("x", &s), Eof(), Err("")))))
}

func TestAction(t *testing.T) {
	var mymatch string
	match := func(match string) bool {
		mymatch = match
		return true
	}

	grammar := And(str{"lang: "}, Action(match, Loop1(rang{'a', 'z'})), eof{})
	if grammar.Parse("lang: go") >= 0 && mymatch != "go" {
		t.Error("action failed")
	}
}

func TestRecurse(t *testing.T) {
	var parens Grammar

	refparens := Ref("parens", &parens)
	parens = Or(
		And(str{"("}, refparens, str{")"}),
		str{"x"},
	)

	test(t, parens, "(((x)))", 7)
	test(t, parens, "(((x))", -1)
}
