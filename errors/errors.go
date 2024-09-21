package errors

import "errors"

var (
	// file errors
	ErrFileDoesNotExist error = errors.New("file does not exist")
	ErrCouldNotReadFromFile error = errors.New("could not read from file")
)