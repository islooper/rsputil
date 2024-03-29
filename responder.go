package rsp_util

import "log"

//// Responder 响应接口
//type Responder interface {
//	SetRet(ret int32)
//	SetMsg(msg string)
//}

type Error interface {
	GetCode() int32
	GetMessage() string
	GetErr() error
}

// ErrInfo 定义错误
type ErrInfo struct {
	Code    int32  // 错误码
	Message string // 展示给用户看的
	Err     error  // 保存内部错误信息
}

func (e *ErrInfo) GetCode() int32 {
	return e.Code
}

func (e *ErrInfo) GetMessage() string {
	return e.Message
}

func (e *ErrInfo) GetErr() error {
	return e.Err
}

func (e *ErrInfo) Is(info *ErrInfo) bool {
	if info == nil {
		return false
	}
	if e.Code == info.Code {
		return true
	} else {
		return false
	}
}

func (e *ErrInfo) IsErrNot() bool {
	return e.Is(ErrNot)
}

var (
	ErrNot    = &ErrInfo{1000, "success", nil}
	ErrUnknow = &ErrInfo{1500, "unknown error", nil}

	ErrParam    = &ErrInfo{1400, "param error", nil}
	ErrDataBase = &ErrInfo{1900, "database error", nil}
)

func NewErrInfo(info *ErrInfo, msg string, err error) *ErrInfo {
	errInfo := &ErrInfo{
		Code:    info.Code,
		Message: info.Message,
		Err:     info.Err,
	}
	if msg != "" {
		errInfo.Message = msg
	}

	if err != nil {
		errInfo.Err = err
	}
	return errInfo
}

func WriteRsp(rspPtr interface{}, svcError Error, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	if svcError == nil {
		panic("error is nil")
	}

	data["Code"] = svcError.GetCode()
	data["Message"] = svcError.GetMessage()

	SetStructVals(rspPtr, data)
}

func WriteRspWithLog(rspPtr interface{}, svcError Error, data map[string]interface{}) {
	WriteRsp(rspPtr, svcError, data)
	go log.Fatal(svcError.GetErr())
}

// WriteRpcRspWithMsg noinspection ALL
func WriteRpcRspWithMsg(rspPtr interface{}, svcError Error, msg string, data map[string]interface{}) {

	if nil == data {
		data = make(map[string]interface{})
	}
	if svcError == nil {
		panic("rpcError is nil")
	}

	data["Code"] = svcError.GetCode()
	data["Message"] = msg

	SetStructVals(rspPtr, data)
}

// WriteRpcRspWithMsgLog noinspection ALL
func WriteRpcRspWithMsgLog(rspPtr interface{}, svcError Error, msg string, data map[string]interface{}) {
	WriteRpcRspWithMsg(rspPtr, svcError, msg, data)
	go log.Fatal(svcError.GetErr())
}
