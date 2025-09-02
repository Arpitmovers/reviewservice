package models

import "time"

// CREATE TABLE Reviews (
// 	review_id BIGINT PRIMARY KEY,
// 	hotel_id BIGINT unsigned NOT NULL,
// 	provider_id INT NOT NULL,
// 	reviewer_id BIGINT unsigned NOT NULL,
// 	rating DECIMAL(3,1),
// 	check_in_month_year VARCHAR(20),
// 	review_date TIMESTAMP,
// 	review_title VARCHAR(200),
// 	review_text TEXT,
// 	response_text TEXT,
// 	room_type VARCHAR(100),
// 	length_of_stay INT,
// 	positives TEXT,
// 	negatives TEXT,
// 	encrypted_data TEXT,
// 	FOREIGN KEY (hotel_id) REFERENCES Hotels(hotel_id),
// 	 FOREIGN KEY (provider_id) REFERENCES Providers(provider_id),
//  FOREIGN KEY (reviewer_id) REFERENCES Reviewers(reviewer_id)
// );

type Review struct {
	ReviewID   uint64 `gorm:"primaryKey;column:review_id"`
	HotelID    uint64 `gorm:"not null;index"`
	ProviderID int    `gorm:"not null;index"`
	ReviewerID int    `gorm:"not null;index"`
	//IsShowReviewResponse    bool
	Rating                  float32 `gorm:"type:decimal(3,1)"`
	CheckInMonthYear        string  `gorm:"size:20"`
	EncryptedReviewData     string  `gorm:"type:text"`
	FormattedRating         string  `gorm:"size:10"`
	FormattedReviewDate     string  `gorm:"size:50"`
	RatingText              string  `gorm:"size:100"`
	ResponderName           string  `gorm:"size:255"`
	ResponseDateText        string  `gorm:"size:50"`
	ResponseText            string  `gorm:"type:text"`
	ResponseTranslateSource string  `gorm:"size:10"`
	ReviewComments          string  `gorm:"type:text"`
	ReviewNegatives         string  `gorm:"type:text"`
	ReviewPositives         string  `gorm:"type:text"`
	ReviewProviderLogo      string  `gorm:"size:255"`
	ReviewProviderText      string  `gorm:"size:100"`
	ReviewTitle             string  `gorm:"size:255"`
	TranslateSource         string  `gorm:"size:10"`
	TranslateTarget         string  `gorm:"size:10"`
	ReviewDate              time.Time
	OriginalTitle           string `gorm:"size:255"`
	OriginalComment         string `gorm:"type:text"`
	FormattedResponseDate   string `gorm:"size:50"`
	RoomType                string `gorm:"size:100"`
	LengthOfStay            int
	PositiveCount           int `gorm:"default:0"`
	NegativeCount           int `gorm:"default:0"`

	// Associations
	Hotel    Hotel    `gorm:"foreignKey:HotelID;references:HotelID"`
	Provider Provider `gorm:"foreignKey:ProviderID;references:ProviderID"`
	Reviewer Reviewer `gorm:"foreignKey:ReviewerID;references:ReviewerID"`
}

func (Review) TableName() string {
	return "Reviews"
}
