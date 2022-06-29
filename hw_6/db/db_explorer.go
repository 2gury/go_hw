package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type ReqValue struct{}

type Body map[string]interface{}

type DbFieldMetaInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type FieldMetaInfo struct {
	FieldName  string
	Type       string
	IsNullable bool
}

type Route struct {
	Method string
	RegExp *regexp.Regexp
	Handler http.HandlerFunc
}

func ConvertDbFieldMetaInfo(dbInfo DbFieldMetaInfo) FieldMetaInfo {
	convInfo := FieldMetaInfo{}

	convInfo.FieldName = dbInfo.Field

	switch {
	case strings.Contains(dbInfo.Type, "int"):
		convInfo.Type = "int"
	case strings.Contains(dbInfo.Type, "varchar"):
		convInfo.Type = "string"
	case strings.Contains(dbInfo.Type, "text"):
		convInfo.Type = "string"
	default:
		convInfo.Type = "unknown"
	}

	switch dbInfo.Null {
	case "YES":
		convInfo.IsNullable = true
	default:
		convInfo.IsNullable = false
	}

	return convInfo
}

func NewDBExplorer(conn *sql.DB) (http.Handler, error) {
	tablesMap := map[string]map[string]FieldMetaInfo{}
	tablesSlice := []string{}

	rows, err := conn.Query(`SHOW TABLES;`)
	if err != nil {
		return nil, err
	}

	var tmpTableName string
	for rows.Next() {
		rows.Scan(&tmpTableName)
		tablesMap[tmpTableName] = map[string]FieldMetaInfo{}
		tablesSlice = append(tablesSlice, tmpTableName)
	}
	rows.Close()
	log.Println(tablesMap)

	for key := range tablesMap {
		stmt, err := conn.Prepare(fmt.Sprintf("SHOW COLUMNS FROM %s;", key))
		if err != nil {
			return nil, err
		}

		rows, err := stmt.Query()
		if err != nil {
			return nil, err
		}
		stmt.Close()
		for rows.Next() {
			tmpMetaInfo := DbFieldMetaInfo{}
			rows.Scan(&tmpMetaInfo.Field, &tmpMetaInfo.Type, &tmpMetaInfo.Null, &tmpMetaInfo.Key, &tmpMetaInfo.Default, &tmpMetaInfo.Extra)

			tablesMap[key][tmpMetaInfo.Field] =  ConvertDbFieldMetaInfo(tmpMetaInfo)
		}
		rows.Close()
	}
	log.Println(tablesMap)

	router := http.NewServeMux()
	routes := []Route{
		NewRoute(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(Body{
				"response": Body{
					"tables": tablesSlice,
				},
			})
		}),
		NewRoute(http.MethodGet, "/([a-z]+)", func(w http.ResponseWriter, r *http.Request) {
			vars := r.Context().Value(ReqValue{}).([]string)
			tableName := vars[0]

			if _, ok := tablesMap[tableName]; !ok {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(Body{
					"error": "unknown table",
				})
				return
			}


		}),
	}
	sort.Slice(routes, func (i, j int) bool {
		return len(routes[i].RegExp.String()) > len(routes[j].RegExp.String())
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			matches := route.RegExp.FindStringSubmatch(r.URL.Path)
			if matches != nil && route.Method == r.Method {
				ctx := context.WithValue(r.Context(), ReqValue{}, matches[1:])
				route.Handler(w, r.WithContext(ctx))
				return
			}
		}
		http.NotFound(w, r)
	})

	return router, nil
}


func NewRoute(method, path string, handler http.HandlerFunc) Route {
	return Route{
		Method: method,
		RegExp: regexp.MustCompile(path),
		Handler: handler,
	}
}
