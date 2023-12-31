package common

import (
	"encoding/xml"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"encoding/json"
	"errors"
)

type (
	Binding interface {
		Bind(*http.Request, interface{}) error
	}

	// JSON binding
	jsonBinding struct{}

	// XML binding
	xmlBinding struct{}

	// form binding
	formBinding struct{}

	// multipart form binding
	multipartFormBinding struct{}
)

const MAX_MEMORY = 1 * 1024 * 1024

var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{} // todo
	MultipartForm = multipartFormBinding{}
)

func (_ jsonBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err == nil {
		return Validate(obj)
	} else {
		log.Printf("errrrrrrrr %v", err)
		return err
	}
}

func (_ xmlBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err == nil {
		return Validate(obj)
	} else {
		return err
	}
}

func (_ formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return Validate(obj)
}

func (_ multipartFormBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(MAX_MEMORY); err != nil {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return Validate(obj)
}

func mapForm(ptr interface{}, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	formStruct := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := formStruct.Field(i)
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get("form")
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			// requirement from @baixiao 2016.11.1 by @liudan
			if structFieldKind == reflect.Struct {
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil {
					return err
				}
			}
			continue
		}

		inputValue, exists := form[inputFieldName]
		if !exists {
			continue
		}
		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
			formStruct.Field(i).Set(slice)
		} else {
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
				return err
			}
		}
	}
	return nil
}

func setIntField(val string, bitSize int, structField reflect.Value) error {
	if val == "" {
		val = "0"
	}

	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		structField.SetInt(intVal)
	}

	return err
}

func setUintField(val string, bitSize int, structField reflect.Value) error {
	if val == "" {
		val = "0"
	}

	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		structField.SetUint(uintVal)
	}

	return err
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		} else {
			structField.SetBool(boolVal)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
	return nil
}

// Don't pass in pointers to bind to. Can lead to bugs. See:
// https://github.com/codegangsta/martini-contrib/issues/40
// https://github.com/codegangsta/martini-contrib/pull/34#issuecomment-29683659
func ensureNotPointer(obj interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		panic("Pointers are not accepted as binding models")
	}
}

func Validate(obj interface{}, parents ...string) error {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			// Allow ignored and unexported fields in the struct
			if len(field.PkgPath) > 0 || field.Tag.Get("form") == "-" {
				continue
			}

			fieldValue := val.Field(i).Interface()
			zero := reflect.Zero(field.Type).Interface()

			if strings.Index(field.Tag.Get("binding"), "required") > -1 {
				fieldType := field.Type.Kind()
				if fieldType == reflect.Struct {
					if reflect.DeepEqual(zero, fieldValue) {
						return errors.New("Required " + field.Name)
					}
					err := Validate(fieldValue, field.Name)
					if err != nil {
						return err
					}
				} else if reflect.DeepEqual(zero, fieldValue) {
					if len(parents) > 0 {
						return errors.New("Required " + field.Name + " on " + parents[0])
					} else {
						return errors.New("Required " + field.Name)
					}
				} else if fieldType == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
					err := Validate(fieldValue)
					if err != nil {
						return err
					}
				}
			} else {
				fieldType := field.Type.Kind()
				if fieldType == reflect.Struct {
					if reflect.DeepEqual(zero, fieldValue) {
						continue
					}
					err := Validate(fieldValue, field.Name)
					if err != nil {
						return err
					}
				} else if fieldType == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
					err := Validate(fieldValue, field.Name)
					if err != nil {
						return err
					}
				}
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fieldValue := val.Index(i).Interface()
			err := Validate(fieldValue)
			if err != nil {
				return err
			}
		}
	default:
		return nil
	}
	return nil
}
