package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const groupResource string = "groups"

func (c *Client) Groups() *GroupsRequestBuilder {
	return &GroupsRequestBuilder{
		c:    c,
		path: joinPath(c.BaseURL, groupResource),
	}
}

type GroupsRequestBuilder struct {
	c    *Client
	path string
}

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

type GetGroupsResponse struct {
	Value []Group `json:"value"`
}

func (r *GroupsRequestBuilder) Get(ctx context.Context) ([]Group, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.path, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %v", err)
	}
	resp, err := r.c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request returned %d: %s", resp.StatusCode, readForError(resp.Body))
	}

	var ret GetGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}

	return ret.Value, nil
}

type GroupItemRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *GroupsRequestBuilder) ById(id string) *GroupItemRequestBuilder {
	return &GroupItemRequestBuilder{
		Id:   id,
		c:    r.c,
		path: joinPath(r.path, id),
	}
}

func (r *GroupItemRequestBuilder) Get(ctx context.Context) (Group, error) {
	var group Group

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.path, nil)
	if err != nil {
		return group, fmt.Errorf("couldn't create request: %v", err)
	}
	resp, err := r.c.do(req)
	if err != nil {
		return group, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return group, fmt.Errorf("request returned %d: %s", resp.StatusCode, readForError(resp.Body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return group, fmt.Errorf("json.Decode: %v", err)
	}

	return group, nil
}

func readForError(r io.Reader) string {
	data, err := io.ReadAll(r)
	if err != nil {
		return "couldn't read body"
	}

	return string(data)
}
