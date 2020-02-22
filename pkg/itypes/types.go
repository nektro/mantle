package itypes

type APIResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}
