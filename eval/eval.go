package eval

import (
	"errors"

	"github.com/sanderhahn/goexpr/expr"
	"github.com/sanderhahn/goexpr/grammar"
	"github.com/sanderhahn/goexpr/stackmachine"
)

type evaluator struct {
	machine stackmachine.StackMachine
	env     map[string]float64
	expr    grammar.Grammar
}

func newEvaluator() (e *evaluator) {
	e = &evaluator{
		machine: stackmachine.StackMachine{},
		env:     map[string]float64{},
	}
	e.expr = expr.Parser(e)
	return e
}

func (e *evaluator) Reset() {
	e.machine.Reset()
}

func (e *evaluator) Number(number float64) {
	e.machine.Push(number)
}

func (e *evaluator) Lookup(identifier string) bool {
	val, ok := e.env[identifier]
	if ok {
		e.machine.Push(val)
	}
	return ok
}

func (e *evaluator) Operator(in string) bool {
	return e.machine.PushOp(in)
}

func (e *evaluator) Assignment(identifier string, operator string) bool {
	val, ok := e.machine.Eval()
	if !ok {
		return false
	}
	if operator != "=" {
		if !e.Lookup(identifier) {
			return false
		}
		e.machine.PushOp(operator[:len(operator)-1])
		val, ok = e.machine.Eval()
		if !ok {
			return false
		}
	}
	e.env[identifier] = val
	return true
}

var e *evaluator = newEvaluator()

func Eval(in string) (val float64, err error) {
	if e.expr.Parse(in) != -1 {
		val, ok := e.machine.Eval()
		if ok {
			return val, nil
		}
	}
	return 0, errors.New("syntax error")
}
