package cdb

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func clean() {
	dbPool = nil
}

func TestInit(t *testing.T) {
	clean()

	err := Init("192.168.0.35", 26257, "root", "", "bluecherry")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	err = Init("", 26257, "root", "", "bluecherry")
	if err == nil {
		t.Fatal("Expected to throw ErrConnect")
	}
}

func TestGetCon(t *testing.T) {
	clean()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	err := Init("192.168.0.35", 26257, "root", "", "bluecherry")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	_, err = getCon(ctx)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	dbPool = nil
	_, err = getCon(ctx)
	if err == nil {
		t.Fatal("Expected to throw ErrNotInstantiated")
	}
}

type StringInterfaceMap map[string]interface{}

func (m *StringInterfaceMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func TestExec(t *testing.T) {
	clean()

	err := Init("192.168.0.35", 26257, "root", "", "bluecherry")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

}

func TestSelect(t *testing.T) {
	email := "daan@dptechnics.com"
	acc := struct {
		Email *string
		Name  *string
	}{Email: &email}

	colNames := map[string]string{
		"Email": "email",
		"Name":  "name",
	}

	fields := []string{"email AS randomEmail", "name"}

	PrepareSelect("account", fields, "a", colNames, acc)
}

func TestUpdate(t *testing.T) {
	email := "daan@dptechnics.com"
	acc := struct {
		Email *string
		Name  *string
	}{Email: &email}

	colNames := map[string]string{
		"Email": "email",
		"Name":  "name",
	}

	selectors := map[string]interface{}{
		"id": "8e17f234-9e0a-49dc-8a7d-c967ca2df1c6",
	}

	_, err := PrepareUpdate("account", colNames, acc, selectors)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	email := "daan@dptechnics.com"
	name := "joske"
	acc := struct {
		Email *string
		Name  *string
	}{Email: &email, Name: &name}

	colNames := map[string]string{
		"Email": "email",
		"Name":  "name",
	}

	_, err := PrepareInsert("account", colNames, acc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	email := "daan@dptechnics.com"
	name := "daan Pape"
	acc := struct {
		Email *string
		Name  *string
	}{Email: &email, Name: &name}

	colNames := map[string]string{
		"Email": "email",
		"Name":  "name",
	}

	PrepareDelete("account", colNames, acc)
}

func TestTransaction(t *testing.T) {
	clean()

	err := Init("192.168.0.35", 26257, "root", "", "bluecherry")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	tx := BeginTx()

	tx.AddStmt(Prepare(`select email, name from account where email = 'daan@dptechnics.com'`))
	tx.AddStmt(Prepare(`select email, name from account where email = 'daan@dptechnics.com'`))
	tx.AddStmt(Prepare(`select email, name from account where email = 'daan@dptechnics.com'`))

	err = tx.Exec()
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
}

func TestTime(t *testing.T) {
	cdbRes := CdbResult{
		data: map[string]interface{}{
			"time":   "2022-07-20 10:44:03.670 +0200",
			"time_2": "2018-12-19 10:19:39 +0100 CET",
		},
	}

	type Tmp struct {
		start *time.Time
		end   *time.Time
	}

	tmp := Tmp{}

	err := cdbRes.Time("time", &tmp.start)
	if err != nil {
		t.Fatal(err)
	}

	err = cdbRes.Time("time_2", &tmp.end)
	if err != nil {
		t.Fatal(err)
	}
}
