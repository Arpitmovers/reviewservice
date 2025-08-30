package services

import (
	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	"github.com/Arpitmovers/reviewservice/internal/repository/models"
)

func mapToHotel(msg *dto.Review) models.Hotel {
	return models.Hotel{
		HotelID:   msg.HotelID,
		Platform:  msg.Platform,
		HotelName: msg.HotelName,
	}
}

func mapToReviewer(msg *dto.Review) models.Reviewer {
	return models.Reviewer{
		//	ReviewerID:          msg.Comment.ReviewerInfo.ReviewerID,   // uint64
		DisplayName:     msg.Comment.ReviewerInfo.DisplayMemberName,
		CountryID:       msg.Comment.ReviewerInfo.CountryID,
		CountryName:     msg.Comment.ReviewerInfo.CountryName,
		FlagName:        msg.Comment.ReviewerInfo.FlagName,
		ReviewGroupID:   msg.Comment.ReviewerInfo.ReviewGroupID,
		ReviewGroupName: msg.Comment.ReviewerInfo.ReviewGroupName,
		RoomTypeID:      msg.Comment.ReviewerInfo.RoomTypeID,
		RoomTypeName:    msg.Comment.ReviewerInfo.RoomTypeName,
		ReviewedCount:   msg.Comment.ReviewerInfo.ReviewerReviewedCount,

		IsExpert:            msg.Comment.ReviewerInfo.IsExpertReviewer,
		IsShowGlobalIcon:    msg.Comment.ReviewerInfo.IsShowGlobalIcon,
		IsShowReviewedCount: msg.Comment.ReviewerInfo.IsShowReviewedCount,
	}
}

func mapToReview(msg *dto.Review) models.Review {
	return models.Review{
		ReviewID:       msg.Comment.HotelReviewID,
		HotelID:        msg.HotelID,
		ReviewerID:     msg.Comment.HotelReviewID,
		Rating:         msg.Comment.Rating,
		ReviewTitle:    msg.Comment.ReviewTitle,
		ReviewComments: msg.Comment.OriginalComment,
	}
}

func mapToProviderScore(msg *dto.Review) models.ProviderScore {
	if len(msg.OverallByProvider) == 0 {
		return models.ProviderScore{}
	}

	provider := msg.OverallByProvider[0]

	score := models.ProviderScore{
		HotelID:      msg.HotelID,
		ProviderID:   provider.ProviderID,
		OverallScore: provider.OverallScore,
		ReviewCount:  provider.ReviewCount,
	}

	// if len(provider.Grades) > 0 {
	// 	if v, ok := provider.Grades["Cleanliness"]; ok {
	// 		score.Cleanliness = v
	// 	}
	// 	if v, ok := provider.Grades["Facilities"]; ok {
	// 		score.Facilities = v
	// 	}
	// 	if v, ok := provider.Grades["Location"]; ok {
	// 		score.Location = v
	// 	}
	// 	if v, ok := provider.Grades["Service"]; ok {
	// 		score.Service = v
	// 	}
	// 	if v, ok := provider.Grades["Value for money"]; ok {
	// 		score.ValueForMoney = v
	// 	}
	// }

	return score
}
