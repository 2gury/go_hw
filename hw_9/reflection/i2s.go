package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	structVal := reflect.ValueOf(out)
	mapVal := reflect.ValueOf(data)

	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	} else {
		return fmt.Errorf("bruh")
	}

	if mapVal.Kind() == reflect.Ptr {
		mapVal = mapVal.Elem()
	}

	return setMapValue(structVal, mapVal)
}

func setMapValue(strcVal, mapVal reflect.Value) (err error) {
	if mapVal.IsValid() && strcVal.IsValid() {
		kindStrcVal := strcVal.Kind()
		kindMapVal := mapVal.Kind()
		if kindMapVal == reflect.Interface {
			kindMapVal = mapVal.Elem().Kind()
		}
		log.Println(kindStrcVal, kindMapVal)
		isEqual := kindStrcVal == kindMapVal
		switch {
		case isEqual && kindStrcVal == reflect.String:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			strcVal.SetString(mapVal.String())
		case isEqual && kindStrcVal == reflect.Int:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			strcVal.SetInt(mapVal.Int())
		case kindStrcVal == reflect.Int && kindMapVal == reflect.Float64:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			intMapVal := int(mapVal.Float())
			strcVal.SetInt(int64(intMapVal))
		case isEqual && kindStrcVal == reflect.Bool:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			strcVal.SetBool(mapVal.Bool())
		case isEqual && kindStrcVal == reflect.Slice:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			strcVal.Set(reflect.MakeSlice(reflect.SliceOf(strcVal.Type().Elem()), mapVal.Len(), mapVal.Len()))
			for i := 0; i < mapVal.Len(); i++ {
				v := mapVal.Index(i)
				// log.Printf("%v", v.Elem().Kind())
				curErr := setMapValue(strcVal.Index(i), v)
				if curErr != nil {
					err = curErr
				}
			}
		case isEqual && kindStrcVal == reflect.Struct:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			for i := 0; i < strcVal.NumField(); i++ {
				curErr := setMapValue(strcVal.Field(i), mapVal.FieldByName(strcVal.Type().Field(i).Name))
				if curErr != nil {
					err = curErr
				}
			}
		case kindStrcVal == reflect.Struct && kindMapVal == reflect.Map:
			if mapVal.Kind() == reflect.Interface {
				mapVal = mapVal.Elem()
			}
			for i := 0; i < strcVal.NumField(); i++ {
				curFieldName := strcVal.Type().Field(i).Name
				field := strcVal.Field(i)
				currMapVal := mapVal.MapIndex(reflect.ValueOf(curFieldName))

				curErr := setMapValue(field, currMapVal)
				if curErr != nil {
					err = curErr
				}
			}
		default:
			err = fmt.Errorf("bruh")
			log.Printf("unknown type: %s, %s, %v", kindStrcVal.String(), kindMapVal.String(), err)
		}
	}

	return err
}

type SimpleStruct struct {
	ID string
}

type MediumStruct struct {
	Val string
}

type ComplexStruct struct {
	SimpStruct SimpleStruct
	MedStruct  []MediumStruct
}

func main() {
	smpl := SimpleStruct{
		ID: "fefe",
	}
	expected := &ComplexStruct{
		SimpStruct: smpl,
		MedStruct: []MediumStruct{
			{
				"lulz",
			},
			{
				"going",
			},
		},
	}

	jsonRaw, _ := json.Marshal(expected)
	log.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := new(ComplexStruct)
	i2s(tmpData, result)

	log.Println(result)
}
