package dto

// This file for unmarshalling in handler

type AccountCreate struct {
	UserName string `json:"user_name"`
}

type AccountUpdate struct {
	Type   *string `json:"type,omitempty"`
	Status string  `json:"status,omitempty"`
}
