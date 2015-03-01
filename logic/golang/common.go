package golang

import (
	"../../model"
	"strings"
)

func (g *generator) toPublic(value string) string {
	return strings.ToUpper(value[:1]) + value[1:]
}

func (g *generator) toCamelcase(field string) string {
	var buf = make([]byte, 0, len(field))
	nextUpper := false
	for i := 0; i < len(field); i++ {
		val := field[i]
		if val == '_' {
			nextUpper = true
			continue
		}
		if nextUpper {
			nextUpper = false
			buf = append(buf, strings.ToUpper(string(val))...)
		} else {
			buf = append(buf, val)
		}
	}
	return string(buf)
}

func (g *generator) toGolangType(fieldType string) string {
	if strings.HasPrefix(fieldType, "varchar") {
		return "string"
	}
	if strings.HasPrefix(fieldType, "text") {
		return "string"
	}
	if strings.HasPrefix(fieldType, "long") {
		return "int64"
	}
	if strings.HasPrefix(fieldType, "float") {
		return "float32"
	}
	if strings.HasPrefix(fieldType, "double") {
		return "float64"
	}
	return fieldType
}

func (g *generator) generateKeyParams(table *model.Table) string {
	body := ""
	for i, id := range table.PrimaryKeys {
		field := table.FieldMap[id]
		if i > 0 {
			body = body + ", "
		}
		body = body + field.Name + " " + field.Type
	}
	return body
}

func (g *generator) generateKeyArgs(table *model.Table) string {
	body := ""
	for i, id := range table.PrimaryKeys {
		field := table.FieldMap[id]
		if i > 0 {
			body = body + ", "
		}
		body = body + field.Name
	}
	return body
}

func (g *generator) generateParams(table *model.Table) string {
	body := ""
	for i, field := range table.Fields {
		if i > 0 {
			body = body + ", "
		}
		nextIndex := i + 1
		if nextIndex < len(table.Fields) {
			next := table.Fields[nextIndex]
			if field.Type == next.Type {
				body = body + field.Name
				continue
			}
		}
		body = body + field.Name + " " + field.Type
	}
	return body
}

func (g *generator) generateArgs(table *model.Table) string {
	body := ""
	for i, field := range table.Fields {
		if i > 0 {
			body = body + ", "
		}
		body = body + field.Name
	}
	return body
}
