package rsp_util

import "reflect"

func SetStructVals(structPtr interface{}, datas map[string]interface{}) {
	if nil == structPtr {
		panic("Fatal error:value of parameters should not be nil")
	}
	if nil == datas {
		return
	}
	rType := reflect.TypeOf(structPtr)
	// structPtr 要为指针才能改变值
	if rType.Kind() != reflect.Ptr {
		panic("not a ptr")
	}

	// 获取指针指向元素类型
	rType = rType.Elem()
	if rType.Kind() != reflect.Struct {
		panic("not struct ptr")
	}
	// 获取指针指向的原始的值
	rVal := reflect.ValueOf(structPtr).Elem()

	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		if data, ok := datas[t.Name]; ok {
			dV := reflect.ValueOf(data)
			if dV.Type().Name() != f.Type().Name() {
				if f.Type().Name() == "" {
					f.Set(dV)
				} else {
					panic("type is not same," + "arrt:" + t.Name + "type is " + t.Type.Name() + ", data type is " + dV.Type().Name())
				}
			} else {
				f.Set(dV)
			}
		}
	}
}

func DeepFields(ifaceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < ifaceType.NumField(); i++ {
		v := ifaceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

// 从 srcStructPtr 复制相同属性值到 dstStructPtr
// ignore 跳过 ignore 值覆
// Deprecated 请使用 https://github.com/jinzhu/copier
func StructCopy(dstStructPtr interface{}, srcStructPtr interface{}, ignoreKey ...string) {
	srcv := reflect.ValueOf(srcStructPtr)
	dstv := reflect.ValueOf(dstStructPtr)
	srct := reflect.TypeOf(srcStructPtr)
	dstt := reflect.TypeOf(dstStructPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := DeepFields(reflect.ValueOf(srcStructPtr).Elem().Type())
	// 把要忽略的转为 set
	set := NewSet()
	if 0 != len(ignoreKey) {
		for _, v := range ignoreKey {
			set.Add(v)
		}
	}

	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		// 需要被忽略的
		if set.Contains(v.Name) {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}
