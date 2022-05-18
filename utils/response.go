package utils

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

/*/ SendErrJSON 有错误发生时，发送错误JSON
func SendErrJSON(msg string, ctx iris.Context) {

	ctx.JSON(iris.Map{
		"errNo": models.ErrorCode.ERROR,
		"msg":   msg,
		"data":  iris.Map{},
	})
}*/

type Response struct {
	ErrNo  int         `json:"status"`
	ErrMsg string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func NewJSONResponse(ctx iris.Context, errno int, msg string, data interface{}) mvc.Response {

	var (
		jsonData Response
	)
	jsonData = Response{
		ErrNo:  errno,
		ErrMsg: msg,
		Data:   data,
	}
	Logger.Info("返回jsonData", zap.Any("response", jsonData))
	return mvc.Response{Object: jsonData}
}

/*
// 状态可能非200
var err error
switch e := data.(type) {

case *models.Err: // 复合错误
	err = e
	jsonData.ErrNo = e.Code()
	jsonData.ErrMsg = e.ErrorErrStack()
case errlib.ErrCode: // 错误码
	err = e
	jsonData.ErrNo = e.Code()
	jsonData.ErrMsg = e.Error()
default: //未知错误
	err, _ = e.(error)
	jsonData.ErrNo = errlib.Unknown.Code()
	jsonData.ErrMsg = err.Error()
}

var httpCode int
switch err {
case errcode.ErrSystemBusy:
	// 明确要求上游重试的错误
	httpCode = http.StatusInternalServerError
default:
	// 200
	httpCode = http.StatusOK
}
logit.AddAllLevel(ctx, logit.Int(ghttp.LogFieldLogStatus, httpCode))
if httpCode != http.StatusOK {
	gdp2log.Warning(ctx, jsonData.ErrMsg, jsonData.ErrNo) // 团队约定的warning日志
}

return ghttp.NewJSONResponse(httpCode, jsonData)
*/
