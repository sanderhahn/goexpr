package ast

import (
	"fmt"
	"strings"
)

type Node interface {
	fmt.Stringer
}

type TokenType string

const (
	Number     TokenType = "number"
	Operator   TokenType = "operator"
	Identifier TokenType = "identifier"
	Paren      TokenType = "paren"
)

type Token struct {
	TokenType TokenType
	Value     string
}

func (t Token) String() string {
	return t.Value
}

type Group []Node

func (g Group) String() string {
	s := []string{}
	for _, i := range g {
		s = append(s, i.String())
	}
	return "(" + strings.Join(s, " ") + ")"
}

type AstBuilder struct {
	nodes []Node
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{}
}

func (a *AstBuilder) Reset() {
	a.nodes = a.nodes[:0]
}

func (a *AstBuilder) pop() (node Node) {
	l := len(a.nodes)
	node = a.nodes[l-1]
	a.nodes = a.nodes[:l-1]
	return
}

func (a *AstBuilder) push(node Node) {
	a.nodes = append(a.nodes, node)
}

func (a *AstBuilder) Token(typ TokenType, value string) {
	a.push(Token{typ, value})
}

func (a *AstBuilder) Group() {
	a.push(Group{Token{Paren, "()"}, a.pop()})
}

func (a *AstBuilder) Evaluate() {
	op := a.pop()
	right := a.pop()
	left := a.pop()
	a.push(Group{op, left, right})
}

func (a *AstBuilder) Root() Node {
	if len(a.nodes) == 1 {
		return a.nodes[0]
	}
	return nil
}
