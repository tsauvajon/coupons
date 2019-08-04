package coupon

// ErrBadRequest : bad user input
type ErrBadRequest struct {
	message string
}

func newErrBadRequest(message string) *ErrBadRequest {
	return &ErrBadRequest{
		message: message,
	}
}

// Error formats the error as a string
func (e *ErrBadRequest) Error() string {
	return e.message
}

// ErrInternal : server or db error
type ErrInternal struct {
	message string
}

func newErrInternal(message string) *ErrInternal {
	return &ErrInternal{
		message: message,
	}
}

// Error formats the error as a string
func (e *ErrInternal) Error() string {
	return e.message
}
