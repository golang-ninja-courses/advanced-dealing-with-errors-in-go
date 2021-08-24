package trace

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// NOTE: Не редактируйте данный файл, чтобы не сбить ExampleGetDeepestStackTrace.

func ExampleGetDeepestStackTrace() {
	st := GetDeepestStackTrace(one())

	// Будет выведен стектрейс самой глубокой ошибки.
	fmt.Printf("%s\n", formatStackTrace(st))

	// Output:
	// deepest-stacktrace.four
	//     trace_test.go:79
	// deepest-stacktrace.three
	//     trace_test.go:75
	// deepest-stacktrace.two
	//     trace_test.go:71
	// deepest-stacktrace.one
	//     trace_test.go:67
	// deepest-stacktrace.ExampleGetDeepestStackTrace
	//     trace_test.go:18
	// main.main
	//     .:45
}

func TestGetDeepestStackTrace_Nil(t *testing.T) {
	st := GetDeepestStackTrace(nil)
	assert.Equal(t, "", formatStackTrace(st))
}

func formatStackTrace(st *sentry.Stacktrace) string {
	if st == nil {
		return ""
	}

	frames := st.Frames
	var b strings.Builder

	for i := len(frames) - 1; i >= 0; i-- {
		b.WriteString(filepath.Base(frames[i].Module))
		b.WriteRune('.')
		b.WriteString(frames[i].Function)
		b.WriteRune('\n')

		b.WriteString("    " + filepath.Base(frames[i].AbsPath))
		b.WriteRune(':')
		b.WriteString(strconv.Itoa(frames[i].Lineno))
		b.WriteRune('\n')
	}

	return b.String()
}

func one() error {
	return errors.Wrap(two(), "one")
}

func two() error {
	return errors.Wrap(three(), "two")
}

func three() error {
	return fmt.Errorf("three: %w", four())
}

func four() error {
	return errors.New("four")
}
