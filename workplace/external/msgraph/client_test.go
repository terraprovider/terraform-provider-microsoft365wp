package msgraph

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient returns a client for srv with short retry waits.
func newTestClient(srv *httptest.Server) Client {
	c := NewClient(VersionBeta)
	c.Endpoint = srv.URL
	c.RetryableClient.RetryWaitMin = 10 * time.Millisecond
	c.RetryableClient.RetryWaitMax = 20 * time.Millisecond
	c.RetryableClient.RetryMax = 5
	return c
}

// notFoundServer returns 404 for the first n requests to each path and 200 afterwards.
func notFoundServer(n int32) (*httptest.Server, func(path string) int32) {
	var mu sync.Mutex
	counts := map[string]*atomic.Int32{}
	counter := func(path string) *atomic.Int32 {
		mu.Lock()
		defer mu.Unlock()
		if counts[path] == nil {
			counts[path] = &atomic.Int32{}
		}
		return counts[path]
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if counter(r.URL.Path).Add(1) <= n {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"code":"Request_ResourceNotFound","message":"not found"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"x"}`))
	}))
	return srv, func(path string) int32 { return counter(path).Load() }
}

func TestConsistencyFailureFuncIsPerRequest(t *testing.T) {
	srv, count := notFoundServer(2)
	defer srv.Close()
	c := newTestClient(srv)

	// Run requests with and without ConsistencyFailureFunc concurrently: each must keep its own retry behaviour.
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			_, status, _, err := c.Get(context.Background(), GetHttpRequestInput{
				Uri:                    Uri{Entity: "/retry/" + string(rune('a'+i))},
				ValidStatusCodes:       []int{http.StatusOK},
				ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
			})
			if err != nil || status != http.StatusOK {
				t.Errorf("request with retry: status %d, err %v", status, err)
			}
		}(i)
		go func(i int) {
			defer wg.Done()
			_, status, _, _ := c.Get(context.Background(), GetHttpRequestInput{
				Uri:              Uri{Entity: "/noretry/" + string(rune('a'+i))},
				ValidStatusCodes: []int{http.StatusOK},
			})
			if status != http.StatusNotFound {
				t.Errorf("request without retry: expected status 404, got %d", status)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		suffix := string(rune('a' + i))
		if got := count("/beta/retry/" + suffix); got != 3 {
			t.Errorf("request with retry %s: expected 3 attempts, got %d", suffix, got)
		}
		if got := count("/beta/noretry/" + suffix); got != 1 {
			t.Errorf("request without retry %s: expected 1 attempt, got %d", suffix, got)
		}
	}
}

func TestRetryOn404ConsistencyFailureFuncForStopsAtDeadline(t *testing.T) {
	srv, count := notFoundServer(1000)
	defer srv.Close()
	c := newTestClient(srv)
	c.RetryableClient.RetryMax = 1000

	start := time.Now()
	_, status, _, err := c.Get(context.Background(), GetHttpRequestInput{
		Uri:                    Uri{Entity: "/gone"},
		ValidStatusCodes:       []int{http.StatusOK},
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFuncFor(200 * time.Millisecond),
	})
	if status != http.StatusNotFound || err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected a 404 error, got status %d, err %v", status, err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("retries did not stop at the deadline (took %s)", elapsed)
	}
	if got := count("/beta/gone"); got < 2 {
		t.Errorf("expected some retries before the deadline, got %d attempt(s)", got)
	}
}

func TestResolveUrl(t *testing.T) {
	cases := []struct {
		endpoint, ref, want string
	}{
		{"https://graph.microsoft.com", "directoryObjects/", "https://graph.microsoft.com/beta/directoryObjects/"},
		{"https://graph.microsoft.us/", "directoryObjects/", "https://graph.microsoft.us/beta/directoryObjects/"},
		{"https://microsoftgraph.chinacloudapi.cn", "/odata/groups/", "https://microsoftgraph.chinacloudapi.cn/odata/groups/"},
		{"https://dod-graph.microsoft.us", "deviceManagement/roleDefinitions('a-b')", "https://dod-graph.microsoft.us/beta/deviceManagement/roleDefinitions('a-b')"},
		{"https://graph.microsoft.us", "https://example.com/x/", "https://example.com/x/"},
	}
	for _, c := range cases {
		client := NewClient(VersionBeta)
		client.Endpoint = c.endpoint
		got, err := client.ResolveUrl(c.ref)
		if err != nil {
			t.Fatalf("ResolveUrl(%q) on %q: %v", c.ref, c.endpoint, err)
		}
		if got != c.want {
			t.Errorf("ResolveUrl(%q) on %q = %q, want %q", c.ref, c.endpoint, got, c.want)
		}
	}
}

// newShortWaitClient returns a client for srv with the default retry limits, but short retry waits.
func newShortWaitClient(srv *httptest.Server) Client {
	c := NewClient(VersionBeta)
	c.Endpoint = srv.URL
	c.RetryableClient.RetryWaitMin = time.Millisecond
	c.RetryableClient.RetryWaitMax = 2 * time.Millisecond
	return c
}

// Eventual consistency in MS Graph can take longer than the default retries (about 15 s), so consistency retries must
// only be limited by their ConsistencyFailureFunc (i.e. its deadline).
func TestConsistencyRetriesAreNotLimitedByDefaultRetryMax(t *testing.T) {
	notFound := int32(defaultRetryMax * 3)
	srv, count := notFoundServer(notFound)
	defer srv.Close()
	c := newShortWaitClient(srv)

	_, status, _, err := c.Get(context.Background(), GetHttpRequestInput{
		Uri:                    Uri{Entity: "/fresh"},
		ValidStatusCodes:       []int{http.StatusOK},
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFuncFor(10 * time.Second),
	})
	if err != nil || status != http.StatusOK {
		t.Fatalf("expected 200, got status %d, err %v", status, err)
	}
	if got := count("/beta/fresh"); got != notFound+1 {
		t.Errorf("expected %d attempts, got %d", notFound+1, got)
	}
}

func TestOtherRetriesStayLimitedByDefaultRetryMax(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"code":"ServiceUnavailable","message":"try again"}}`))
	}))
	defer srv.Close()
	c := newShortWaitClient(srv)

	_, status, _, _ := c.Get(context.Background(), GetHttpRequestInput{
		Uri:                    Uri{Entity: "/busy"},
		ValidStatusCodes:       []int{http.StatusOK},
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc, // must not extend the retries of other responses
	})
	if status != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", status)
	}
	if got := attempts.Load(); got != int32(defaultRetryMax+1) {
		t.Errorf("expected %d attempts, got %d", defaultRetryMax+1, got)
	}
}

