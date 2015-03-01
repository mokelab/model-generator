# ModelGenerator

Code generator.

## How to use

1. Write table definitions in table.json

    {
        "tables":[
            {
                "name" : "user",
                "fields" : [
                    {
                        "name" : "id",
                        "type" : "varchar(32)"
                    },
                    {
                        "name" : "login_name",
                        "type" : "text"
                    },
                    {
                        "name" : "password",
                        "type" : "varchar(32)"
                    }
                ],
                "primary_keys" : [
                    "id"
                ]
            }
        ],
        "db" : {
            "type" : "mysql"
        },
        "test" : {
            "db_type" : "mysql",
            "db_dsn" : "testUser:pass@/testDatabase"
        }
    }

2. Execute the following command.

    $ go run main.go -file=table.json

3. Copy generated codes(output folder) to your project.