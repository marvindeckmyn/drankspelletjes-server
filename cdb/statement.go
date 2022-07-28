package cdb

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Statement represents the statement which is used for querying the database.
type Statement struct {
	Query  string
	params map[string]interface{}
}

// AppendReturningID appends the returning id clause to the statement
func (s *Statement) AppendReturningID() error {
	s.Query = fmt.Sprintf("%s\nRETURNING id", s.Query)
	return nil
}

// appendWhere appends the where clause to the statement.
func (s *Statement) appendWhere(colNames map[string]string, label string, obj interface{}) error {
	var v = reflect.Value{}

	v = reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typeOfS := v.Type()
	containsWhere := strings.Contains(strings.ToLower(s.Query), "where")

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			continue
		}

		name, ok := colNames[typeOfS.Field(i).Name]
		if !ok {
			return &ErrNoSuchKey{name}
		}

		action := "AND"
		if !containsWhere {
			action = "WHERE"
			containsWhere = true
		}

		if label == "" {
			s.Query = fmt.Sprintf("%s\n%s %s = :%s:", s.Query, action, name, name)
		} else {
			s.Query = fmt.Sprintf("%s\n%s %s.%s = :%s.%s:", s.Query, action, label, name, label, name)
		}
	}

	return nil
}

// Bind adds the key value to the statement parameters.
func (s *Statement) Bind(key string, val interface{}) {
	s.params[key] = val
}

// BindMap adds the key values from the map to the statement parameters.
func (s *Statement) BindMap(params map[string]interface{}) {
	for key, value := range params {
		s.params[key] = value
	}
}

// BindObject dynamically adds the non nil values from the given struct to the statement parameters.
func (s *Statement) BindObject(colNames map[string]string, label string, obj interface{}) error {
	var v = reflect.Value{}

	v = reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if s.params == nil {
		s.params = map[string]interface{}{}
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct || v.Field(i).IsNil() {
			continue
		}

		name, ok := colNames[v.Type().Field(i).Name]
		if !ok {
			if name == "" {
				continue
			}

			return &ErrNoSuchKey{name}
		}

		elem := v.Field(i).Elem()
		val := elem.Interface()

		if label == "" {
			s.params[name] = val
		} else {
			s.params[label+"."+name] = val
		}
	}

	return nil
}

// CreateFields creates the fields to query from the db
func CreateFields(colNames map[string]string) []string {
	fields := []string{}
	for _, field := range colNames {
		fields = append(fields, field)
	}
	return fields
}

// GetString returns the query string from the statement with it's linked values.
func (s *Statement) GetString() (string, map[string]interface{}) {
	return s.Query, s.params
}

func (s *Statement) getRequiredValues() ([]string, error) {
	r := regexp.MustCompile(`:.*?:`)
	matches := r.FindAllString(s.Query, -1)

	requiredValues := []string{}

	for _, match := range matches {
		key := strings.ReplaceAll(match, ":", "")
		if key == "" {
			continue
		}

		requiredValues = append(requiredValues, key)
	}

	return requiredValues, nil
}

// GetPgQuery returns the query which is ready to be executed on the database.
func (s *Statement) getPgQuery() (string, []interface{}, error) {
	r := regexp.MustCompile(`:.*?:`)
	matches := r.FindAllString(s.Query, -1)

	pgQuery := s.Query
	values := []interface{}{}

	for idx, match := range matches {
		pgQuery = strings.Replace(pgQuery, match, "$"+fmt.Sprint(idx+1), 1)
		key := strings.ReplaceAll(match, ":", "")
		if key == "" {
			continue
		}

		val, ok := s.params[key]
		if !ok {
			return "", nil, &ErrMissingValue{key}
		}

		values = append(values, val)
	}

	return pgQuery, values, nil
}

// Prepare creates a statement with the given query.
func Prepare(query string) Statement {
	return Statement{
		Query:  query,
		params: map[string]interface{}{},
	}
}

// PrepareSelect builds a select statement with the given fields.
func PrepareSelect(table string, fields []string, label string, colNames map[string]string, obj interface{}) Statement {
	query := "SELECT "

	for i, field := range fields {
		if label == "" {
			query = fmt.Sprintf("%s %s", query, field)
		} else {
			query = fmt.Sprintf("%s %s.%s", query, label, field)
		}

		if i != len(fields)-1 {
			query = query + ","
		}
	}

	query = fmt.Sprintf("%s \nFROM %s %s", query, table, label)

	stmt := Prepare(query)
	stmt.appendWhere(colNames, label, obj)
	stmt.BindObject(colNames, label, obj)
	return stmt
}

// PrepareInsert builds an insert statement with the given fields.
func PrepareInsert(table string, colNames map[string]string, obj interface{}) (Statement, error) {
	names := ""
	values := ""

	var v = reflect.Value{}

	v = reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct || v.Field(i).IsNil() {
			continue
		}

		name, ok := colNames[typeOfS.Field(i).Name]
		if !ok {
			if name == "" {
				continue
			}

			return Statement{}, &ErrNoSuchKey{name}
		}

		if names == "" {
			names = name
			values = fmt.Sprintf(":%v:", name)
		} else {
			names = fmt.Sprintf("%s, %s", names, name)
			values = fmt.Sprintf("%s, :%v:", values, name)
		}
	}

	stmt := Prepare(fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES(%s)
	`, table, names, values))
	stmt.BindObject(colNames, "", obj)

	return stmt, nil
}

// PrepareUpdate builds an update statement with the given fields.
func PrepareUpdate(table string, colNames map[string]string, obj interface{}, selectors map[string]interface{}) (Statement, error) {
	query := fmt.Sprintf("UPDATE %s \nSET ", table)
	var v = reflect.Value{}

	v = reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fields := []string{}
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct || v.Field(i).IsNil() {
			continue
		}

		name, ok := colNames[typeOfS.Field(i).Name]
		if !ok {
			if name == "" {
				continue
			}

			return Statement{}, &ErrNoSuchKey{name}
		}

		fields = append(fields, name)
	}

	i := 0
	for _, field := range fields {
		query = fmt.Sprintf("%s %s = :%s:", query, field, field)

		if i < len(fields)-1 {
			query = query + ", "
		}

		i++
	}

	stmt := Prepare(query)

	containsWhere := strings.Contains(strings.ToLower(stmt.Query), "where")
	for key, val := range selectors {
		action := "AND"
		if !containsWhere {
			action = "WHERE"
			containsWhere = true
		}

		stmt.Query = fmt.Sprintf("%s\n%s %s = :%s:", stmt.Query, action, key, "UPDATE_"+key)
		stmt.Bind("UPDATE_"+key, val)
	}

	stmt.BindObject(colNames, "", obj)

	return stmt, nil
}

// PrepareDelete builds a delete statement with the given fields.
func PrepareDelete(table string, colNames map[string]string, obj interface{}) Statement {
	query := fmt.Sprintf("DELETE FROM %s", table)
	stmt := Prepare(query)
	stmt.appendWhere(colNames, "", obj)
	stmt.BindObject(colNames, "", obj)

	return stmt
}
