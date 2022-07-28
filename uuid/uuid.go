// The UUID package provides a UUID type which is compatible with any 16-byte
// long UUID. By default the package works with v4 and variant 1 UUIDs.
package uuid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
)

// A UUID is a byte array which contains 16 bytes.
type UUID [16]byte

// UnmarshalJSON defines how to convert a byte slice to the UUID type.
func (instance *UUID) UnmarshalJSON(b []byte) error {
	uuidStr := string(b)
	uuidStr = strings.ReplaceAll(uuidStr, "\"", "")
	uuid, err := FromString(uuidStr)
	if err != nil {
		return err
	}

	*instance = uuid
	return nil
}

// MarshalJSON defines how to convert a UUID to a byte slice.
func (u *UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// Nil returns an empty all zero UUID.
func Nil() UUID {
	return UUID{}
}

// UUIDv4 generates a random uuid which conforms to the UUID version 4 standard.
func UUIDv4() UUID {
	data := UUID{}
	rand.Read(data[:])

	data[6] &= 0x0F
	data[6] |= 0x40
	data[8] &= 0x3F
	data[8] |= 0x80

	return data
}

// String returns the string representation of a UUID.
func (u UUID) String() string {
	res := make([]byte, 36)

	hex.Encode(res[:8], u[:4])
	res[8] = '-'
	hex.Encode(res[9:13], u[4:6])
	res[13] = '-'
	hex.Encode(res[14:18], u[6:8])
	res[18] = '-'
	hex.Encode(res[19:23], u[8:10])
	res[23] = '-'
	hex.Encode(res[24:], u[10:])

	return string(res)
}

// String is a helper function which prints out a uint8 slice into a UUID string. When the incoming
// byte array is not a valid UUID a nil UUID will be printed.
func String(bytes []uint8) string {
	uuid, err := FromBytes(bytes)
	if err != nil {
		return Nil().String()
	}

	return uuid.String()
}

// FromBytes tries to read a uint8 array into a UUID object. The contents of the
// bytes are not checked for maximum flexibility.
func FromBytes(bytes []uint8) (UUID, error) {
	uuid := UUID{}

	if len(bytes) != 16 {
		return uuid, &ErrMalformed{fmt.Errorf("must be 16B long, got %dB", len(bytes))}
	}

	copy(uuid[:], bytes)
	return uuid, nil
}

// FromString parses a string into a UUID. This function will accept any kind of
// UUID for flexibility.
func FromString(str string) (UUID, error) {
	uuid := UUID{}

	if len(str) != 36 {
		return uuid, &ErrMalformed{fmt.Errorf("string UUID must be 36 characters long")}
	}

	if str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		return uuid, &ErrMalformed{fmt.Errorf("missing dashes in UUID string")}
	}

	b := []byte(str)

	_, err := hex.Decode(uuid[:4], b[:8])
	if err != nil {
		return uuid, &ErrMalformed{fmt.Errorf("invalid characters in UUID string")}
	}

	_, err = hex.Decode(uuid[4:6], b[9:13])
	if err != nil {
		return uuid, &ErrMalformed{fmt.Errorf("invalid characters in UUID string")}
	}

	_, err = hex.Decode(uuid[6:8], b[14:18])
	if err != nil {
		return uuid, &ErrMalformed{fmt.Errorf("invalid characters in UUID string")}
	}

	_, err = hex.Decode(uuid[8:10], b[19:23])
	if err != nil {
		return uuid, &ErrMalformed{fmt.Errorf("invalid characters in UUID string")}
	}

	_, err = hex.Decode(uuid[10:], b[24:])
	if err != nil {
		return uuid, &ErrMalformed{fmt.Errorf("invalid characters in UUID string")}
	}

	return uuid, nil
}
