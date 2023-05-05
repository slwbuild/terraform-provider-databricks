package access

// REST API: https://docs.databricks.com/dev-tools/api/latest/ip-access-list.html#operation/create-list

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/databricks/databricks-sdk-go/apierr"
	"github.com/databricks/terraform-provider-databricks/qa"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	TestingID               = "234567"
	TestingLabel            = "Naughty"
	TestingListTypeString   = "BLOCK"
	TestingListType         = "BLOCK"
	TestingEnabled          = true
	TestingIPAddresses      = []string{"1.2.3.4", "1.2.4.0/24"}
	TestingIPAddressesState = []any{"1.2.3.4", "1.2.4.0/24"}
)

func TestIPACLCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodPost,
				Resource: "/api/2.0/ip-access-lists",
				ExpectedRequest: createIPAccessListRequest{
					Label:       TestingLabel,
					ListType:    TestingListType,
					IPAddresses: TestingIPAddresses,
				},
				Response: IpAccessListStatusWrapper{
					IPAccessList: IpAccessListStatus{
						ListID:        TestingID,
						Label:         TestingLabel,
						ListType:      TestingListType,
						IPAddresses:   TestingIPAddresses,
						AddressCount:  2,
						CreatedAt:     87939234,
						CreatorUserID: 1234556,
						UpdatedAt:     87939234,
						UpdatorUserID: 1234556,
						Enabled:       TestingEnabled,
					},
				},
			},
			{
				Method:   http.MethodGet,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: IpAccessListStatusWrapper{
					IPAccessList: IpAccessListStatus{
						ListID:        TestingID,
						Label:         TestingLabel,
						ListType:      TestingListType,
						IPAddresses:   TestingIPAddresses,
						AddressCount:  2,
						CreatedAt:     87939234,
						CreatorUserID: 1234556,
						UpdatedAt:     87939234,
						UpdatorUserID: 1234556,
						Enabled:       TestingEnabled,
					},
				},
			},
		},
		Resource: ResourceIPAccessList(),
		State: map[string]any{
			"label":        TestingLabel,
			"list_type":    TestingListTypeString,
			"ip_addresses": TestingIPAddressesState,
		},
		Create: true,
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, TestingID, d.Id())
	assert.Equal(t, TestingLabel, d.Get("label"))
	assert.Equal(t, TestingListTypeString, d.Get("list_type"))
	assert.Equal(t, TestingEnabled, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("ip_addresses.#"))
}

func TestAPIACLCreate_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodPost,
				Resource: "/api/2.0/ip-access-lists",
				Response: apierr.APIErrorBody{
					ErrorCode: "RESOURCE_ALREADY_EXISTS",
					Message:   "IP access list with type (" + TestingListTypeString + ") and label (" + TestingLabel + ") already exists",
				},
				Status: 400,
			},
		},
		Resource: ResourceIPAccessList(),
		State: map[string]any{
			"label":        TestingLabel,
			"list_type":    TestingListTypeString,
			"ip_addresses": TestingIPAddressesState,
		},
		Create: true,
	}.Apply(t)
	assert.Error(t, err)
	qa.AssertErrorStartsWith(t, err, "IP access list with type")
	assert.Equal(t, "", d.Id(), "Id should be empty for error creates")
}

func TestIPACLUpdate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodGet,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: IpAccessListStatusWrapper{
					IPAccessList: IpAccessListStatus{
						ListID:        TestingID,
						Label:         TestingLabel,
						ListType:      TestingListType,
						IPAddresses:   TestingIPAddresses,
						AddressCount:  2,
						CreatedAt:     87939234,
						CreatorUserID: 1234556,
						UpdatedAt:     87939234,
						UpdatorUserID: 1234556,
						Enabled:       TestingEnabled,
					},
				},
			},
			{
				Method:   http.MethodPut,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: IpAccessListStatusWrapper{
					IPAccessList: IpAccessListStatus{
						ListID:        TestingID,
						Label:         TestingLabel,
						ListType:      TestingListType,
						IPAddresses:   TestingIPAddresses,
						AddressCount:  2,
						CreatedAt:     87939234,
						CreatorUserID: 1234556,
						UpdatedAt:     87939234,
						UpdatorUserID: 1234556,
						Enabled:       TestingEnabled,
					},
				},
			},
		},
		Resource: ResourceIPAccessList(),
		State: map[string]any{
			"label":        TestingLabel,
			"list_type":    TestingListTypeString,
			"ip_addresses": TestingIPAddressesState,
		},
		Update: true,
		ID:     TestingID,
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, TestingID, d.Id())
	assert.Equal(t, TestingLabel, d.Get("label"))
	assert.Equal(t, TestingListTypeString, d.Get("list_type"))
	assert.Equal(t, TestingEnabled, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("ip_addresses.#"))
}

