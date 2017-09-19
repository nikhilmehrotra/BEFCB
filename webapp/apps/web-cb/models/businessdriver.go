package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type BusinessDriverModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            		bson.ObjectId ` bson:"_id" , json:"_id" `
	BusinessDriverId 	string ` bson:"Id" , json:"Id" `
	Name 				string ` bson:"Name" , json:"Name" `
	Seq 				int ` bson:"Seq" , json:"Seq" `
	Type 				string ` bson:"Type" , json:"Type" `
}

func NewBusinessDriverModel() *BusinessDriverModel {
	m := new(BusinessDriverModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *BusinessDriverModel) RecordID() interface{} {
	return e.Id
}

func (m *BusinessDriverModel) TableName() string {
	return "MasterBusinessDriver"
}
