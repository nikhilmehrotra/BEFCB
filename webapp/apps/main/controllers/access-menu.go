package controllers

import (
	"eaciit/scb-apps/webapp/apps/main/models"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/knot/knot.v1"
	"gopkg.in/mgo.v2/bson"
)

func (c *AccessController) GetAccessMenu(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	csr, err := acl.Find(new(models.AccessMenuModel), nil, nil)
	defer csr.Close()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]models.AccessMenuModel, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(data)
}

func (c *AccessController) GetAccessMenuForCurrentUser(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	// if NoLogin mode, then return all available access for dev purpose
	// if c.GetCurrentUsername(k) == "" && c.NoLogin {
	if c.GetCurrentUsername(k) == "" {
		return c.GetAccessMenu(k)
	}

	return c.SetResultOK(LoginDataAccessMenu)
}

func (c *AccessController) SelectAccessMenu(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.AccessMenuModel)
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

func (c *AccessController) DeleteAccessMenu(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.AccessMenuModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	err = acl.Delete(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	// ====== inject the new changes to grants "ADMIN"

	err = c.FixAdminGroupAccessMenu()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(nil)
}

func (c *AccessController) SaveAccessMenu(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.AccessMenuModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	if payload.ID == "" {
		payload.ID = bson.NewObjectId().Hex()
	}

	err = acl.Save(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	// ====== inject the new changes to grants "ADMIN"

	err = c.FixAdminGroupAccessMenu()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(payload)
}

func (c *AccessController) GetAccessMenuByApplication(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := struct{ ApplicationID string }{""}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	dataFlat, err := c.GetAccessMenuByApplicationID(payload.ApplicationID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(dataFlat)
}

func (c *AccessController) FixAdminGroupAccessMenu() error {
	group := new(models.GroupModel)
	group.ID = "admin"
	err := acl.FindByID(group, group.ID)
	if err != nil {
		return err
	}

	csrAccessMenu, err := acl.Find(new(models.AccessMenuModel), nil, nil)
	defer csrAccessMenu.Close()
	if err != nil {
		return err
	}
	dataAccessMenu := make([]models.AccessMenuModel, 0)
	err = csrAccessMenu.Fetch(&dataAccessMenu, 0, false)
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
	err = acl.Save(group)
	if err != nil {
		return err
	}

	return nil
}
