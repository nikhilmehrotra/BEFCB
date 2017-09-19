package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type InitiativeOwnerModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	LoginID       string
	Name          string
	Created_By    string
	Created_Date  time.Time
	Updated_By    string
	Updated_Date  time.Time
	IsDeleted     bool
}

func NewInitiativeOwner() *InitiativeOwnerModel {
	m := new(InitiativeOwnerModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *InitiativeOwnerModel) RecordID() interface{} {
	return e.Id
}

func (m *InitiativeOwnerModel) TableName() string {
	return "InitiativeOwner"
}
