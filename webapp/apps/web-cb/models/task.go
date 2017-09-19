package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TaskModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	Id               bson.ObjectId `bson:"_id" , json:"_id" `
	Name             string
	Owner            string
	Statement        string
	Description      string
	DateCreated      time.Time
	DateUpdated      time.Time
	LifeCycleId      string
	BusinessDriverId string
	TaskType         string
	SCCategory       string
	IsGlobal         bool     `bson:"IsGlobal" , json:"IsGlobal" `
	Region           []string `bson:"Region" , json:"Region" `
	Country          []string `bson:"Country" , json:"Country" `
}

func NewTaskModel() *TaskModel {
	m := new(TaskModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *TaskModel) RecordID() interface{} {
	return e.Id
}

func (s *TaskModel) TableName() string {
	return "Task"
}
