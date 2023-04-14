package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	GormModel
	Message			string	`gorm:"not null" json:"message" form:"message" valid:"required~Your message is required"`
	UserID 			uint		`json:"user_id"`
	User				*User		`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	PhotoID 		uint		`json:"photo_id"`
	Photo				*Photo 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo,omitempty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}
	err = nil
	return
}