package types

import (
	"reflect"
	"strconv"
)

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func UInt64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func StringToInt(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// Fill 通过反射将对象2的值填充给对象1
func Fill(obj1 interface{}, obj2 interface{}) {
	//1.通过反射获取两个结构的字段
	v1 := reflect.ValueOf(obj1).Elem()
	v2 := reflect.ValueOf(obj2).Elem()

	//2.循环填充
	for i := 0; i < v1.NumField(); i++ {
		//2.1获取结构1字段详细信息
		fieldInfo1 := v1.Type().Field(i)
		field1Name := fieldInfo1.Name
		field1Type := fieldInfo1.Type

		//2.2 循环结构2的字段
		for i2 := 0; i2 < v2.NumField(); i2++ {
			//2.2.1获取解构2的详细信息
			fieldInfo2 := v2.Type().Field(i2)
			field2Name := fieldInfo2.Name
			field2Type := fieldInfo2.Type

			//2.2.2如果两个结构的字段名相等，而且值类型相等且有值，将结构2的值赋给结构1,
			if field1Name == field2Name && field1Type == field2Type {

				//2.2.2.1 判断是否有值
				if v2.FieldByName(fieldInfo2.Name).IsValid() {
					switch v2.FieldByName(fieldInfo2.Name).Type().String() {
					case "int":
						if v2.FieldByName(fieldInfo2.Name).Int() == 0 {
							continue
						}
					case "string":
						if v2.FieldByName(fieldInfo2.Name).String() == "" {
							continue
						}
					}
				}

				//2.2.2.1 设置值
				newValue := v2.FieldByName(field2Name)
				if newValue.IsValid(){
					v1.FieldByName(field1Name).Set(newValue)
				}
			}
		}
	}
}
