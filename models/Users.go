package models

type Users struct {
	Id       int    `gorm:"primaryKey" json:"id"`
	NamaUser string `gorm:"type:varchar(191)" json:"nama_user"`
	Email    string `gorm:"type:varchar(191)" json:"email"`
	Password string `gorm:"type:varchar(191)" json:"password"`
}
