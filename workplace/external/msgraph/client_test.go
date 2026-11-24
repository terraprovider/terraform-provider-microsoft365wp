package msgraph

import (
	"context"
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
