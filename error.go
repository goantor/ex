package ex

import "fmt"

var (
	errorsMapping = Errors{
		UNKNOWN_ERRNO: "unknown error",
	}

	UNKNOWN_ERRNO FatalErrno = -9999
)

func ErrorsRegister(errors map[IErrno]string) {
	for no, err := range errors {
		ErrorSave(no, err)
	}
}

func ErrorsInit(errors map[IErrno]string) {
	if errorsMapping == nil {
		errorsMapping = errors
		errorsMapping[UNKNOWN_ERRNO] = "unknown error"
	}
}

func ErrorSave(no IErrno, message string) {
	if _, ok := errorsMapping[no]; ok {
		panic(fmt.Sprintf("error is exist [%d]:%s", no, errorsMapping[no]))
	}

	errorsMapping[no] = message
}

type Errors map[IErrno]string

type IErrno interface {
	Level() string
}

type Errno int

func (e Errno) Level() string {
	return "base"
}

func (e Errno) Error() string {
	return errorMessage(e)
}

type InfoErrno Errno

func (e InfoErrno) Level() string {
	return "info"
}

func (e InfoErrno) Error() string {
	return errorMessage(e)
}

type TraceErrno Errno

func (e TraceErrno) Level() string {
	return "trace"
}

func (e TraceErrno) Error() string {
	return errorMessage(e)
}

type DebugErrno Errno

func (e DebugErrno) Level() string {
	return "debug"
}

func (e DebugErrno) Error() string {
	return errorMessage(e)
}

type WarnErrno Errno

func (e WarnErrno) Level() string {
	return "warning"
}

func (e WarnErrno) Error() string {
	return errorMessage(e)
}

type ErrorErrno Errno

func (e ErrorErrno) Level() string {
	return "error"
}

func (e ErrorErrno) Error() string {
	return errorMessage(e)
}

type FatalErrno Errno

func (e FatalErrno) Level() string {
	return "fatal"
}

func (e FatalErrno) Error() string {
	return errorMessage(e)
}

type PanicErrno Errno

func (e PanicErrno) Level() string {
	return "panic"
}

func (e PanicErrno) Error() string {
	return errorMessage(e)
}

func GetErrorMessage(no IErrno) string {
	if message, ok := errorsMapping[no]; ok {
		return message
	}

	return "unknown error"
}

type Error struct {
	Code    IErrno // 错误码
	Data    interface{}
	Message string // 内容
}

//func (e Error) Level() string {
//	return "error"
//}
//
//func (e Error) Error() string {
//	return e.Message
//}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    ErrorErrno(code),
		Message: message,
	}
}

func QuickThrowError(code int, message string) {
	e := NewError(code, message)
	throw(e)
}

func New(code IErrno, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

func ThrowError(error interface{}, args ...interface{}) {
	var (
		data interface{} = nil
	)

	if len(args) > 0 {
		data = args[0]
	}

	e := getError(error)

	e.Data = data
	//log.Auto(e.Code, e.Message, data)
	throw(e)
}

func getError(error interface{}) *Error {
	var (
		e     *Error
		errno IErrno
	)

	switch error.(type) {
	case IErrno:

		errno = error.(IErrno)
		e = &Error{
			Code:    errno,
			Message: errorMessage(errno),
		}

		break
	case Error:
		ie := error.(Error)
		e = &ie

		break
	case *Error:
		e = error.(*Error)
		break
	default:
		e = &Error{
			Code:    UNKNOWN_ERRNO,
			Message: "unknown error",
			Data:    error,
		}
	}

	return e
}

func throw(e *Error) {
	panic(e)
}

func errorMessage(errno IErrno, def ...string) string {
	err, ok := errorsMapping[errno]
	if ok {
		return err
	}

	if len(def) == 0 {
		return errorsMapping[UNKNOWN_ERRNO]
	}

	return def[0]
}

func Recover(fn func(r interface{}, message string)) {
	if r := recover(); r != nil {
		switch r.(type) {
		case IErrno:
			fn(r, errorsMapping[r.(IErrno)])
			break
		case *Error:
			fn(r, "")
			break
		default:
			fn(UNKNOWN_ERRNO, fmt.Sprintf("%+v", r))
		}
	}
}
