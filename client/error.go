package client

type HttpError struct {
	err    string
	Status int
}

func NewHttpError(err string, status int) *HttpError {
	return &HttpError{
		err:    err,
		Status: status,
	}
}

func (e *HttpError) Error() string {
	return e.err
}

// *HttpError implements error
var _ error = (*HttpError)(nil)
