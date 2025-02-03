package result

const Success code = "0x0000"

const (
	ErrInternalServerError code = "0x0010"
	ErrBadRequest          code = "0x0011"
	ErrRequestTimeout      code = "0x0012"
	ErrValidationFailed    code = "0x0013"
	ErrDataBaseOperation   code = "0x0014"

	ErrUnauthorized  code = "0x0020"
	ErrAuthForbidden code = "0x0021"

	ErrResourceNotFound     code = "0x0030"
	ErrResourceAlreadyExist code = "0x0031"
)

var MapCodeMessage = map[code]string{
	Success:                 "success",
	ErrInternalServerError:  "internal server error",
	ErrBadRequest:           "bad request",
	ErrRequestTimeout:       "request timeout",
	ErrValidationFailed:     "validation failed",
	ErrDataBaseOperation:    "database operation failed",
	ErrUnauthorized:         "unauthorized",
	ErrAuthForbidden:        "authorization forbidden",
	ErrResourceNotFound:     "resource not found",
	ErrResourceAlreadyExist: "resource already exist",
}

type code string

func (c code) Code() code {
	return c
}

func (c code) Data() any {
	return nil
}

func (c code) Error() error {
	return nil
}

// With catch error
func (c code) With(err error) Result {
	return &result{
		code: c,
		err:  err,
	}
}
