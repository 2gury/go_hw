package main

import (
	"log"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	structVal := reflect.ValueOf(out)
	mapVal := reflect.ValueOf(data)

	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}

	if mapVal.Kind() == reflect.Ptr {
		mapVal = mapVal.Elem()
	}
	if mapVal.Kind() != reflect.Map {
		log.Println("not map")
		return nil
	}

	for i := 0; i < structVal.NumField(); i ++ {
		curFieldName := structVal.Type().Field(i).Name
		field := structVal.Field(i)
		val :=  mapVal.MapIndex(reflect.ValueOf(curFieldName))

		if val.IsValid() {
			log.Printf("%s: %v", curFieldName, val)
			if field.IsValid() && field.CanSet() {
				log.Println("sheesh")
			}
		}
	}

	return nil
}

type Kek struct {
	ID       int
	Username string
	Active   bool
}

func main() {
	data := map[string]interface{}{
		"Lol":  1,
		"fefe": "pepe",
		"ID": "going",
	}

	out := new(Kek)
	i2s(&data, out)
}
