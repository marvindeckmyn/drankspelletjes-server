package cdb

// ErrConnect is thrown when the application could not connect to the dbms.
type ErrConnect struct {
	Cause error
}

func (e *ErrConnect) Error() string {
	if e.Cause != nil {
		return "failed to connect: " + e.Cause.Error()
	}

	return "failed to connect"
}

// ErrNotInstantiated is returned when the pool was not yet connected.
type ErrNotInstantiated struct{}

func (e *ErrNotInstantiated) Error() string {
	return "pool is not yet instantiated"
}

// ErrQuery is returned when a generic postgress query was erroneous.
type ErrQuery struct {
	Cause error
}

func (e *ErrQuery) Error() string {
	if e.Cause != nil {
		return "query error: " + e.Cause.Error()
	}

	return "query error"
}

// ErrMissingValue is returned when someone tried to execute a query without the required values.
type ErrMissingValue struct {
	Value string
}

func (e *ErrMissingValue) Error() string {
	return "Missing value: " + e.Value
}

// ErrMissingValue is returned when someone tried to execute a query without the required values.
type ErrFailedToParseQuery struct {
	Cause error
}

func (e *ErrFailedToParseQuery) Error() string {
	if e.Cause != nil {
		return "Failed to parse query" + e.Cause.Error()
	}

	return "Failed to parse Query"
}

// ErrMissingValue is returned when someone tried to execute a result without the required values.
type ErrFailedToParseResult struct {
	Cause error
}

func (e *ErrFailedToParseResult) Error() string {
	if e.Cause != nil {
		return "Failed to parse result" + e.Cause.Error()
	}

	return "Failed to parse Query"
}

// ErrInvalidParseMethod is thrown when the destination to parse to has an invalid parse method.
type ErrInvalidParseMethod struct {
}

func (e *ErrInvalidParseMethod) Error() string {
	return "Invalid parse method"
}

// ErrParseResult is thrown when the parse function returned an error.
type ErrParseResult struct {
	Cause error
}

func (e *ErrParseResult) Error() string {
	if e.Cause != nil {
		return "Failed to parse: " + e.Cause.Error()
	}

	return "Failed to parse"
}

// ErrMalformed is returned when an object in de database was malformed.
type ErrMalformed struct {
	Cause error
}

func (e *ErrMalformed) Error() string {
	if e.Cause != nil {
		return "malformed object: " + e.Cause.Error()
	}

	return "malformed object"
}

// ErrNoSuchKey is returned when a certain key isn't present in a result.
type ErrNoSuchKey struct {
	Key string
}

func (e *ErrNoSuchKey) Error() string {
	return "key '" + e.Key + "' doesn't exist"
}

// ErrMissingResult is returned when a statement returned no results.
type ErrMissingResult struct{}

func (e *ErrMissingResult) Error() string {
	return "No results found"
}

// ErrCreateStmt is returned when a prepare function failed.
type ErrCreateStmt struct {
	Cause error
}

func (e *ErrCreateStmt) Error() string {
	if e.Cause != nil {
		return "Error creating statement: " + e.Cause.Error()
	}

	return "Error creating statement"
}
