package stacktrace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE: Не редактируйте данный файл, чтобы не сбить пример и тесты ниже.

func ExampleTrace() {
	stack := Trace()
	fmt.Println(stack[:1])

	// Output:
	// handmade-stacktrace.ExampleTrace
	// handmade-stacktrace/stack_trace_test.go:14
}

func TestTrace(t *testing.T) {
	t.Run("simple call", func(t *testing.T) {
		frames := Trace()
		assert.Len(t, frames, 3)

		re := regexp.MustCompile(`handmade-stacktrace\.TestTrace\.func1
handmade-stacktrace\/stack_trace_test\.go:24
testing\.tRunner
testing\/testing\.go:\d{1,4}`)
		trace := frames[:2].String()
		assert.True(t, re.MatchString(trace), trace)
	})

	t.Run("depth = 100", func(t *testing.T) {
		frames := dive(100) // На самом деле чуть больше.
		assert.Len(t, frames, 32)
	})
}

func dive(depth int) StackTrace {
	if depth == 0 {
		return Trace()
	}
	return dive(depth - 1)
}
