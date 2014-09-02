package expr

import (
	"log"

	"github.com/sanderhahn/goexpr/ast"
	. "github.com/sanderhahn/goexpr/grammar"
)

type ExpressionParser interface {
	Parse(input string) ast.Node
}

type expressionParser struct {
	grammar Grammar
	builder ast.AstBuilder
	opstack []string
}

func Parser() ExpressionParser {

	e := expressionParser{}

	reset := func(_ string) bool {
		e.builder.Reset()
		e.opstack = e.opstack[:0]
		return true
	}

	token := func(typ ast.TokenType, rule Grammar) Grammar {
		return Action(func(value string) bool {
			e.builder.Token(typ, value)
			return true
		}, rule)
	}

	operator := func(operator string) bool {
		return e.pushOp(operator)
	}

	paren := func(paren string) bool {
		e.pushOp(paren)
		return true
	}

	finish := func(_ string) bool {
		for len(e.opstack) > 0 {
			e.popOp()
		}
		return true
	}

	// Grammar

	ws := Opt(Loop1(Group(" \t")))
	op := Alts("+= -= **= *= /= ||= |= &&= &= %= ^= + - ** * / || | && & <= << < >= >> > != == % ^ =")

	digit := Group("0123456789")
	alpha := Or(Rang('a', 'z'), Rang('A', 'Z'))
	num := Or(
		And(Str("."), Loop1(digit)),
		And(Loop1(digit), Opt(And(Str("."), Loop(digit)))))

	id := And(alpha, Loop(Or(alpha, digit)))

	var expr Grammar
	refexpr := Ref("expr", &expr)

	term := And(
		ws,
		Or(
			token(ast.Number, num),
			token(ast.Identifier, id),
			And(Action(paren, Str("(")), refexpr, ws, Action(paren, Str(")")))))

	expr = Sep1(term, And(ws, Action(operator, op)))

	statement := And(
		Action(reset, Epsilon()),
		ws,
		expr,
		ws,
		Action(finish, Eof()))

	e.grammar = statement
	return &e
}

func (e *expressionParser) Parse(input string) ast.Node {
	n := e.grammar.Parse(input)
	if n >= 0 {
		return e.builder.Root()
	}
	return nil
}

func (e *expressionParser) peekOp() string {
	return e.opstack[len(e.opstack)-1]
}

func (e *expressionParser) pushOp(op string) bool {

	if op == "(" {
		e.opstack = append(e.opstack, op)
		return true
	} else if op == ")" {
		for e.peekOp() != "(" {
			e.popOp()
		}
		e.popOp()
		return true
	}

	for len(e.opstack) > 0 {
		assoc := operatorAssoc(e.peekOp(), op)
		if assoc == non {
			return false
		} else if assoc == left {
			e.popOp()
		} else {
			break
		}
	}

	e.opstack = append(e.opstack, op)
	return true
}

func operatorAssoc(last, operator string) assoc {
	prioOp, assocOp := operatorInfo(operator)
	prioLast, assocLast := operatorInfo(last)

	// Operators of equal precedence that are non-associative cannot be used next to each other
	if prioLast == prioOp && assocLast == non && assocOp == non {
		return non
	}

	if (assocOp == left && prioOp >= prioLast) ||
		(assocOp == right && prioOp > prioLast) ||
		(assocOp == non && prioOp > prioLast) {
		return left
	}
	return right
}

func (e *expressionParser) popOp() {
	l := len(e.opstack) - 1
	op := e.opstack[l]
	e.opstack = e.opstack[:l]
	if op != "(" {
		e.builder.Token(ast.Operator, op)
		e.builder.Evaluate()
	} else {
		e.builder.Group()
	}
}

type assoc int

const (
	left assoc = iota
	right
	non
)

func operatorInfo(op string) (prio int, assoc assoc) {

	switch op {
	case "*", "/", "%", "<<", ">>", "&":
		return 1, left
	case "**":
		return 1, right

	case "+", "-", "|", "^":
		return 2, left

	case "==", "!=", "<", "<=", ">", ">=":
		return 3, left

	case "&&":
		return 4, left

	case "||":
		return 5, left

	case "=", "+=", "-=", "**=", "*=", "/=", "||=", "|=", "&&=", "&=", "%=", "^=":
		return 9, non

	case "(":
		return 10, non
	}

	log.Fatalf("operator %s has no priority", op)
	return -1, non
}
