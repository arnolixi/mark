package ioc

import (
	"log"
	"reflect"
)

type BeanFactoryImpl struct {
	bm BeanMapper
}

func NewBFI() *BeanFactoryImpl {
	return &BeanFactoryImpl{bm: make(BeanMapper)}
}

var BFI *BeanFactoryImpl

const injectTag = "inject"

func init() {
	BFI = NewBFI()
}

func (impl *BeanFactoryImpl) Set(args ...interface{}) {
	if args == nil || len(args) == 0 {
		return
	}
	for _, v := range args {
		impl.bm.set(v)
	}
}

func (impl *BeanFactoryImpl) Get(key interface{}) interface{} {
	if key == nil {
		return nil
	}
	val := impl.bm.get(key)
	if val.IsValid() {
		return val.Interface()
	}
	return nil
}

func (impl *BeanFactoryImpl) Apply(bean interface{}) {
	if bean == nil {
		return
	}
	val := reflect.ValueOf(bean)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fVal, ok := field.Tag.Lookup(injectTag)
		if val.Field(i).CanSet() && ok {
			if fVal == "" || fVal == "-" {
				if getV := impl.Get(field.Type); getV != nil {
					val.Field(i).Set(reflect.ValueOf(getV))
					impl.Apply(getV)
				} else {
					log.Printf(
						"inject failed     Error: %s;  ObjectName: %s   field: %s     InjectName: %s",
						"not found field struct",
						val.Type().Name(),
						field.Name,
						field.Type.Name(),
					)
				}
			} else {
				// TODO 表达式注入
				continue
			}
		}
	}

}

// ApplyAll 递归注入
func (impl *BeanFactoryImpl) ApplyAll() {
	for t, v := range impl.bm {
		if t.Elem().Kind() == reflect.Struct {
			impl.Apply(v.Interface())
		}
	}
}
