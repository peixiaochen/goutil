package assert

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/stdutil"
)

// Nil asserts that the given is a nil value
func Nil(t TestingT, give any, fmtAndArgs ...any) bool {
	if stdutil.IsNil(give) {
		return true
	}

	t.Helper()
	return fail(t, fmt.Sprintf("Expected nil, but got: %#v", give), fmtAndArgs)
}

// NotNil asserts that the given is a not nil value
func NotNil(t TestingT, give any, fmtAndArgs ...any) bool {
	if !stdutil.IsNil(give) {
		return true
	}

	t.Helper()
	return fail(t, "Should not nil value", fmtAndArgs)
}

// True asserts that the given is a bool true
func True(t TestingT, give bool, fmtAndArgs ...any) bool {
	if !give {
		t.Helper()
		return fail(t, "Result should be True", fmtAndArgs)
	}
	return true
}

// False asserts that the given is a bool false
func False(t TestingT, give bool, fmtAndArgs ...any) bool {
	if give {
		t.Helper()
		return fail(t, "Result should be False", fmtAndArgs)
	}
	return true
}

// Empty asserts that the give should be empty
func Empty(t TestingT, give any, fmtAndArgs ...any) bool {
	empty := stdutil.IsEmpty(give)
	if !empty {
		t.Helper()
		return fail(t, fmt.Sprintf("Should be empty, but was:\n%#v", give), fmtAndArgs)
	}

	return empty
}

// NotEmpty asserts that the give should not be empty
func NotEmpty(t TestingT, give any, fmtAndArgs ...any) bool {
	nEmpty := !stdutil.IsEmpty(give)
	if !nEmpty {
		t.Helper()
		return fail(t, fmt.Sprintf("Should not be empty, but was:\n%#v", give), fmtAndArgs)
	}

	return nEmpty
}

// PanicRunFunc define
type PanicRunFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func runPanicFunc(f PanicRunFunc) (didPanic bool, message interface{}, stack string) {
	didPanic = true
	defer func() {
		message = recover()
		if didPanic {
			stack = string(debug.Stack())
		}
	}()

	// call the target function
	f()
	didPanic = false

	return
}

// Panics asserts that the code inside the specified func panics.
func Panics(t TestingT, fn PanicRunFunc, fmtAndArgs ...any) bool {
	if hasPanic, panicVal, _ := runPanicFunc(fn); !hasPanic {
		t.Helper()

		return fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", fn, panicVal), fmtAndArgs)
	}

	return true
}

// NotPanics asserts that the code inside the specified func NOT panics.
func NotPanics(t TestingT, fn PanicRunFunc, fmtAndArgs ...any) bool {
	if hasPanic, panicVal, stackMsg := runPanicFunc(fn); hasPanic {
		t.Helper()

		return fail(t, fmt.Sprintf(
			"func %#v should not panic\n\tPanic value:\t%#v\n\tPanic stack:\t%s",
			fn, panicVal, stackMsg,
		), fmtAndArgs,
		)
	}

	return true
}

// PanicsMsg should panic and with a value
func PanicsMsg(t TestingT, fn PanicRunFunc, wantVal interface{}, fmtAndArgs ...any) bool {
	hasPanic, panicVal, stackMsg := runPanicFunc(fn)
	if !hasPanic {
		t.Helper()
		return fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", fn, panicVal), fmtAndArgs)
	}

	if panicVal != wantVal {
		t.Helper()
		return fail(t, fmt.Sprintf(
			"func %#v should panic.\n\tWant  value:\t%#v\n\tPanic value:\t%#v\n\tPanic stack:\t%s",
			fn, wantVal, panicVal, stackMsg),
			fmtAndArgs,
		)
	}

	return true
}

// PanicsErrMsg should panic and with error message
func PanicsErrMsg(t TestingT, fn PanicRunFunc, errMsg string, fmtAndArgs ...any) bool {
	hasPanic, panicVal, stackMsg := runPanicFunc(fn)
	if !hasPanic {
		t.Helper()
		return fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", fn, panicVal), fmtAndArgs)
	}

	err, ok := panicVal.(error)
	if !ok {
		t.Helper()
		return fail(t, fmt.Sprintf("func %#v should panic and is error type,\nbut type was: %T", fn, panicVal), fmtAndArgs)
	}

	if err.Error() != errMsg {
		t.Helper()
		return fail(t, fmt.Sprintf(
			"func %#v should panic.\n\tWant  error:\t%#v\n\tPanic value:\t%#v\n\tPanic stack:\t%s",
			fn, errMsg, panicVal, stackMsg),
			fmtAndArgs,
		)
	}

	return true
}

