package models

type ProviderScore struct {
	ID            int   `gorm:"primaryKey;autoIncrement"`
	HotelID       int64 `gorm:"index"`
	ProviderID    int
	ProviderName  string `gorm:"size:100"`
	OverallScore  float64
	ReviewCount   int
	Cleanliness   float64
	Facilities    float64
	Location      float64
	Service       float64
	ValueForMoney float64

	// Hotel Hotel `gorm:"foreignKey:HotelID"`
}

func (ProviderScore) TableName() string {
	return "provider_scores"
}
