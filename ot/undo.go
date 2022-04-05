package ot

// When you use ctrl-z to undo your latest changes, you expect the program not
// to undo every single keystroke but to undo your last sentence you wrote at
// a stretch or the deletion you did by holding the backspace key down. This
// This can be implemented by composing operations on the undo stack. This
// method can help decide whether two operations should be composed. It
// returns true if the operations are consecutive insert operations or both
// operations delete text at the same position. You may want to include other
// factors like the time since the last change in your decision.
func (o *Operation) shouldBeComposedWith(other *Operation) bool {
	if o.isNoop() || other.isNoop() {
		return true
	}

	startA := getStartIndex(o)
	startB := getStartIndex(other)
	simpleA := getSimpleOp(o)
	simpleB := getSimpleOp(other)

	if simpleA == nil || simpleB == nil {
		return false
	}

	if simpleA.IsInsert() && simpleB.IsInsert() {
		return startA+len(simpleA.StrVal) == startB
	}

	if simpleA.IsDelete() && simpleB.IsDelete() {
		// there are two possibilities to delete: with backspace and with the
		// delete key.
		return startB-int(simpleB.Val) == startA || startA == startB
	}
	return false
}

// Decides whether two operations should be composed with each other
// if they were inverted, that is
// `shouldBeComposedWith(a, b) = shouldBeComposedWithInverted(b^{-1}, a^{-1})`.
func (o *Operation) shouldBeComposedWithInverted(other *Operation) bool {
	if o.isNoop() || other.isNoop() {
		return true
	}

	startA := getStartIndex(o)
	startB := getStartIndex(other)
	simpleA := getSimpleOp(o)
	simpleB := getSimpleOp(other)

	if simpleA == nil || simpleB == nil {
		return false
	}

	if simpleA.IsInsert() && simpleB.IsInsert() {
		return startA+len(simpleA.StrVal) == startB || startA == startB
	}

	if simpleA.IsDelete() && simpleB.IsDelete() {
		return startB-int(simpleB.Val) == startA
	}
	return false
}