func TestConsistencyRetriesWithoutDeadlineStopAtConsistencyRetryMax(t *testing.T) {
	srv, count := notFoundServer(1000)
	defer srv.Close()
	c := newShortWaitClient(srv)

	_, status, _, _ := c.Get(context.Background(), GetHttpRequestInput{
		Uri:                    Uri{Entity: "/gone"},
		ValidStatusCodes:       []int{http.StatusOK},
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
	})
	if status != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", status)
	}
	if got := count("/beta/gone"); got != int32(consistencyRetryMax+1) {
		t.Errorf("expected %d attempts, got %d", consistencyRetryMax+1, got)
	}
}

func TestRetryBackoff(t *testing.T) {
	respWith := func(state *retryState) *http.Response {
		req, _ := http.NewRequest(http.MethodGet, "https://example.invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), retryContextKey{}, retryContext{state: state}))
		return &http.Response{StatusCode: http.StatusNotFound, Request: req}
	}
	min, max := time.Second, 30*time.Second

	for _, tc := range []struct {
		name  string
		state *retryState
		want  time.Duration
	}{
		{"first consistency retry", &retryState{consistencyRetries: 1, lastWasConsistency: true}, time.Second},
		{"later consistency retries are capped", &retryState{consistencyRetries: 10, lastWasConsistency: true}, consistencyRetryWaitMax},
		{"other retries count their own attempts", &retryState{consistencyRetries: 10, defaultRetries: 2}, 2 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryBackoff(min, max, 11, respWith(tc.state)); got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

// TestDescribeDefaultRetry checks the log message of retries by the default retry policy, which do not show up in the
// request/response logs (#247).
func TestDescribeDefaultRetry(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "https://graph.microsoft.com/beta/policies/authenticationStrengthPolicies", nil)
	for _, tc := range []struct {
		name string
		resp *http.Response
		err  error
		want string
	}{
		{"5xx", &http.Response{StatusCode: http.StatusBadGateway, Request: req}, nil,
			"MS Graph returned 502 for POST https://graph.microsoft.com/beta/policies/authenticationStrengthPolicies, retrying"},
		{"connection error", nil, errors.New(`Post "https://graph.microsoft.com/beta/x": EOF`),
			`request failed, retrying: Post "https://graph.microsoft.com/beta/x": EOF`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := describeDefaultRetry(tc.resp, tc.err); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
