package graph

type Identity struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

type IdentitySet struct {
	User        Identity `json:"user"`
	Application Identity `json:"application"`
}

type Recipient struct {
	EmailAddress EmailAddress `json:"emailAddress"`
}

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
