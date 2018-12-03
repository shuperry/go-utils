package utils

import (
	"reflect"
	"strings"
)

func Keys(m map[interface{}]interface{}) []interface{} {
	keys := []interface{}{}
	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}

func IsMapEmpty(m map[interface{}]interface{}) bool {
	return len(Keys(m)) == 0
}

func IsArrayEmpty(a []interface{}) bool {
	return len(a) == 0
}

func IsMapIncludeKey(m map[interface{}]interface{}, key interface{}) bool {
	flag := false

	for k, _ := range m {
		if k == key {
			flag = true
		}
	}

	return flag
}

func IsArrayInclude(a []interface{}, val interface{}) bool {
	flag := false

	for _, v := range a {
		if v == val {
			flag = true
		}
	}

	return flag
}

func GetStrWithSingleQuoteForArray(a []string) string {
	temp_a := make([]string, len(a))
	for i, v := range a {
		temp_a[i] = strings.Join([]string{"'", v, "'"}, "")
	}

	return strings.Join(temp_a, ",")
}

/**
  e.g: in: ["aaa,bbb", "ccc,ddd", "eee"], out: ["aaa", "bbb", "ccc", "ddd", "eee"].
*/
func TransArrayWithCommaToFullArray(a []string) []string {
	temp_a := []string{}

	for _, v := range a {
		temp_b := strings.Split(v, ",")
		for _, v1 := range temp_b {
			temp_a = append(temp_a, v1)
		}
	}

	return temp_a
}

func ArrayHasDuplicateItems(a []string, val string) bool {
	count := 0

	for _, v := range a {
		if v == val {
			count++
		}
	}

	return count > 1
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	for _, v1 := range a {
		if len(ret) == 0 {
			ret = append(ret, v1)
		} else {
			for k2, v2 := range ret {
				if v1 == v2 || v1 == "" {
					break
				} else if k2 == len(ret)-1 {
					ret = append(ret, v1)
				}
			}
		}
	}
	return
}

func IsMapKeyOutRange(m map[interface{}]interface{}, rangeKeysArr []interface{}) bool {
	flag := false

loop:
	for key, _ := range m {
		if !IsArrayInclude(rangeKeysArr, key) {
			flag = true
			break loop
		}
	}

	return flag
}

func Type(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func FilterOutRangeFields(m map[interface{}]interface{}, rangeKeysArr []interface{}) {
	for key, _ := range m {
		if !IsArrayInclude(rangeKeysArr, key) {
			delete(m, key)
		}
	}
}
