package controllers

import (
	"eaciit/scb-apps/webapp/apps/main/models"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/knot/knot.v1"
	"gopkg.in/mgo.v2/bson"
)

func (c *AccessController) GetGroup(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	csr, err := acl.Find(new(models.GroupModel), nil, nil)
	defer csr.Close()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]models.GroupModel, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(data)
}

func (c *AccessController) SelectGroup(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.GroupModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	err = acl.FindByID(payload, payload.ID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(payload)
}

func (c *AccessController) DeleteGroup(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.GroupModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	if payload.IsImportant {
		return c.SetResultError("Cannot delete this group", nil)
	}

	err = acl.Delete(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(nil)
}

func (c *AccessController) SaveGroup(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.GroupModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	if payload.IsImportant {
		return c.SetResultError("Cannot update this group", nil)
	}

	if payload.ID == "" {
		payload.ID = bson.NewObjectId().Hex()
		payload.Enable = true
	}

	err = acl.Save(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	err = c.FixGroupAccessMenuByGroupID(payload.ID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(payload)
}

func (c *AccessController) SaveGroups(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new([]models.GroupModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	for _, each := range *payload {
		if each.IsImportant {
			return c.SetResultError("Cannot update this group", nil)
		}

		if each.ID == "" {
			each.ID = bson.NewObjectId().Hex()
			each.Enable = true
		}

		err = acl.Save(&each)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}

		err = c.FixGroupAccessMenuByGroupID(each.ID)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}
	}

	return c.SetResultOK(nil)
}

func (c *AccessController) FixGroupAccessMenuByGroupID(groupID string) error {
	group := new(models.GroupModel)
	group.ID = groupID
	err := acl.FindByID(group, group.ID)
	if err != nil {
		return err
	}

	for _, app := range group.Applications {
		dataAccessMenu, err := c.GetAccessMenuByApplicationID(app)
		if err != nil {
			return err
		}

		group.Grants = make([]models.AccessGrant, 0)
		for _, each := range dataAccessMenu {
			group.Grants = append(group.Grants, models.AccessGrant{
				AccessID:      each.ID,
				AccessValue:   15,
				ApplicationID: each.ApplicationID,
			})
		}
	}

	err = acl.Save(group)
	if err != nil {
		return err
	}

	return nil
}
