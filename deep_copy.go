package deepcopy

import (
	"reflect"
	"time"
)

// DeepCopy interface 可以自己实现深拷贝函数
type DeepCopy interface {
	DeepCopy() interface{}
}

func Copy(i interface{}) interface{} {
	if i == nil {
		return i
	}

	v := reflect.ValueOf(i)
	// 创建 i 的副本，注意要获取源类型，需要调用 .Elem()
	newVal := reflect.New(v.Type()).Elem()
	parse(v, newVal)
	return newVal.Interface()
}

func parse(oldVal, newVal reflect.Value) {
	// 如果实现了 DeepCopy interface 直接返回
	if oldVal.CanInterface() {
		if dp, ok := oldVal.Interface().(DeepCopy); ok {
			newVal.Set(reflect.ValueOf(dp.DeepCopy()))
			return
		}
	}

	switch oldVal.Kind() {
	case reflect.Invalid:
		return

	case reflect.Array:
		// 创建数组副本，并递归解析
		newVal.Set(reflect.New(oldVal.Type()).Elem())
		for i := 0; i < oldVal.Len(); i++ {
			parse(oldVal.Index(i), newVal.Index(i))
		}

	case reflect.Interface:
		if oldVal.IsNil() {
			return
		}
		// 注意：需要创建 oldVal 源类型的副本，因为这块需要解析 oldVal 的真正类型
		cpy := reflect.New(oldVal.Elem().Type()).Elem()
		parse(oldVal.Elem(), cpy)
		// 递归解析真正类型后，再赋值
		newVal.Set(cpy)
	case reflect.Map:
		if oldVal.IsNil() {
			return
		}
		newVal.Set(reflect.MakeMap(oldVal.Type()))
		for _, key := range oldVal.MapKeys() {
			newVal.SetMapIndex( // 直接调用 Copy 比较方便
				reflect.ValueOf(Copy(key.Interface())),
				reflect.ValueOf(Copy(oldVal.MapIndex(key).Interface())),
			)

		}

	case reflect.Ptr:
		// 非法的类型直接返回
		if !oldVal.Elem().IsValid() {
			return
		}
		newVal.Set(reflect.New(oldVal.Type().Elem()))
		parse(oldVal.Elem(), newVal.Elem())
	case reflect.Slice:
		if oldVal.IsNil() {
			return
		}
		newVal.Set(reflect.MakeSlice(oldVal.Type(), oldVal.Len(), oldVal.Cap()))
		for i := 0; i < oldVal.Len(); i++ {
			parse(oldVal.Index(i), newVal.Index(i))
		}

	case reflect.Struct:
		// 因为 time 包比较常用，所以这里特殊判断下，直接设置源值即可，其他的结构体不可以先全复制，在复制可到出字段
		// 因为不可导出字段也有可能需要深拷贝
		t, ok := oldVal.Interface().(time.Time)
		if ok {
			newVal.Set(reflect.ValueOf(t))
			return
		}

		for i := 0; i < oldVal.NumField(); i++ {
			// 只复制可导出的字段
			if oldVal.Type().Field(i).IsExported() {
				parse(oldVal.Field(i), newVal.Field(i))
			}

		}

	default:
		// 其他类型都可以直接复制
		newVal.Set(oldVal)
	}

}
