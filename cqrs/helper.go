package cqrs

import "reflect"

func extractName(i interface{}) (name string) {
    name = reflect.TypeOf(i).String()

    if reflect.ValueOf(i).Kind() == reflect.Ptr{
        name = name[1:]
    }

    return
}
