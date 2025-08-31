package models

import "time"

type Review struct {
	ReviewID int64 `gorm:"primaryKey"`

	HotelID          int64 `gorm:"index"`
	ProviderID       int
	Rating           float64
	ReviewTitle      string `gorm:"size:255"`
	ReviewComments   string `gorm:"type:text"`
	ReviewPositives  string `gorm:"type:text"`
	ReviewNegatives  string `gorm:"type:text"`
	CheckinMonthYear string `gorm:"size:20"`
	ReviewDate       time.Time
	ResponseText     string `gorm:"type:text"`
	ResponseDate     *time.Time
	ResponderName    string `gorm:"size:100"`
	TranslateSource  string `gorm:"size:10"`
	TranslateTarget  string `gorm:"size:10"`

	ReviewerID int64

	Hotel Hotel `gorm:"foreignKey:HotelID"`
	//Reviewer Reviewer `gorm:"foreignKey:ReviewerID"`
}

func (Review) TableName() string {
	return "reviews"
}
