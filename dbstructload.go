package dbstructload

import (
	"database/sql"
	"errors"
	"reflect"
)

const (
	tagKey = "queryField"
)

var (
	// ErrNotStructPtr is returns when bad arguments are provided to Scan.
	ErrNotStructPtr = errors.New("arguments must be struct pointers")

	// ErrMissingField is returns when the structs provided to Scan don't match
	// the columns in the query results.
	ErrMissingField = errors.New("no provided structs contain field matching column")
)

// Rows is a wrapper around *sql.Rows that loads query results into structs.
type Rows struct {
	*sql.Rows
}

// Query uses the provided *sql.DB to make a query.
func Query(db *sql.DB, query string, args ...interface{}) (*Rows, error) {
	sqlRows, err := db.Query(query, args...)
	return &Rows{sqlRows}, err
}

// Load loads the current row into the provided structs.
func (r *Rows) Load(destStructs ...interface{}) error {
	cols, err := r.Columns()
	if err != nil {
		return err
	}

	destValues := make([]interface{}, 0, len(cols))

	var (
		tmpPtr       reflect.Value
		tmpStruct    reflect.Value
		tmpSructType reflect.Type
		tmpTagVal    string
		ok           bool
	)

ColumnIteration:
	for _, col := range cols {
		for _, destStruct := range destStructs {
			tmpPtr = reflect.ValueOf(destStruct)

			if tmpPtr.Kind() != reflect.Ptr {
				return ErrNotStructPtr
			}

			tmpStruct = tmpPtr.Elem()
			if tmpStruct.Kind() != reflect.Struct {
				return ErrNotStructPtr
			}

			tmpSructType = tmpStruct.Type()
			for i := 0; i < tmpSructType.NumField(); i++ {
				tmpTagVal, ok = tmpSructType.Field(i).Tag.Lookup(tagKey)
				if ok && tmpTagVal == col {
					destValues = append(destValues, tmpStruct.Field(i).Interface())
					continue ColumnIteration
				}
			}
		}

		return ErrMissingField
	}

	return r.Scan(destValues...)
}
