package handlers

type ProcessReviewRequest struct {
	PathPrefix string `json:"pathPrefix" validate:"required"`
}
