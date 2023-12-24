package domains

import "time"

type (
	CreateUserRequest struct {
		Name string `json:"name"`
	}

	CreateUserResponse struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	User struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		AccountIDs []string  `json:"account_ids,omitempty"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	GetUsersResponse struct {
		Users []*User `json:"users"`
	}
)
