package utils

import (
	// "fmt"
	// "reflect"
	// "unsafe"

	"github.com/sanity-io/litter"
)

// func PrintContextInternals(ctx interface{}, inner bool) {
//     contextValues := reflect.ValueOf(ctx).Elem()
//     contextKeys := reflect.TypeOf(ctx).Elem()
//
//     if !inner {
//         fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
//     }
//
//     if contextKeys.Kind() == reflect.Struct {
//         for i := 0; i < contextValues.NumField(); i++ {
//             reflectValue := contextValues.Field(i)
//             reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()
//
//             reflectField := contextKeys.Field(i)
//
//             if reflectField.Name == "Context" {
//                 PrintContextInternals(reflectValue.Interface(), true)
//             } else {
//                 fmt.Printf("field name: %+v\n", reflectField.Name)
//                 fmt.Printf("value: %+v\n", reflectValue.Interface())
//             }
//         }
//     } else {
//         fmt.Printf("context is empty (int)\n")
//     }
// }

func PrettyPrint(i interface{}) {
	if IsDev() {
		litter.Dump(i)
	}
}

var IsDevMode *bool

func IsDev() bool {
	if IsDevMode == nil {
		return true
	}
	return *IsDevMode
}
