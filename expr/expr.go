package expr

import (
	"log"
	"strconv"

	"github.com/sanderhahn/goexpr/grammar"
	. "github.com/sanderhahn/goexpr/grammar"
)

type Calculator interface {
	Reset()
	Number(float64)
	Lookup(identifier string) bool
	Operator(operator string) bool
	Assignment(identifier string, operator string) bool
}

func Parser(calc Calculator) grammar.Grammar {

	// Map grammar actions to calculator interface

	reset := func(_ string) bool {
		calc.Reset()
		return true
	}

	number := func(number string) bool {
		value, err := strconv.ParseFloat(number, 64)
		if err != nil {
			log.Fatal(err)
		}
		calc.Number(value)
		return true
	}

	lookup := func(identifier string) bool {
		return calc.Lookup(identifier)
	}

	operator := func(operator string) bool {
		return calc.Operator(operator)
	}

	var assignIdentifier string
	var assignOperator string

	assignment := func(_ string) bool {
		return calc.Assignment(assignIdentifier, assignOperator)
	}

	setAssignIdentifier := func(identifier string) bool {
		assignIdentifier = identifier
		return true
	}

	setAssignOperator := func(operator string) bool {
		assignOperator = operator
		return true
	}

	// Grammar

	ws := Opt(Loop1(Group(" \t")))
	op := Alts("+ - ** * / || | && & <= << < >= >> > != == % ^")
	assignOp := Or(And(Str("="), Notahead(Str("="))), Alts("+= -= **= *= /= ||= |= &&= &= %= ^="))

	digit := Group("0123456789")
	alpha := Or(Rang('a', 'z'), Rang('A', 'Z'))
	num := Or(
		And(Str("."), Loop1(digit)),
		And(Loop1(digit), Opt(And(Str("."), Loop(digit)))))

	id := And(alpha, Loop(Or(alpha, digit)))

	lparens := Loop(Action(operator, Str("(")))
	rparens := Loop(Action(operator, Str(")")))

	term := And(
		ws,
		lparens,
		ws,
		Or(
			Action(number, num),
			Action(lookup, id),
		),
		ws,
		rparens)

	expr := Sep1(term, And(ws, Action(operator, op)))

	assign := And(
		Ahead(And(id, ws, assignOp)),
		Action(assignment,
			And(
				Action(setAssignIdentifier, id),
				ws,
				Action(setAssignOperator, assignOp),
				expr)))

	statement := And(
		Action(reset, Epsilon()),
		ws,
		Or(
			assign,
			expr,
		),
		ws,
		Eof())

	return statement
}
