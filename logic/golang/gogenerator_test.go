package golang

import (
	"../../model"
	"testing"
)

func TestGenerator_0000_GetSignatue_1(t *testing.T) {
	g := &generator{}
	table := &model.Table{
		FieldMap: map[string]*model.Field{
			"id": &model.Field{
				Name: "id",
				Type: "string",
			},
		},
		PrimaryKeys: []string{
			"id",
		},
	}

	sig := g.generateGetSignature("User", table)
	if sig != "Get(id string) (*User, error)" {
		t.Errorf("Unexpected signature")
	}
}

func TestGenerator_0001_GetSignatue_2(t *testing.T) {
	g := &generator{}
	table := &model.Table{
		FieldMap: map[string]*model.Field{
			"id": &model.Field{
				Name: "id",
				Type: "string",
			},
			"user_id": &model.Field{
				Name: "userId",
				Type: "string",
			},
		},
		PrimaryKeys: []string{
			"id",
			"user_id",
		},
	}

	sig := g.generateGetSignature("Inventory", table)
	if sig != "Get(id string, userId string) (*Inventory, error)" {
		t.Errorf("Unexpected signature %s", sig)
	}
}
