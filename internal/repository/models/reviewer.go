package models

type Reviewer struct {
	ReviewerID          uint64 `gorm:"primaryKey;"`
	DisplayName         string `gorm:"size:100"`
	CountryID           int    `gorm:"index"`
	CountryName         string `gorm:"size:100"`
	FlagName            string `gorm:"size:10"`
	ReviewGroupID       int
	ReviewGroupName     string `gorm:"size:50"`
	RoomTypeID          int
	RoomTypeName        string `gorm:"size:100"`
	ReviewedCount       int
	IsExpert            bool
	IsShowGlobalIcon    bool
	IsShowReviewedCount bool

	// Reviews []Review `gorm:"foreignKey:ReviewerID"`
}

func (Reviewer) TableName() string {
	return "reviewers"
}
