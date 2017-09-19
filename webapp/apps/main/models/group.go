package models

import (
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/toolkit"
)

type AccessGrant struct {
	AccessID      string
	AccessValue   int
	ApplicationID string
}

type GroupModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string `json:"_id",bson:"_id"`
	Title         string
	Grants        []AccessGrant
	Applications  []string
	IsImportant   bool

	Owner      string
	Enable     bool
	GroupType  acl.GroupTypeEnum
	GroupConf  toolkit.M
	MemberConf toolkit.M
}

func (g *GroupModel) TableName() string {
	return "acl_groups"
}

func (g *GroupModel) RecordID() interface{} {
	return g.ID
}

func (g *GroupModel) PreSave() error {
	if g.ID == "" {
		g.ID = toolkit.RandomString(32)
	}
	return nil
}
