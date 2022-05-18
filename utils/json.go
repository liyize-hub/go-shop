package utils

import (
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

func setField(obj interface{}, name string, value interface{}) error {
	structData := reflect.ValueOf(obj).Elem()
	fieldValue := structData.FieldByName(name)
	if !fieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj ", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value ", name)
	}

	fieldType := fieldValue.Type()
	val := reflect.ValueOf(value)
	valTypeStr := val.Type().String()
	fieldTypeStr := fieldType.String()
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType) //Convert将v持有的值转换为类型为t的值，并返回该值的Value封装
	} else if fieldType != val.Type() { //json中不能表达的类型
		return errors.New("Provided value type " + valTypeStr + " didn't match obj field type " + fieldTypeStr)
	}
	fieldValue.Set(val) //将v的持有值修改为x的持有值。如果v.CanSet()返回假，会panic
	return nil
}

// SetStructByJSON 由json对象生成 struct
func SetStructByJSON(obj interface{}, mapData map[string]interface{}) {
	for key, value := range mapData {
		if err := setField(obj, key, value); err != nil {
			Logger.Error("SetStructByJSON failed", zap.Any("error", err))
		}
	}
}
