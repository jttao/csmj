package api

const (
	//代理错误代码
	errorCodeShareBase = 50000 + iota
	errorCodeShareAlreadyGet
	errorCodeShareNotFinish
)

var (
	errorMap = map[int]string{
		errorCodeShareAlreadyGet: "已经领取过了",
		errorCodeShareNotFinish:  "还没有分享",
	}
)

type codeError struct {
	code int
}

func (e *codeError) Code() int {
	return e.code
}

func (e *codeError) Error() string {
	return errorMap[e.code]
}

var (
	errorShareAlreadyGet = &codeError{errorCodeShareAlreadyGet}
	errorShareNotFinish  = &codeError{errorCodeShareNotFinish}
)
