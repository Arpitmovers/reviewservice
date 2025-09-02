package models

type Hotel struct {
	HotelID   uint64 `gorm:"primaryKey;column:hotel_id"`
	HotelName string `gorm:"size:500;not null"`
	Platform  string `gorm:"size:100;not null"`
}

func (Hotel) TableName() string {
	return "Hotels"
}
