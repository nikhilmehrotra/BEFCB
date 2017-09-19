package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type BEFSponsorModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Name          string
	Created_By    string
	Created_Date  time.Time
	Updated_By    string
	Updated_Date  time.Time
	IsDeleted     bool
}

func NewBEFSponsor() *BEFSponsorModel {
	m := new(BEFSponsorModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *BEFSponsorModel) RecordID() interface{} {
	return e.Id
}

func (m *BEFSponsorModel) TableName() string {
	return "BEFSponsor"
}
