package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SharedAgendaModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Name          string
	RAG           string
	Leads         []string
	CreatedDate   time.Time
	CreatedBy     string
	UpdatedDate   time.Time
	UpdatedBy     string
	// Reference
	SCId               string
	SCName             string
	BDId               string
	BusinessDriverName string
	Seq                float64
	IsDeleted          bool
}

func NewSharedAgendaModel() *SharedAgendaModel {
	m := new(SharedAgendaModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SharedAgendaModel) RecordID() interface{} {
	return e.Id
}

func (s *SharedAgendaModel) TableName() string {
	return "SharedAgenda"
}
