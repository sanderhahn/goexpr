package eval

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/sanderhahn/goexpr/ast"
	"github.com/sanderhahn/goexpr/expr"
)

type Environment struct {
	Variables map[string]float64
	Parser    expr.ExpressionParser
}

func NewEnvironment() *Environment {
	return &Environment{
		Variables: map[string]float64{},
		Parser:    expr.Parser(),
	}
}

func (env *Environment) evalExpression(group ast.Group) (value float64, err error) {
	operator := group[0].(ast.Token)
	left, err := env.evalNode(group[1])
	if err != nil {
		return
	}
	right, err := env.evalNode(group[2])
	if err != nil {
		return
	}
	value = opEval(left, operator.Value, right)
	return
}

func (env *Environment) evalNode(node ast.Node) (value float64, err error) {
	switch node := node.(type) {

	case ast.Token:
		switch node.TokenType {

		case ast.Number:
			value, err = strconv.ParseFloat(node.Value, 64)

		case ast.Identifier:
			var ok bool
			value, ok = env.Variables[node.Value]
			if !ok {
				err = errors.New("variable " + node.Value + " is undefined")
			}
		}

	case ast.Group:
		operator := node[0].(ast.Token)
		switch operator.TokenType {

		case ast.Paren:
			value, err = env.evalNode(node[1])

		case ast.Operator:
			if isAssignment(operator.Value) {
				identifier := node[1].(ast.Token)
				if identifier.TokenType != ast.Identifier {
					err = errors.New("no variable at left hand side")
					return
				}
				if operator.Value != "=" {
					shortcut := operator.Value[:len(operator.Value)-1]
					value, err = env.evalExpression(ast.Group{ast.Token{ast.Operator, shortcut}, identifier, node[2]})
				} else {
					value, err = env.evalNode(node[2])
				}
				env.Variables[identifier.Value] = value
			} else {
				value, err = env.evalExpression(node)
			}
		}
	}
	return
}

func (env *Environment) Eval(in string) (value float64, err error) {
	var node ast.Node = env.Parser.Parse(in)
	fmt.Printf("ast: %v\n", node)
	if node != nil {
		value, err = env.evalNode(node)
	} else {
		err = errors.New("syntax error")
	}
	return
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

func opEval(l float64, op string, r float64) (v float64) {
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
	default:
		log.Fatalf("operator %s not implemented", op)
	}
	return v
}

func isAssignment(operator string) bool {
	switch operator {
	case "=", "+=", "-=", "**=", "*=", "/=", "||=", "|=", "&&=", "&=", "%=", "^=":
		return true
	}
	return false
}
