package logic

import (
	"../model"
	"./golang"
)

type Generator interface {
	Generate(tables []*model.Table, testInfo *model.TestInfo)
}

func CreateGenerator() Generator {
	return golang.CreateGenerator()
}
