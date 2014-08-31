package main

import "testing"

func test(t *testing.T, grammar grammar, input string, ok int) {
	if n := grammar.parse(input); n != ok {
		t.Errorf("input %s (%d != %d)\ngrammar %s\n", input, n, ok, grammar)
	}
}

func Test(t *testing.T) {
	test(t, eof{}, "", 0)
	test(t, eof{}, "1", -1)
	test(t, str{"1"}, "1", 1)
	test(t, str{"1"}, "2", -1)
	test(t, str{"12"}, "12", 2)
	test(t, and{str{"1"}, str{"2"}, eof{}}, "12", 2)
	test(t, and{str{"1"}, str{"2"}, eof{}}, "11", -1)
	test(t, or{str{"1"}, str{"2"}, eof{}}, "1", 1)
	test(t, or{str{"1"}, str{"2"}, eof{}}, "3", -1)
	test(t, or{and{str{"1"}, str{"1"}}, and{str{"1"}, str{"2"}}}, "12", 2)
	test(t, loop{str{"1"}}, "", 0)
	test(t, loop{str{"1"}}, "111", 3)
	test(t, and{loop{str{"1"}}, err{"woeps"}}, "111", -1)
	test(t, rang{'a', 'z'}, "a", 1)
	test(t, rang{'a', 'z'}, "A", -1)
	test(t, group{"0123456789"}, "9", 1)
	test(t, group{"0123456789"}, "a", -1)
	test(t, str{"\u2318"}, "\u2318", 3)
	test(t, group{"\u2318"}, "\u2318", 3)
	test(t, loop1(str{"1"}), "", -1)
	test(t, loop1(str{"1"}), "1", 1)
	test(t, loop1(str{"1"}), "11", 2)
	test(t, and{str{"lang: "}, loop1(rang{'a', 'z'}), eof{}}, "lang: go", 8)
	test(t, opt{str{"1"}}, "", 0)
	test(t, opt{str{"1"}}, "1", 1)
	test(t, sep1(str{"1"}, str{","}), "1,1,1", 5)
	test(t, and{sep1(str{"1"}, str{","}), eof{}}, "1,,1", -1)

	// the first succesful or is chosen and no backtracking is done
	test(t, and{or{str{"*"}, str{"**"}}, eof{}}, "**", -1)
	// loop consumes greedy and this will always fail
	test(t, and{loop{str{"*"}}, str{"*"}, eof{}}, "***", -1)

	test(t, not{str{"1"}}, "1", -1)
	test(t, not{str{"1"}}, "2", 0)
	test(t, ahead{str{"1"}}, "1", 0)

	test(t, and{str{"1"}, notahead(str{"1"})}, "12", 1)
	test(t, and{str{"1"}, notahead(str{"1"})}, "11", -1)
}

func TestAction(t *testing.T) {
	var mymatch string
	match := func(match string) bool {
		mymatch = match
		return true
	}

	grammar := and{str{"lang: "}, action{match, loop1(rang{'a', 'z'})}, eof{}}
	if grammar.parse("lang: go") >= 0 && mymatch != "go" {
		t.Error("action failed")
	}
}
