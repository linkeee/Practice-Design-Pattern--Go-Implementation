package db

import (
	"reflect"
	"strings"
)

type record struct {
	primaryKey interface{}
	fields     map[string]int
	values     []interface{}
}

func recordFrom(key interface{}, value interface{}) (r record, e error) {
	defer func() {
		if err := recover(); err != nil {
			r = record{}
			e = ErrRecordTypeInvalid
		}
	}()
	vType := reflect.TypeOf(value)
	vVal := reflect.ValueOf(value)
	if vVal.Type().Kind() == reflect.Pointer {
		vType = vType.Elem()
		vVal = vVal.Elem()
	}
	record := record{
		primaryKey: key,
		fields:     make(map[string]int, vVal.NumField()),
		values:     make([]interface{}, vVal.NumField()),
	}
	for i := 0; i < vVal.NumField(); i++ {
		fieldType := vType.Field(i)
		fieldVal := vVal.Field(i)
		name := strings.ToLower(fieldType.Name)
		record.fields[name] = i
		record.values[i] = fieldVal.Interface()
	}
	return record, nil
}

func (r record) convertByValue(result interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = ErrRecordTypeInvalid
		}
	}()
	rType := reflect.TypeOf(result)
	rVal := reflect.ValueOf(result)
	if rType.Kind() == reflect.Pointer {
		rType = rType.Elem()
		rVal = rVal.Elem()
	}
	for i := 0; i < rType.NumField(); i++ {
		field := rVal.Field(i)
		field.Set(reflect.ValueOf(r.values[i]))
	}
	return nil
}

func (r record) convertByType(rType reflect.Type) (result interface{}, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = ErrRecordTypeInvalid
		}
	}()
	if rType.Kind() == reflect.Pointer {
		rType = rType.Elem()
	}
	rVal := reflect.New(rType)
	return rVal, nil
}
