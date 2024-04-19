package ioc

import "reflect"

type BeanMapper map[reflect.Type]reflect.Value

func (b BeanMapper) set(bean interface{}) {
	t := reflect.TypeOf(bean)
	if t.Kind() != reflect.Ptr {
		panic("require ptr object")
	}
	b[t] = reflect.ValueOf(bean)
}

func (b BeanMapper) get(bean interface{}) reflect.Value {
	var t reflect.Type
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		t = reflect.TypeOf(bean)
	}

	if v, ok := b[t]; ok {
		return v
	}
	for k, v := range b {
		if k.Kind() == reflect.Interface && k.Implements(t) {
			return v
		}
	}
	return reflect.Value{}
}
