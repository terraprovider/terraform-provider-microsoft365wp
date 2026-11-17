package generic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// RawContains reports whether the raw MS Graph value got contains the raw value want (e.g. what has just been written):
//   - maps: got must contain every key of want with a matching value (keys starting with "@odata." are ignored)
//   - lists: same length, and every element of want must match a distinct element of got (in any order)
//   - strings are compared case-insensitively (MS Graph may normalize the casing of enum values)
//   - nil in want matches anything, and empty lists/maps in want also match nil in got
//
// It is meant to detect stale reads after a write, not to validate the write.
func RawContains(want, got any) bool {
	switch w := want.(type) {
	case nil:
		return true
	case map[string]any:
		if got == nil {
			return len(w) == 0
		}
		g, ok := got.(map[string]any)
		if !ok {
			return false
		}
		for k, wv := range w {
			if strings.HasPrefix(k, "@odata.") {
				continue
			}
			gv, ok := g[k]
			if !ok {
				if wv == nil {
					continue
				}
				return false
			}
			if !RawContains(wv, gv) {
				return false
			}
		}
		return true
	case []any:
		if got == nil {
			return len(w) == 0
		}
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return false
		}
		used := make([]bool, len(g))
		for _, wv := range w {
			found := false
			for i, gv := range g {
				if !used[i] && RawContains(wv, gv) {
					used[i] = true
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	case string:
		g, ok := got.(string)
		return ok && strings.EqualFold(w, g)
	default:
		return want == got
	}
}

type WaitForRawValuesOptions struct {
	Timeout          time.Duration
	PollInterval     time.Duration
	ConsecutiveReads int // reads may be served by different replicas, so require several consecutive matching reads
}

var DefaultWaitForRawValuesOptions = WaitForRawValuesOptions{
	Timeout:          30 * time.Second,
	PollInterval:     2 * time.Second,
	ConsecutiveReads: 2,
}

// WaitForRawValues polls uri until MS Graph returns values that contain want (see RawContains). This helps for
// entities that MS Graph serves stale reads for after a write, which would otherwise lead to differences in a plan run
// right after an apply. It only adds a warning if this does not happen within the timeout, since the write itself has
// succeeded.
func WaitForRawValues(ctx context.Context, diags *diag.Diagnostics, client *msgraph.Client, uri msgraph.Uri,
	want map[string]any, opts WaitForRawValuesOptions) {

	if len(want) == 0 {
		return
	}

	deadline := time.Now().Add(opts.Timeout)
	matches := 0
	for {
		readDiags := diag.Diagnostics{}
		got := ReadRaw2(ctx, &readDiags, client, uri, nil, nil, true) // 404 just means "not visible yet"
		if !readDiags.HasError() && got != nil && RawContains(want, got) {
			matches++
			if matches >= opts.ConsecutiveReads {
				return
			}
		} else {
			matches = 0
		}
		if time.Now().After(deadline) {
			diags.AddWarning("MS Graph does not return the written values yet",
				fmt.Sprintf("%q has been written successfully, but MS Graph still returned other values after %s "+
					"(eventual consistency). A plan run right now might therefore still show differences, which should "+
					"disappear after a while.", uri.Entity, opts.Timeout))
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(opts.PollInterval):
		}
	}
}
