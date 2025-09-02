package models

type Reviewer struct {
	ReviewerID            int    `gorm:"primaryKey;autoIncrement;column:reviewer_id"`
	DisplayName           string `gorm:"size:100"`
	CountryID             int
	CountryName           string `gorm:"size:100"`
	ReviewGroupID         int
	ReviewGroupName       string `gorm:"size:100"`
	ReviewerReviewedCount int
	IsExpertReviewer      bool
}

func (Reviewer) TableName() string {
	return "Reviewers"
}
