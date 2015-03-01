package golang

import (
	"../../model"
)

func (g *generator) convertTable(table *model.Table) *model.Table {
	var goFields []*model.Field
	var goFieldMap = make(map[string]*model.Field)
	for _, field := range table.Fields {
		goField := &model.Field{
			Name: g.toCamelcase(field.Name),
			Type: g.toGolangType(field.Type),
		}
		goFields = append(goFields, goField)
		goFieldMap[field.Name] = goField
	}
	return &model.Table{
		Name:          table.Name,
		Fields:        goFields,
		FieldMap:      goFieldMap,
		PrimaryKeys:   table.PrimaryKeys,
		PrimaryKeySet: table.PrimaryKeySet,
	}
}
