package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SponsorModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Name          string
}

func NewSponsor() *SponsorModel {
	m := new(SponsorModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *SponsorModel) RecordID() interface{} {
	return e.Id
}

func (m *SponsorModel) TableName() string {
	return "Sponsor"
}
