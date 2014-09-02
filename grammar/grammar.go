package grammar

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type Grammar interface {
	Parse(string) int
	fmt.Stringer
}

var epsilon Grammar = str{}

func Epsilon() Grammar {
	return epsilon
}

type eof struct{}

func Eof() Grammar {
	return &eof{}
}

func (_ eof) Parse(s string) int {
	if len(s) == 0 {
		return 0
	}
	return -1
}

func (_ eof) String() string {
	return "<eof>"
}

type str struct {
	in string
}

func Str(in string) Grammar {
	return str{in}
}

func (str str) Parse(in string) int {
	if strings.HasPrefix(in, str.in) {
		return len(str.in)
	}
	return -1
}

func (str str) String() string {
	return fmt.Sprintf("%q", str.in)
}

type and struct {
	seq []Grammar
}

func And(seq ...Grammar) Grammar {
	return and{seq}
}

func (and and) Parse(in string) int {
	pos := 0
	for _, next := range and.seq {
		if n := next.Parse(in[pos:]); n >= 0 {
			pos += n
		} else {
			return -1
		}
	}
	return pos
}

func (and and) String() string {
	return join(and.seq, " ")
}

type or struct {
	alts []Grammar
}

func Or(alts ...Grammar) Grammar {
	return or{alts}
}

func (or or) Parse(in string) int {
	pos := 0
	for _, alt := range or.alts {
		if n := alt.Parse(in[pos:]); n >= 0 {
			return n
		}
	}
	return -1
}

func (or or) String() string {
	return join(or.alts, " / ")
}

type loop struct {
	Grammar
}

func Loop(rule Grammar) Grammar {
	return loop{rule}
}

func (loop loop) Parse(in string) int {
	pos := 0
	for pos < len(in) {
		n := loop.Grammar.Parse(in[pos:])
		if n <= 0 {
			return pos
		}
		pos += n
	}
	return pos
}

func (loop loop) String() string {
	return "( " + loop.Grammar.String() + " ) *"
}

type err struct {
	msg string
}

func Err(msg string) Grammar {
	return err{msg}
}

func (err err) Parse(_ string) int {
	// fmt.Println(err.msg)
	return -1
}

func (_ err) String() string {
	return "<err>"
}

type rang struct {
	start, end rune
}

func Rang(start, end rune) Grammar {
	return rang{start, end}
}

func (rang rang) Parse(in string) int {
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

func Group(match string) Grammar {
	return group{match}
}

func (group group) Parse(in string) int {
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

type opt struct {
	rule Grammar
}

func Opt(rule Grammar) Grammar {
	return opt{rule}
}

func (opt opt) Parse(in string) int {
	n := opt.rule.Parse(in)
	if n >= 0 {
		return n
	}
	return 0
}

func (opt opt) String() string {
	return "( " + opt.rule.String() + " ) ?"
}

type not struct {
	rule Grammar
}

func Not(rule Grammar) Grammar {
	return not{rule}
}

func (not not) Parse(in string) int {
	n := not.rule.Parse(in)
	if n == -1 {
		return 0
	}
	return -1
}

func (not not) String() string {
	return "! ( " + not.rule.String() + " )"
}

type ahead struct {
	rule Grammar
}

func Ahead(rule Grammar) Grammar {
	return ahead{rule}
}

func (ahead ahead) Parse(in string) int {
	n := ahead.rule.Parse(in)
	if n >= 0 {
		return 0
	}
	return -1
}

func (ahead ahead) String() string {
	return "> ( " + ahead.rule.String() + " )"
}

type action struct {
	act  func(match string) bool
	rule Grammar
}

func Action(act func(match string) bool, rule Grammar) Grammar {
	return action{act, rule}
}

func (action action) Parse(in string) int {
	n := action.rule.Parse(in)
	if n >= 0 {
		if action.act(in[:n]) {
			return n
		}
	}
	return -1
}

func (action action) String() string {
	return "@ ( " + action.rule.String() + " )"
}

type ref struct {
	name string
	rule *Grammar
}

func Ref(name string, rule *Grammar) *ref {
	return &ref{name, rule}
}

func (ref ref) Parse(in string) int {
	return (*ref.rule).Parse(in)
}

func (ref ref) String() string {
	return "<" + ref.name + ">"
}

func Notahead(rule Grammar) Grammar {
	return ahead{not{rule}}
}

func Loop1(item Grammar) Grammar {
	return And(item, loop{item})
}

func Sep1(item Grammar, sep Grammar) Grammar {
	return Loop1(And(item, opt{sep}))
}

func Alts(choices string) Grammar {
	alts := []Grammar{}
	for _, alt := range strings.Split(choices, " ") {
		alts = append(alts, str{alt})
	}
	return Or(alts...)
}

func join(rule []Grammar, sep string) string {
	more := []string{}
	for _, s := range rule {
		more = append(more, s.String())
	}
	return "( " + strings.Join(more, sep) + " )"
}
