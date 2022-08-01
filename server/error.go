package server

// ErrNil is thrown when a function was executed on a nil pointer.
type ErrNil struct{}

func (e *ErrNil) Error() string {
	return "cannot execute function on a nil pointer"
}

// ErrInvalidURL is thrown when a given URL is not valid.
type ErrInvalidURL struct {
	url string
}

func (e *ErrInvalidURL) Error() string {
	return "invalid URL '" + e.url + "'"
}

// ErrListen is thrown when the server could not start listening.
type ErrListen struct {
	Cause error
}

func (e *ErrListen) Error() string {
	if e.Cause != nil {
		return "could not start listening: " + e.Cause.Error()
	}

	return "could not start listening"
}

// ErrCookieNotFound is thrown when a requested cookie could not be found.
type ErrCookieNotFound struct {
	name string
}

func (e *ErrCookieNotFound) Error() string {
	return "cookie '" + e.name + "' is not found"
}

// ErrMarshaling is thrown when the content could'nt be marshaled to JSON.
type ErrMarshaling struct{}

func (e *ErrMarshaling) Error() string {
	return "Failed to marshal JSON"
}
