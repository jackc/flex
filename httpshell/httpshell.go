package httpshell

import (
	"fmt"
	"net/http"

	"github.com/jackc/flex"
)

type ParamsBuilder func(*http.Request) (map[string]interface{}, error)
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// BuildParamsError is wrapped around an error returned by BuildParams before passing to HandleError.
type BuildParamsError struct {
	commandName string
	err         error
}

func (e *BuildParamsError) Unwrap() error {
	return e.err
}

func (e *BuildParamsError) Error() string {
	return fmt.Sprintf("%s: build params: %v", e.commandName, e.err)
}

func (e *BuildParamsError) CommandName() string {
	return e.commandName
}

// BuildParamsError is wrapped around an error returned by Exec(JSON) before passing to HandleError.
type ExecError struct {
	commandName string
	err         error
}

func (e *ExecError) Unwrap() error {
	return e.err
}

func (e *ExecError) Error() string {
	return fmt.Sprintf("%s: exec: %v", e.commandName, e.err)
}

func (e *ExecError) CommandName() string {
	return e.commandName
}

// WriteError is wrapped around an error returned by Write before passing to HandleError.
type WriteError struct {
	commandName string
	err         error
}

func (e *WriteError) Unwrap() error {
	return e.err
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("%s: write: %v", e.commandName, e.err)
}

func (e *WriteError) CommandName() string {
	return e.commandName
}

type JSONHandler struct {
	Shell       *flex.Shell
	CommandName string
	BuildParams ParamsBuilder
	HandleError ErrorHandler
}

func (h *JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := h.BuildParams(r)
	if err != nil {
		h.HandleError(w, r, &BuildParamsError{commandName: h.CommandName, err: err})
		return
	}

	jsonBytes, err := h.Shell.ExecJSON(r.Context(), h.CommandName, params)
	if err != nil {
		h.HandleError(w, r, &ExecError{commandName: h.CommandName, err: err})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		h.HandleError(w, r, &WriteError{commandName: h.CommandName, err: err})
		return
	}
}
