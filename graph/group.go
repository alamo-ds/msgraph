package graph

import "time"

type Group struct {
	ID                            string      `json:"id"`
	DeletedDateTime               time.Time   `json:"deletedDateTime"`
	Classification                string      `json:"classification"`
	CreatedDateTime               time.Time   `json:"createdDateTime"`
	Description                   string      `json:"description"`
	DisplayName                   string      `json:"displayName"`
	ExpirationDateTime            time.Time   `json:"expirationDateTime"`
	GroupTypes                    []string    `json:"groupTypes"`
	IsAssignableToRole            bool        `json:"isAssignableToRole"`
	Mail                          string      `json:"mail"`
	MailEnabled                   bool        `json:"mailEnabled"`
	MailNickname                  string      `json:"mailNickname"`
	MembershipRule                string      `json:"membershipRule"`
	MembershipRuleProcessingState string      `json:"membershipRuleProcessingState"`
	OnPremisesDomainName          string      `json:"onPremisesDomainName"`
	OnPremisesLastSyncDateTime    time.Time   `json:"onPremisesLastSyncDateTime"`
	OnPremisesNetBiosName         string      `json:"onPremisesNetBiosName"`
	OnPremisesSamAccountName      string      `json:"onPremisesSamAccountName"`
	OnPremisesSecurityIdentifier  string      `json:"onPremisesSecurityIdentifier"`
	OnPremisesSyncEnabled         bool        `json:"onPremisesSyncEnabled"`
	PreferredDataLocation         string      `json:"preferredDataLocation"`
	PreferredLanguage             string      `json:"preferredLanguage"`
	ProxyAddresses                []string    `json:"proxyAddresses"`
	RenewedDateTime               time.Time   `json:"renewedDateTime"`
	ResourceBehaviorOptions       []string    `json:"resourceBehaviorOptions"`
	ResourceProvisioningOptions   []string    `json:"resourceProvisioningOptions"`
	SecurityEnabled               bool        `json:"securityEnabled"`
	SecurityIdentifier            string      `json:"securityIdentifier"`
	Theme                         string      `json:"theme"`
	UniqueName                    string      `json:"uniqueName"`
	Visibility                    string      `json:"visibility"`
	OnPremisesProvisioningErrors  []string    `json:"onPremisesProvisioningErrors"`
	ServiceProvisioningErrors     []time.Time `json:"serviceProvisioningErrors"`
}
