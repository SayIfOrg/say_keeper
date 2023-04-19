package models

import "gorm.io/gorm"

//MigrateSchema migrates the schema with preserving data
// use it in early development stage and then use migration files
func MigrateSchema(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Comment{})
	return err
}
