package main

import (
	"log"
	"strconv"
)

type calculator interface {
	reset()
	number(float64)
	lookup(identifier string) bool
	operator(operator string) bool
	assignment(identifier string, operator string) bool
}

func expr(calc calculator) grammar {

	// Map grammar actions to calculator interface

	reset := func(_ string) bool {
		calc.reset()
		return true
	}

	number := func(number string) bool {
		value, err := strconv.ParseFloat(number, 64)
		if err != nil {
			log.Fatal(err)
		}
		calc.number(value)
		return true
	}

	lookup := func(identifier string) bool {
		return calc.lookup(identifier)
	}

	operator := func(operator string) bool {
		return calc.operator(operator)
	}

	var assignIdentifier string
	var assignOperator string

	assignment := func(_ string) bool {
		return calc.assignment(assignIdentifier, assignOperator)
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

	ws := opt{loop1(group{" \t"})}
	op := alts("+ - ** * / || | && & <= << < >= >> > != == % ^")
	assignOp := or{and{str{"="}, notahead(str{"="})}, alts("+= -= **= *= /= ||= |= &&= &= %= ^=")}

	digit := group{"0123456789"}
	alpha := or{rang{'a', 'z'}, rang{'A', 'Z'}}
	num := or{
		and{str{"."}, loop1(digit)},
		and{loop1(digit), opt{and{str{"."}, loop{digit}}}}}

	id := and{alpha, loop{or{alpha, digit}}}

	lparens := loop{action{operator, str{"("}}}
	rparens := loop{action{operator, str{")"}}}

	term := and{
		ws,
		lparens,
		ws,
		or{
			action{number, num},
			action{lookup, id},
		},
		ws,
		rparens}

	expr := sep1(term, and{ws, action{operator, op}})

	assign := and{
		ahead{and{id, ws, assignOp}},
		action{assignment,
			and{
				action{setAssignIdentifier, id},
				ws,
				action{setAssignOperator, assignOp},
				expr}}}

	statement := and{
		action{reset, epsilon},
		ws,
		or{
			assign,
			expr,
		},
		ws,
		eof{}}

	return statement
}
