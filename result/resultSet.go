package result

import (
	"fmt"
	"strconv"
)

// Set helps managing a set of Result
type Set struct {
	results []Result
	message string
}

// NewSet constructor.
func NewSet(results []Result, message string) Set {
	return Set{results, message}
}

// IsSuccess checks if all the results are in success.
func (results Set) IsSuccess() bool {
	for _, element := range results.results {
		if !element.isSuccess {
			return false
		}
	}
	return true
}

// OverallResult gets a result object for the all set.
func (results Set) OverallResult() Result {
	return New(results.IsSuccess(), results.Message())
}

// IsEmpty checks if the set contains NO result.
func (results Set) IsEmpty() bool {
	return len(results.results) == 0
}

// Size gets the number of results (both success and failures).
func (results Set) Size() int {
	return len(results.results)
}

// Print prints a result.
func (results Set) Print() {
	for _, res := range results.results {
		res.Print()
	}
	if !results.IsEmpty() {
		fmt.Println()
	}
	results.OverallResult().Print()
}

// Add a new result to the result set.
// The result set is returned for convinience.
func (results *Set) Add(res Result) Set {
	results.results = append(results.results, res)
	return *results
}

// Message getter.
// Returns a default message if no message is set.
func (results Set) Message() string {
	if results.message != "" {
		return results.message
	}
	return results.DefaultMessage()
}

// SetMessage setter.
func (results *Set) SetMessage(message string) {
	results.message = message
}

// DefaultMessage computes a default message based on the object state.
func (results Set) DefaultMessage() string {
	tot, success, failures := results.statistics()

	if tot == 0 {
		return "No result available"
	}

	if success == tot {
		return fmt.Sprintf("All %s elements executed with success", strconv.Itoa(tot))
	}

	if failures == tot {
		return fmt.Sprintf("All %s elements executed with error", strconv.Itoa(tot))
	}

	return fmt.Sprintf("%s success, %s failures", strconv.Itoa(success), strconv.Itoa(failures))
}

func (results Set) statistics() (int, int, int) {
	success := 0
	failures := 0
	for _, element := range results.results {
		if element.isSuccess {
			success++
		} else {
			failures++
		}
	}
	return len(results.results), success, failures
}
