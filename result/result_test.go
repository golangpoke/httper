package result_test

import (
	"fmt"
	"github.com/golangpoke/httper/result"
	"os"
	"testing"
)

func TestResult(t *testing.T) {
	err := os.ErrNotExist
	resultImplement := result.ErrInternalServerError.Wrap(err)
	HandleResult(resultImplement)
	codeImplement := result.ErrBadRequest
	HandleResult(codeImplement)
	mapImplement := result.Map{
		"hello": "world",
	}
	HandleResult(mapImplement)
}

func HandleResult(rs result.Result) {
	if rs.Error() != nil || rs.Code() != result.Success {
		fmt.Printf("fail code: %s, error: %v\n", rs.Code(), rs.Error())
		return
	}
	fmt.Printf("success code: %s, data: %v\n", rs.Code(), rs.Data())
}
