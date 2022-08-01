package validator

// ErrNilValue is thrown when a value is nil.
type ErrNil struct {
	Key string
}

func (e *ErrNil) Error() string {
	return e.Key + " is nill"
}

// ErrInvalidJSON is thrown when an error occurred while marshaling the content.
type ErrInvalidJSON struct {
	Cause error
}

func (e *ErrInvalidJSON) Error() string {
	return "Invalid JSON: " + e.Cause.Error()
}

// ErrInvalidContent is thrown when an url or body contains the wrong content for the request.
type ErrInvalidContent struct {
	Cause string
}

func (e *ErrInvalidContent) Error() string {
	return "Invalid content: " + e.Cause
}
