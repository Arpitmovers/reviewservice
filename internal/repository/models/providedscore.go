package models

type ProviderScore struct {
	ID            int   `gorm:"primaryKey;autoIncrement"`
	HotelID       int64 `gorm:"index"`
	ProviderID    int
	ProviderName  string `gorm:"size:100"`
	OverallScore  float64
	ReviewCount   int
	Cleanliness   float32
	Facilities    float32
	Location      float32
	Service       float32
	ValueForMoney float32

	Hotel Hotel `gorm:"foreignKey:HotelID"`
}

func (ProviderScore) TableName() string {
	return "provider_scores"
}
