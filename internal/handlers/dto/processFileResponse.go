package handlers

type APIResponse struct {
	ErrorMsg string `json:"errorMsg,omitempty"`
	Success  bool   `json:"success"`
}
