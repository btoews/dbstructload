package model

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

// Rows is a wrapper around *sql.Rows that loads query results into model
// structs.
type Rows struct {
	sqlRows *sql.Rows
	err     error
}

// Query uses the provided *sql.DB to make a query.
func Query(db *sql.DB, query string, args ...interface{}) (*Rows, error) {
	sqlRows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return &Rows{sqlRows, nil}, nil
}

// Next prepares the next result row for reading with the Scan method. It
// returns true on success, or false if there is no next result row or an error
// happened while preparing it. Err should be consulted to distinguish between
// the two cases.
//
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (r *Rows) Next() bool {
	if err := r.Err(); err != nil {
		return false
	}

	return r.sqlRows.Next()
}

// Scan loads the current row into the provided model structs.
func (r *Rows) Scan(destStructs ...interface{}) error {
	cols, err := r.sqlRows.Columns()
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

	return r.sqlRows.Scan(destValues...)
}

// Err returns the error, if any, that was encountered during iteration. Err may
// be called after an explicit or implicit Close.
func (r *Rows) Err() error {
	if r.err != nil {
		return r.err
	}

	r.err = r.sqlRows.Err()
	return r.err
}

// Close closes the Rows, preventing further enumeration. If Next is called and
// returns false and there are no further result sets, the Rows are closed
// automatically and it will suffice to check the result of Err. Close is
// idempotent and does not affect the result of Err.
func (r *Rows) Close() error {
	return r.sqlRows.Close()
}
