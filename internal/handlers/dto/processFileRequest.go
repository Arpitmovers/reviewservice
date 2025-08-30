package handlers

type ProcessReviewRequest struct {
	Bucket     string `json:"bucket" validate:"required,min=3,max=63"`
	PathPrefix string `json:"pathPrefix" validate:"required"`
	Force      bool   `json:"force"`
}
