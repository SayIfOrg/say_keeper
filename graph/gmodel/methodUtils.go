package gmodel

import (
	"github.com/SayIfOrg/say_keeper/models"
	"github.com/SayIfOrg/say_keeper/utils"
	"strconv"
)

func FromDBComment(dc *models.Comment) *Comment {
	return &Comment{
		ID:        strconv.Itoa(int(dc.ID)),
		UserID:    strconv.Itoa(int(dc.UserID)),
		ReplyToID: utils.RUintToString(dc.ReplyToId),
		Content:   dc.Content,
		Agent:     dc.Agent,
	}
}

func FromDBUser(du *models.User) *User {
	return &User{
		ID:   strconv.Itoa(int(du.ID)),
		Name: du.Name,
	}
}
