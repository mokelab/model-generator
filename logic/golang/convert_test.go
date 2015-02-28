package golang

import (
	"../../model"
	"testing"
)

func TestConvert_0000_convertTable(t *testing.T) {
	g := &generator{}
	table := &model.Table{
		Name: "user",
		Fields: []*model.Field{
			&model.Field{
				Name: "id",
				Type: "varchar(32)",
			},
			&model.Field{
				Name: "user_name",
				Type: "text",
			},
		},
	}
	out := g.convertTable(table)
	if out.Name != "user" {
		t.Errorf("Unexpected table name : %s", out.Name)
	}
	if len(out.Fields) != 2 {
		t.Errorf("Unexpected field count : %d", len(out.Fields))
	}
	field := out.Fields[0]
	if field.Name != "id" {
		t.Errorf("Unexpected field name[0] : %s", field.Name)
	}
	if field.Type != "string" {
		t.Errorf("Unexpected field type[0] : %s", field.Type)
	}
	field = out.Fields[1]
	if field.Name != "userName" {
		t.Errorf("Unexpected field name[1] : %s", field.Name)
	}
	if field.Type != "string" {
		t.Errorf("Unexpected field type[1] : %s", field.Type)
	}
	field = out.FieldMap["user_name"]
	if field.Name != "userName" {
		t.Errorf("Unexpected field name[user_name] : %s", field.Name)
	}
	if field.Type != "string" {
		t.Errorf("Unexpected field type[user_name] : %s", field.Type)
	}
}
