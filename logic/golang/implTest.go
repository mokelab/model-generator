package golang

import (
	"../../model"
	"fmt"
	"os"
)

func (g *generator) generateImplTest(table, goTable *model.Table) {
	typeName := g.toPublic(g.toCamelcase(table.Name))

	fileName := fmt.Sprintf("output/model/impl/%s_test.go", table.Name)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to open %s : %s\n", fileName, err)
		return
	}

	body := fmt.Sprintf("package impl\n\n"+
		"import (\n"+
		"\tm \"../\"\n"+
		"\t\"testing\"\n"+
		"\t\"database/sql\"\n"+
		")\n\n%s\n%s",
		g.generateTestUtils(table, goTable, table.Name, typeName),
		g.generateAllTest(table, goTable, table.Name, typeName),
	)

	_, err = file.WriteString(body)
	if err != nil {
		fmt.Printf("Failed to Write : %s\n", err)
		return
	}
	file.Close()
}

func (g *generator) generateTestUtils(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func create%sDAO(db *sql.DB) *%sDAO{\n"+
		"\treturn New%sDAO(NewConnection(db))\n"+
		"}\n\n",
		typeName, name, typeName)
	body = body + fmt.Sprintf("func assert%s(t *testing.T, item *m.%s, %s) {\n",
		typeName, typeName, g.generateParams(goTable))
	for _, field := range goTable.Fields {
		body = body + fmt.Sprintf("\tif item.%s != %s {\n"+
			"\t\tt.Errorf(\"%s must be %%s but %%s\", %s, item.%s)\n"+
			"\t}\n",
			g.toPublic(field.Name), field.Name,
			field.Name, field.Name, g.toPublic(field.Name))
	}
	body = body + "}\n\n"
	// HardDelete
	body = body + fmt.Sprintf("func hardDelete%s(db *sql.DB, %s) {\n"+
		"\ts, _ := db.Prepare(\"DELETE FROM %s WHERE %s\")\n"+
		"\tdefer s.Close()\n"+
		"\ts.Exec(%s)\n"+
		"}\n",
		typeName, g.generateKeyParams(goTable),
		name, g.generateKeyWhere(table),
		g.generateKeyArgs(goTable),
	)

	return body
}

func (g *generator) generateAllTest(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func Test%s_All(t *testing.T) {\n",
		typeName)
	body = body + "\tdb, err := connect()\n" +
		"\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to connect\")\n" +
		"\t\treturn\n" +
		"\t}\n" +
		"\tdefer db.Close()\n\n"
	// create DAO
	body = body + fmt.Sprintf("\tdao := create%sDAO(db)\n\n", typeName)
	// Create
	body = body + g.generateTestValues(goTable)
	body = body + fmt.Sprintf("\thardDelete%s(db, %s)\n",
		typeName, g.generateKeyArgs(goTable))
	body = body + fmt.Sprintf("\titem, err := dao.Create(%s)\n",
		g.generateArgs(goTable))
	body = body + "\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to Create : %s\", err)\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + fmt.Sprintf("\tassert%s(t, &item, %s)\n\n",
		typeName, g.generateTestAssertValues(goTable))
	// Get
	body = body + fmt.Sprintf("\titem2, err := dao.Get(%s)\n",
		g.generateKeyArgs(goTable),
	)
	body = body + "\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to Get : %s\", err)\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + fmt.Sprintf("\tassert%s(t, &item2, %s)\n\n",
		typeName, g.generateTestAssertValues(goTable))
	// Update
	body = body + g.generateTestUpdateValues(goTable)
	body = body + fmt.Sprintf("\titem3, err := dao.Update(%s)\n",
		g.generateArgs(goTable),
	)
	body = body + "\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to Update : %s\", err)\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + fmt.Sprintf("\tassert%s(t, &item3, %s)\n\n",
		typeName, g.generateTestAssertUpdateValues(goTable))
	// Update check
	body = body + fmt.Sprintf("\titem4, err := dao.Get(%s)\n",
		g.generateKeyArgs(goTable),
	)
	body = body + "\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to Get : %s\", err)\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + fmt.Sprintf("\tassert%s(t, &item4, %s)\n\n",
		typeName, g.generateTestAssertUpdateValues(goTable))
	// Delete
	body = body + fmt.Sprintf("\titem5, err := dao.Delete(%s)\n",
		g.generateKeyArgs(goTable),
	)
	body = body + "\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to Delete : %s\", err)\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + "\tif !item5.IsEmpty() {\n" +
		"\t\tt.Errorf(\"Returned item must be empty\")\n" +
		"\t\treturn\n" +
		"\t}\n"
	// Delete check
	body = body + fmt.Sprintf("\titem6, err := dao.Get(%s)\n",
		g.generateKeyArgs(goTable),
	)
	body = body + "\tif err != nil {\n" +
		"\t\tt.Errorf(\"Failed to Get : %s\", err)\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + "\tif !item6.IsEmpty() {\n" +
		"\t\tt.Errorf(\"Item must be empty\")\n" +
		"\t\treturn\n" +
		"\t}\n"
	body = body + "}\n"
	return body
}

func (g *generator) generateTestValues(table *model.Table) string {
	body := ""
	for _, field := range table.Fields {
		body = body + fmt.Sprintf("\tvar %s %s = %s\n",
			field.Name, field.Type, g.generateTestValue(field))
	}
	return body
}

func (g *generator) generateTestUpdateValues(table *model.Table) string {
	body := ""
	for _, field := range table.Fields {
		if table.PrimaryKeySet[field.Name] {
			continue
		}
		body = body + fmt.Sprintf("\t%s = %s\n",
			field.Name, g.generateTestUpdateValue(field))
	}
	return body
}

func (g *generator) generateTestAssertValues(table *model.Table) string {
	body := ""
	for i, field := range table.Fields {
		if i > 0 {
			body = body + ", "
		}
		body = body + g.generateTestValue(field)
	}
	return body
}

func (g *generator) generateTestAssertUpdateValues(table *model.Table) string {
	body := ""
	for i, field := range table.Fields {
		if i > 0 {
			body = body + ", "
		}
		if table.PrimaryKeySet[field.Name] {
			body = body + g.generateTestValue(field)
		} else {
			body = body + g.generateTestUpdateValue(field)
		}
	}
	return body
}

func (g *generator) generateTestValue(field *model.Field) string {
	if field.Type == "string" {
		return fmt.Sprintf("\"%s\"", field.Name)
	}
	if field.Type == "int" {
		return "100"
	}
	if field.Type == "long" {
		return "2000"
	}
	if field.Type == "float32" {
		return "12.34"
	}
	if field.Type == "float64" {
		return "345.6789"
	}
	return ""
}

func (g *generator) generateTestUpdateValue(field *model.Field) string {
	if field.Type == "string" {
		return fmt.Sprintf("\"New%s\"", field.Name)
	}
	if field.Type == "int" {
		return "200"
	}
	if field.Type == "long" {
		return "4000"
	}
	if field.Type == "float32" {
		return "34.56"
	}
	if field.Type == "float64" {
		return "678.5432"
	}
	return ""
}
