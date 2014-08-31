package main

import (
	"log"
	"math"
)

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
	}
	return v
}

type assoc int

const (
	left assoc = iota
	right
	non
)

func opInfo(op string) (prio int, assoc assoc) {

	switch op {
	case "*":
		return 1, left
	case "/":
		return 1, left
	case "%":
		return 1, left
	case "<<":
		return 1, left
	case ">>":
		return 1, left
	case "&":
		return 1, left
	case "**":
		return 1, right

	case "+":
		return 2, left
	case "-":
		return 2, left
	case "|":
		return 2, left
	case "^":
		return 2, left

	case "==":
		return 3, non
	case "!=":
		return 3, non
	case "<":
		return 3, non
	case "<=":
		return 3, non
	case ">":
		return 3, non
	case ">=":
		return 3, non

	case "&&":
		return 4, left

	case "||":
		return 5, left

	case "(":
		return 10, non
	}

	log.Fatalf("operator %s has no priority", op)
	return -1, non
}
