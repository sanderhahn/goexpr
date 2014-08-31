package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type grammar interface {
	parse(string) int
	fmt.Stringer
}

var epsilon grammar = str{}

type eof struct{}

func (_ eof) parse(s string) int {
	if len(s) == 0 {
		return 0
	}
	return -1
}

func (_ eof) String() string {
	return "<eof>"
}

type str struct{ in string }

func (str str) parse(in string) int {
	if strings.HasPrefix(in, str.in) {
		return len(str.in)
	}
	return -1
}

func (str str) String() string {
	return fmt.Sprintf("%q", str.in)
}

type and []grammar

func (and and) parse(in string) int {
	pos := 0
	for _, next := range and {
		if n := next.parse(in[pos:]); n >= 0 {
			pos += n
		} else {
			return -1
		}
	}
	return pos
}

func (and and) String() string {
	return join(and, " ")
}

type or []grammar

func (or or) parse(in string) int {
	pos := 0
	for _, alt := range or {
		if n := alt.parse(in[pos:]); n >= 0 {
			return n
		}
	}
	return -1
}

func (or or) String() string {
	return join(or, " / ")
}

type loop struct {
	grammar
}

func (loop loop) parse(in string) int {
	pos := 0
	for pos < len(in) {
		// fmt.Printf("in %s\n", in[pos:])
		n := loop.grammar.parse(in[pos:])
		if n <= 0 {
			return pos
		}
		pos += n
	}
	return pos
}

func (loop loop) String() string {
	return "( " + loop.grammar.String() + " ) *"
}

type err struct {
	msg string
}

func (err err) parse(_ string) int {
	// fmt.Println(err.msg)
	return -1
}

func (_ err) String() string {
	return "<err>"
}

type rang struct {
	start, end rune
}

func (rang rang) parse(in string) int {
	if len(in) > 0 {
		r, width := utf8.DecodeRuneInString(in)
		if r == 1 {
			log.Fatal("invalid utf8")
		}
		if r >= rang.start && r <= rang.end {
			return width
		}
	}
	return -1
}

func (rang rang) String() string {
	return fmt.Sprintf("[%q-%q]", rang.start, rang.end)
}

type group struct {
	group string
}

func (group group) parse(in string) int {
	if len(in) > 0 {
		r, width := utf8.DecodeRuneInString(in)
		if r == 1 {
			log.Fatal("invalid utf8")
		}
		if strings.ContainsRune(group.group, r) {
			return width
		}
	}
	return -1
}

func (group group) String() string {
	return fmt.Sprintf("[%q]", group.group)
}

type opt struct{ grammar grammar }

func (opt opt) parse(in string) int {
	n := opt.grammar.parse(in)
	if n >= 0 {
		return n
	}
	return 0
}

func (opt opt) String() string {
	return "( " + opt.grammar.String() + " ) ?"
}

type not struct{ grammar grammar }

func (not not) parse(in string) int {
	n := not.grammar.parse(in)
	if n == -1 {
		return 0
	}
	return -1
}

func (not not) String() string {
	return "! ( " + not.grammar.String() + " )"
}

type ahead struct{ grammar grammar }

func (ahead ahead) parse(in string) int {
	n := ahead.grammar.parse(in)
	if n >= 0 {
		return 0
	}
	return -1
}

func (ahead ahead) String() string {
	return "> ( " + ahead.grammar.String() + " )"
}

type action struct {
	act     func(match string) bool
	grammar grammar
}

func (action action) parse(in string) int {
	n := action.grammar.parse(in)
	if n >= 0 {
		if action.act(in[:n]) {
			return n
		}
	}
	return -1
}

func (action action) String() string {
	return "@ ( " + action.grammar.String() + " )"
}

type ref struct {
	name string
	rule *grammar
}

func Ref(name string, rule *grammar) *ref {
	return &ref{name, rule}
}

func (ref ref) parse(in string) int {
	return (*ref.rule).parse(in)
}

func (ref ref) String() string {
	return "<" + ref.name + ">"
}

func notahead(grammar grammar) grammar {
	return ahead{not{grammar}}
}

func loop1(item grammar) grammar {
	return and{item, loop{item}}
}

func sep1(item grammar, sep grammar) grammar {
	return loop1(and{item, opt{sep}})
}

func alts(choices string) (or or) {
	for _, alt := range strings.Split(choices, " ") {
		or = append(or, str{alt})
	}
	return or
}

func join(grammar []grammar, sep string) string {
	more := []string{}
	for _, s := range grammar {
		more = append(more, s.String())
	}
	return "( " + strings.Join(more, sep) + " )"
}
