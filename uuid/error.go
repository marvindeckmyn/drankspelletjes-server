package uuid

// ErrNotInstantiated is returned when the etcd lib was not yet initialized.
type ErrMalformed struct {
	Cause error
}

func (e *ErrMalformed) Error() string {
	if e.Cause != nil {
		return "malformed UUID: " + e.Cause.Error()
	}

	return "malformed UUID"
}
