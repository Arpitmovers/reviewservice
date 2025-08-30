package handlers

import (
	"encoding/json"
	"fmt"

	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	services "github.com/Arpitmovers/reviewservice/internal/service"
)

func (h *ReviewHandler) ConsumeReview() func([]byte) error {

	return func(body []byte) error {
		fmt.Printf("Received review message: %s\n", string(body))
		var review *dto.Review

		if err := json.Unmarshal(body, &review); err != nil {
			return fmt.Errorf("failed to unmarshal review json: %w", err)
		}
		// TODO: decode, validate, save to DB etc.

		services.SaveReview(review)
		// if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 	// Save Review
		// 	if err := tx.Create(&review).Error; err != nil {
		// 		return fmt.Errorf("failed to save review: %w", err)
		// 	}
		// 	return nil
		// }); err != nil {
		// 	return err
		// }
		// return nil
	}
}
