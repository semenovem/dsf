package schema

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator(t *testing.T) {

	t.Run("addErr", func(t *testing.T) {
		resp := newResponse()
		assert.NoError(t, resp.error())
		resp.addErr(fmt.Errorf(""))
		assert.Equal(t, 1, len(resp.errs))
	})

	t.Run("uuidParser", func(t *testing.T) {
		resp := newResponse()

		id, ok := resp.uuidParser("")
		assert.False(t, ok)
		assert.Equal(t, "", id)
		assert.Equal(t, 1, len(resp.errs))

		id2, err := uuid.NewUUID()
		assert.NoError(t, err)

		id, ok = resp.uuidParser(id2)
		assert.True(t, ok)
		assert.Equal(t, id2.String(), id)
		assert.Equal(t, 1, len(resp.errs))

		id, ok = resp.uuidParser(id2.String())
		assert.True(t, ok)
		assert.Equal(t, id2.String(), id)
		assert.Equal(t, 1, len(resp.errs))
	})

	t.Run("checkString", func(t *testing.T) {
		resp := newResponse()

		t.Run("ok", func(t *testing.T) {
			ok := resp.checkString("")
			assert.True(t, ok)
			assert.Equal(t, 0, len(resp.errs))
		})

		t.Run("error", func(t *testing.T) {
			set := []interface{}{
				1, 1.1, map[string]string{}, []string{}, false, true,
			}

			for i, v := range set {
				ok := resp.checkString(v)
				assert.False(t, ok)
				assert.Equal(t, i+1, len(resp.errs))
			}
		})
	})

	t.Run("addPath", func(t *testing.T) {
		resp := newResponse()

		path := "sdfasdf234123"
		path2 := "sdfasdf234123"

		assert.Equal(t, "", resp.addPath())
		assert.Equal(t, path, resp.addPath(path))

		expect := fmt.Sprintf("%s%s%s", path, pathSeparator1, path2)
		assert.Equal(t, expect, resp.addPath(path, path2))
		field := "3432234"

		resp.field = field
		assert.Equal(t, field, resp.addPath())

		expect = fmt.Sprintf("%s%s%s", field, pathSeparator1, path)
		assert.Equal(t, expect, resp.addPath(path))

		expect = fmt.Sprintf("%s%s%s%s%s", field, pathSeparator1, path, pathSeparator1, path2)
		assert.Equal(t, expect, resp.addPath(path, path2))
	})

	t.Run("child", func(t *testing.T) {
		resp := newResponse()
		path := "sdfasdf234123"
		clone := resp.child(path)
		assert.Equal(t, path, clone.field)
	})

	t.Run("isOptional", func(t *testing.T) {
		r := newResponse()
		assert.True(t, r.isOptional("?sfsfsadf"))
		assert.False(t, r.isOptional("sfsfsadf"))
	})

	t.Run("validateMap", func(t *testing.T) {
		r := newResponse()

		m := map[string]interface{}{
			"first":   1,
			"?second": 1,
		}

		set := []interface{}{
			1, 1.1, map[string]string{}, []string{}, false, true,
		}
		for i, v := range set {
			r.validateMap(m, v)
			assert.Equal(t, i+1, len(r.errs))
		}

		m2 := map[string]interface{}{
			"first": 1,
		}

		r = newResponse()

		r.validateMap(m, m2)
		assert.NoError(t, r.error())
	})

}

