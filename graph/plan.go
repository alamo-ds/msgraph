package graph

import "time"

type Plan struct {
	OdataEtag       string        `json:"@odata.etag"`
	CreatedDateTime time.Time     `json:"createdDateTime"`
	Owner           string        `json:"owner"`
	Title           string        `json:"title"`
	ID              string        `json:"id"`
	CreatedBy       IdentitySet   `json:"createdBy"`
	Container       PlanContainer `json:"container"`
}

type PlanContainer struct {
	ContainerID string `json:"containerId"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}
