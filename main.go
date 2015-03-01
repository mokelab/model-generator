package main

import (
	"./logic"
	"./model"
	"flag"
	"fmt"
	rj "github.com/fkmhrk-go/rawjson"
	"os"
)

func main() {
	var fileName *string = flag.String("file", "", "JSON file name")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Failed to open file : %s", err)
		return
	}
	fmt.Println(flag.NArg())
	fmt.Printf("File=" + *fileName)
	obj, err := rj.ObjectFromReader(file)
	if err != nil {
		fmt.Printf("Failed to parse JSON : %s", err)
		return
	}
	array, err := obj.Array("tables")
	if err != nil {
		fmt.Printf("Failed to array field \"tables\" : %s", err)
		return
	}
	testInfoObj, err := obj.Object("test")
	if err != nil {
		fmt.Printf("Failed to object field \"test\" : %s", err)
		return
	}
	var tables []*model.Table
	for i := 0; i < len(array); i++ {
		obj, err := array.Object(i)
		if err != nil {
			continue
		}
		table, err := model.ParseTable(obj)
		if err != nil {
			continue
		}
		tables = append(tables, table)
	}
	testInfo, err := model.ParseTestInfo(testInfoObj)
	if err != nil {
		fmt.Printf("Failed to parse \"test\" field : %s", err)
		return
	}
	g := logic.CreateGenerator()
	g.Generate(tables, testInfo)
}
