package model

import (
	rj "github.com/fkmhrk-go/rawjson"
)

type TestInfo struct {
	DBType string
	DBDSN  string
}

func ParseTestInfo(json rj.RawJsonObject) (*TestInfo, error) {
	dbType, err := json.String("db_type")
	if err != nil {
		return nil, err
	}
	dbDSN, err := json.String("db_dsn")
	if err != nil {
		return nil, err
	}
	return &TestInfo{
		DBType: dbType,
		DBDSN:  dbDSN,
	}, nil
}
