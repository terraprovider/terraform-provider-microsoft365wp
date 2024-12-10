package retryablehttputil

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func ConfigureClientRetryLimitsAndBackoff(retryableClient *retryablehttp.Client) {

	// set retries / adjust from default values to avoid problems due to throttling
	retryableClient.RetryWaitMin = 3 * time.Second
	retryableClient.RetryWaitMax = 60 * time.Second
	retryableClient.RetryMax = 10

	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		backoffTime := retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
		// MS Graph often returns `Retry-After: 0` but we want (and have) to delay a bit to become successful at some
		// point. Therefore we ignore the value from `Retry-After` in case it's below our min and call DefaultBackoff
		// again to get the calculated default delay instead.
		if backoffTime < min {
			backoffTime = retryablehttp.DefaultBackoff(min, max, attemptNum, nil)
		}
		return backoffTime
	}
}
