package models

import (
	// . "eaciit/scb-bef/webapp/bef/models"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	// "time"
)

type LogCommentModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	IdInitiative  bson.ObjectId
	Comment       Comment
	Action        string
}

func NewLogCommentModel() *LogCommentModel {
	m := new(LogCommentModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *LogCommentModel) RecordID() interface{} {
	return e.Id
}

func (s *LogCommentModel) TableName() string {
	return "LogComment"
}
