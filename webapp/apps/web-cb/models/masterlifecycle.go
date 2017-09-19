package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type LifeCycleModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	LifeCycleId   string        ` bson:"Id" , json:"Id" `
	Name          string        ` bson:"Name" , json:"Name" `
	Seq           int           ` bson:"Seq" , json:"Seq" `
	SubLC         []SubLifeCycleModel
}

type SubLifeCycleModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string `bson:"Id" , json:"Id" `
	Name          string ` bson:"Name" , json:"Name" `
}

func NewLifeCycleModelModel() *LifeCycleModel {
	m := new(LifeCycleModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *LifeCycleModel) RecordID() interface{} {
	return e.Id
}

func (m *LifeCycleModel) TableName() string {
	return "MasterLifeCycle"
}
