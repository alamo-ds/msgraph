package graph

type Bucket struct {
	OdataEtag string `json:"@odata.etag"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	OrderHint string `json:"orderHint"`
	PlanID    string `json:"planId"`
}
