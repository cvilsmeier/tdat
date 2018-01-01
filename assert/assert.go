/*
Package assert provides assertion functions for tests.
We define our own little set of assert functions because
we don't want to depend on 3rd party packages.

As Rob Pike says: "A little copying is better than a little dependency".
*/
package assert

import (
	"fmt"
	"testing"
)

// True asserts a condition.
func True(t *testing.T, cond bool) {
	t.Helper()
	Truef(t, cond, "")
}

// True asserts a condition with a formatted message.
func Truef(t *testing.T, cond bool, format string, args ...interface{}) {
	t.Helper()
	if !cond {
		msg := "assertion failed"
		msg2 := fmt.Sprintf(format, args...)
		if msg2 != "" {
			msg = msg + ": " + msg2
		}
		t.Fatalf(msg)
	}
}

// EqStr asserts that a string equals another.
func EqStr(t *testing.T, exp, act string) {
	t.Helper()
	EqStrf(t, exp, act, "")
}

// EqStr asserts that a string equals another, with a formattd message.
func EqStrf(t *testing.T, exp, act string, format string, args ...interface{}) {
	t.Helper()
	if exp != act {
		msg := fmt.Sprintf("exp [%q] but was [%q]", exp, act)
		msg2 := fmt.Sprintf(format, args...)
		if msg2 != "" {
			msg = msg + ": " + msg2
		}
		t.Fatalf(msg)
	}
}

// EqInt asserts that two ints are equal.
func EqInt(t *testing.T, exp, act int) {
	t.Helper()
	EqIntf(t, exp, act, "")
}

// EqInt asserts that two ints are equal, with a formattd message.
func EqIntf(t *testing.T, exp, act int, format string, args ...interface{}) {
	t.Helper()
	if exp != act {
		msg := fmt.Sprintf("exp %d but was %d", exp, act)
		msg2 := fmt.Sprintf(format, args...)
		if msg2 != "" {
			msg = msg + ": " + msg2
		}
		t.Fatalf(msg)
	}
}
