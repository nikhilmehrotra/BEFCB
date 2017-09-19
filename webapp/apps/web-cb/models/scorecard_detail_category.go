package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ScorecardDetailCategoryModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	SdcId         string
	Name          string
	SCId          string
	SCName        string
	Updated_Date  time.Time
	Updated_By    string
	Seq           float64
}

func NewScorecardDetailCategory() *ScorecardDetailCategoryModel {
	m := new(ScorecardDetailCategoryModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *ScorecardDetailCategoryModel) RecordID() interface{} {
	return e.Id
}

func (m *ScorecardDetailCategoryModel) TableName() string {
	return "ScorecardDetailCategory"
}
