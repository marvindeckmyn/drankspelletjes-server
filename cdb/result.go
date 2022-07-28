package cdb

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/marvindeckmyn/drankspelletjes-server/log"
	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
)

// DbResult is a map[string]interface{} with convenient data selection functions.
type CdbResult struct {
	// The raw db result set.
	data map[string]interface{}

	// Slice which contains
	errs []error
}

func (rs *CdbResult) GetData() map[string]interface{} {
	return rs.data
}

// addError appends an error to a given CdbResult, the error is then passed on as the return
// value.
func (rs *CdbResult) addError(err error) error {
	if rs == nil {
		return &ErrMissingResult{}
	}

	rs.errs = append(rs.errs, err)
	return err
}

// ClearErrors will remove any pending errors from a CdbResult.
func (rs *CdbResult) ClearErrors() {
	if rs == nil {
		return
	}

	rs.errs = nil
}

// HasErrors returns true when one of the select functions caused an error.
func (rs *CdbResult) HasErrors() bool {
	if rs == nil {
		return true
	}

	return len(rs.errs) > 0
}

// HasErrorsLog returns true when one of the select functions cased an error. In case there are
// errors they will be logged using the log package with the given tag and prefix.
func (rs *CdbResult) HasErrorsLog(tag string, prefix string) bool {
	if rs == nil {
		return true
	}

	if !rs.HasErrors() {
		return false
	}

	for _, err := range rs.Errors() {
		log.Error("%s%s", prefix, err.Error())
	}

	return true
}

// Get al list of all errors which have occurred while selecting data from the CdbResult.
func (rs *CdbResult) Errors() []error {
	if rs == nil {
		return []error{&ErrMissingResult{}}
	}

	return rs.errs
}

// time tries to select a key as a time object. An error will occur when the key doesn't exist.
func (rs *CdbResult) Time(key string, dest **time.Time) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	success := false

	layout := "2006-01-02 15:04:05.000 -0700"
	t, err := time.Parse(layout, fmt.Sprintf("%v", raw))
	if err == nil {
		success = true
	}

	layout = "2006-01-02 15:04:05 -0700 MST"
	t, err = time.Parse(layout, fmt.Sprintf("%v", raw))
	if err == nil {
		success = true
	}

	if !success {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to time.Time", key, raw)})
	}

	*dest = &t
	return nil
}

// Opttime tries to select a key as a time object.
func (rs *CdbResult) OptTime(key string, dest **time.Time) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	if raw == nil {
		return nil
	}

	success := false

	layout := "2006-01-02 15:04:05.000 -0700"
	t, err := time.Parse(layout, fmt.Sprintf("%v", raw))
	if err == nil {
		success = true
	}

	layout = "2006-01-02 15:04:05 -0700 MST"
	t, err = time.Parse(layout, fmt.Sprintf("%v", raw))
	if err == nil {
		success = true
	}

	if !success {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to opt time.Time", key, raw)})
	}

	*dest = &t
	return nil
}

// Iface tries to select a key as an interface. An error will occur when the key doesn't exist.
func (rs *CdbResult) Iface(key string, dest **interface{}) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	val, present := rs.data[key]
	if !present {
		*dest = nil
		return rs.addError(&ErrNoSuchKey{key})
	}

	*dest = &val
	return nil
}

// OptIface tries to select a key as an interface but doesn't error out when the key doesn't exist.
func (rs *CdbResult) OptIface(key string, dest **interface{}) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	val, present := rs.data[key]
	if present {
		*dest = &val
	}

	return nil
}

// Uint8 tries to select a key as an uint8. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an uint8.
func (rs *CdbResult) Uint8(key string, dest **uint8) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseUint(fmt.Sprintf("%v", raw), 10, 8)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int8", key, raw)})
	}

	num8 := uint8(num)
	*dest = &num8
	return nil
}

// Uint16 tries to select a key as an uint16. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an uint16.
func (rs *CdbResult) Uint16(key string, dest **uint16) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseUint(fmt.Sprintf("%v", raw), 10, 16)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int16", key, raw)})
	}

	num16 := uint16(num)
	*dest = &num16
	return nil
}

