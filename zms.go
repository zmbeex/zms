package zms

import (
	"github.com/zmbeex/gkit"
	"reflect"
)

type Zms struct {
	clientParams *Params
	clientResult *Result
	UserId       int64 `title:"用户ID"`
}

func InitZms(params *Params, result *Result) *Zms {
	z := new(Zms)
	z.clientParams = params
	z.clientResult = result
	return z
}

func (z *Zms) GetParams(params interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			gkit.Info("请求入参：" + gkit.SetJson(z.clientParams))
			panic(r)
		}
	}()
	data := make(map[string]interface{})
	err := gkit.GetJson(z.clientParams.Params, &data)
	gkit.CheckPanic(err, "参数异常")
	gkit.Info("请求入参：" + gkit.SetJson(&data))

	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)
	if t.Kind() != reflect.Ptr {
		gkit.Panic("入参必须是一个结构体指针")
	}
	t = t.Elem()
	v = v.Elem()
	if t.Kind() != reflect.Struct {
		gkit.Panic("入参必须是一个结构体指针")
	}
	for i := 0; i < t.NumField(); i++ {
		t0 := t.Field(i)
		v0 := v.Field(i)
		key := t0.Name
		check := t0.Tag.Get("check")
		title := t0.Tag.Get("title")
		val := data[gkit.StringFirstLower(key)]
		err := gkit.CheckValue(check, title, gkit.ToString(val))
		gkit.CheckPanic(err, "参数异常")
		if val == "" {
			val = t0.Tag.Get("defaultValue")
		}
		v0.Set(reflect.ValueOf(gkit.GetReflectValue(t0.Type.Kind(), val)))
	}
}

// 设置状态
func (z *Zms) SetStatus(status int) {
	z.clientResult.Status = status
}

// 设置提示
func (z *Zms) SetNote(note string) {
	z.clientResult.Note = note
}

// 设置结果
func (z *Zms) SetResult(result interface{}) {
	t := reflect.TypeOf(result)
	val := reflect.ValueOf(result)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		sliceLen := val.Len()
		for i := 0; i < sliceLen; i++ {
			z.clientResult.List = append(z.clientResult.List, val.Index(i).Interface())
		}
	}
	if t.Kind() == reflect.Map || t.Kind() == reflect.Struct {
		z.clientResult.List = append(z.clientResult.List, result)
	}
}

// 设置值
func (z *Zms) SetData(key string, v interface{}) {
	if z.clientResult.Data == nil {
		z.clientResult.Data = make(map[string]interface{})
	}
	z.clientResult.Data[key] = v
}
