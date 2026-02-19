package graph

type Identity struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

type IdentitySet struct {
	User        Identity `json:"user"`
	Application Identity `json:"application"`
}
