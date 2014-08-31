package main

import (
	"errors"
	"log"
	"math"
	"strconv"
)

type eval struct {
	count      int
	stack      []float64
	ops        []string
	assign     string
	assignOp   string
	assignment bool
	env        map[string]float64
}

func newEval() *eval {
	return &eval{
		env: map[string]float64{},
	}
}

func (eval *eval) reset() bool {
	eval.count = 0
	eval.stack = eval.stack[:0]
	eval.ops = []string{}
	return true
}

func (eval *eval) op(in string) bool {
	inPrio := eval.opPrio(in)
	for len(eval.ops) > 0 {
		top := eval.ops[len(eval.ops)-1]
		topPrio := eval.opPrio(top)

		assoc := eval.opAssoc(in)
		topAssoc := eval.opAssoc(top)
		if assoc == non && topAssoc == non {
			return false
		}

		if (assoc == left && inPrio >= topPrio) ||
			(assoc == right && inPrio > topPrio) ||
			(assoc == non && inPrio > topPrio) {
			eval.evalOnce()
		} else {
			break
		}
	}

	eval.ops = append(eval.ops, in)
	return true
}

func (eval *eval) num(in string) bool {
	f, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Fatal(err)
	}
	eval.stack = append(eval.stack, f)
	return true
}

func (eval *eval) id(in string) bool {
	v, ok := eval.env[in]
	eval.stack = append(eval.stack, v)
	return ok
}

func btof(b bool) float64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func ftob(f float64) bool {
	if f != 0 {
		return true
	} else {
		return false
	}
}

func (eval *eval) evalOnce() {
	// log.Printf("eval: %v %v\n", eval.stack, eval.ops)

	op := eval.ops[len(eval.ops)-1]
	eval.ops = eval.ops[:len(eval.ops)-1]

	if op == "(" || op == ")" {
		return
	}

	l := eval.stack[len(eval.stack)-2]
	r := eval.stack[len(eval.stack)-1]
	eval.stack = eval.stack[:len(eval.stack)-2]

	var v float64
	switch op {
	case "+":
		v = l + r
	case "*":
		v = l * r
	case "-":
		v = l - r
	case "/":
		v = l / r
	case "%":
		v = math.Mod(l, r)
	case "**":
		v = math.Pow(l, r)
	case "<<":
		v = float64(int64(l) << uint(r))
	case ">>":
		v = float64(int64(l) >> uint(r))
	case "&":
		v = float64(int64(l) & int64(r))
	case "|":
		v = float64(int64(l) | int64(r))
	case "^":
		v = float64(int64(l) ^ int64(r))
	case "==":
		v = btof(l == r)
	case "!=":
		v = btof(l != r)
	case "<":
		v = btof(l < r)
	case "<=":
		v = btof(l <= r)
	case ">":
		v = btof(l > r)
	case ">=":
		v = btof(l >= r)
	case "&&":
		v = btof(ftob(l) && ftob(r))
	case "||":
		v = btof(ftob(l) != ftob(r))
	}
	eval.stack = append(eval.stack, v)
}

func (eval *eval) opPrio(op string) int {

	switch op {
	case "*":
		return 1
	case "/":
		return 1
	case "%":
		return 1
	case "<<":
		return 1
	case ">>":
		return 1
	case "&":
		return 1
	case "**":
		return 1

	case "+":
		return 2
	case "-":
		return 2
	case "|":
		return 2
	case "^":
		return 2

	case "==":
		return 3
	case "!=":
		return 3
	case "<":
		return 3
	case "<=":
		return 3
	case ">":
		return 3
	case ">=":
		return 3

	case "&&":
		return 4

	case "||":
		return 5

	case "(":
		return 10
	}
	log.Fatalf("operator %s has no priority", op)
	return -1
}

type assoc int

const (
	left assoc = iota
	right
	non
)

func (eval *eval) opAssoc(op string) assoc {
	if op == "**" {
		return right
	}
	if op == "==" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">=" {
		return non
	}
	return left
}

func (eval *eval) paren(in string) bool {
	if in == "(" {
		eval.ops = append(eval.ops, "(")
		eval.count++
	} else if in == ")" {
		eval.ops = append(eval.ops, ")")
		eval.count--
		for len(eval.ops) > 0 && eval.ops[len(eval.ops)-1] != "(" {
			eval.evalOnce()
		}
		if len(eval.ops) > 0 {
			eval.ops = eval.ops[:len(eval.ops)-1]
		}
	}
	return eval.count >= 0
}

func (eval *eval) resetAssign() bool {
	eval.assignment = false
	return true
}

func (eval *eval) setAssign(in string) bool {
	eval.assign = in
	return true
}

func (eval *eval) setAssignOp(in string) bool {
	eval.assignOp = in
	return true
}

func (eval *eval) useAssign(in string) bool {
	eval.assignment = true
	return true
}

func (eval *eval) balanced(_ string) bool {
	for len(eval.ops) > 0 {
		eval.evalOnce()
	}
	if eval.assignment {
		val := eval.stack[len(eval.stack)-1]
		if eval.assignOp != "=" {
			prev, ok := eval.env[eval.assign]
			if !ok {
				return false
			}
			eval.ops = append(eval.ops, eval.assignOp[:len(eval.assignOp)-1])
			eval.stack[len(eval.stack)-1] = prev
			eval.stack = append(eval.stack, val)
			eval.evalOnce()
			val = eval.stack[len(eval.stack)-1]
		}
		eval.env[eval.assign] = val
	}
	return eval.count == 0 && len(eval.stack) == 1 && len(eval.ops) == 0
}

func (eval *eval) eval() float64 {
	return eval.stack[len(eval.stack)-1]
}

var statement grammar
var calc *eval = newEval()

func init() {
	ws := loop{group{" \t"}}
	op := alts("+ - ** * / || | && & <= << < >= >> > != == % ^")
	assignOp := or{and{str{"="}, notahead{str{"="}}}, alts("+= -= **= *= /= ||= |= &&= &= %= ^=")}
	digit := group{"0123456789"}
	alpha := or{rang{'a', 'z'}, rang{'A', 'Z'}}
	num := loop1(digit)
	id := and{alpha, loop{or{alpha, digit}}}

	lparens := loop{action{nil, str{"("}, calc.paren}}
	rparens := loop{action{nil, str{")"}, calc.paren}}
	term := and{ws, lparens, ws, or{action{nil, num, calc.num}, action{nil, id, calc.id}}, ws, rparens}
	expr := sep1(term, and{ws, action{nil, op, calc.op}})
	assign := opt{action{calc.resetAssign,
		and{ws, action{nil, id, calc.setAssign},
			ws, action{nil, assignOp, calc.setAssignOp}},
		calc.useAssign}}
	statement = action{calc.reset, and{assign, expr, ws, eof{}}, calc.balanced}
}

func Eval(in string) (float64, error) {
	if statement.parse(in) != -1 {
		return calc.eval(), nil
	}
	return 0, errors.New("syntax error")
}
