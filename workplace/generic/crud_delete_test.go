package generic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// deleteTestServer answers DELETE and GET on /beta/policies/p with the given status sequences (the last status
// repeats) and counts the requests.
type deleteTestServer struct {
	mu      sync.Mutex
	deletes []int
	gets    []int
	nDelete int
	nGet    int
}

const trustedLocationBody = `{"error":{"code":"BadRequest","message":"The server could not process the request because it is malformed or incorrect.","innerError":{"message":"1177: Named location: 'p' cannot be deleted because it is marked as a Trusted location."}}}`

// retryTrustedLocation is a WriteRetryFunc for the error in trustedLocationBody.
func retryTrustedLocation(resp *http.Response, o *odata.OData) bool {
	return resp.StatusCode == http.StatusBadRequest && o != nil && o.Error != nil && strings.Contains(o.Error.String(), "1177:")
}

const checkAccessBody = `{"error":{"code":"AccessDenied","message":"Access denied. CheckAccess returned 'NotAllowed' for action 'microsoft.directory/conditionalAccessPolicies/delete'."}}`

func (s *deleteTestServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		if req.URL.Path != "/beta/policies/p" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
		}
		var status int
		switch req.Method {
		case http.MethodDelete:
			status = s.deletes[min(s.nDelete, len(s.deletes)-1)]
			s.nDelete++
		case http.MethodGet:
			status = s.gets[min(s.nGet, len(s.gets)-1)]
			s.nGet++
		default:
			t.Errorf("unexpected method %s", req.Method)
		}
		if status != http.StatusNoContent {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		switch status {
		case http.StatusOK:
			_, _ = io.WriteString(w, `{"id":"p"}`)
		case http.StatusNotFound:
			_, _ = io.WriteString(w, `{"error":{"code":"ResourceNotFound","message":"ConditionalAccessPolicy with id p is null/invalid"}}`)
		case http.StatusForbidden:
			_, _ = io.WriteString(w, checkAccessBody)
		case http.StatusBadRequest:
			_, _ = io.WriteString(w, trustedLocationBody)
		}
	}
}

func TestDeleteEntityRaw(t *testing.T) {
	oldPoll, oldResend := deleteVerifyPollInterval, deleteVerifyResendInterval
	deleteVerifyPollInterval, deleteVerifyResendInterval = time.Millisecond, 30*time.Millisecond
	defer func() { deleteVerifyPollInterval, deleteVerifyResendInterval = oldPoll, oldResend }()

	withOptions := WriteOptions{DeleteRetryNotFoundFor: 300 * time.Millisecond, DeleteVerifyFor: 300 * time.Millisecond}
	withWriteRetry := WriteOptions{WriteRetryFor: 300 * time.Millisecond, WriteRetryFunc: retryTrustedLocation}

	for _, tc := range []struct {
		name        string
		options     WriteOptions
		deletes     []int
		gets        []int
		wantDeletes int // -1: don't check, -2: at least 2 (deleted again)
		wantGets    int // -1: don't check
		wantWarning bool
		wantError   string
	}{
		{"no options: like DeleteRaw", WriteOptions{}, []int{204}, []int{200}, 1, 0, false, ""},
		{"no options: 404 is accepted as gone", WriteOptions{}, []int{404}, []int{200}, 1, 0, true, ""},
		{"deleted and gone", withOptions, []int{204}, []int{404}, 1, 3, false, ""},
		{"404 right after create is retried", withOptions, []int{404, 404, 204}, []int{404}, 3, 3, false, ""},
		{"403 CheckAccess right after create is retried", withOptions, []int{403, 204}, []int{404}, 2, 3, false, ""},
		{"404 reads must be consecutive", withOptions, []int{204}, []int{404, 200, 404}, 1, -1, false, ""},
		{"lost 204 is deleted again", withOptions, []int{204}, append(repeat(200, 50), 404), -2, -1, false, ""},
		{"404 for the whole window and gone", withOptions, []int{404}, []int{404}, -1, 3, true, ""},
		{"404 for the whole window, but still exists", withOptions, []int{404}, []int{200}, -1, -1, false, "still exists"},
		{"404 for the whole window, then deleted again", withOptions, []int{404, 404, 404, 404, 404, 404, 404, 404, 204}, append(repeat(200, 20), 404), -2, -1, false, ""},
		{"still exists after verify window", withOptions, []int{204}, []int{200}, -1, -1, false, "still exists"},
		{"other DELETE errors fail at once", withOptions, []int{400}, []int{404}, 1, 0, false, "Error deleting"},
		{"WriteRetryFunc: accepted error is retried", withWriteRetry, []int{400, 400, 204}, []int{404}, 3, 0, false, ""},
		{"WriteRetryFunc: accepted error fails after the window", withWriteRetry, []int{400}, []int{404}, -2, 0, false, "1177"},
		{"WriteRetryFunc: 404 is not retried", withWriteRetry, []int{404, 204}, []int{404}, 1, 0, true, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &deleteTestServer{deletes: tc.deletes, gets: tc.gets}
			server := httptest.NewServer(s.handler(t))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.RetryableClient.RetryWaitMin = time.Millisecond
			client.RetryableClient.RetryWaitMax = 5 * time.Millisecond
			client.RetryableClient.RetryMax = 1000

			aps := &AccessParams{BaseUri: "/policies", WriteOptions: tc.options, graphClient: &client}
			diags := diag.Diagnostics{}
			aps.deleteEntityRaw(context.Background(), &diags, "", "p", nil)

			if tc.wantError == "" && diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if tc.wantError != "" && (!diags.HasError() || !strings.Contains(diagsString(diags), tc.wantError)) {
				t.Fatalf("want error containing %q, got %v", tc.wantError, diags)
			}
			if gotWarning := diags.WarningsCount() > 0; gotWarning != tc.wantWarning {
				t.Errorf("warning: got %v, want %v (diags %v)", gotWarning, tc.wantWarning, diags)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			if tc.wantDeletes == -2 && s.nDelete < 2 {
				t.Errorf("got %d DELETEs, want at least 2", s.nDelete)
			}
			if tc.wantDeletes >= 0 && s.nDelete != tc.wantDeletes {
				t.Errorf("got %d DELETEs, want %d", s.nDelete, tc.wantDeletes)
			}
			if tc.wantGets >= 0 && s.nGet != tc.wantGets {
				t.Errorf("got %d GETs, want %d", s.nGet, tc.wantGets)
			}
		})
	}
}

func repeat(v int, n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = v
	}
	return r
}

func diagsString(diags diag.Diagnostics) string {
	var sb strings.Builder
	for _, d := range diags {
		fmt.Fprintf(&sb, "%s: %s\n", d.Summary(), d.Detail())
	}
	return sb.String()
}
