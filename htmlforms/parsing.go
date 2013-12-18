package htmlforms

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/astaxie/beego"

	"github.com/francoishill/goutils/stringutils"
)

//Some obtained from: https://github.com/beego/wetalk

func ParseForm(form interface{}, values url.Values, dateFormat string, dateTimeFormat string) error {
	val := reflect.ValueOf(form)
	elm := reflect.Indirect(val)

	if !isStructPointer(val) {
		return fmt.Errorf("%s must be a struct pointer", val.Type().Name())
	}

outFor:
	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		fT := elm.Type().Field(i)

		fName := fT.Name

		for _, v := range strings.Split(fT.Tag.Get("form"), ";") {
			v = strings.TrimSpace(v)
			if v == "-" {
				continue outFor
			} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
				tN := v[:i]
				v = strings.TrimSpace(v[i+1 : len(v)-1])
				switch tN {
				case "name":
					fName = v
				}
			}
		}

		value := ""
		var vs []string
		if v, ok := values[fName]; ok {
			vs = v
			if len(v) > 0 {
				value = v[0]
			}
		}
		if v, ok := values[fName+"[]"]; ok {
			vs = v
			if len(v) > 0 {
				value = v[0]
			}
		}

		switch fT.Type.Kind() {
		case reflect.Bool:
			b, _ := stringutils.StrTo(value).Bool()
			f.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, _ := stringutils.StrTo(value).Int64()
			f.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, _ := stringutils.StrTo(value).Uint64()
			f.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, _ := stringutils.StrTo(value).Float64()
			f.SetFloat(x)
		case reflect.Struct:
			if fT.Type.String() == "time.Time" {
				if len(value) > 10 {
					t, err := beego.DateParse(value, dateTimeFormat)
					if err != nil {
						continue
					}
					f.Set(reflect.ValueOf(t))
				} else {
					t, err := beego.DateParse(value, dateFormat)
					if err != nil {
						continue
					}
					f.Set(reflect.ValueOf(t))
				}
			}
		case reflect.String:
			f.SetString(value)
		case reflect.Slice:
			f.Set(reflect.ValueOf(vs))
		}
	}

	return nil
}

// assert an object must be a struct pointer
func isStructPointer(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		return true
	}
	return false
}
