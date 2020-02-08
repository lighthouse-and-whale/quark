package store

type update int
type query int

var (
	Update update
	Query  query
)

// update $:{} Rename,Set,Unset,Inc,Mul,Min,Max

func (*update) Rename(v string) map[string]interface{} {
	// 字段重命名
	return map[string]interface{}{"$rename": v}
}
func (*update) Set(v interface{}) map[string]interface{} {
	// 字段值更新
	return map[string]interface{}{"$set": v}
}
func (*update) Unset(v interface{}) map[string]interface{} {
	// 删除字段
	return map[string]interface{}{"$unset": v}
}
func (*update) Inc(v interface{}) map[string]interface{} {
	// 数值字段增减
	return map[string]interface{}{"$inc": v}
}
func (*update) Mul(v interface{}) map[string]interface{} {
	// 数值字段乘积
	return map[string]interface{}{"$mul": v}
}
func (*update) Min(v interface{}) map[string]interface{} {
	// 指定的值小于原本值就更新
	return map[string]interface{}{"$min": v}
}
func (*update) Max(v interface{}) map[string]interface{} {
	// 指定的值大于原本值就更新
	return map[string]interface{}{"$max": v}
}

// update $:{array:?} AddToSet,Pull,PullAll,PopHead,PopTail,Append

func (*update) AddToSet(v interface{}) map[string]interface{} {
	// 数组字段为空就更新
	return map[string]interface{}{"$addToSet": v}
}
func (*update) Pull(v interface{}) map[string]interface{} {
	// 删除数组中指定值
	return map[string]interface{}{"$pull": v}
}
func (*update) PullAll(v interface{}) map[string]interface{} {
	// 删除数组中多个指定值
	return map[string]interface{}{"$pullAll": v}
}
func (*update) PopHead(array string) map[string]interface{} {
	// 删除数组头
	return map[string]interface{}{"$pop": map[string]interface{}{array: -1}}
}
func (*update) PopTail(array string) map[string]interface{} {
	// 删除数组尾
	return map[string]interface{}{"$pop": map[string]interface{}{array: 1}}
}
func (*update) Append(array string, v interface{}) map[string]interface{} {
	// 数组内容追加
	return map[string]interface{}{"$push": map[string]interface{}{array: map[string]interface{}{"$each": v}}}
}

// query field:{$:?} Gt,Gte,Lt,Lte,Ne,Not,Exists,Regex

func (*query) Gt(v int) map[string]interface{} {
	// '> ' 大于
	return map[string]interface{}{"$gt": v}
}
func (*query) Gte(v int) map[string]interface{} {
	// '>=' 大于等于
	return map[string]interface{}{"$gte": v}
}
func (*query) Lt(v int) map[string]interface{} {
	// '< ' 小于
	return map[string]interface{}{"$lt": v}
}
func (*query) Lte(v int) map[string]interface{} {
	// '<=' 小于等于
	return map[string]interface{}{"$lte": v}
}
func (*query) Ne(v interface{}) map[string]interface{} {
	// '!=' 不等于
	return map[string]interface{}{"$ne": v}
}
func (*query) Not(v interface{}) map[string]interface{} {
	// '! ' 取反
	return map[string]interface{}{"$not": v}
}
func (*query) Exists(v bool) map[string]interface{} {
	// 匹配具有指定字段的文档
	return map[string]interface{}{"$exists": v}
}
func (*query) Regex(v string) map[string]interface{} {
	// 指定字段正则匹配
	return map[string]interface{}{"$regex": v, "$options": "i"}
}

// query field:{$:[]} In,Nin,Mod

func (*query) In(v interface{}) map[string]interface{} {
	// 一个字段的多个匹配
	return map[string]interface{}{"$in": v}
}
func (*query) Nin(v interface{}) map[string]interface{} {
	// 一个字段的多个匹配并取反
	return map[string]interface{}{"$nin": v}
}
func (*query) Mod(v interface{}) map[string]interface{} {
	// 一个字段的余数匹配
	return map[string]interface{}{"$mod": v}
}

// query $:[] Or,And,Nor

func (*query) Or(v interface{}) map[string]interface{} {
	// 多组条件满足一个
	return map[string]interface{}{"$or": v}
}
func (*query) And(v interface{}) map[string]interface{} {
	// 多组条件同时满足
	return map[string]interface{}{"$and": v}
}
func (*query) Nor(v interface{}) map[string]interface{} {
	// 多组条件同时满足并取反
	return map[string]interface{}{"$nor": v}
}

// query array:{$:?} ElemMatch,Size,Slice

func (*query) ElemMatch(v interface{}) map[string]interface{} {
	// 匹配数组的第一个元素
	return map[string]interface{}{"$elemMatch": v}
}
func (*query) Size(v interface{}) map[string]interface{} {
	// 匹配数组指定大小
	return map[string]interface{}{"$size": v}
}
func (*query) Slice(v interface{}) map[string]interface{} {
	// 数组切片
	return map[string]interface{}{"$slice": v}
}

// query array:{$:[]} All

func (*query) All(v interface{}) map[string]interface{} {
	// 数组字段具有多个内容的匹配
	return map[string]interface{}{"$all": v}
}
