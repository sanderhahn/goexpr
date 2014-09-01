package stackmachine

type StackMachine struct {
	stack  []float64
	ops    []string
	groups int
}

func (s *StackMachine) Reset() {
	s.stack = s.stack[:0]
	s.ops = []string{}
	s.groups = 0
}

func (s *StackMachine) valid() bool {
	return s.groups == 0 && len(s.stack) == 1 && len(s.ops) == 0
}

func (s *StackMachine) peek() float64 {
	return s.stack[len(s.stack)-1]
}

func (s *StackMachine) Push(val float64) {
	s.stack = append(s.stack, val)
}

func (s *StackMachine) pop() float64 {
	l := len(s.stack)
	val := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return val
}

func (s *StackMachine) peekOp() string {
	return s.ops[len(s.ops)-1]
}

func (s *StackMachine) PushOp(op string) bool {

	if op == "(" {
		s.ops = append(s.ops, op)
		s.groups++
		return true
	} else if op == ")" {
		if s.groups > 0 {
			for s.peekOp() != "(" {
				s.evalOnce()
			}
			s.popOp()
		}
		s.groups--
		return s.groups >= 0
	}

	for len(s.ops) > 0 {
		opPrio, opAssoc := opInfo(op)
		topPrio, topAssoc := opInfo(s.peekOp())

		if opAssoc == non && topAssoc == non {
			return false
		}

		if (opAssoc == left && opPrio >= topPrio) ||
			(opAssoc == right && opPrio > topPrio) ||
			(opAssoc == non && opPrio > topPrio) {
			s.evalOnce()
		} else {
			break
		}
	}

	s.ops = append(s.ops, op)
	return true
}

func (s *StackMachine) popOp() string {
	l := len(s.ops)
	op := s.ops[l-1]
	s.ops = s.ops[:l-1]
	return op
}

func (s *StackMachine) evalOnce() {
	right := s.pop()
	left := s.pop()
	s.Push(opEval(left, s.popOp(), right))
}

func (s *StackMachine) Eval() (val float64, ok bool) {
	for len(s.ops) > 0 {
		s.evalOnce()
	}
	return s.peek(), s.valid()
}
