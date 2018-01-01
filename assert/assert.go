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

func True(t *testing.T, cond bool) {
	t.Helper()
	Truef(t, cond, "")
}

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

func EqStr(t *testing.T, exp, act string) {
	t.Helper()
	EqStrf(t, exp, act, "")
}

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

func EqInt(t *testing.T, exp, act int) {
	t.Helper()
	EqIntf(t, exp, act, "")
}

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
