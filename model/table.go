package model

import (
	rj "github.com/fkmhrk-go/rawjson"
)

type Table struct {
	Name          string
	Fields        []*Field
	FieldMap      map[string]*Field
	PrimaryKeys   []string
	PrimaryKeySet map[string]bool
}

type Field struct {
	Name string
	Type string
}

func ParseTable(json rj.RawJsonObject) (*Table, error) {
	name, err := json.String("name")
	if err != nil {
		return nil, err
	}
	array, err := json.Array("fields")
	if err != nil {
		return nil, err
	}
	keyArray, err := json.Array("primary_keys")
	if err != nil {
		return nil, err
	}
	var fields []*Field
	var fieldMap map[string]*Field = make(map[string]*Field)
	for i := 0; i < len(array); i++ {
		obj, err := array.Object(i)
		if err != nil {
			continue
		}
		name, err := obj.String("name")
		if err != nil {
			continue
		}
		typeName, err := obj.String("type")
		if err != nil {
			continue
		}
		f := &Field{
			Name: name,
			Type: typeName,
		}
		fields = append(fields, f)
		fieldMap[name] = f
	}
	// primary key
	var primaryKeys []string
	var primaryKeySet = make(map[string]bool)
	for i := 0; i < len(keyArray); i++ {
		key, err := keyArray.String(i)
		if err != nil {
			continue
		}
		primaryKeys = append(primaryKeys, key)
		primaryKeySet[key] = true
	}
	return &Table{
		Name:          name,
		Fields:        fields,
		FieldMap:      fieldMap,
		PrimaryKeys:   primaryKeys,
		PrimaryKeySet: primaryKeySet,
	}, nil
}