func TestValidatorSchemeString(t *testing.T) {

	t.Run("string", func(t *testing.T) {
		sets := []string{"", "sdfsdfsdaf"}

		for _, v := range sets {
			r := newResponse()
			assert.True(t, r.schemeString(stringType, v))
			assert.NoError(t, r.error())
		}
	})

	t.Run("bool", func(t *testing.T) {
		sets := []bool{true, false}

		for _, v := range sets {
			r := newResponse()
			assert.True(t, r.schemeString(boolType, v))
			assert.NoError(t, r.error())
		}
	})

	t.Run("int", func(t *testing.T) {
		sets := []interface{}{
			uint(1), uint8(1), uint16(1), uint32(1), uint64(1), 1, int8(1), int16(1), int32(1), int64(1),
		}

		for _, v := range sets {
			r := newResponse()
			assert.True(t, r.schemeString(intType, v))
			assert.NoError(t, r.error())
		}
	})

	t.Run("float64", func(t *testing.T) {
		sets := []interface{}{
			uint(1), uint8(1), uint16(1), uint32(1), uint64(1), 1, int8(1), int16(1), int32(1), int64(1), float32(1.1), 1.1,
		}

		for _, v := range sets {
			r := newResponse()
			assert.True(t, r.schemeString(float64Type, v))
			assert.NoError(t, r.error())
		}
	})

	t.Run("uuid", func(t *testing.T) {
		sets := []interface{}{
			uuid.New(), uuid.NewString(),
		}

		for _, v := range sets {
			r := newResponse()
			assert.True(t, r.schemeString(uuidType, v))
			assert.NoError(t, r.error())
		}
	})

	t.Run("string-exactly", func(t *testing.T) {
		sets := []string{
			"", "0000", "12121212", "asdfasdfdsf",
		}

		for _, v := range sets {
			r := newResponse()
			assert.True(t, r.schemeString(v, v))
			assert.NoError(t, r.error())
		}
	})
}

func TestValidatorEngine(t *testing.T) {
	engine := func(a, b interface{}) error {
		r := newResponse()
		r.engine(a, b)
		return r.error()
	}

	t.Run("string", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			assert.NoError(t, engine("string", ""))
			assert.NoError(t, engine("string", "sdfsdfs"))
			assert.NoError(t, engine("string", "1111"))
		})

		t.Run("error", func(t *testing.T) {
			assert.Error(t, engine("string", 1))
			assert.Error(t, engine("string", 1.1))
			assert.Error(t, engine("string", []string{}))
			assert.Error(t, engine("string", true))
			assert.Error(t, engine("string", false))
		})
	})
}

func TestValidatorArray(t *testing.T) {
	validateArray := func(a []interface{}, b interface{}) error {
		r := newResponse()
		r.validateArray(a, b)
		return r.error()
	}

	t.Run("string1", func(t *testing.T) {
		schema := []interface{}{"string"}
		input := []interface{}{"sdfsdfg", "sdfsdfgdf"}

		assert.NoError(t, validateArray(schema, input))
	})

	t.Run("string2", func(t *testing.T) {
		schema := []interface{}{"asdfsadfsdf"}
		input := []interface{}{"asdfsadfsdf", "asdfsadfsdf"}

		assert.NoError(t, validateArray(schema, input))
	})

	t.Run("string3", func(t *testing.T) {
		schema := []interface{}{"11111111", "222222222"}
		input := []interface{}{"11111111", "222222222"}

		assert.NoError(t, validateArray(schema, input))
	})

	//schema := []map[string]interface{}{
	//  {
	//    "org": map[string]interface{}{
	//      "onpremise_enabled":   "bool",
	//      "record_enabled":      "bool",
	//      "change_name_enabled": "bool",
	//    },
	//
	//    "conf_cloud": map[string]interface{}{
	//      "guest_permission_enabled": "bool",
	//      "trusted_groups_enabled":   "bool",
	//      "auth_required":            "bool",
	//    },
	//
	//    "conf_onprem": map[string]interface{}{
	//      "guest_permission_enabled": "bool",
	//      "trusted_groups_enabled":   "bool",
	//      "auth_required":            "bool",
	//    },
	//  },
	//}
	//
	//response := map[string]interface{}{
	//  "org": map[string]interface{}{
	//    "onpremise_enabled":   false,
	//    "record_enabled":      true,
	//    "change_name_enabled": true,
	//  },
	//}

	schema := []interface{}{"string"}
	input := []interface{}{"sdfsdfg", "sdfsdfgdf"}

	t.Run("string", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {

			assert.NoError(t, validateArray(schema, input))

			//assert.NoError(t, engine("string", "sdfsdfs"))
			//assert.NoError(t, engine("string", "1111"))
		})

		//t.Run("error", func(t *testing.T) {
		//	assert.Error(t, engine("string", 1))
		//	assert.Error(t, engine("string", 1.1))
		//	assert.Error(t, engine("string", []string{}))
		//	assert.Error(t, engine("string", true))
		//	assert.Error(t, engine("string", false))
		//})
	})
}
