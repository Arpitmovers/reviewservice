package handlers

type Review struct {
	HotelID   int64  `json:"hotelId"`
	Platform  string `json:"platform"`
	HotelName string `json:"hotelName"`
	//Comment           Comment             `json:"comment"`
	Comment           Comment             `gorm:"-"`
	OverallByProvider []OverallByProvider `gorm:"-"`
	//OverallByProvider []OverallByProvider `json:"overallByProviders"`
}

type Comment struct {
	IsShowReviewResponse    bool    `json:"isShowReviewResponse"`
	HotelReviewID           int64   `json:"hotelReviewId"`
	ProviderID              int     `json:"providerId"`
	Rating                  float64 `json:"rating"`
	CheckInDateMonthAndYear string  `json:"checkInDateMonthAndYear"`
	EncryptedReviewData     string  `json:"encryptedReviewData"`
	FormattedRating         string  `json:"formattedRating"`
	FormattedReviewDate     string  `json:"formattedReviewDate"`
	RatingText              string  `json:"ratingText"`
	ResponderName           string  `json:"responderName"`
	ResponseDateText        string  `json:"responseDateText"`
	ResponseTranslateSource string  `json:"responseTranslateSource"`
	ReviewComments          string  `json:"reviewComments"`
	ReviewNegatives         string  `json:"reviewNegatives"`
	ReviewPositives         string  `json:"reviewPositives"`
	ReviewProviderLogo      string  `json:"reviewProviderLogo"`
	ReviewProviderText      string  `json:"reviewProviderText"`
	ReviewTitle             string  `json:"reviewTitle"`
	TranslateSource         string  `json:"translateSource"`
	TranslateTarget         string  `json:"translateTarget"`
	ReviewDate              string  `json:"reviewDate"`
	//	ReviewerInfo            ReviewerInfo `json:"reviewerInfo"`
	ReviewerInfo          ReviewerInfo `gorm:"type:json"`
	OriginalTitle         string       `json:"originalTitle"`
	OriginalComment       string       `json:"originalComment"`
	FormattedResponseDate string       `json:"formattedResponseDate"`
}

type ReviewerInfo struct {
	CountryName           string `json:"countryName"`
	DisplayMemberName     string `json:"displayMemberName"`
	FlagName              string `json:"flagName"`
	ReviewGroupName       string `json:"reviewGroupName"`
	RoomTypeName          string `json:"roomTypeName"`
	CountryID             int    `json:"countryId"`
	LengthOfStay          int    `json:"lengthOfStay"`
	ReviewGroupID         int    `json:"reviewGroupId"`
	RoomTypeID            int    `json:"roomTypeId"`
	ReviewerReviewedCount int    `json:"reviewerReviewedCount"`
	IsExpertReviewer      bool   `json:"isExpertReviewer"`
	IsShowGlobalIcon      bool   `json:"isShowGlobalIcon"`
	IsShowReviewedCount   bool   `json:"isShowReviewedCount"`
}

type OverallByProvider struct {
	ProviderID   int     `json:"providerId"`
	Provider     string  `json:"provider"`
	OverallScore float64 `json:"overallScore"`
	ReviewCount  int     `json:"reviewCount"`
	//Grades       map[string]float64 `json:"grades"`
}
