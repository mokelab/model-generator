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

func (g *generator) Generate(tables []*model.Table, testInfo *model.TestInfo) {
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
	err = g.generateOutputFolders()
	if err != nil {
		return
	}

	g.generateConnection()
	g.generateTestConnection(testInfo)
	for _, table := range tables {
		// convert
		goTable := g.convertTable(table)
		g.generateType(goTable)
		g.generateImpl(table, goTable)
		g.generateImplTest(table, goTable)
	}
}

func (g *generator) generateOutputFolders() error {
	err := os.Mkdir("output/model", 0755)
	if err != nil {
		fmt.Printf("Failed to create output/model folder : %s", err)
		return err
	}
	err = os.Mkdir("output/model/impl", 0755)
	if err != nil {
		fmt.Printf("Failed to create output/model/impl folder : %s", err)
		return err
	}
	return nil
}

func (g *generator) generateConnection() {
	fileName := "output/model/impl/connection.go"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to open %s : %s\n", fileName, err)
		return
	}

	body := "package impl\n\n" +
		"import (\n" +
		"\t\"database/sql\"\n" +
		"\t_ \"github.com/go-sql-driver/mysql\"\n" +
		")\n\n" +
		"type Connection struct {\n" +
		"\tdb *sql.DB\n" +
		"}\n\n" +
		"func NewConnection(db *sql.DB) *Connection {\n" +
		"\treturn &Connection{\n" +
		"\t\tdb: db,\n" +
		"\t}\n" +
		"}\n\n" +
		"func (c *Connection) Connect() *sql.DB {\n" +
		"\treturn c.db\n" +
		"}\n\n" +
		"func (c *Connection) Begin() (*sql.Tx, error) {\n" +
		"\treturn c.db.Begin()\n" +
		"}\n"
	_, err = file.WriteString(body)
	if err != nil {
		fmt.Printf("Failed to Write : %s\n", err)
		return
	}
	file.Close()
}

func (g *generator) generateTestConnection(testInfo *model.TestInfo) {
	fileName := "output/model/impl/connection_test.go"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to open %s : %s\n", fileName, err)
		return
	}

	body := fmt.Sprintf("package impl\n\n"+
		"import (\n"+
		"\t\"database/sql\"\n"+
		"\t_ \"github.com/go-sql-driver/mysql\"\n"+
		")\n\n"+
		"func connect() (*sql.DB, error) {\n"+
		"\treturn sql.Open(\"%s\", \"%s\")\n"+
		"}\n\n"+
		"func begin() (*sql.DB, *sql.Tx, error) {\n"+
		"\tdb, err := connect()\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, nil, err\n"+
		"\t}\n\n"+
		"\ttx, err := db.Begin()\n"+
		"\tif err != nil {\n"+
		"\t\tdb.Close()\n"+
		"\t\treturn nil, nil, err\n"+
		"\t}\n\n"+
		"\treturn db, tx, nil\n"+
		"}\n",
		testInfo.DBType, testInfo.DBDSN)
	_, err = file.WriteString(body)
	if err != nil {
		fmt.Printf("Failed to Write : %s\n", err)
		return
	}
	file.Close()
}
