package zms

// 参数
type Params struct {
	Code   string `title:"服务编码"`
	Token  string `title:"认证签名"`
	Params string `title:"请求入参"`
	Uuid   string `title:"唯一标识"`
}

// 返回数据
type Result struct {
	Uuid   string                 `title:"唯一标识"`
	Code   string                 `title:"服务编码"`
	Status int                    `title:"状态"`
	Note   string                 `title:"提示信息"`
	Data   map[string]interface{} `title:"返回结果/任意类型"`
	List   []interface{}          `title:"返回结果，数组"`
}

func InitZms(params *Params, result *Result) *Zms {
	z := new(Zms)
	z.clientParams = params
	z.clientResult = result
	return z
}
