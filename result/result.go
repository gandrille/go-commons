package result

import (
	"fmt"
	"time"
)

// Result type
type Result struct {
	isSuccess bool
	message   string
}

// New constructs a Result object.
func New(isSuccess bool, message string) Result {
	return Result{isSuccess, message}
}

// Run executes a treatement and append the time spent to the result message
func Run(runner func() Result) Result {
	t0 := time.Now()
	result := runner()
	t1 := time.Now()
	duration := t1.Sub(t0).Truncate(time.Second).String()
	result.SetMessage(fmt.Sprintf("%s (%s)", result.Message(), duration))
	return result
}

// Success constructs a Result success object.
func Success(message string) Result {
	return Result{true, message}
}

// Failure constructs a Result failure object.
func Failure(message string) Result {
	return Result{false, message}
}

// IsSuccess function
func (result Result) IsSuccess() bool {
	return result.isSuccess
}

// SetMessage setter
func (result *Result) SetMessage(message string) {
	result.message = message
}

// Message getter
func (result Result) Message() string {
	return result.message
}

// Print prints a result.
func (result Result) Print() {
	if result.isSuccess {
		PrintOK(result.message)
	} else {
		PrintError(result.message)
	}
}
