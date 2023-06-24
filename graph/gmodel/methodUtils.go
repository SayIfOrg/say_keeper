package gmodel

import (
	"encoding/json"
	"github.com/SayIfOrg/say_keeper/models"
	"github.com/SayIfOrg/say_keeper/utils"
	"strconv"
)

func FromDBComment(dc *models.Comment) *Comment {
	return &Comment{
		ID:        strconv.Itoa(int(dc.ID)),
		UserID:    strconv.Itoa(int(dc.UserID)),
		ReplyToID: utils.SqlItptS(dc.ReplyToId),
		Content:   dc.Content,
		Agent:     dc.Agent,
	}
}

func FromDBUser(du *models.User) *User {
	return &User{
		ID: strconv.Itoa(int(du.ID)),
	}
}

func UnmarshalComment(b []byte) (*Comment, error) {
	obj := &Comment{}
	err := json.Unmarshal(b, obj)
	return obj, err
}
