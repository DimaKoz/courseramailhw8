package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func i2s(data interface{}, out interface{}) error {

	incoming := make(map[string]reflect.Value)
	// See what the map has now
	fmt.Printf("mp is now: %+v\n", data)
	vData := reflect.ValueOf(data)
	if vData.Kind() == reflect.Map {
		for _, key := range vData.MapKeys() {
			strct := vData.MapIndex(key)
			str := fmt.Sprintf("%v", key.Interface())
			incoming[str] = strct
			fmt.Println(key.Interface(), strct.Interface())
		}
	} else {
		return errors.New("unknown type")
	}

	indirect := reflect.ValueOf(out)
	v := indirect.Elem()

	s := v.Field(0)
	fmt.Printf("out type is now: %+v\n", s)
	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		values[i] = field.Interface()
		name := v.Type().Field(i).Name
		fmt.Printf("v.Field(%d) is now: %+v type: %s name: %s \n", i, field, field.Type(), name)
		if found, ok := incoming[name]; ok {
			kind := field.Type().Kind()
			if kind == reflect.Int {
				strValue := fmt.Sprintf("%v", found.Interface())
				intValue, err := strconv.Atoi(strValue)
				if err == nil {
					field.SetInt(int64(intValue))
				}
			} else if kind == reflect.String {
				strValue := fmt.Sprintf("%v", found.Interface())
				field.SetString(strValue)
			} else if kind == reflect.Bool {
				boolValue := found.Interface().(interface{}).(bool)
				field.SetBool(boolValue)
			} else {
				log.Println("unknown type:", kind)
			}

		}

	}
	return nil
}
