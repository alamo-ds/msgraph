package graph

import "time"

type Task struct {
	OdataEtag                string                `json:"@odata.etag"`
	PlanID                   string                `json:"planId"`
	BucketID                 string                `json:"bucketId"`
	Title                    string                `json:"title"`
	OrderHint                string                `json:"orderHint"`
	AssigneePriority         string                `json:"assigneePriority"`
	PercentComplete          int                   `json:"percentComplete"`
	StartDateTime            time.Time             `json:"startDateTime"`
	CreatedDateTime          time.Time             `json:"createdDateTime"`
	DueDateTime              time.Time             `json:"dueDateTime"`
	HasDescription           bool                  `json:"hasDescription"`
	PreviewType              string                `json:"previewType"`
	CompletedDateTime        time.Time             `json:"completedDateTime"`
	ReferenceCount           int                   `json:"referenceCount"`
	ChecklistItemCount       int                   `json:"checklistItemCount"`
	ActiveChecklistItemCount int                   `json:"activeChecklistItemCount"`
	ConversationThreadID     string                `json:"conversationThreadId"`
	Priority                 int                   `json:"priority"`
	ID                       string                `json:"id"`
	CreatedBy                IdentitySet           `json:"createdBy"`
	CompletedBy              IdentitySet           `json:"completedBy"`
	AppliedCategories        map[string]bool       `json:"appliedCategories"`
	Assignments              map[string]Assignment `json:"assignments"`
}

type Assignment struct {
	OrderHint        string      `json:"orderHint"`
	AssignedBy       IdentitySet `json:"assignedBy"`
	AssignedDateTime time.Time   `json:"assignedDateTime"`
}

type TaskDetails struct {
	OdataEtag   string                       `json:"@odata.etag"`
	ID          string                       `json:"id"`
	Description string                       `json:"description"`
	PreviewType string                       `json:"previewType"`
	Checklist   map[string]ChecklistItem     `json:"checklist"`
	References  map[string]ExternalReference `json:"references"`
}

type ChecklistItem struct {
	Title                string      `json:"title"`
	IsChecked            bool        `json:"isChecked"`
	OrderHint            string      `json:"orderHint"`
	LastModifiedBy       IdentitySet `json:"lastModifiedBy"`
	LastModifiedDateTime time.Time   `json:"lastModifiedByDateTime"`
}

type ExternalReference struct {
	Alias                string      `json:"alias"`
	LastModifiedBy       IdentitySet `json:"lastModifiedBy"`
	LastModifiedDateTime time.Time   `json:"lastModifiedByDateTime"`
	PreviewPriority      string      `json:"previewPriority"`
	Type                 string      `json:"type"`
}
