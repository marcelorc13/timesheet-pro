package domain

type (
	DBResponse struct {
		Success bool
		Message string
		Data    any
	}

	HttpResponse struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}
)
