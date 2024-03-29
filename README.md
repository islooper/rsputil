# rsp-util 

## 错误码自定义文言

### 概述
由于通用错误码（如：1400参数错误）在多个场景下可能涉及到不同的参数，导致调用方在遇到错误时难以理解具体是哪个参数出了问题，同时我们在排查问题时也会很麻烦。

### 需求
- 支持对于同一个错误码，根据出错的具体场景自定义不同的错误信息，当然，很明确的码值就不建议随意变更文言了。
- 使调用方和开发者能够更直观、更快速地定位和理解错误原因。

## 反射机制支持

### 概述
为了提高代码的灵活性和扩展性，需要引入反射机制。

### 需求
- 实现反射机制，支持直接将信息反射到响应（rsp）对象中。
- 可参考测试用例来了解具体的实现方法。

## 自动记录日志

### 概述
可以自动写日志，省去项目中繁多的log.xxx 影响美观

目前只能记录最终返回的log日志，但后续可以支持一次context下的所有error以链表形式自动记录（不一定写得出来）

## 显隐error

### 概述
支持对外返回一个error 对内用一个error。 对内的就可以记录更多的信息

```markdown
```golang

const (
ParamErr = 1400
)

// 基础测试
type CreateSharingRsp struct {
Code int32
Message string
Data string
}

// 指针
type CarsharingDetail struct {
Code int32   `json:"code"`
Message string  `json:"message"`
Detail  *Detail `json:"data"`
}

type Detail struct {
ShareId string
Vin string
}

// 切片指针
type CarsharingListRsp struct {
Code int32   `json:"code"`
Message string  `json:"message"`
List    []*List `json:"data"`
}

type List struct {
ShareId string
Vin string
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

```
```
