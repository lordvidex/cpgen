package cond

// If - imitates ternary operator to allow writing code in the format <cond> ? <true> : <false>
func If[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}
