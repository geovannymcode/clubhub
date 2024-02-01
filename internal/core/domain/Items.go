package domain

// Items holds a collection of ItemServe, primarily for JSON serialization.
type Items struct {
	Items []ItemServe `json:"items"`
}
