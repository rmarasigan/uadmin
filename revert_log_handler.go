package uadmin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// RevertLogHandler !
func revertLogHandler(w http.ResponseWriter, r *http.Request) {
	log := Log{}
	Get(&log, "id = ?", r.FormValue("log_id"))
	if log.Action == log.Action.Deleted() {
		s := db.NewScope(models[log.TableName])
		tableName := s.TableName()
		sql := fmt.Sprintf("update %s set deleted_at = null where id = %d", tableName, log.TableID)
		db.Exec(sql)
	} else {
		now := time.Now()
		DType := reflect.TypeOf(now)
		var langParser map[string]string
		err := json.Unmarshal([]byte(log.Activity), &langParser)
		if err != nil {
			fmt.Println(err.Error())
		}
		modelType := reflect.TypeOf(models[log.TableName]).Elem()
		newType := reflect.New(modelType)
		Get(newType.Elem().Addr().Interface(), "id = ?", log.TableID)
		model, ok := NewModel(log.TableName, true)
		if !ok {
		}
		var t reflect.Type
		t = reflect.TypeOf(model.Interface()).Elem()
		for index := 0; index < t.NumField(); index++ {
			if t.Field(index).Type.Kind() == reflect.Int {
				_v, _ := strconv.Unquote(string(langParser[t.Field(index).Name]))
				_v = fmt.Sprintf("%+v", _v)
				if len(langParser[t.Field(index).Name]) >= 2 {
				}
				i, _ := strconv.ParseInt(_v, 10, 64)
				newType.Elem().FieldByName(t.Field(index).Name).SetInt(i)
			} else if t.Field(index).Type.Kind() == reflect.String {
				// Check if Multilingual
				val := ""
				if t.Field(index).Tag.Get("multilingual") == cTRUE {
					tVal := map[string]string{}
					langs := []Language{}
					Filter(&langs, "`active` = ?", true)
					for _, lang := range langs {
						tVal[lang.Code] = fmt.Sprint(langParser[lang.Code+"-"+t.Field(index).Name])
					}
					b, _ := json.Marshal(tVal)
					val = string(b)
				} else {
					val = string(langParser[t.Field(index).Name])
				}

				newType.Elem().FieldByName(t.Field(index).Name).SetString(val)
			} else if t.Field(index).Type.Kind() == reflect.Bool {
				var val bool
				val = false
				if string(langParser[t.Field(index).Name]) == "true" {
					val = true
				}
				newType.Elem().FieldByName(t.Field(index).Name).SetBool(val)
			} else if t.Field(index).Type.Kind() == reflect.Uint {
				_v := string(langParser[t.Field(index).Name])
				i, _ := strconv.ParseInt(_v, 10, 64)
				val := uint(i)
				newType.Elem().FieldByName(t.Field(index).Name).Set(reflect.ValueOf(val))
			} else if t.Field(index).Type.Kind() == reflect.Float64 {
				_v := string(langParser[t.Field(index).Name])
				i, _ := strconv.ParseFloat(_v, 10)
				newType.Elem().FieldByName(t.Field(index).Name).Set(reflect.ValueOf(i))
			} else if t.Field(index).Type.Kind() == reflect.Ptr {
				if fmt.Sprint(langParser[t.Field(index).Name]) != "" {
					tm, _ := time.Parse("2006-01-02 15:04", string(langParser[t.Field(index).Name]))
					newType.Elem().FieldByName(t.Field(index).Name).Set(reflect.ValueOf(&tm))
				}
			} else if t.Field(index).Type == DType {
				_v := string(langParser[t.Field(index).Name])
				tm, _ := time.Parse("2006-01-02 15:04", _v)
				newType.Elem().FieldByName(t.Field(index).Name).Set(reflect.ValueOf(tm))
			} else {
			}
		}

		Save(newType.Elem().Addr().Interface())
	}
}
