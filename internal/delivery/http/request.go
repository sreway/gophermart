package http

type (
	credentialsRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	withdrawRequest struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}
)
