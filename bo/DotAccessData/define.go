package DotAccessData

import (
	"encoding/json"
	"github.com/lamgor666/goboot-common/util/castx"
	"github.com/lamgor666/goboot-common/util/stringx"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

type Bo struct {
	data map[string]interface{}
}

func FromMap(arg0 interface{}) Bo {
	if arg0 == nil {
		return DotAccessData{}
	}
	
	map1 := castx.ToStringMap(arg0)
	
	if len(map1) < 1 {
		return Bo{}
	}

	return Bo{data: map1}
}

func FromJson(arg0 interface{}) Bo {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok && len(_buf) > 0 {
		buf = _buf
	} else if reader, ok := arg0.(io.Reader); ok {
		buf, _ = ioutil.ReadAll(reader)
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		buf = []byte(s1)
	}

	var data map[string]interface{}

	if len(buf) < 1 || json.Unmarshal(buf, &data) != nil {
		return Bo{}
	}

	return Bo{data: data}
}

func FromYaml(arg0 interface{}) Bo {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok && len(_buf) > 0 {
		buf = _buf
	} else if reader, ok := arg0.(io.Reader); ok {
		buf, _ = ioutil.ReadAll(reader)
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		if stat, err := os.Stat(s1); err == nil {
			if !stat.IsDir() {
				buf, _ = ioutil.ReadFile(s1)
			}
		} else {
			buf = []byte(s1)
		}
	}

	var data map[string]interface{}

	if len(buf) < 1 || yaml.Unmarshal(buf, &data) != nil {
		return Bo{}
	}

	return Bo{data: data}
}

func (bo Bo) IsEmpty() bool {
	return len(bo.data) < 1
}

func (bo Bo) GetMap(path string) map[string]interface{} {
	if bo.IsEmpty() {
		return map[string]interface{}{}
	}

	if !strings.Contains(path, ".") {
		return castx.ToStringMap(bo.getValueInternal(path))
	}

	lastKey := stringx.SubstringAfterLast(path, ".")
	keys := strings.Split(stringx.SubstringBeforeLast(path, "."), ".")
	var map1 map[string]interface{}

	for idx, key := range keys {
		if idx == 0 {
			map1 = castx.ToStringMap(bo.getValueInternal(key))
			continue
		}

		if len(map1) < 1 {
			break
		}

		map1 = castx.ToStringMap(bo.getValueInternal(key, map1))
	}

	if len(map1) < 1 {
		return map[string]interface{}{}
	}

	return castx.ToStringMap(bo.getValueInternal(lastKey, map1))
}

func (bo Bo) GetStringMap(path string) map[string]string {
	return castx.ToStringMapString(bo.GetMap(path))
}

func (bo Bo) GetSlice(path string) []interface{} {
	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return []interface{}{}
		}

		return castx.ToSlice(bo.getValueInternal(path))
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToSlice(bo.getValueInternal(key, map1))
}

func (bo Bo) GetStringSlice(path string) []string {
	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return []string{}
		}

		return castx.ToStringSlice(bo.getValueInternal(path))
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToStringSlice(bo.getValueInternal(key, map1))
}

func (bo Bo) GetIntSlice(path string) []int {
	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return []int{}
		}

		return castx.ToIntSlice(bo.getValueInternal(path))
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToIntSlice(bo.getValueInternal(key, map1))
}

func (bo Bo) GetMapSlice(path string) []map[string]interface{} {
	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return []map[string]interface{}{}
		}

		return castx.ToMapSlice(bo.getValueInternal(path))
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToMapSlice(bo.getValueInternal(key, map1))
}

func (bo Bo) GetString(path string, defaultValue ...string) string {
	var _defaultValue string

	if len(defaultValue) > 0 && defaultValue[0] != "" {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return _defaultValue
		}

		if s1, err := castx.ToStringE(bo.getValueInternal(path)); err == nil && s1 != "" {
			return s1
		}

		return _defaultValue
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if s1, err := castx.ToStringE(bo.getValueInternal(key, map1)); err == nil && s1 != "" {
		return s1
	}

	return _defaultValue
}

func (bo Bo) GetInt(path string, defaultValue ...int) int {
	_defaultValue := math.MinInt32

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return _defaultValue
		}

		if n1, err := castx.ToIntE(bo.getValueInternal(path)); err == nil {
			return n1
		}

		return _defaultValue
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if n1, err := castx.ToIntE(bo.getValueInternal(key, map1)); err == nil {
		return n1
	}

	return _defaultValue
}

func (bo Bo) GetInt64(path string, defaultValue ...int64) int64 {
	_defaultValue := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return _defaultValue
		}

		if n1, err := castx.ToInt64E(bo.getValueInternal(path)); err == nil {
			return n1
		}

		return _defaultValue
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if n1, err := castx.ToInt64E(bo.getValueInternal(key, map1)); err == nil {
		return n1
	}

	return _defaultValue
}

func (bo Bo) GetFloat(path string, defaultValue ...float64) float64 {
	_defaultValue := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return _defaultValue
		}

		if n1, err := castx.ToFloat64E(bo.getValueInternal(path)); err == nil {
			return n1
		}

		return _defaultValue
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if n1, err := castx.ToFloat64E(bo.getValueInternal(key, map1)); err == nil {
		return n1
	}

	return _defaultValue
}

func (bo Bo) GetBoolean(path string, defaultValue ...bool) bool {
	var _defaultValue bool

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if bo.IsEmpty() {
			return _defaultValue
		}

		if b1, err := castx.ToBoolE(bo.getValueInternal(path)); err == nil {
			return b1
		}

		return _defaultValue
	}

	map1 := bo.GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if b1, err := castx.ToBoolE(bo.getValueInternal(key, map1)); err == nil {
		return b1
	}

	return _defaultValue
}

func (bo Bo) GetDataSize(path string) int64 {
	return castx.ToDataSize(bo.GetString(path))
}

func (bo Bo) GetDuration(path string) time.Duration {
	return castx.ToDuration(bo.GetString(path))
}

func (bo Bo) getValueInternal(key string, source ...map[string]interface{}) interface{} {
	var data map[string]interface{}

	if len(source) > 0 {
		data = source[0]
	} else {
		data = bo.data
	}

	if len(data) < 1 {
		return nil
	}

	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ToLower(key)

	for compkey, value := range data {
		compkey = strings.ReplaceAll(compkey, "-", "")
		compkey = strings.ReplaceAll(compkey, "_", "")
		compkey = strings.ToLower(compkey)

		if compkey == key {
			return value
		}
	}

	return nil
}
