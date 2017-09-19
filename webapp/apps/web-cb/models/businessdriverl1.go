package models

import (
	"github.com/eaciit/orm"
	"time"
)

type BusinessDriverL1Model struct {
	orm.ModelBase  `bson:"-",json:"-"`
	Id             int ` bson:"_id" , json:"_id" `
	Idx            string
	Name           string
	Description    string
	UpdatedBy      string
	UpdatedDate    time.Time
	BusinessMetric []BusinessMetric
	Seq            int
}

func NewBusinessDriverL1Model() *BusinessDriverL1Model {
	m := new(BusinessDriverL1Model)
	return m
}

func (e *BusinessDriverL1Model) RecordID() interface{} {
	return e.Id
}

func (m *BusinessDriverL1Model) TableName() string {
	return "BusinessDriverL1"
}
