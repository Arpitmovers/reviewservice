package services

import (
	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	"github.com/Arpitmovers/reviewservice/internal/repository/models"
)

func MapHotelReviewDTOToModels(dto *dto.HotelReviewDTO) (*models.Hotel,
	*models.Provider, *models.Reviewer, *models.Review, *models.ProviderSummary) {
	hotel := &models.Hotel{
		HotelID:   dto.HotelId,
		HotelName: dto.HotelName,
		Platform:  dto.Platform,
	}

	// Assuming one provider per record from comment
	provider := &models.Provider{
		ProviderID:   dto.Comment.ProviderId,
		ProviderName: "",
	}

	// Map reviewer info
	reviewerDTO := dto.Comment.ReviewerInfo
	reviewer := &models.Reviewer{
		DisplayName:           reviewerDTO.DisplayMemberName,
		CountryID:             reviewerDTO.CountryId,
		CountryName:           reviewerDTO.CountryName,
		ReviewGroupID:         reviewerDTO.ReviewGroupId,
		ReviewGroupName:       reviewerDTO.ReviewGroupName,
		ReviewerReviewedCount: reviewerDTO.ReviewerReviewedCount,
		IsExpertReviewer:      reviewerDTO.IsExpertReviewer,
	}

	review := &models.Review{
		ReviewID:                dto.Comment.HotelReviewId,
		HotelID:                 dto.HotelId,
		ProviderID:              dto.Comment.ProviderId,
		Rating:                  dto.Comment.Rating,
		CheckInMonthYear:        dto.Comment.CheckInDateMonthYear,
		EncryptedReviewData:     dto.Comment.EncryptedReviewData,
		FormattedRating:         dto.Comment.FormattedRating,
		FormattedReviewDate:     dto.Comment.FormattedReviewDate,
		RatingText:              dto.Comment.RatingText,
		ResponderName:           dto.Comment.ResponderName,
		ResponseDateText:        dto.Comment.ResponseDateText,
		ResponseText:            "", // This can be set if response content exists
		ResponseTranslateSource: dto.Comment.ResponseTranslateSource,
		ReviewComments:          dto.Comment.ReviewComments,
		ReviewNegatives:         dto.Comment.ReviewNegatives,
		ReviewPositives:         dto.Comment.ReviewPositives,
		ReviewProviderLogo:      dto.Comment.ReviewProviderLogo,
		ReviewProviderText:      dto.Comment.ReviewProviderText,
		ReviewTitle:             dto.Comment.ReviewTitle,
		TranslateSource:         dto.Comment.TranslateSource,
		TranslateTarget:         dto.Comment.TranslateTarget,
		ReviewDate:              dto.Comment.ReviewDate,
		OriginalTitle:           dto.Comment.OriginalTitle,
		OriginalComment:         dto.Comment.OriginalComment,
		FormattedResponseDate:   dto.Comment.FormattedResponseDate,
		RoomType:                reviewerDTO.RoomTypeName,
		LengthOfStay:            reviewerDTO.LengthOfStay,
	}

	// Map provider summary using first overall provider for simplicity
	var providerSummary *models.ProviderSummary
	if len(dto.OverallByProviders) > 0 {
		p := dto.OverallByProviders[0]
		providerSummary = &models.ProviderSummary{
			HotelID:       dto.HotelId,
			ProviderID:    p.ProviderId,
			OverallScore:  p.OverallScore,
			ReviewCount:   p.ReviewCount,
			Cleanliness:   p.Grades["Cleanliness"],
			Facilities:    p.Grades["Facilities"],
			Location:      p.Grades["Location"],
			Service:       p.Grades["Service"],
			ValueForMoney: p.Grades["Value for money"],
			RoomComfort:   p.Grades["Room comfort and quality"],
		}
	}

	return hotel, provider, reviewer, review, providerSummary
}
