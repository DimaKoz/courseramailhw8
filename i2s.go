package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func i2s(data interface{}, out interface{}) error {

	vData := reflect.ValueOf(data)
	kind := vData.Kind()
	indirect := reflect.ValueOf(out)
	var v reflect.Value
	if indirect.Kind() != reflect.Ptr {
		v = reflect.ValueOf(out)
	} else {
		v = indirect.Elem()
	}

	switch kind {
	case reflect.Map:
		if v.Kind() != reflect.Struct {
			return errors.New("wrong type: out.Kind() != reflect.Struct")
		}
		incoming := make(map[string]reflect.Value)
		for _, key := range vData.MapKeys() {
			strct := vData.MapIndex(key)
			str := fmt.Sprintf("%v", key.Interface())
			incoming[str] = strct
			fmt.Println(key.Interface(), strct.Interface())
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)

			name := v.Type().Field(i).Name
			fmt.Printf("v.Field(%d) is now: %+v type: %s name: %s \n", i, field, field.Type(), name)
			if found, ok := incoming[name]; ok {

				if !field.CanAddr() {
					return errors.New("wrong type: need ptr")
				}
				err := i2s(found.Interface(), field.Addr().Interface())
				if err != nil {
					return err
				}
			}
		}

	case reflect.Float64:
		if v.Kind() != reflect.Int {
			return errors.New("wrong type: out.Kind() != reflect.Int")
		}
		strValue := fmt.Sprintf("%v", vData.Interface())
		intValue, err := strconv.Atoi(strValue)
		if err == nil {
			v.SetInt(int64(intValue))
		}

	case reflect.String:
		if v.Kind() != reflect.String {
			return errors.New("wrong type: out.Kind() != reflect.String")
		}
		strValue := fmt.Sprintf("%v", vData.Interface())
		v.SetString(strValue)

	case reflect.Bool:
		if v.Kind() != reflect.Bool {
			return errors.New("wrong type: out.Kind() != reflect.Bool")
		}
		boolValue := vData.Interface().(interface{}).(bool)
		v.SetBool(boolValue)

	case reflect.Struct:
		err := i2s(vData.Interface(), v.Addr().Interface())
		if err != nil {
			log.Println("err:", err)
		}

	case reflect.Slice:
		sliceIncome := reflect.ValueOf(vData.Interface())
		if v.Type().Kind() != reflect.Slice {
			return errors.New("wrong type")
		}
		outSlice := reflect.MakeSlice(v.Type(), sliceIncome.Len(), sliceIncome.Len())
		for ii := 0; ii < sliceIncome.Len(); ii++ {
			err := i2s(sliceIncome.Index(ii).Interface(), outSlice.Index(ii).Addr().Interface())
			if err != nil {
				return err
			}
		}
		v.Set(outSlice)

	default:
		log.Println("unknown type:", kind)
	}
	return nil
}
