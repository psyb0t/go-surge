package surge_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-surge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey      = "EXAMPLE-DO-NOT-USE"
	testAccountID   = "acct_example"
	testPhoneNumber = "+12025550123"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *surge.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := surge.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)
	client, err := surge.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBearerAuth(t *testing.T, req *http.Request) {
	t.Helper()
	assert.Equal(t, "Bearer "+testAPIKey, req.Header.Get("Authorization"))
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write([]byte(body))
	assert.NoError(t, err)
}

func TestListAccountsPagination(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/accounts", r.URL.Path)
		assertBearerAuth(t, r)

		switch r.URL.Query().Get("after") {
		case "":
			writeJSON(t, w, http.StatusOK, `{
				"data":[{"id":"acct_first","name":"First",
				"brand_name":null,"time_zone":null,"organization":{}}],
				"pagination":{"next_cursor":"page+2/=","previous_cursor":null}
			}`)
		case "page+2/=":
			writeJSON(t, w, http.StatusOK, `{
				"data":[{"id":"acct_second","name":"Second",
				"brand_name":null,"time_zone":null,"organization":{}}],
				"pagination":{"next_cursor":"page3","previous_cursor":"page1"}
			}`)
		case "page3":
			writeJSON(t, w, http.StatusOK, `{
				"data":[{"id":"acct_third","name":"Third",
				"brand_name":null,"time_zone":null,"organization":{}}],
				"pagination":{"next_cursor":null,"previous_cursor":"page2"}
			}`)
		default:
			writeJSON(t, w, http.StatusOK, `{
				"data":[],"pagination":{"next_cursor":null}}
			`)
		}
	})

	params := &surge.ListAccountsParams{}

	ids := make([]string, 0, 3)

	for range 3 {
		res, err := client.ListAccountsWithResponse(t.Context(), params)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())
		require.NotNil(t, res.JSON200)
		require.Len(t, res.JSON200.Data, 1)
		assert.Nil(t, res.JSON200.Data[0].BrandName)
		ids = append(ids, res.JSON200.Data[0].Id)
		params.After = res.JSON200.Pagination.NextCursor
	}

	assert.Equal(t, []string{"acct_first", "acct_second", "acct_third"}, ids)
	assert.Nil(t, params.After)

	end := "past-end"
	res, err := client.ListAccountsWithResponse(
		t.Context(), &surge.ListAccountsParams{After: &end},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Empty(t, res.JSON200.Data)
	assert.Nil(t, res.JSON200.Pagination.NextCursor)
}

func TestCreateContact(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/accounts/acct_example/contacts", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearerAuth(t, r)

		var body surge.ContactParams
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, testPhoneNumber, body.PhoneNumber)

		writeJSON(t, w, http.StatusCreated, `{
			"id":"ctc_example","phone_number":"+12025550123",
			"first_name":"Sample","outbound_disabled":false
		}`)
	})

	res, err := client.CreateContactWithResponse(
		t.Context(), testAccountID,
		surge.ContactParams{PhoneNumber: testPhoneNumber},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())
	require.NotNil(t, res.JSON201)
	assert.Equal(t, "ctc_example", res.JSON201.Id)
	assert.Equal(t, testPhoneNumber, res.JSON201.PhoneNumber)
	require.NotNil(t, res.JSON201.FirstName)
	assert.Equal(t, "Sample", *res.JSON201.FirstName)
}

func TestCreateMessage(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/accounts/acct_example/messages", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearerAuth(t, r)

		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, testPhoneNumber, body["to"])
		assert.Equal(t, "Sample message", body["body"])
		assert.NotContains(t, body, "conversation")

		writeJSON(t, w, http.StatusCreated, `{
			"id":"msg_example","status":"pending","body":null,
			"blast_id":null,"metadata":{},"attachments":[]
		}`)
	})

	text := "Sample message"

	var body surge.MessageParams

	err := body.FromMessageParams1(surge.MessageParams1{
		To: testPhoneNumber, Body: &text,
	})
	require.NoError(t, err)

	res, err := client.CreateMessageWithResponse(
		t.Context(), testAccountID, body,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())
	require.NotNil(t, res.JSON201)
	require.NotNil(t, res.JSON201.Id)
	assert.Equal(t, "msg_example", *res.JSON201.Id)
	require.NotNil(t, res.JSON201.Status)
	assert.Equal(t, surge.MessageStatusPending, *res.JSON201.Status)
	assert.Nil(t, res.JSON201.Body)
	assert.Nil(t, res.JSON201.BlastId)
}

func TestAuthenticationFailures(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		auth   string
		status int
	}{
		{name: "missing", status: http.StatusUnauthorized},
		{name: "wrong", auth: "Bearer wrong", status: http.StatusForbidden},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/contacts/ctc_example", r.URL.Path)
				assert.Equal(t, tc.auth, r.Header.Get("Authorization"))
				writeJSON(t, w, tc.status, `{
					"error":{"type":"unauthorized","message":"Access denied"}}
				`)
			}
			client := newTestClient(t, handler)
			res, err := client.GetContactWithResponse(
				t.Context(), "ctc_example",
				func(_ context.Context, req *http.Request) error {
					req.Header.Set("Authorization", tc.auth)

					return nil
				},
			)
			require.NoError(t, err)
			assert.Equal(t, tc.status, res.StatusCode())
			assert.Nil(t, res.JSON200)
			require.NotNil(t, res.JSONDefault)
			assert.Equal(t, "unauthorized", res.JSONDefault.Error.Type)
		})
	}
}

func TestMalformedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusOK, `{"id":`)
	})
	_, err := client.GetContactWithResponse(t.Context(), "ctc_example")

	var syntaxError *json.SyntaxError
	require.ErrorAs(t, err, &syntaxError)
}

func TestCanceledRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("canceled request reached the server")
	})
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := client.GetContactWithResponse(ctx, "ctc_example")
	require.ErrorIs(t, err, context.Canceled)
}

func TestWebhookModel(t *testing.T) {
	t.Parallel()

	var event surge.MessageSentEvent

	err := json.Unmarshal([]byte(`{
		"type":"message.sent","account_id":"acct_example",
		"data":{"id":"msg_example","body":null,"status":"sent"}
	}`), &event)
	require.NoError(t, err)
	assert.Equal(t, surge.MessageSentEventTypeMessageSent, event.Type)
	assert.True(t, event.Type.Valid())
	assert.Nil(t, event.Data.Body)
}
