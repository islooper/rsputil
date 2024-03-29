package rsp_util

import (
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
	"testing"
)

const (
	ParamErr = 1400
)

// 基础测试
type CreateSharingRsp struct {
	Code    int32
	Message string
	Data    string
}

// 指针
type CarsharingDetail struct {
	Code    int32   `json:"code"`
	Message string  `json:"message"`
	Detail  *Detail `json:"data"`
}

type Detail struct {
	ShareId string
	Vin     string
}

// 切片指针
type CarsharingListRsp struct {
	Code    int32   `json:"code"`
	Message string  `json:"message"`
	List    []*List `json:"data"`
}

type List struct {
	ShareId string
	Vin     string
}

func TestWriteDefaultSuccessRsp(t *testing.T) {
	// 基础返回
	rsp := new(CreateSharingRsp)

	WriteRsp(rsp, ErrNot, map[string]interface{}{
		"Data": ksuid.New().String(),
	})

	fmt.Println(fmt.Sprintf("%+v", rsp))
	WriteRpcRspWithMsg(rsp, ErrNot, "test message", map[string]interface{}{
		"Data": ksuid.New().String(),
	})
	fmt.Println(fmt.Sprintf("%+v", rsp))

	// 指针返回
	detailRsp := new(CarsharingDetail)
	detail := &Detail{
		ShareId: ksuid.New().String(),
		Vin:     ksuid.New().String(),
	}
	WriteRsp(detailRsp, ErrNot, map[string]interface{}{
		"Detail": detail,
	})

	fmt.Println(fmt.Sprintf("%+v", detailRsp))

	// 切片指针
	rspSlice := new(CarsharingListRsp)
	lists := make([]*List, 0)

	l1 := &List{
		ShareId: ksuid.New().String(),
		Vin:     ksuid.New().String(),
	}

	l2 := &List{
		ShareId: ksuid.New().String(),
		Vin:     ksuid.New().String(),
	}

	lists = append(lists, l1, l2)

	WriteRsp(rspSlice, ErrNot, map[string]interface{}{
		"List": lists,
	})

	fmt.Println(fmt.Sprintf("%+v", rspSlice))

	if len(rspSlice.List) != 2 {
		t.Error("rsp.list is not right")
	}

	// nil test

	//var nilRsp CarsharingListRsp
	//
	//WriteRsp(nilRsp, ErrNot, map[string]interface{}{
	//	"List": lists,
	//})

}

func TestWriteRspNil(t *testing.T) {
	rsp := new(CreateSharingRsp)

	WriteRsp(rsp, ErrNot, nil)

	fmt.Println(rsp)
	if rsp.Code != 1000 {
		t.Error("rsp.Code is not 1000")
	}

	if rsp.Message != "success" {
		t.Error("rsp.Message is not success")
	}
}

func TestCustomError(t *testing.T) {
	// e.g 参数错误
	paramErr := NewErrInfo(&ErrInfo{
		Code:    ParamErr,
		Message: "参数错误（默认message）",
		Err:     errors.New("不会返回的error"),
	},
		"优先级最高message", errors.New("优先级最高error"))

	// 基础返回
	rsp := new(CreateSharingRsp)
	WriteRsp(rsp, paramErr, map[string]interface{}{
		"Data": "",
	})

	// 完全不定义结构体
	WriteRsp(rsp, paramErr, map[string]interface{}{
		"Data": "",
	})

	fmt.Println(fmt.Sprintf("%+v", rsp))
	WriteRpcRspWithMsg(rsp, ErrNot, "test message", map[string]interface{}{
		"Data": ksuid.New().String(),
	})
	fmt.Println(fmt.Sprintf("%+v", rsp))

}

func TestErrorLog(t *testing.T) {
	// e.g 各种参数错误的情况
	paramErr := NewErrInfo(&ErrInfo{
		Code:    ParamErr,
		Message: "参数错误（默认message）",
		Err:     errors.New("不会返回的error"),
	},
		"优先级最高message", nil)

	// 基础返回
	rsp := new(CreateSharingRsp)
	WriteRspWithLog(rsp, paramErr, map[string]interface{}{
		"Data": "",
	})

	fmt.Println(fmt.Sprintf("%+v", rsp))
}
