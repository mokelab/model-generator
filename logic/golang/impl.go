package golang

import (
	"../../model"
	"fmt"
	"os"
)

func (g *generator) generateImpl(table, goTable *model.Table) {
	typeName := g.toPublic(table.Name)

	fileName := fmt.Sprintf("output/model/impl/%s.go", table.Name)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to open %s : %s\n", fileName, err)
		return
	}

	body := fmt.Sprintf("package impl\n\n"+
		"import (\n"+
		"\tm \"../\"\n"+
		"\t\"database/sql\"\n"+
		")\n\n%s\n%s\n%s\n%s\n%s\n%s",
		g.generateConstructor(table.Name, typeName),
		g.generateCreateImpl(table, goTable, table.Name, typeName),
		g.generateGetImpl(table, goTable, table.Name, typeName),
		g.generateUpdateImpl(table, goTable, table.Name, typeName),
		g.generateDeleteImpl(table, goTable, table.Name, typeName),
		g.generateScan(table, goTable, table.Name, typeName))

	_, err = file.WriteString(body)
	if err != nil {
		fmt.Printf("Failed to Write : %s\n", err)
		return
	}
	file.Close()
}

func (g *generator) generateConstructor(name, typeName string) string {
	return fmt.Sprintf("type %sDAO struct {\n"+
		"\tconnection *Connection\n"+
		"}\n\n"+
		"func New%sDAO(connection *Connection) *%sDAO {\n"+
		"\treturn &%sDAO{\n"+
		"\t\tconnection: connection,\n"+
		"\t}\n"+
		"}\n",
		name, typeName, name, name)
}

func (g *generator) generateCreateImpl(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func (d *%sDAO) %s {\n",
		name, g.generateCreateSignature("m."+typeName, goTable))
	body = body + g.generateBegin()
	body = body + fmt.Sprintf("\n\tst, err := tr.Prepare(\"INSERT INTO %s(%s) VALUES(%s)\")\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n"+
		"\tdefer st.Close()\n\n",
		name, g.generateRows(table), g.generatePlaceholder(table))
	// exec
	body = body + fmt.Sprintf("\t_, err = st.Exec(%s)\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n\n",
		g.generateKeyArgs(goTable))
	// commit
	body = body + "\ttr.Commit()\n\n"
	// return
	body = body + g.generateReturn(typeName, goTable, ", nil\n")
	body = body + "}\n"
	return body
}

func (g *generator) generateGetImpl(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func (d *%sDAO) %s {\n",
		name, g.generateGetSignature("m."+typeName, goTable))
	body = body + g.generateConnect()
	body = body + fmt.Sprintf("\tst, err := db.Prepare(\"SELECT %s FROM %s WHERE %s AND deleted <> 1\")\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n"+
		"\tdefer st.Close()\n\n",
		g.generateRows(table), name, g.generateKeyWhere(table))
	// execute
	body = body + fmt.Sprintf("\trows, err := st.Query(%s)\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n"+
		"\tdefer rows.Close()\n\n",
		g.generateKeyArgs(goTable))
	body = body + "\tif !rows.Next() {\n" +
		"\t\treturn nil, nil\n" +
		"\t}\n\n"
	// return
	body = body + "\treturn d.scan(rows), nil\n"
	body = body + "}\n"
	return body
}

func (g *generator) generateUpdateImpl(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func (d *%sDAO) %s {\n",
		name, g.generateUpdateSignature("m."+typeName, goTable))
	body = body + g.generateConnect()
	// query
	body = body + fmt.Sprintf("\tst, err := db.Prepare(\"UPDATE %s SET %s WHERE %s AND deleted <> 1\")\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n"+
		"\tdefer st.Close()\n\n",
		name, g.generateUpdateSet(table), g.generateKeyWhere(table))
	// exec
	body = body + fmt.Sprintf("\t_, err = st.Exec(%s, %s)\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n\n",
		g.generateArgs(goTable), g.generateKeyArgs(goTable))
	// return
	body = body + g.generateReturn(typeName, goTable, ", nil\n")
	body = body + "}\n"
	return body
}

func (g *generator) generateDeleteImpl(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func (d *%sDAO) %s {\n",
		name, g.generateDeleteSignature("m."+typeName, goTable))
	body = body + g.generateConnect()
	// query
	body = body + fmt.Sprintf("\tst, err := db.Prepare(\"UPDATE %s SET deleted=1 WHERE %s AND deleted <> 1\")\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n"+
		"\tdefer st.Close()\n\n",
		name, g.generateKeyWhere(table))
	// exec
	body = body + fmt.Sprintf("\t_, err = st.Exec(%s)\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, err\n"+
		"\t}\n\n",
		g.generateKeyArgs(goTable))
	// return
	body = body + "\treturn nil, nil\n"
	body = body + "}\n"
	return body
}

func (g *generator) generateScan(table, goTable *model.Table, name, typeName string) string {
	body := fmt.Sprintf("func (d *%sDAO) scan(rows *sql.Rows) *m.%s {\n",
		name, typeName)
	for _, field := range goTable.Fields {
		body = body + fmt.Sprintf("\tvar %s %s\n", field.Name, field.Type)
	}
	// args
	args := ""
	for i, field := range goTable.Fields {
		if i > 0 {
			args = args + ", "
		}
		args = args + "&" + field.Name
	}
	body = body + fmt.Sprintf("\trows.Scan(%s)\n", args)
	// return
	body = body + g.generateReturn(typeName, goTable, "\n")
	body = body + "}\n"
	return body
}

//----

func (g *generator) generateConnect() string {
	return "\tdb := d.connection.Connect()\n"
}

func (g *generator) generateBegin() string {
	return "\ttr, err := d.connection.Begin()\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tdefer tr.Rollback()\n"
}

func (g *generator) generateRows(table *model.Table) string {
	body := ""
	for i, field := range table.Fields {
		if i > 0 {
			body = body + ","
		}
		body = body + field.Name
	}
	return body
}

func (g *generator) generatePlaceholder(table *model.Table) string {
	body := ""
	for i, _ := range table.Fields {
		if i > 0 {
			body = body + ","
		}
		body = body + "?"
	}
	return body
}

func (g *generator) generateReturn(typeName string, table *model.Table, postfix string) string {
	body := fmt.Sprintf("\treturn &m.%s {\n", typeName)
	for _, field := range table.Fields {
		body = body + fmt.Sprintf("\t\t%s: %s,\n",
			g.toPublic(field.Name), field.Name)
	}
	body = body + "\t}" + postfix
	return body
}

func (g *generator) generateUpdateSet(table *model.Table) string {
	body := ""
	for i, field := range table.Fields {
		if i > 0 {
			body = body + ", "
		}
		body = body + field.Name + "=?"
	}
	return body
}

func (g *generator) generateKeyWhere(table *model.Table) string {
	body := ""
	for i := 0; i < len(table.PrimaryKeys); i++ {
		field := table.FieldMap[table.PrimaryKeys[i]]
		if i > 0 {
			body = body + " AND "
		}
		body = body + field.Name + "=?"
	}
	return body
}
