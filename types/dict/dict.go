package dict

// EachFunc 遍历Key回调函数
type EachFunc func(key string, val interface{}) bool

// Dict 接口封装
type Dict interface {
	Get(key string) (val interface{}, exists bool)
	Len() int
	Exists(key string) (exists bool)
	Put(key string, val interface{}) (result int)
	Inc(key string, val int) (result int)
	PutIfAbsent(key string, val interface{}) (result int)
	PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	ForEach(eachFn EachFunc)
	ForScanKeys(eachFn EachFunc, start int, count int) [][]byte
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
