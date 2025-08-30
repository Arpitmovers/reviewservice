package models

type Hotel struct {
	HotelID   int64  `gorm:"primaryKey;column:hotel_id"`
	Platform  string `gorm:"size:50;index"`
	HotelName string `gorm:"size:255"`

	Reviews        []Review        `gorm:"foreignKey:HotelID"`
	ProviderScores []ProviderScore `gorm:"foreignKey:HotelID"`
}

func (Hotel) TableName() string {
	return "hotels"
}
