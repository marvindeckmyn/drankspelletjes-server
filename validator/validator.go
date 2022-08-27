package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
)

type V map[string]func(item interface{}) bool

func IsNotNil(values map[string]interface{}) error {
	nilValues := []string{}

	for key, value := range values {
		if value == nil || reflect.ValueOf(value).IsNil() {
			nilValues = append(nilValues, key)
		}
	}

	if len(nilValues) != 0 {
		return &ErrNil{Key: strings.Join(nilValues, ", ")}
	}

	return nil
}

func IsTime(item interface{}) bool {
	if item == nil {
		return false
	}

	layout := "2006-01-02T15:04:05.000Z"
	_, ok := time.Parse(layout, fmt.Sprintf("%v", item))
	return ok != nil
}

func IsOptTime(item interface{}) bool {
	return true
}

func IsOptUUID(item interface{}) bool {
	return true
}

func IsEmail(item interface{}) bool {
	if item == nil {
		return false
	}

	if str, ok := item.(string); ok {
		if str != "" {
			return true
		}
	}

	return false
}

func isUUID(item interface{}) bool {
	str, ok := item.(string)
	if !ok {
		return true
	}

	_, err := uuid.FromString(str)
	if err != nil {
		return true
	}

	return false
}

func IsUUIDSlice(item interface{}) bool {
	values, ok := item.([]string)
	if !ok {
		return true
	}

	for _, str := range values {
		_, err := uuid.FromString(str)
		if err != nil {
			return true
		}
	}

	return false
}

func IsBool(item interface{}) bool {
	if item == nil {
		return false
	}

	if _, ok := item.(bool); ok {
		return true
	}

	return false
}

func IsString(item interface{}) bool {
	if item == nil {
		return false
	}

	if _, ok := item.(string); ok {
		return true
	}

	return false
}

func IsOptString(item interface{}) bool {
	if item == nil {
		return true
	}

	if _, ok := item.(string); ok {
		return true
	}

	return false
}

func IsUint(item interface{}) bool {
	if item == nil {
		return false
	}

	switch t := item.(type) {
	case int64:
		if t < 0 {
			return false
		}

		return true

	case float32:
		if t < 0 {
			return false
		}

		return true

	case float64:
		if t < 0 {
			return false
		}

		return true

	case int:
		if t < 0 {
			return false
		}

		return true

	case string:
		int, err := strconv.ParseInt(t, 10, 64)
		if int < 0 {
			return false
		}

		return err == nil

	default:
		return false
	}
}

func IsInt(item interface{}) bool {
	if item == nil {
		return false
	}

	switch t := item.(type) {
	case int64:
		return true

	case float32:
		return true

	case float64:
		return true

	case int:
		return true

	case string:
		_, err := strconv.ParseInt(t, 10, 64)
		return err == nil

	default:
		return false
	}
}

func IsUUIDV4(item interface{}) bool {
	if item == nil {
		return false
	}

	if val, ok := item.(string); ok {
		r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
		return r.MatchString(val)
	} else {
		fmt.Print("Failed to parse string")
	}

	return false
}

func IsMapStrInterface(item interface{}) bool {
	if item == nil {
		return false
	}

	if _, ok := item.(map[string]interface{}); !ok {
		return false
	}

	return true
}

func IsMapStrStr(item interface{}) bool {
	if item == nil {
		return false
	}

	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return false
	}

	for _, value := range itemMap {
		if _, ok := value.(string); !ok {
			return false
		}
	}

	return true
}

func MarshalBody(requestBody io.Reader, val interface{}) error {
	body, _ := io.ReadAll(requestBody)
	jsonData := map[string]interface{}{}

	if err := json.Unmarshal(body, &jsonData); err != nil {
		return &ErrInvalidJSON{Cause: err}
	}

	encoded, _ := json.Marshal(jsonData)
	err := json.Unmarshal(encoded, &val)
	if err != nil {
		return &ErrInvalidJSON{Cause: err}
	}

	return nil
}

func (instance V) ValidateAndMarshalBody(requestBody io.Reader, val interface{}) error {
	jsonData := map[string]interface{}{}
	errs := []string{}

	if err := json.NewDecoder(requestBody).Decode(&jsonData); err != nil {
		if err != nil {
			return &ErrInvalidJSON{Cause: err}
		}
	}

	for key, validationFunc := range instance {
		valid := validationFunc(jsonData[key])
		if !valid {
			errs = append(errs, key)
		}
	}

	if len(errs) == 0 {
		encoded, _ := json.Marshal(jsonData)
		err := json.Unmarshal(encoded, &val)

		if err != nil {
			return &ErrInvalidJSON{Cause: err}
		}

		return nil
	}

	return &ErrInvalidContent{Cause: strings.Join(errs, ", ")}
}

func (instance V) ValidateAndMarshalURL(r *http.Request, val interface{}) error {
	errs := []string{}
	urlParams := map[string]string{}

	for key, validationFunc := range instance {
		val := r.URL.Query().Get(key)
		valid := validationFunc(val)
		if !valid {
			errs = append(errs, key)
		}

		urlParams[key] = val
	}

	if len(errs) == 0 {
		encoded, _ := json.Marshal(urlParams)
		err := json.Unmarshal(encoded, &val)

		if err != nil {
			return &ErrInvalidJSON{Cause: err}
		}

		return nil
	}

	return &ErrInvalidContent{Cause: strings.Join(errs, ", ")}
}
