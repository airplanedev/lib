package api

import "fmt"

// TaskMissingError implements an explainable error.
type TaskMissingError struct {
	AppURL string
	Slug   string
}

// Error implementation.
func (err TaskMissingError) Error() string {
	return fmt.Sprintf("task with slug %q does not exist", err.Slug)
}

// ExplainError implementation.
func (err TaskMissingError) ExplainError() string {
	return fmt.Sprintf(
		"Follow the URL below to create the task:\n%s",
		err.AppURL+"/tasks/new",
	)
}

// AppMissingError implements an explainable error.
type AppMissingError struct {
	AppURL string
	Slug   string
}

// Error implementation.
func (err AppMissingError) Error() string {
	return fmt.Sprintf("app with slug %q does not exist", err.Slug)
}

// ExplainError implementation.
func (err AppMissingError) ExplainError() string {
	return fmt.Sprintf(
		"Follow the URL below to create the app:\n%s",
		err.AppURL+"/apps/new",
	)
}
