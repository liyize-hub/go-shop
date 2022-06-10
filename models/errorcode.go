package models

type errorCode struct {
	SUCCESS    int
	ERROR      int
	NotFound   int
	LoginError int
}

// ErrorCode 错误码
var ErrorCode = errorCode{
	SUCCESS:    0,
	ERROR:      1,
	NotFound:   2,
	LoginError: -1, //用户名或密码错误
}
