package options

import (
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

func resolve(obj interface{}, hget HandleGet) (err error) {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := reflect.Indirect(val).FieldByName(field.Name)

		switch fieldVal.Kind() {
		case reflect.Struct:
			err = resolve(fieldVal.Addr().Interface(), hget)
		case reflect.Ptr:
			if !fieldVal.IsNil() {
				err = resolve(filedVal, hget)
			}
		default:
			err = assign(field, fieldVal, hget)
		}

		if err != nil {
			return
		}
	}

	return nil
}

func assign(field reflect.StructField, val reflect.Value, hget HandleGet) (err error) {
	key := field.Tag.Get("options")
	if len(key) == 0 {
		return
	}

	value := hget(key, DefaultVal)
	if value == DefaultVal {
		return
	}

	switch val.Kind() {
	case reflect.String:
		val.SetString(value)
	case reflect.Slice, reflect.Array:
		val2 := reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr()))
		if err := json.Unmarshal([]byte(value), val2.Interface()); err != nil {
			return err
		}
	case reflect.Map:
		val2 := reflect.New(val.Type())
		if err := json.Unmarshal([]byte(value), val2.Interface()); err != nil {
			return err
		}

		val.Set(reflect.Indirect(val2))
	default:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		val.SetInt(v)
	}

	return nil
}
