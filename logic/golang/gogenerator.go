package golang

import (
	"../../model"
	"fmt"
	"os"
)

type generator struct {
}

func CreateGenerator() *generator {
	return &generator{}
}

func (g *generator) Generate(tables []*model.Table) {
	// create output
	err := os.RemoveAll("output")
	if err != nil {
		if os.IsExist(err) {
			return
		}
		fmt.Printf("Failed to remove output folder : %s", err)
	}
	err = os.Mkdir("output", 0755)
	if err != nil {
		fmt.Printf("Failed to create output folder : %s", err)
		return
	}
	// generate model folder
	err = os.Mkdir("output/model", 0755)
	if err != nil {
		fmt.Printf("Failed to create output/model folder : %s", err)
		return
	}

	// convert
	for _, table := range tables {
		goTable := g.convertTable(table)
		g.generateType(goTable)
	}
}
