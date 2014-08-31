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
	return join(or, " | ")
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
	return loop.grammar.String() + "*"
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
	return fmt.Sprintf("[%s]", group.group)
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
	return fmt.Sprintf("%s ?", opt.grammar)
}

type notahead struct{ grammar grammar }

func (notahead notahead) parse(in string) int {
	n := notahead.grammar.parse(in)
	if n >= 0 {
		return -1
	}
	return 0
}

func (notahead notahead) String() string {
	return fmt.Sprintf("q(?! %s )", notahead.grammar)
}

type action struct {
	start   func() bool
	grammar grammar
	end     func(match string) bool
}

func (action action) parse(in string) int {
	if action.start != nil {
		if !action.start() {
			return -1
		}
	}
	n := action.grammar.parse(in)
	if n >= 0 {
		if action.end == nil {
			return n
		}
		if action.end(in[0:n]) {
			return n
		}
	}
	return -1
}

func (action action) String() string {
	return fmt.Sprintf("@(%s)", action.grammar)
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
	return "(" + strings.Join(more, sep) + ")"
}
