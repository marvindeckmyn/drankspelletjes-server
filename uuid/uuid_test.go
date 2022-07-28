package uuid

import (
	"fmt"
	"regexp"
	"testing"
)

func TestNil(t *testing.T) {
	uuid := Nil()
	str := uuid.String()

	if str != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("UUID '%s' is not valid nil UUID", str)
	}
}

func TestUUIDv4(t *testing.T) {
	uuid := UUIDv4()
	str := uuid.String()

	fmt.Println(str)

	want := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-" +
		"[8|9|a|b][a-f0-9]{3}-[a-f0-9]{12}$")

	if !want.MatchString(str) {
		t.Fatalf("UUID '%s' is not valid UUIDv4", str)
	}
}

func TestByteParse(t *testing.T) {
	data := []uint8{0xaa, 0x6f, 0xce, 0xf7, 0x36, 0x08, 0x48, 0x95, 0xa4, 0xc9, 0x21, 0x7b, 0x51, 0x72, 0x68, 0xef}

	uuid, err := FromBytes(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(uuid.String())
	fmt.Println(String(data))

	data = []uint8{0xaa, 0x6f, 0xce, 0xf7, 0x36, 0x08, 0x48, 0x95, 0xa4, 0xc9, 0x21, 0x7b, 0x51, 0x72, 0x68}

	uuid, err = FromBytes(data)
	if err == nil {
		t.Fatalf(err.Error())
	}

	if String(data) != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("Invalid byte array should return nil UUID")
	}

	fmt.Println(uuid.String())
}

func TestStringParse(t *testing.T) {
	uuid, err := FromString("1b1b47e6-c1e7-4d38-91a7-dfb08f1d3446")
	if err != nil {
		t.Fatalf("Parsing the UUID failed: %s", err.Error())
	}

	if uuid.String() != "1b1b47e6-c1e7-4d38-91a7-dfb08f1d3446" {
		t.Fatalf("Parse UUID doesn't match given UUID")
	}

	_, err = FromString("hello")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}

	_, err = FromString("1b1b47e6+c1e7-4d38-91a7-dfb08f1d3446")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}

	_, err = FromString("gb1b47e6-c1e7-4d38-91a7-dfb08f1d3446")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}

	_, err = FromString("1b1b47e6-g1e7-4d38-91a7-dfb08f1d3446")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}

	_, err = FromString("1b1b47e6-c1e7-gd38-91a7-dfb08f1d3446")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}

	_, err = FromString("1b1b47e6-c1e7-4d38-g1a7-dfb08f1d3446")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}

	_, err = FromString("1b1b47e6-c1e7-4d38-91a7-gfb08f1d3446")
	if err == nil {
		t.Fatalf("Parsing should of have failed")
	}
}
