package dict

// SimpleDict wraps a map, it is not thread safe
type SimpleDict struct {
	m map[string]interface{}
}

// MakeSimple 构造基础字典
func MakeSimple() *SimpleDict {
	return &SimpleDict{
		m: make(map[string]interface{}),
	}
}

func (dict *SimpleDict) Inc(key string, val int) (result int) {
	return 0
}

// Get 返回字典Key对应值
func (dict *SimpleDict) Get(key string) (val interface{}, exists bool) {
	val, ok := dict.m[key]
	return val, ok
}

// Len 返回Dict元素数量
func (dict *SimpleDict) Len() int {
	if dict.m == nil {
		panic("m is nil")
	}
	return len(dict.m)
}

func (dict *SimpleDict) Exists(key string) (exists bool) {
	_, exists = dict.m[key]
	return
}

// Put 设置Key的值
func (dict *SimpleDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m[key]
	dict.m[key] = val
	if existed {
		return 0
	}
	return 1
}

// PutIfAbsent 设置数据(不存在Key才设置)
func (dict *SimpleDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := dict.m[key]
	if existed {
		return 0
	}
	dict.m[key] = val
	return 1
}

// PutIfExists 设置数据(如果存在Key)
func (dict *SimpleDict) PutIfExists(key string, val interface{}) (result int) {
	_, existed := dict.m[key]
	if existed {
		dict.m[key] = val
		return 1
	}
	return 0
}

// Remove 删除Key数据
func (dict *SimpleDict) Remove(key string) (result int) {
	_, existed := dict.m[key]
	delete(dict.m, key)
	if existed {
		return 1
	}
	return 0
}

// Keys 返回所有的Hash keys slice
func (dict *SimpleDict) Keys() []string {
	result := make([]string, len(dict.m))
	i := 0
	for k := range dict.m {
		result[i] = k
		i++
	}
	return result
}

func (dict *SimpleDict) ForScanKeys(eachFn EachFunc, start int, count int) [][]byte {
	//TODO implement me
	panic("implement me")
}

// ForEach 遍历所有的Keys
func (dict *SimpleDict) ForEach(eachFu EachFunc) {
	for k, v := range dict.m {
		if !eachFu(k, v) {
			break
		}
	}
}

// RandomKeys 按数量返回随机的Keys, 可能包含重复的Key
func (dict *SimpleDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		for k := range dict.m {
			result[i] = k
			break
		}
	}
	return result
}

// RandomDistinctKeys 按数量返回随机的Keys, 不包含重复的Key
func (dict *SimpleDict) RandomDistinctKeys(limit int) []string {
	size := limit
	if size > len(dict.m) {
		size = len(dict.m)
	}
	result := make([]string, size)
	i := 0
	for k := range dict.m {
		if i == limit {
			break
		}
		result[i] = k
		i++
	}
	return result
}

// Clear 清除所有的Keys
func (dict *SimpleDict) Clear() {
	*dict = *MakeSimple()
}
