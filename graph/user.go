package graph

import (
	"strings"
	"time"
)

// TODO: add other fields. For simplicity, I only added what I needed for
// a specific project.
type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Mail        string `json:"mail"`
	// NOTE: below is what should be used to search User by email, not "mail"
	UserPrincipalName     string    `json:"userPrincipalName"`
	CreationType          string    `json:"creationType"`
	Department            string    `json:"department"`
	GivenName             string    `json:"givenName"`
	AccountEnabled        bool      `json:"accountEnabled"`
	CreatedDateTime       time.Time `json:"createdDateTime,omitzero"`
	DeletedDateTime       time.Time `json:"deletedDateTime,omitzero"`
	EmployeeLeaveDateTime time.Time `json:"employeeLeaveDateTime,omitzero"`
	// TODO: the key matches, but for some reason not getting values. Investigate why.
	AssignedLicenses []AssignedLicense `json:"assignedLicenses"`
}

type AssignedLicense struct {
	SkuID         string   `json:"skuId"`
	DisabledPlans []string `json:"disabledPlans"`
}

func userSelectParams(params []string) string {
	var sb strings.Builder

	params = append([]string{
		"id",
		"displayName",
		"mail",
		"userPrincipalName",
		"accountEnabled",
		"createdDateTime",
		"deletedDateTime",
		"employeeLeaveDateTime",
	}, params...)

	sb.WriteString(strings.Join(params, ","))
	return sb.String()
}