// Uint32 tries to select a key as an uint32. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an uint32.
func (rs *CdbResult) Uint32(key string, dest **uint32) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseUint(fmt.Sprintf("%v", raw), 10, 32)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int32", key, raw)})
	}

	num32 := uint32(num)
	*dest = &num32
	return nil
}

// OptUint32 tries to select a key as an uint32. An error will occur when the value cannot be parsed
// to an uint32.
func (rs *CdbResult) OptUint32(key string, dest **uint32) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	num, err := strconv.ParseUint(fmt.Sprintf("%v", raw), 10, 32)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int32", key, raw)})
	}

	num32 := uint32(num)
	*dest = &num32
	return nil
}

// OptUint64 tries to select a key as an uint64. An error will occur when the value cannot be parsed
// to an uint64.
func (rs *CdbResult) OptUint64(key string, dest **uint64) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	num, err := strconv.ParseUint(fmt.Sprintf("%v", raw), 10, 64)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int64", key, raw)})
	}

	num64 := uint64(num)
	*dest = &num64
	return nil
}

// Int32 tries to select a key as an int32. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an int32.
func (rs *CdbResult) Int32(key string, dest **int32) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseInt(fmt.Sprintf("%v", raw), 10, 32)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s %v cannot be parsed to int32", key, raw)})
	}

	num32 := int32(num)
	*dest = &num32
	return nil
}

// OptInt32 tries to select a key as an int32. An error will occur when the value cannot be parsed
// to an int32.
func (rs *CdbResult) OptInt32(key string, dest **int32) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	if raw == nil {
		return nil
	}

	num, err := strconv.ParseInt(fmt.Sprintf("%v", raw), 10, 32)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int32", key, raw)})
	}

	num32 := int32(num)
	*dest = &num32
	return nil
}

// Uint64 tries to select a key as an uint64. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an uint64.
func (rs *CdbResult) Uint64(key string, dest **uint64) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseUint(fmt.Sprintf("%v", raw), 10, 64)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int64", key, raw)})
	}

	num64 := uint64(num)
	*dest = &num64
	return nil
}

// Int64 tries to select a key as an int64. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an int64.
func (rs *CdbResult) Int64(key string, dest **int64) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseInt(fmt.Sprintf("%v", raw), 10, 64)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int64", key, raw)})
	}

	num64 := int64(num)
	*dest = &num64
	return nil
}

// OptInt64 tries to select a key as an int64. An error will occur when the value cannot be parsed
// to an int64.
func (rs *CdbResult) OptInt64(key string, dest **int64) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	num, err := strconv.ParseInt(fmt.Sprintf("%v", raw), 10, 64)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to int64", key, raw)})
	}

	*dest = &num
	return nil
}

// Str tries to select a key as a string. An error will occur when the key doesn't exist or when the
// value cannot be parsed to a string.
func (rs *CdbResult) Str(key string, dest **string) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	str, ok := raw.(string)
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to string", key, raw)})
	}

	*dest = &str
	return nil
}

// OptStr tries to select a key as a string. An error will occur when the value cannot be parsed to
// a string.
func (rs *CdbResult) OptStr(key string, dest **string) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present || raw == nil {
		*dest = nil
		return nil
	}

	if raw == nil {
		return nil
	}

	str, ok := raw.(string)
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%s %v cannot be parsed to string", key, raw)})
	}

	*dest = &str
	return nil
}

// StrSlice gets an array of strings. An error will occur when the key doesn't exist or when the
// value could not be parsed to a slice of strings.
func (rs *CdbResult) StrSlice(key string, dest **[]string) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	slice, ok := raw.([]interface{})
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to []string", raw)})
	}

	ss := []string{}

	for idx, entry := range slice {
		rawStr, ok := entry.(string)
		if !ok {
			return rs.addError(&ErrMalformed{fmt.Errorf("%v on idx %d cannot be parsed to a string",
				rawStr, idx)})
		}

		ss = append(ss, rawStr)
	}

	*dest = &ss
	return nil
}

