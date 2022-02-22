package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type DataType int

const (
	T_RAW DataType = iota
	T_ROWS
)

type Data interface {
	Type() DataType
	Columns() []string
	Rows() [][]interface{}
	Raw() []byte
}

type payload struct {
	t    DataType
	cols []string
	rows [][]interface{}
	raw  []byte
}

func (d payload) Type() DataType {
	return d.t
}

func (d payload) Columns() []string {
	return d.cols
}

func (d payload) Rows() [][]interface{} {
	return d.rows
}

func (d payload) Raw() []byte {
	return d.raw
}

func FromJson(data []byte) (Data, error) {
	var js interface{}
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, err
	}

	var cols []string
	var rows [][]interface{}
	switch v := js.(type) {
	case []map[string]interface{}:
		for i, data := range v {
			if i == 0 {
				for k, _ := range data {
					cols = append(cols, k)
				}
				sort.Strings(cols)
			}

			row := createRowFromMap(cols, data)
			rows = append(rows, row)
		}
	case map[string]interface{}:
		for k, _ := range v {
			cols = append(cols, k)
		}
		sort.Strings(cols)
		rows = append(rows, createRowFromMap(cols, v))
	}

	return &payload{
		t:    T_ROWS,
		cols: cols,
		rows: rows,
	}, nil
}

func FromRawBytes(data []byte) Data {
	return &payload{
		t:   T_RAW,
		raw: data,
	}
}

func FromStruct(data interface{}) (Data, error) {
	dval := reflect.ValueOf(data)
	dtyp := dval.Type()

	if dtyp.Kind() == reflect.Ptr {
		dtyp = dtyp.Elem()
		dval = dval.Elem()
	}

	if dtyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}

	cols, _ := getCols(dtyp)
	row, _ := getRow(cols, dval)

	return &payload{
		t:    T_ROWS,
		cols: cols,
		rows: [][]interface{}{row},
	}, nil
}

func FromSliceOfStruct(data interface{}) (Data, error) {
	dval := reflect.ValueOf(data)
	dtyp := dval.Type()

	if dtyp.Kind() != reflect.Slice {
		return nil, fmt.Errorf("data is not a slice")
	}

	var cols []string
	var rows [][]interface{}
	for i := 0; i < dval.Len(); i++ {
		elval := dval.Index(i)
		if elval.Type().Kind() == reflect.Ptr {
			elval = elval.Elem()
		}

		if i == 0 {
			if elval.Type().Kind() != reflect.Struct {
				return nil, fmt.Errorf("value is not a struct")
			}

			cols, _ = getCols(elval.Type())
		}

		row, _ := getRow(cols, elval)
		rows = append(rows, row)

	}

	return &payload{
		t:    T_ROWS,
		cols: cols,
		rows: rows,
	}, nil

}

func FromHeadersAndRows(header []string, rows [][]interface{}) Data {
	return &payload{
		t:    T_ROWS,
		cols: header,
		rows: rows,
	}
}

func getCols(dtyp reflect.Type) ([]string, error) {
	if dtyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type is not a struct")
	}
	var cols []string
	for i := 0; i < dtyp.NumField(); i++ {
		field := dtyp.Field(i)
		tag, ok := field.Tag.Lookup("col")
		if !ok {
			tag = field.Name
		}

		cols = append(cols, tag)
	}

	sort.Strings(cols)
	return cols, nil
}

func getRow(cols []string, dval reflect.Value) ([]interface{}, error) {
	if dval.Type().Kind() != reflect.Struct {
		return nil, fmt.Errorf("value is not a struct")
	}

	dtyp := dval.Type()
	var colmap = make(map[string]int)
	for i, c := range cols {
		colmap[c] = i
	}

	var row = make([]interface{}, len(cols))
	for i := 0; i < dval.NumField(); i++ {
		val := dval.Field(i)
		field := dtyp.Field(i)
		tag, ok := field.Tag.Lookup("col")
		if !ok {
			tag = field.Name
		}

		if idx, ok := colmap[tag]; ok {
			row[idx] = val.Interface()
		}
	}

	return row, nil
}

func createRowFromMap(cols []string, data map[string]interface{}) []interface{} {
	var row []interface{}
	for _, c := range cols {
		var val interface{}
		v, ok := data[c]
		if ok {
			val = v
		} else {
			val = nil
		}

		row = append(row, val)
	}

	return row
}
