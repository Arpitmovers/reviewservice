package models

type ProviderSummary struct {
	HotelID       uint64  `gorm:"primaryKey;column:hotel_id"`
	ProviderID    int     `gorm:"primaryKey;column:provider_id"`
	OverallScore  float32 `gorm:"type:decimal(3,1)"`
	ReviewCount   int
	Cleanliness   float32 `gorm:"type:decimal(3,1)"`
	Facilities    float32 `gorm:"type:decimal(3,1)"`
	Location      float32 `gorm:"type:decimal(3,1)"`
	Service       float32 `gorm:"type:decimal(3,1)"`
	ValueForMoney float32 `gorm:"type:decimal(3,1)"`
	RoomComfort   float32 `gorm:"type:decimal(3,1)"`

	// Associations
	Hotel    Hotel    `gorm:"foreignKey:HotelID;references:HotelID"`
	Provider Provider `gorm:"foreignKey:ProviderID;references:ProviderID"`
}

func (ProviderSummary) TableName() string {
	return "ProviderSummary"
}