// UUID tries to select a key as a UUID. An error will occur when the key doesn't exist or when the
// value cannot be parsed to a UUID.
func (rs *CdbResult) UUID(key string, dest **uuid.UUID) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	if reflect.TypeOf(raw) == nil {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to uuid", raw)})
	}

	if reflect.TypeOf(raw).String() == "string" {
		uuid, err := uuid.FromString(raw.(string))
		if err != nil {
			return rs.addError(
				&ErrMalformed{fmt.Errorf("%v cannot be parsed to uuid", raw)})
		}

		if dest == nil {
			uuidPtr := &uuid
			dest = &uuidPtr
			return nil
		}

		*dest = &uuid
		return nil
	}

	uintSlice, ok := raw.([16]uint8)
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to uuid", raw)})
	}

	uuid, err := uuid.FromBytes(uintSlice[:])
	if err != nil {
		return err
	}

	*dest = &uuid
	return nil
}

// OptUUID tries to select a key as a UUID. An error will occur when the value cannot be parsed to a
// UUID.
func (rs *CdbResult) OptUUID(key string, dest **uuid.UUID) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	if reflect.TypeOf(raw) == nil {
		return nil
	}

	if reflect.TypeOf(raw).String() == "string" {
		uuid, err := uuid.FromString(raw.(string))
		if err != nil {
			return rs.addError(
				&ErrMalformed{fmt.Errorf("%v cannot be parsed to uuid", raw)})
		}

		if dest == nil {
			uuidPtr := &uuid
			dest = &uuidPtr
			return nil
		}

		*dest = &uuid
		return nil
	}

	uintSlice, ok := raw.([16]uint8)
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to uuid", raw)})
	}

	uuid, err := uuid.FromBytes(uintSlice[:])
	if err != nil {
		return err
	}

	*dest = &uuid
	return nil
}

// Object tries to select a key as an object. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an object.
func (rs *CdbResult) ObjectSlice(key string, dest *[]*CdbResult) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	objSlice, ok := raw.([]interface{})
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to a slice of CdbResults", raw)})
	}

	for _, obj := range objSlice {
		data, ok := obj.(map[string]interface{})
		if !ok {
			return rs.addError(
				&ErrMalformed{fmt.Errorf("%v cannot be parsed to a CdbResult", obj)})
		}

		*dest = append(*dest, &CdbResult{
			data: data,
		})
	}

	return nil
}

// Object tries to select a key as an object. An error will occur when the key doesn't exist or when
// the value cannot be parsed to an object.
func (rs *CdbResult) Object(key string, dest **CdbResult) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	obj, ok := raw.(map[string]interface{})
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to a CdbResult", raw)})
	}

	*dest = &CdbResult{
		data: obj,
	}
	return nil
}

// OptObject tries to select a key as an object. An error will occur when the value cannot be parsed
// to an object.
func (rs *CdbResult) OptObject(key string, dest **CdbResult) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		*dest = nil
		return nil
	}

	if raw == nil {
		return nil
	}

	obj, ok := raw.(map[string]interface{})
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to a CdbResult", raw)})
	}

	*dest = &CdbResult{
		data: obj,
	}

	return nil
}

// Bool tries to select a key as a boolean. An error will occur when the key doesn't exist or when
// the value cannot be parsed to a boolean.
func (rs *CdbResult) Bool(key string, dest **bool) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	res, ok := raw.(bool)
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s %v cannot be parsed to a bool", key, raw)})
	}

	*dest = &res
	return nil
}

// OptBool tries to select a key as a boolean. An error will occur when the value cannot be parsed
// to a boolean.
func (rs *CdbResult) OptBool(key string, dest **bool) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	res, ok := raw.(bool)
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a bool", raw)})
	}

	*dest = &res
	return nil
}

