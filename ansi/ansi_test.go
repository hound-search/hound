package ansi

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	printTests = false
)

func makeReal(s string) string {  //nolint
	return strings.Replace(s, "~", "\x1b", -1)
}

func makeFake(s string) string {
	return strings.Replace(s, "\x1b", "~", -1)
}

func assertEqual(t *testing.T, got string, exp string) {
	if printTests {
		fmt.Println(got)
	}

	exp = strings.Replace(exp, "~", "\x1b", -1)
	if got != exp {
		t.Errorf("mismatch: %s & %s", makeFake(got), makeFake(exp))
	}
}

func TestEnabled(t *testing.T) {
	a := Colorer{true}
	assertEqual(t,
		a.FgBg("x", Black, Normal, Colorless, Normal),
		"~[30mx~[0m")
	assertEqual(t,
		a.FgBg("x", Red, Normal, Colorless, Normal),
		"~[31mx~[0m")
	assertEqual(t,
		a.FgBg("x", Red, Intense, Colorless, Normal),
		"~[91mx~[0m")
	assertEqual(t,
		a.FgBg("x", Green, Bold|Blink|Underline|Invert, Colorless, Normal),
		"~[1;5;4;7;32mx~[0m")
	assertEqual(t,
		a.FgBg("x", Green, Bold|Blink|Underline|Invert|Intense, Colorless, Normal),
		"~[1;5;4;7;92mx~[0m")

	assertEqual(t,
		a.FgBg("x", Green, Bold|Blink|Underline|Intense, Magenta, Normal),
		"~[1;5;4;92;45mx~[0m")
	assertEqual(t,
		a.FgBg("x", Yellow, Bold|Blink|Underline|Intense, Cyan, Intense),
		"~[1;5;4;93;106mx~[0m")
}

func TestDisabled(t *testing.T) {
	a := Colorer{false}
	assertEqual(t,
		a.FgBg("x", Black, Normal, Colorless, Normal),
		"x")
	assertEqual(t,
		a.FgBg("foo", Red, Normal, Colorless, Normal),
		"foo")
	assertEqual(t,
		a.FgBg("butter", Red, Intense, Colorless, Normal),
		"butter")
	assertEqual(t,
		a.FgBg("x", Green, Bold|Blink|Underline|Invert, Colorless, Normal),
		"x")
	assertEqual(t,
		a.FgBg("x", Green, Bold|Blink|Underline|Invert|Intense, Colorless, Normal),
		"x")

	assertEqual(t,
		a.FgBg("x", Green, Bold|Blink|Underline|Intense, Magenta, Normal),
		"x")
	assertEqual(t,
		a.FgBg("x", Yellow, Bold|Blink|Underline|Intense, Cyan, Intense),
		"x")
}

func TestIsTerminal(t *testing.T) {
	// just make sure we can call this thing.
	isTTY(os.Stdout.Fd())
}
