package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type ApplicationModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string `bson:"_id",json:"_id"`
	Name          string `bson:"Name",json:"Name"`
	LandingURL    string `bson:"LandingURL",json:"LandingURL"`
}

func NewApplicationModel() *ApplicationModel {
	m := new(ApplicationModel)
	m.ID = bson.NewObjectId().Hex()
	return m
}

func (e *ApplicationModel) RecordID() interface{} {
	return e.ID
}

func (s *ApplicationModel) TableName() string {
	return "applications"
}
