package models

type Product struct {
	Id          int    `gorm:"primaryKey" json:"id"`
	NamaProduct string `gorm:"type:varchar(191)" json:"nama_product"`
	Deskripsi   string `gorm:"type:text" json:"deskripsi"`
}