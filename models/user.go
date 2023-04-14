package models

import (
	"final-project-dts-go/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username	string		`gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Your username is required"`
	Email			string		`gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required,email~Invalid email format"`
	Age				uint				`gorm:"not null" json:"age" form:"age" valid:"required~Your age is required,in(8|)~Input age more than equal 8"`	
	Password 	string		`gorm:"not null" json:"password,omitempty" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Photos 		[]Photo 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPass(u.Password)
	err = nil
	return
}

