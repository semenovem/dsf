package schema

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const (
	boolType       = "bool"
	stringType     = "string"
	intType        = "int"
	float64Type    = "float64"
	uuidType       = "uuid"
	pathSeparator1 = "/"
	optionalPrefix = "?"
)

type Validator struct {
	nameSpace string
}

type response struct {
	errs           []error
	field          string
	optionalPrefix string
	pathSeparator  string
	parent         *response
}

func New() *Validator {
	return &Validator{}
}

func (r *Validator) Validate(scheme, response interface{}, path ...string) error {
	return newResponse().validate(scheme, response, path...)
}

func newResponse() *response {
	r := _newResponse()
	r.errs = make([]error, 0)
	return r
}

func _newResponse() *response {
	return &response{
		optionalPrefix: optionalPrefix,
		pathSeparator:  pathSeparator1,
	}
}

func (r *response) child(path string) *response {
	child := _newResponse()
	child.field = r.addPath(path)
	child.parent = r
	return child
}

func (r *response) addPath(additionalPath ...string) string {
	if len(additionalPath) != 0 {
		if r.field == "" {
			return strings.Join(additionalPath, r.pathSeparator)
		} else {
			a := make([]string, len(additionalPath)+1)
			a[0] = r.field
			copy(a[1:], additionalPath)

			return strings.Join(a, r.pathSeparator)
		}
	}

	return r.field
}

func (r *response) isOptional(s string) bool {
	return strings.HasPrefix(s, r.optionalPrefix)
}

func (r *response) validate(scheme, response interface{}, additionalPath ...string) error {
	r.field = r.addPath(additionalPath...)
	r.engine(scheme, response)
	return r.error()
}

func (r *response) addErr(err error) {
	if err == nil {
		return
	}

	if r.parent == nil {
		r.errs = append(r.errs, err)
	} else {
		r.parent.addErr(err)
	}
}

func (r *response) error() error {
	if len(r.errs) == 0 {
		return nil
	}

	ret := make([]string, len(r.errs))
	for i, v := range r.errs {
		ret[i] = v.Error()
	}

	return fmt.Errorf("%+v", strings.Join(ret, "\n"))
}

func (r *response) schemeTypeErr(b, c interface{}) {
	err := fmt.Errorf("поле=%s, ожидаемый тип=%T, факт=%T, значение=%+v", r.field, b, c, c)
	r.addErr(err)
}

func (r *response) schemeValueErr(b, c interface{}) {
	err := fmt.Errorf("поле=%s, ожидаемое значение='%v', факт='%v'", r.field, b, c)
	r.addErr(err)
}

func (r *response) uuidParser(v interface{}) (string, bool) {
	if u, ok := v.(uuid.UUID); ok {
		return u.String(), true
	}

	str, ok := v.(string)
	if !ok {
		err := schemeTypeErr(r.field, "", v)
		r.addErr(errors.WithMessage(err, "ожидается строковое представление uuid.UUID"))
		return "", false
	}

	id, err := uuid.Parse(str)
	if err != nil {
		err = errors.WithMessage(err,
			fmt.Sprintf("поле=%s, не валидный формат строкового представления uuid.UUID='%v'", r.field, str))

		r.addErr(err)
		return "", false
	}

	return id.String(), true
}

func (r *response) checkString(val interface{}) bool {
	_, ok := val.(string)
	if !ok {
		r.schemeTypeErr("", val)
	}
	return ok
}

func (r *response) engine(scheme, response interface{}) {

	switch t := scheme.(type) {

	case string:
		r.schemeString(t, response)

	case bool: // в ответе ждем bool значение, равное указанному
		if r.schemeString(boolType, response) && t != response {
			r.schemeValueErr(t, response)
		}

	case int: // в ответе ждем только int, равное указанному
		if r.schemeString(intType, response) && t != response {
			r.schemeValueErr(t, response)
		}

	case float64: // в ответе ждем только float, равное указанному
		if r.schemeString(float64Type, response) && t != response {
			r.schemeValueErr(t, response)
		}

	case uuid.UUID:
		if id, ok := r.uuidParser(response); ok && t.String() != id {
			r.schemeValueErr(t, response)
		}

	case *regexp.Regexp: // валидация строки регулярным выражением
		if r.checkString(response) && !t.MatchString(response.(string)) {
			err := fmt.Errorf("поле=%s, значение='%v' не соответствует регулярному выражению='%v'", r.field, response, t.String())
			r.addErr(err)
		}

	case map[string]interface{}:
		r.validateMap(t, response)

	case []interface{}:
		r.validateArray(t, response)

	default:
		panic(fmt.Errorf("поле=%s, тип не поддерживается во входящих данных=%T, значение='%+v'. Ожидаемый тип=%T, ожидаемое значение='%+v'",
			r.field, response, response, scheme, scheme))
	}
}

func (r *response) schemeString(scheme string, response interface{}) bool {
	ok := false

	switch scheme {

	case stringType:
		ok = r.checkString(response)

	case boolType:
		if _, ok = response.(bool); !ok {
			r.schemeTypeErr(true, response)
		}

	case intType:
		switch response.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			ok = true
		default:
			r.schemeTypeErr(0, response)
		}

	case float64Type:
		switch response.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			ok = true
		default:
			r.schemeTypeErr(1.1, response)
		}

	case uuidType:
		_, ok = r.uuidParser(response)

	// TODO делать time.Time
	//case "time":

	default:
		ok = r.checkString(response)
		if ok && scheme != response {
			r.schemeValueErr(scheme, response)
			ok = false
		}
	}

	return ok
}

func (r *response) validateMap(scheme map[string]interface{}, response interface{}) {
	var (
		ok              bool
		field, nextPath string
		resp            map[string]interface{}
		fact, expect    interface{}
	)

	if resp, ok = response.(map[string]interface{}); !ok {
		r.schemeTypeErr(scheme, response)
		return
	}

	for field, expect = range scheme {
		nextPath = fmt.Sprintf("%s/%s", r.field, field)

		fact, ok = resp[field]

		if !ok {
			if !r.isOptional(field) {
				r.addErr(fmt.Errorf("нет поля %s", nextPath))
			}
			continue
		}

		r.child(field).engine(expect, fact)
	}
}

func (r *response) validateArray(schema []interface{}, response interface{}) {
	var (
		ok           bool
		i            int
		resp         []interface{}
		fact, expect interface{}
	)

	if resp, ok = response.([]interface{}); !ok {
		r.schemeTypeErr(schema, response)
		return
	}

	if len(schema) == 0 {
		panic(fmt.Errorf("поле=%s, массив типа=%T в схеме не может быть пустым", r.field, schema))
	}

	if len(schema) == 1 {
		expect = schema[0]

		for i, fact = range resp {
			r.child(fmt.Sprintf("%s[%d]", r.field, i)).engine(expect, fact)
		}

		return
	}

	if len(schema) > 1 {
		for i, expect = range schema {
			if len(resp) <= i {
				r.addErr(fmt.Errorf("поле:%s, индекс[%d]", r.field, i))
			}

			r.child(fmt.Sprintf("%s, индекс[%d]", r.field, i)).engine(expect, resp[i])
		}
	}
}
