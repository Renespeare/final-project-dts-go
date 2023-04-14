package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	GormModel
	Title			string	`gorm:"not null" json:"title" form:"title" valid:"required~Your title is required"`
	Caption		string	`json:"caption" form:"caption"`
	PhotoUrl	string	`gorm:"not null" json:"photo_url" form:"photo_url" valid:"required~Your photo is required"`
	UserID		uint		`json:"user_id,omitempty"`
	User			*User		`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}
	err = nil
	return
}