package golang

import (
	"../../model"
	"fmt"
	"os"
)

func (g *generator) generateType(table *model.Table) {
	typeName := g.toPublic(table.Name)

	fileName := fmt.Sprintf("output/model/%s.go", table.Name)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to open %s : %s\n", fileName, err)
		return
	}

	body := fmt.Sprintf("package model\n\n%s\n%s",
		g.generateInterface(typeName, table),
		g.generateEntity(typeName, table),
	)

	_, err = file.WriteString(body)
	if err != nil {
		fmt.Printf("Failed to Write : %s\n", err)
		return
	}
	file.Close()
}

func (g *generator) generateInterface(typeName string, table *model.Table) string {
	body := fmt.Sprintf("type %sDAO interface{\n", typeName)
	body = body + "\t" + g.generateCreateSignature(typeName, table) + "\n"
	body = body + "\t" + g.generateGetSignature(typeName, table) + "\n"
	body = body + "\t" + g.generateUpdateSignature(typeName, table) + "\n"
	body = body + "\t" + g.generateDeleteSignature(typeName, table) + "\n"
	body = body + "}\n"
	return body
}

func (g *generator) generateCreateSignature(typeName string, table *model.Table) string {
	return "Create(" + g.generateParams(table) + ") " +
		"(*" + typeName + ", error)"
}

func (g *generator) generateGetSignature(typeName string, table *model.Table) string {
	return "Get(" + g.generateKeyParams(table) + ") " +
		"(*" + typeName + ", error)"
}

func (g *generator) generateUpdateSignature(typeName string, table *model.Table) string {
	return "Update(" + g.generateParams(table) + ") " +
		"(*" + typeName + ", error)"
}

func (g *generator) generateDeleteSignature(typeName string, table *model.Table) string {
	return "Delete(" + g.generateKeyParams(table) + ") " +
		"(*" + typeName + ", error)"
}

func (g *generator) generateEntity(typeName string, table *model.Table) string {
	body := fmt.Sprintf("type %s struct{\n", typeName)
	for _, field := range table.Fields {
		body = body + "\t" + g.toPublic(field.Name) + " " + field.Type + "\n"
	}
	body = body + "}\n"
	return body
}