// Contains asserts that the given data(string,slice,map) should contain element
//
// TIP: only support types: string, map, array, slice
// 	map         - check key exists
// 	string      - check sub-string exists
// 	array,slice - check sub-element exists
func Contains(t TestingT, src, elem any, fmtAndArgs ...any) bool {
	valid, found := stdutil.CheckContains(src, elem)
	if valid && found {
		return true
	}

	t.Helper()

	// src invalid
	if !valid {
		return fail(t, fmt.Sprintf("%#v could not be applied builtin len()", src), fmtAndArgs)
	}

	// not found
	return fail(t, fmt.Sprintf("%#v\nShould contain: %#v", src, elem), fmtAndArgs)
}

// NotContains asserts that the given data(string,slice,map) should not contain element
//
// TIP: only support types: string, map, array, slice
// 	map         - check key exists
// 	string      - check sub-string exists
// 	array,slice - check sub-element exists
func NotContains(t TestingT, src, elem any, fmtAndArgs ...any) bool {
	valid, found := stdutil.CheckContains(src, elem)
	if valid && !found {
		return true
	}

	t.Helper()

	// src invalid
	if !valid {
		return fail(t, fmt.Sprintf("%#v could not be applied builtin len()", src), fmtAndArgs)
	}

	// found
	return fail(t, fmt.Sprintf("%#v\nShould not contain: %#v", src, elem), fmtAndArgs)
}

// ContainsKey asserts that the given map is contains key
func ContainsKey(t TestingT, mp, key any, fmtAndArgs ...any) bool {
	if !maputil.HasKey(mp, key) {
		t.Helper()
		return fail(t,
			fmt.Sprintf("Map data should contains the key: %#v\nMap data: %#v", key, mp),
			fmtAndArgs,
		)
	}

	return true
}

// ContainsKeys asserts that the map is contains all given keys
func ContainsKeys(t TestingT, mp any, keys []interface{}, fmtAndArgs ...any) bool {
	ok, noKey := maputil.HasAllKeys(mp, keys...)
	if !ok {
		t.Helper()
		return fail(t,
			fmt.Sprintf("Map data should contains the key: %#v\nMap data: %#v", noKey, mp),
			fmtAndArgs,
		)
	}

	return true
}

// StrContains asserts that the given strings is contains sub-string
func StrContains(t TestingT, s, sub string, fmtAndArgs ...any) bool {
	if !strings.Contains(s, sub) {
		t.Helper()
		return fail(t,
			fmt.Sprintf("String value check fail:\nGiven string: %#v\nNot contains: %#v", s, sub),
			fmtAndArgs,
		)
	}

	return true
}

//
// -------------------- error --------------------
//

// NoErr asserts that the given is a nil error
func NoErr(t TestingT, err error, fmtAndArgs ...any) bool {
	if err != nil {
		t.Helper()
		return fail(t, fmt.Sprintf("Received unexpected error:\n%+v", err), fmtAndArgs)
	}
	return true
}

// Err asserts that the given is a not nil error
func Err(t TestingT, err error, fmtAndArgs ...any) bool {
	if err == nil {
		t.Helper()
		return fail(t, "An error is expected but got nil.", fmtAndArgs)
	}
	return true
}

// ErrMsg asserts that the given is a not nil error and error message equals wantMsg
func ErrMsg(t TestingT, err error, wantMsg string, fmtAndArgs ...any) bool {
	if err == nil {
		t.Helper()
		return fail(t, "An error is expected but got nil.", fmtAndArgs)
	}

	errMsg := err.Error()
	if errMsg != wantMsg {
		t.Helper()
		return fail(t, fmt.Sprintf("Error message not equal:\n"+
			"expect: %q\n"+
			"actual: %q", wantMsg, errMsg), fmtAndArgs)
	}

	return true
}

// ErrSubMsg asserts that the given is a not nil error and the error message contains subMsg
func ErrSubMsg(t TestingT, err error, subMsg string, fmtAndArgs ...any) bool {
	if err == nil {
		t.Helper()
		return fail(t, "An error is expected but got nil.", fmtAndArgs)
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, subMsg) {
		t.Helper()
		return fail(t, fmt.Sprintf("Error message check fail:\n"+
			"error  message : %q\n"+
			"should contains: %q", errMsg, subMsg), fmtAndArgs)
	}

	return true
}

//
// -------------------- Len --------------------
//

