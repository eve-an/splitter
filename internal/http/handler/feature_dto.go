package handler

type variantPayload struct {
	Name   string `json:"name"`
	Weight uint8  `json:"weight"`
}

type featureRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Active      bool             `json:"active"`
	Variants    []variantPayload `json:"variants"`
}

type eventRequest struct {
	UserID  string `json:"user_id"`
	Variant string `json:"variant"`
	Type    string `json:"type"`
}
