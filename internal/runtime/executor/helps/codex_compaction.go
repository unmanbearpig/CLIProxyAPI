package helps

import (
	"strings"

	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// IsCodexCompactionRequest reads per-message metadata because one downstream
// websocket can carry both compactions and normal turns under the same headers.
func IsCodexCompactionRequest(payload []byte) bool {
	raw := strings.TrimSpace(gjson.GetBytes(payload, "client_metadata.x-codex-turn-metadata").String())
	return strings.EqualFold(strings.TrimSpace(gjson.Get(raw, "request_kind").String()), "compaction")
}

// StripCodexWebsocketEnvelopeType removes the websocket-only envelope before
// HTTP translation and payload rules; /responses rejects this top-level field.
func StripCodexWebsocketEnvelopeType(req *cliproxyexecutor.Request, opts *cliproxyexecutor.Options) {
	if stripped, err := sjson.DeleteBytes(req.Payload, "type"); err == nil {
		req.Payload = stripped
	}
	if len(opts.OriginalRequest) > 0 {
		if stripped, err := sjson.DeleteBytes(opts.OriginalRequest, "type"); err == nil {
			opts.OriginalRequest = stripped
		}
	}
}
