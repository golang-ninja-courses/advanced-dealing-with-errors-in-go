package tmpl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"text/template"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsExecUnexportedFieldError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		assert.False(t, IsExecUnexportedFieldError(nil))
	})

	t.Run("not `unexported field` error", func(t *testing.T) {
		err := parseAndExecuteTemplate(bytes.NewBuffer(nil), "example",
			`{{ with .Name }}`,
			struct{ Name unsafe.Pointer }{},
		)
		require.Error(t, err)
		assert.False(t, IsExecUnexportedFieldError(err))
	})

	t.Run("partial assertion is not working", func(t *testing.T) {
		assert.False(t, IsExecUnexportedFieldError(errors.New("is an unexported field of struct type")))
		assert.False(t, IsExecUnexportedFieldError(errors.New("template")))
	})

	t.Run("not template.Exec error", func(t *testing.T) {
		lie := `template: example:1:3: executing "example" at <.name>: name is an unexported field of struct type`
		assert.False(t, IsExecUnexportedFieldError(errors.New(lie)))
	})

	t.Run("`unexported field` error", func(t *testing.T) {
		err := parseAndExecuteTemplate(bytes.NewBuffer(nil), "example",
			`{{ .name }}`,
			struct{ name string }{name: "Bob"},
		)
		require.Error(t, err)
		assert.True(t, IsExecUnexportedFieldError(err))
	})
}

func TestIsFunctionNotDefinedError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		assert.False(t, IsFunctionNotDefinedError(nil))
	})

	t.Run("not `function not defined` error", func(t *testing.T) {
		err := parseAndExecuteTemplate(bytes.NewBuffer(nil), "example", `{{ with`, nil)
		require.Error(t, err)
		assert.False(t, IsFunctionNotDefinedError(err))
	})

	t.Run("partial assertion is not working", func(t *testing.T) {
		assert.False(t, IsFunctionNotDefinedError(errors.New("function")))
		assert.False(t, IsFunctionNotDefinedError(errors.New("not defined")))
		assert.False(t, IsFunctionNotDefinedError(errors.New("template")))
	})

	t.Run("`function not defined` error", func(t *testing.T) {
		err := parseAndExecuteTemplate(bytes.NewBuffer(nil), "example", `{{ call XXX }}`, nil)
		require.Error(t, err)
		assert.True(t, IsFunctionNotDefinedError(err))
	})
}

func parseAndExecuteTemplate(wr io.Writer, name, text string, data interface{}) error {
	t, err := template.New(name).Parse(text)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}

	if err := t.Execute(wr, data); err != nil {
		return fmt.Errorf("cannot execute template: %w", err)
	}
	return nil
}
