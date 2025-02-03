package result

type Result interface {
	Code() code   // service code
	Data() any    // response data
	Error() error // api/service error
}

type result struct {
	code code
	data any
	err  error
}

func (r *result) Code() code {
	return r.code
}

func (r *result) Data() any {
	return r.data
}

func (r *result) Error() error {
	return r.err
}

type Map map[string]any

func (m Map) Code() code {
	return Success
}

func (m Map) Data() any {
	return m
}

func (m Map) Error() error {
	return nil
}