func TestIPACLUpdate_Error(t *testing.T) {
	_, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodPut,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: apierr.APIErrorBody{
					ErrorCode: "SERVER_ERROR",
					Message:   "Something unexpected happened",
				},
				Status: 500,
			},
		},
		Resource: ResourceIPAccessList(),
		Update:   true,
		ID:       TestingID,
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Something unexpected")
}

func TestIPACLRead(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodGet,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: IpAccessListStatusWrapper{
					IPAccessList: IpAccessListStatus{
						ListID:        TestingID,
						Label:         TestingLabel,
						ListType:      TestingListType,
						IPAddresses:   TestingIPAddresses,
						AddressCount:  2,
						CreatedAt:     87939234,
						CreatorUserID: 1234556,
						UpdatedAt:     87939234,
						UpdatorUserID: 1234556,
						Enabled:       TestingEnabled,
					},
				},
			},
		},
		Resource: ResourceIPAccessList(),
		Read:     true,
		New:      true,
		ID:       TestingID,
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, TestingID, d.Id())
	assert.Equal(t, TestingLabel, d.Get("label"))
	assert.Equal(t, TestingListTypeString, d.Get("list_type"))
	assert.Equal(t, TestingEnabled, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("ip_addresses.#"))
}

func TestIPACLRead_NotFound(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodGet,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: apierr.APIErrorBody{
					ErrorCode: "RESOURCE_DOES_NOT_EXIST",
					Message:   "Can't find an IP access list with id: " + TestingID + ".",
				},
				Status: 404,
			},
		},
		Resource: ResourceIPAccessList(),
		Read:     true,
		Removed:  true,
		ID:       TestingID,
	}.ApplyNoError(t)
}

func TestIPACLRead_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodGet,
				Resource: "/api/2.0/ip-access-lists/" + TestingID,
				Response: apierr.APIErrorBody{
					ErrorCode: "SERVER_ERROR",
					Message:   "Something unexpected happened",
				},
				Status: 500,
			},
		},
		Resource: ResourceIPAccessList(),
		Read:     true,
		ID:       TestingID,
	}.Apply(t)
	assert.Error(t, err)
	qa.AssertErrorStartsWith(t, err, "Something unexpected happened")
	assert.Equal(t, TestingID, d.Id(), "Id should not be empty for error reads")
}

func TestIPACLDelete(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodDelete,
				Resource: fmt.Sprintf("/api/2.0/ip-access-lists/%s?", TestingID),
			},
		},
		Resource: ResourceIPAccessList(),
		Delete:   true,
		ID:       TestingID,
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, TestingID, d.Id())
}

func TestIPACLDelete_Error(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   http.MethodDelete,
				Resource: fmt.Sprintf("/api/2.0/ip-access-lists/%s?", TestingID),
				Response: apierr.APIErrorBody{
					ErrorCode: "INVALID_STATE",
					Message:   "Something went wrong",
				},
				Status: 400,
			},
		},
		Resource: ResourceIPAccessList(),
		Delete:   true,
		Removed:  true,
		ID:       TestingID,
	}.ExpectError(t, "Something went wrong")
}

func TestListIpAccessLists(t *testing.T) {
	client, server, err := qa.HttpFixtureClient(t, []qa.HTTPFixture{
		{
			Method:   "GET",
			Resource: "/api/2.0/ip-access-lists",
			Response: map[string]any{},
		},
	})
	defer server.Close()
	require.NoError(t, err)

	ctx := context.Background()
	ipLists, err := NewIPAccessListsAPI(ctx, client).List()
	require.NoError(t, err)
	assert.Equal(t, len(ipLists.ListIPAccessListsResponse), 0)
}
