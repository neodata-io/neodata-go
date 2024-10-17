package model

// Ability defines the permissions for a user on the frontend.
type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
}