// Float32 tries to select a key as a float32. An error will occur when the key doesn't exist or
// when the value cannot be parsed to an float32.
func (rs *CdbResult) Float32(key string, dest **float32) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	rawNum, err := strconv.ParseFloat(fmt.Sprintf("%v", raw), 10)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a float32", raw)})
	}

	num := float32(rawNum)
	*dest = &num
	return nil
}

// OptFloat32 tries to select a key as a float32. An error will occur when the value cannot be
// parsed to a float32.
func (rs *CdbResult) OptFloat32(key string, dest **float32) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	num, ok := raw.(float32)
	if !ok {
		return rs.addError(
			&ErrMalformed{fmt.Errorf("%v cannot be parsed to a float32", raw)})
	}

	*dest = &num
	return nil
}

// Float64 tries to select a key as a float64. An error will occur when the key doesn't exist or
// when the value cannot be parsed to an float64.
func (rs *CdbResult) Float64(key string, dest **float64) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	num, err := strconv.ParseFloat(fmt.Sprintf("%v", raw), 10)
	if err != nil {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a float64", raw)})
	}

	*dest = &num
	return nil
}

// OptFloat64 tries to select a key as a float64. An error will occur when the value cannot be
// parsed to a float64.
func (rs *CdbResult) OptFloat64(key string, dest **float64) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	num, ok := raw.(float64)
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a float64", raw)})
	}

	*dest = &num
	return nil
}

// OptMapStrStr tries to select a key as a map[string]string. An error will occur when the value cannot be
// parsed to a map[string]string.
func (rs *CdbResult) OptMapStrStr(key string, dest **map[string]string) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	if raw == nil {
		return nil
	}

	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a map[string]string", raw)})
	}

	mapStrStr := map[string]string{}

	for key, value := range rawMap {
		str, ok := value.(string)
		if !ok {
			return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to string", key, str)})
		}

		mapStrStr[key] = str
	}

	*dest = &mapStrStr
	return nil
}

// MapStrStr tries to select a key as a map[string]string. An error will occur when the value cannot be
// parsed to a map[string]string.
func (rs *CdbResult) MapStrStr(key string, dest **map[string]string) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a map[string]string", raw)})
	}

	mapStrStr := map[string]string{}

	for key, value := range rawMap {
		str, ok := value.(string)
		if !ok {
			return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to string", key, str)})
		}

		mapStrStr[key] = str
	}

	*dest = &mapStrStr
	return nil
}

// MapStrInterface tries to select a key as a map[string]interface{}. An error will occur when the
// value cannot be parsed to a map[string]interface{}.
func (rs *CdbResult) MapStrInterface(key string, dest **map[string]interface{}) error {
	if rs == nil || rs.data == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a map[string]interface{}", raw)})
	}

	*dest = &rawMap
	return nil
}

// MapStrInterface tries to select a key as a map[string]interface{}. An error will occur when the
// value cannot be parsed to a map[string]interface{}.
func (rs *CdbResult) OptMapStrInterface(key string, dest **map[string]interface{}) error {
	if rs == nil || rs.data == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	if raw == nil {
		return nil
	}

	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%v cannot be parsed to a map[string]interface{}", raw)})
	}

	*dest = &rawMap
	return nil
}

func OptCustomType[T interface{}](rs *CdbResult, key string, dest **T, cb func(s string) T) error {
	if rs == nil {
		return nil
	}

	raw, present := rs.data[key]
	if !present {
		return nil
	}

	if raw == nil {
		return nil
	}

	str, ok := raw.(string)
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to string", key, raw)})
	}

	res := cb(str)

	*dest = &res

	return nil
}

func CustomType[T interface{}](rs *CdbResult, key string, dest **T, cb func(s string) T) error {
	if rs == nil {
		return &ErrMissingResult{}
	}

	raw, present := rs.data[key]
	if !present {
		return rs.addError(&ErrNoSuchKey{key})
	}

	str, ok := raw.(string)
	if !ok {
		return rs.addError(&ErrMalformed{fmt.Errorf("%s: %v cannot be parsed to string", key, raw)})
	}

	res := cb(str)

	*dest = &res

	return nil
}
