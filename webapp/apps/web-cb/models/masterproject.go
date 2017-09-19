package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type ProjectModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	ProjectId     int           ` bson:"Id" , json:"Id" `
	ProjectName   string        ` bson:"ProjectName" , json:"ProjectName" `
	Desc          string        ` bson:"Desc" , json:"Desc" `
}

func NewProjectModel() *ProjectModel {
	m := new(ProjectModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *ProjectModel) RecordID() interface{} {
	return e.Id
}

func (m *ProjectModel) TableName() string {
	return "MasterProject"
}
