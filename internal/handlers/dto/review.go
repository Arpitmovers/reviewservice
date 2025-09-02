package handlers

import "time"

type ReviewerInfoDTO struct {
	CountryName           string `json:"countryName"`
	DisplayMemberName     string `json:"displayMemberName"`
	FlagName              string `json:"flagName"`
	ReviewGroupName       string `json:"reviewGroupName"`
	RoomTypeName          string `json:"roomTypeName"`
	CountryId             int    `json:"countryId"`
	LengthOfStay          int    `json:"lengthOfStay"`
	ReviewGroupId         int    `json:"reviewGroupId"`
	ReviewerReviewedCount int    `json:"reviewerReviewedCount"`
	IsExpertReviewer      bool   `json:"isExpertReviewer"`
	IsShowReviewedCount   bool   `json:"isShowReviewedCount"`
}

type CommentDTO struct {
	IsShowReviewResponse    bool            `json:"isShowReviewResponse"`
	HotelReviewId           uint64          `json:"hotelReviewId"`
	ProviderId              int             `json:"providerId"`
	Rating                  float32         `json:"rating"`
	CheckInDateMonthYear    string          `json:"checkInDateMonthAndYear"`
	EncryptedReviewData     string          `json:"encryptedReviewData"`
	FormattedRating         string          `json:"formattedRating"`
	FormattedReviewDate     string          `json:"formattedReviewDate"`
	RatingText              string          `json:"ratingText"`
	ResponderName           string          `json:"responderName"`
	ResponseDateText        string          `json:"responseDateText"`
	ResponseTranslateSource string          `json:"responseTranslateSource"`
	ReviewComments          string          `json:"reviewComments"`
	ReviewNegatives         string          `json:"reviewNegatives"`
	ReviewPositives         string          `json:"reviewPositives"`
	ReviewProviderLogo      string          `json:"reviewProviderLogo"`
	ReviewProviderText      string          `json:"reviewProviderText"`
	ReviewTitle             string          `json:"reviewTitle"`
	TranslateSource         string          `json:"translateSource"`
	TranslateTarget         string          `json:"translateTarget"`
	ReviewDate              time.Time       `json:"reviewDate"`
	ReviewerInfo            ReviewerInfoDTO `json:"reviewerInfo"`
	OriginalTitle           string          `json:"originalTitle"`
	OriginalComment         string          `json:"originalComment"`
	FormattedResponseDate   string          `json:"formattedResponseDate"`
}

type OverallByProviderDTO struct {
	ProviderId   int                `json:"providerId"`
	Provider     string             `json:"provider"`
	OverallScore float32            `json:"overallScore"`
	ReviewCount  int                `json:"reviewCount"`
	Grades       map[string]float32 `json:"grades"`
}

type HotelReviewDTO struct {
	HotelId            uint64                 `json:"hotelId"`
	Platform           string                 `json:"platform"`
	HotelName          string                 `json:"hotelName"`
	Comment            CommentDTO             `json:"comment"`
	OverallByProviders []OverallByProviderDTO `json:"overallByProviders"`
}
