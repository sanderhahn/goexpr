package main

type stackmachine struct {
	stack  []float64
	ops    []string
	groups int
}

func (s *stackmachine) reset() {
	s.stack = s.stack[:0]
	s.ops = []string{}
	s.groups = 0
}

func (s *stackmachine) valid() bool {
	return s.groups == 0 && len(s.stack) == 1 && len(s.ops) == 0
}

func (s *stackmachine) peek() float64 {
	return s.stack[len(s.stack)-1]
}

func (s *stackmachine) push(val float64) {
	s.stack = append(s.stack, val)
}

func (s *stackmachine) pop() float64 {
	l := len(s.stack)
	val := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return val
}

func (s *stackmachine) peekOp() string {
	return s.ops[len(s.ops)-1]
}

func (s *stackmachine) pushOp(op string) bool {

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

func (s *stackmachine) popOp() string {
	l := len(s.ops)
	op := s.ops[l-1]
	s.ops = s.ops[:l-1]
	return op
}

func (s *stackmachine) evalOnce() {
	right := s.pop()
	left := s.pop()
	s.push(opEval(left, s.popOp(), right))
}

func (s *stackmachine) eval() (val float64, ok bool) {
	for len(s.ops) > 0 {
		s.evalOnce()
	}
	return s.peek(), s.valid()
}