func Len(t TestingT, give any, wantLn int, fmtAndArgs ...any) bool {
	gln := reflects.Len(reflect.ValueOf(give))
	if gln < 0 {
		t.Helper()
		return fail(t, fmt.Sprintf("\"%s\" could not be calc length", give), fmtAndArgs)
	}

	if gln != wantLn {
		t.Helper()
		return fail(t, fmt.Sprintf("\"%s\" should have %d item(s), but has %d", give, wantLn, gln), fmtAndArgs)
	}
	return false
}

func LenGt(t TestingT, give any, minLn int, fmtAndArgs ...any) bool {
	gln := reflects.Len(reflect.ValueOf(give))
	if gln < 0 {
		t.Helper()
		return fail(t, fmt.Sprintf("\"%s\" could not be calc length", give), fmtAndArgs)
	}

	if gln < minLn {
		t.Helper()
		return fail(t, fmt.Sprintf("\"%s\" should less have %d item(s), but has %d", give, minLn, gln), fmtAndArgs)
	}
	return false
}

//
// -------------------- compare --------------------
//

// Eq asserts that the want should equal to the given
func Eq(t TestingT, want, give any, fmtAndArgs ...any) bool {
	t.Helper()

	if err := checkEqualArgs(want, give); err != nil {
		return fail(t, fmt.Sprintf("Cannot compare: %#v == %#v (%s)",
			want, give, err), fmtAndArgs)
	}

	if !reflects.IsEqual(want, give) {
		// TODO diff := diff(want, give)
		want, give = formatUnequalValues(want, give)
		return fail(t, fmt.Sprintf("Not equal: \n"+
			"expect: %s\n"+
			"actual: %s", want, give), fmtAndArgs)
	}

	return true
}

// Neq asserts that the want should not be equal to the given.
// alias of NotEq()
func Neq(t TestingT, want, give any, fmtAndArgs ...any) bool {
	t.Helper()
	return NotEq(t, want, give, fmtAndArgs...)
}

// NotEq asserts that the want should not be equal to the given
func NotEq(t TestingT, want, give any, fmtAndArgs ...any) bool {
	t.Helper()

	if err := checkEqualArgs(want, give); err != nil {
		return fail(t, fmt.Sprintf("Cannot compare: %#v == %#v (%s)",
			want, give, err), fmtAndArgs)
	}

	if reflects.IsEqual(want, give) {
		return fail(t, fmt.Sprintf("Given should not be: %#v\n", give), fmtAndArgs)
	}
	return true
}

func Lt(t TestingT, give, max int, fmtAndArgs ...any) bool {
	gInt, err := mathutil.ToInt(give)
	if err == nil && gInt <= max {
		return true
	}

	return fail(t, fmt.Sprintf("Given should later than or equal %d(but was %d)", max, gInt), fmtAndArgs)

}

func Gt(t TestingT, give, min int, fmtAndArgs ...any) bool {
	gInt, err := mathutil.ToInt(give)
	if err == nil && gInt >= min {
		return true
	}

	return fail(t, fmt.Sprintf("Given should gater than or equal %d(but was %d)", min, gInt), fmtAndArgs)
}

//
// -------------------- fail --------------------
//

// Fail reports a failure through
func Fail(t TestingT, failMsg string, fmtAndArgs ...any) bool {
	t.Helper()
	return fail(t, failMsg, fmtAndArgs)
}

type failNower interface {
	FailNow()
}

// FailNow fails test
func FailNow(t TestingT, failMsg string, fmtAndArgs ...any) bool {
	t.Helper()
	fail(t, failMsg, fmtAndArgs)

	if fnr, ok := t.(failNower); ok {
		fnr.FailNow()
	}
	return false
}

//
// -------------------- render error --------------------
//

var (
	// ShowFullPath on show error trace
	ShowFullPath = true
	// EnableColor on show error trace
	EnableColor = true
)

// DisableColor render
func DisableColor() {
	EnableColor = false
}

// HideFullPath render
func HideFullPath() {
	ShowFullPath = false
}

// fail reports a failure through
func fail(t TestingT, failMsg string, fmtAndArgs []any) bool {
	t.Helper()

	tName := t.Name()
	if EnableColor {
		tName = color.Red.Sprint(tName)
	}

	labeledTexts := []labeledText{
		{"Test Name", tName},
		{"Error At", strings.Join(callerInfos(), "\n")},
		{"Error Msg", failMsg},
	}

	// user custom message
	if userMsg := formatTplAndArgs(fmtAndArgs...); len(userMsg) > 0 {
		labeledTexts = append(labeledTexts, labeledText{"User Msg", userMsg})
	}

	t.Error("\n" + formatLabeledTexts(labeledTexts))
	return false
}
