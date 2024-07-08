package handle

import (
	"application-web/logger"

	"application-web/pkg/buserr"
	"application-web/pkg/dto"
)

func Result(code int, result interface{}, msg string) dto.Result {
	return dto.Result{
		Code:   code,
		Msg:    msg,
		Result: result,
	}
}

func Ok() dto.Result {
	return Result(buserr.Ok, nil, "ok")
}
func OkWithMsg(msg string) dto.Result {
	return Result(buserr.Ok, nil, msg)
}

func Success(result interface{}) dto.Result {
	return Result(buserr.Ok, result, "ok")
}

func Error(error string) dto.Result {
	return dto.Result{
		Code: buserr.Err,
		Msg:  error,
	}
}

func Fail(code int, msg string) dto.Result {
	return Result(code, nil, msg)
}

func FailWithMsg(code int, msg string) dto.Result {
	return Result(code, nil, msg)
}

func HandleError(err error, codes ...interface{}) {
	if err != nil {
		if len(codes) > 0 {
			logger.Error("%v\n%s", codes, err.Error())
		} else {
			// panic(err.Error())
		}
	}
}
