package client

type HttpError struct {
	err        string
	StatusCode int
}

func NewHttpError(err string, statusCode int) *HttpError {
	return &HttpError{
		err:        err,
		StatusCode: statusCode,
	}
}

func (e *HttpError) Error() string {
	return e.err
}

// *HttpError implements error
var _ error = (*HttpError)(nil)
