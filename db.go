package gosqlgo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func execQuery(query string) *sql.Rows {
	os.Setenv("USER_DATABASE", "test")
	os.Setenv("PASS_DATABASE", "12345678")
	os.Setenv("HOST_DATABASE", "localhost")
	os.Setenv("NAME_DATABASE", "test")
	usuario := os.Getenv("USER_DATABASE")
	password := os.Getenv("PASS_DATABASE")
	host := os.Getenv("HOST_DATABASE")
	database := os.Getenv("NAME_DATABASE")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, password, host, database))
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Query(query)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}

	return res
}

func StructScan(rows *sql.Rows, model interface{}) error {

	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination") // @todo add new error message
	}

	v = reflect.Indirect(v)
	t := v.Type()

	// Find out if what we are using is a Struct ( Query One ) or a Slice with Structs ( QueryAll )
	isStruct, isSlice := false, false
	if t.Kind() == reflect.Slice {
		isSlice = true
	} else if t.Kind() == reflect.Struct {
		isStruct = true
	}

	// Ensure we only get the column information one time!
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	var m []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		x := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			x[colName] = *val
		}
		m = append(m, x)

		// If  we are dealing with a struct. There is no point in looping over all the results, if they are more then one!
		if isStruct == true {
			break
		}
	}

	if isStruct == true {
		if len(m) > 0 { // Ensure we have data fields!
			changeStruct(v, t, m[0])
		}
	}

	if isSlice == true {
		if len(m) > 0 { // Ensure we have data in the slice!
			var elem reflect.Value
			for _, d := range m {

				typ := v.Type().Elem()
				elem = reflect.New(typ).Elem()

				changeStruct(elem, typ, d)
				v.Set(reflect.Append(v, elem))
			}
		}
	}

	return nil
}

func changeStruct(v reflect.Value, t reflect.Type, m map[string]interface{}) {
	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("db"), ",")[0]

		if item, ok := m[field]; ok {
			if v.Field(i).CanSet() {
				if item != nil {
					switch v.Field(i).Kind() {
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						uinttext := (item.(interface{}).([]uint8))
						if intVal, err := strconv.ParseUint(string(uinttext), 10, 64); err == nil {
							v.Field(i).SetUint(uint64(intVal))
						}
					case reflect.String:
						v.Field(i).SetString(item.(string)) //	s := bytesToString(item.([]uint8))
					case reflect.Float32, reflect.Float64:
						v.Field(i).SetFloat(item.(float64))
					case reflect.Int, reflect.Int32, reflect.Int64:
						v.Field(i).SetInt(item.(int64))
					case reflect.Bool:
						v.Field(i).Set(reflect.ValueOf(!(item.(int64) == 0)))
					case reflect.Ptr:
						if reflect.ValueOf(item).Kind() == reflect.Bool {
							itemBool := item.(bool)
							v.Field(i).Set(reflect.ValueOf(&itemBool))
						}
					case reflect.Struct:
						v.Field(i).Set(reflect.ValueOf(item))
					default:
						fmt.Println(t.Field(i).Name, ": ", v.Field(i).Kind(), " - > - ", reflect.ValueOf(item).Kind()) // @todo remove after test out the Get methods
					}
				}
			}
		}
	}
}
