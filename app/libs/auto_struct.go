package libs

import "reflect"

type (
	Builder struct {
		fieldId []reflect.StructField // 存储结构体字段
	}

	Struct struct {
		typ   reflect.Type
		index map[string]int
	}
)

func NewBuilder() *Builder {
	return &Builder{}
}

// 新增字段
func (b *Builder) AddField(field string, typ reflect.Type) *Builder {
	b.fieldId = append(b.fieldId, reflect.StructField{
		Name: field,
		Type: typ,
	})

	return b
}

func (b *Builder) Build() *Struct {
	stu := reflect.StructOf(b.fieldId)
	index := make(map[string]int)

	for i := 0; i < stu.NumField(); i++ {
		index[stu.Field(i).Name] = i
	}

	return &Struct{
		typ:   stu,
		index: index,
	}
}

func (b *Builder) AddString(val string) *Builder {
	b.AddField(val, reflect.TypeOf(""))
	return b
}
