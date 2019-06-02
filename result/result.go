package result

import (
	"fmt"
	"strings"
	"time"
)

// Result type
type Result struct {
	status  Status
	message string
}

// Status
type Status int

const (
	Created Status = iota
	Updated
	Unchanged
	Removed
	Info
	Error
)

func (s Status) String() string {
	return toString[s]
}

var toString = map[Status]string{
	Created:   "Created",
	Updated:   "Updated",
	Unchanged: "Unchanged",
	Removed:   "Removed",
	Info:      "Info",
	Error:     "Error",
}

var toID = map[string]Status{
	"Created":   Created,
	"Updated":   Updated,
	"Unchanged": Unchanged,
	"Removed":   Removed,
	"Info":      Info,
	"Error":     Error,
}

// New constructs a Result object.
func New(status Status, message string) Result {
	return Result{status, message}
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

// Message getter
func (result Result) Message() string {
	return result.message
}

// Status getter
func (result Result) Status() Status {
	return result.status
}

// StandardizeMessage constructs a Result object with a standardized message
func (result Result) StandardizeMessage(name, value string) Result {
	if result.IsCreated() {
		return NewCreated(name + " created with value " + value)
	}
	if result.IsUpdated() {
		return NewUpdated(name + " updated. Value is now " + value)
	}
	if result.IsUnchanged() {
		return NewUnchanged(name + " has value " + value + " (unchanged)")
	}
	if result.IsRemoved() {
		return NewRemoved(name + " has been deleted")
	}
	if result.IsRemoved() {
		return NewInfo(name + " has value " + value + " (info)")
	}

	// we want to update the original message
	result.SetMessage(name + ": " + result.Message())
	return result
}

// =============================================

// NewCreated constructs a Created Result object
func NewCreated(message string) Result {
	return Result{Created, message}
}

// NewUpdated constructs an Updated Result object
func NewUpdated(message string) Result {
	return Result{Updated, message}
}

// NewUnchanged constructs an Unchanged Result object
func NewUnchanged(message string) Result {
	return Result{Unchanged, message}
}

// NewRemoved constructs a Removed Result object
func NewRemoved(message string) Result {
	return Result{Removed, message}
}

// NewInfo constructs an Info Result object
func NewInfo(message string) Result {
	return Result{Info, message}
}

// NewError constructs an Error Result object
func NewError(message string) Result {
	return Result{Error, message}
}

// =============================================

// IsSuccess function
func (result Result) IsSuccess() bool {
	return result.status != Error
}

// IsFailure function
func (result Result) IsFailure() bool {
	return !result.IsSuccess()
}

// =============================================

// IsCreated function
func (result Result) IsCreated() bool {
	return result.status == Created
}

// IsUpdated function
func (result Result) IsUpdated() bool {
	return result.status == Updated
}

// IsUnchanged function
func (result Result) IsUnchanged() bool {
	return result.status == Unchanged
}

// IsRemoved function
func (result Result) IsRemoved() bool {
	return result.status == Removed
}

// IsInfo function
func (result Result) IsInfo() bool {
	return result.status == Info
}

// IsError function
func (result Result) IsError() bool {
	return result.status == Error
}

// =============================================

// SetMessage setter
func (result Result) SetMessage(message string) {
	result.message = message
}

// Print prints a result.
func (result Result) Print() {
	tag := strings.ToUpper("[" + result.status.String() + "]")
	if result.IsSuccess() {
		fmt.Printf("%s %s\n", green(tag), result.message)
	} else {
		fmt.Printf("%s %s\n", red(tag), result.message)
	}
}
