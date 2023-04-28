package models

import (
	"database/sql"
	"github.com/SayIfOrg/say_keeper/commenting"
)

func (c *Comment) PopulateIdentifier(outerID string) {
	if c.Agent != commenting.WebAgent && outerID == "" {
		panic("for none web agent `outerId` should be provided")
	}
	if c.Agent == commenting.WebAgent && outerID != "" {
		panic("only for none web agent `outerId` should be provided")
	}
	if c.Agent != commenting.WebAgent && outerID != "" {
		c.Identifier = sql.NullString{String: c.Agent + "-" + outerID}
	}
}
