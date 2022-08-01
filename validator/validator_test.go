package validator

import (
	"fmt"
	"testing"
)

func TestIsNil(t *testing.T) {
	type Tmp struct {
		Test string
	}

	var tmp *Tmp = nil

	numb := 5
	ptr := &numb
	ptr = nil

	err := IsNotNil(map[string]interface{}{
		"tmp":     tmp,
		"random":  nil,
		"nog een": nil,
		"ptr":     ptr,
	})
	if err == nil {
		t.Fatal("expected to throw error")
	}

	fmt.Println(err)
}
