package maps

import "reflect"

// RefToMap 将结构体中的字段根据标签转换为map
// data: 需要转换的结构体或结构体指针
// tag: 结构体字段上的标签名称，用于在map中作为键
// 返回值: 一个map，其中键是结构体字段上的标签值，值是对应字段的值
func RefToMap(data any, tag string) map[string]any {
	// 初始化一个空的map用于存放结果
	maps := map[string]any{}
	// 获取data的类型
	t := reflect.TypeOf(data)
	// 获取data的值
	v := reflect.ValueOf(data)

	// 遍历结构体的每个字段
	for i := 0; i < t.NumField(); i++ {
		// 获取当前字段的信息
		field := t.Field(i)
		// 通过tag获取当前字段的标签值
		getTag, ok := field.Tag.Lookup(tag)
		// 如果字段没有指定的标签，则跳过当前字段
		if !ok {
			continue
		}
		// 获取当前字段的值
		val := v.Field(i)
		// 如果字段值是零值，则跳过当前字段
		if val.IsZero() {
			continue
		}

		// 如果当前字段是结构体，则递归转换该结构体为map
		if field.Type.Kind() == reflect.Struct {
			newMaps := RefToMap(val.Interface(), tag)
			maps[getTag] = newMaps
			continue
		}

		// 如果当前字段是指向结构体的指针，则递归转换该指针指向的结构体为map
		if field.Type.Kind() == reflect.Ptr {
			if field.Type.Elem().Kind() == reflect.Struct {
				newMaps := RefToMap(val.Elem().Interface(), tag)
				maps[getTag] = newMaps
				continue
			}
			// 如果当前字段是指向非结构体的指针，则直接取指针的值放入map
			maps[getTag] = val.Elem().Interface()
			continue
		}
		// 如果当前字段既不是结构体也不是指针，则直接将值放入map
		maps[getTag] = val.Interface()

	}
	// 返回转换后的map
	return maps
}
