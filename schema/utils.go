package schema

import "fmt"

func schemeTypeErr(a, b, c interface{}) error {
	return fmt.Errorf("поле=%s, ожидаемый тип=%T, факт=%T, значение=%+v", a, b, c, c)
}
