package couchbase

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/couchbase/gocb.v1"
)

const (
	readUserDelay        = 100
	userNameProperty     = "user_name"
	userPasswordProperty = "user_password"
	userRolesProperty    = "user_roles"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: createUser,
		Read:   readUser,
		Update: updateUser,
		Delete: deleteUser,
		Schema: map[string]*schema.Schema{
			userNameProperty: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			userPasswordProperty: {
				Type:     schema.TypeString,
				Required: true,
			},
			userRolesProperty: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createUser(data *schema.ResourceData, meta interface{}) (err error) {
	doCreateUser(data, meta)
	return doReadUser(data, meta, readUserDelay)
}

func readUser(data *schema.ResourceData, meta interface{}) (err error) {
	doReadUser(data, meta, 0)
	return
}

func updateUser(data *schema.ResourceData, meta interface{}) (err error) {
	doCreateUser(data, meta)
	return doReadUser(data, meta, readUserDelay)
}

func deleteUser(data *schema.ResourceData, meta interface{}) (err error) {
	return doDropUser(data, meta)
}

func doReadUser(data *schema.ResourceData, meta interface{}, delay int) (err error) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	_, manager, err := connect(meta.(*Config))
	if err != nil {
		return
	}
	userName := data.Get(userNameProperty).(string)
	_, err = manager.GetUser(gocb.LocalDomain, userName)
	if err != nil {
		log.Printf("[WARN] Can not find user %q", userName)
		data.SetId("")
	}
	return
}

func doDropUser(data *schema.ResourceData, meta interface{}) (err error) {
	_, manager, err := connect(meta.(*Config))
	if err != nil {
		return
	}
	userName := data.Get(userNameProperty).(string)
	err = manager.RemoveUser(gocb.LocalDomain, userName)
	if err == nil {
		log.Printf("[INFO] User with the name %q was removed", userName)
		data.SetId("")
	}
	return
}

func doCreateUser(data *schema.ResourceData, meta interface{}) (err error) {
	_, manager, err := connect(meta.(*Config))
	if err != nil {
		return
	}
	userRoles := strings.Split(data.Get(userRolesProperty).(string), ",")
	roles := make([]gocb.UserRole, len(userRoles))
	for i, role := range userRoles {
		tokens := strings.Split(role, ":")
		roles[i] = gocb.UserRole{tokens[0], tokens[1]}
	}
	userName := data.Get(userNameProperty).(string)
	userPassword := data.Get(userPasswordProperty).(string)
	err = manager.UpsertUser(gocb.LocalDomain, userName, &gocb.UserSettings{Password: userPassword, Roles: roles})
	if err == nil {
		log.Printf("[INFO] User with the name %q was created", userName)
		data.SetId(userName)
	}
	return
}
