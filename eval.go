package main

import "errors"

type evaluator struct {
	calculator
	machine stackmachine
	env     map[string]float64
	expr    grammar
}

func newEvaluator() (e *evaluator) {
	e = &evaluator{
		machine: stackmachine{},
		env:     map[string]float64{},
	}
	e.expr = expr(e)
	return e
}

func (e *evaluator) reset() {
	e.machine.reset()
}

func (e *evaluator) number(number float64) {
	e.machine.push(number)
}

func (e *evaluator) lookup(identifier string) bool {
	val, ok := e.env[identifier]
	if ok {
		e.machine.push(val)
	}
	return ok
}

func (e *evaluator) operator(in string) bool {
	return e.machine.pushOp(in)
}

func (e *evaluator) assignment(identifier string, operator string) bool {
	val, ok := e.machine.eval()
	if !ok {
		return false
	}
	if operator != "=" {
		if !e.lookup(identifier) {
			return false
		}
		e.machine.pushOp(operator[:len(operator)-1])
		val, ok = e.machine.eval()
		if !ok {
			return false
		}
	}
	e.env[identifier] = val
	return true
}

var e *evaluator = newEvaluator()

func Eval(in string) (val float64, err error) {
	if e.expr.parse(in) != -1 {
		val, ok := e.machine.eval()
		if ok {
			return val, nil
		}
	}
	return 0, errors.New("syntax error")
}
