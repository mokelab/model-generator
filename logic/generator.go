package logic

import (
	"../model"
	"./golang"
)

type Generator interface {
	Generate(tables []*model.Table)
}

func CreateGenerator() Generator {
	return golang.CreateGenerator()
}
