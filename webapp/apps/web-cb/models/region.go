package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type RegionModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Major_Region  string        ` bson:"Major_Region" , json:"Major_Region" `
	Region        string        ` bson:"Region" , json:"Region" `
	Country       string        ` bson:"Country" , json:"Country" `
	CountryCode   string        ` bson:"CountryCode" , json:"CountryCode" `
	IsDeleted     bool
}

func NewRegion() *RegionModel {
	m := new(RegionModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *RegionModel) RecordID() interface{} {
	return e.Id
}

func (m *RegionModel) TableName() string {
	return "Region"
}
