package response

import "bytes"

var (
	Success	= &BaseErr{Code: 0, Message: "Success"}
	ServerError = &BaseErr{Code: 10001, Message: "Server Error"}
	FileExist = &BaseErr{Code: 10002, Message: "The file already exists"}
)

type BaseErr struct {
	Code int
	Message string
}

func (err BaseErr) Error() string {
	return err.Message
}

type AllErr struct {
	Code int
	Message string
	Err string
}

func New(errno *BaseErr, err string) *AllErr {
	return &AllErr{Code: errno.Code, Message: errno.Message, Err:err}
}

func (err *BaseErr) Add(message string) error {
	err.Message += " " + message

	return err
}

func (err *BaseErr) AddParam(param string) error {
	var buffer bytes.Buffer
	buffer.WriteString(err.Message)
	buffer.WriteString(" ")
	buffer.WriteString(param)
	err.Message = buffer.String()

	return err
}

func (err *AllErr) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(err.Message)
	buffer.WriteString(err.Err)

	return buffer.String()
}

func formatErr(err error) (int, string) {
	if err == nil {
		return Success.Code, Success.Message
	}

	switch typed := err.(type) {
	case *AllErr:
		return typed.Code, typed.Error()
	case *BaseErr:
		return typed.Code, typed.Message
	}

	return ServerError.Code, err.Error()
}