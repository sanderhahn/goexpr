package stackmachine

import "testing"

func TestOperator(t *testing.T) {
	prio, assoc := opInfo("+")
	if prio != 2 || assoc != left {
		t.Error("operator + wrong info")
	}
	if opEval(1, "+", 2) != 3 {
		t.Error("operator + not working")
	}
}
