package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SummaryBusinessDriverModel struct {
	orm.ModelBase   `bson:"-",json:"-"`
	Id              bson.ObjectId    ` bson:"_id" , json:"_id" `
	Idx             string           ` bson:"Id" , json:"Id" `
	Name            string           ` bson:"Name" , json:"Name" `
	Seq             int32            ` bson:"Seq" , json:"Seq" `
	Type            string           ` bson:"Type" , json:"Type" `
	Completion      float64          ` bson:"Completion" , json:"Completion" `
	BusinessMetrics []BusinessMetric ` bson:"BusinessMetrics,omitempty" , json:"BusinessMetrics" `
	Parentid        string
	Parentname      string
	Category        string
	SeqParent       float64
}

func NewSummaryBusinessDriverModel() *SummaryBusinessDriverModel {
	m := new(SummaryBusinessDriverModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *SummaryBusinessDriverModel) RecordID() interface{} {
	return e.Id
}

func (m *SummaryBusinessDriverModel) TableName() string {
	return "SummaryBusinessDriver"
}
