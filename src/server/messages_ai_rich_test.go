package server

import (
	"encoding/json"
	"testing"

	"go.mau.fi/whatsmeow/proto/waAICommonDeprecated"
)

func TestBuildAIRichResponseMessage(t *testing.T) {
	msg := buildAIRichResponseMessage([]*waAICommonDeprecated.AIRichResponseSubMessage{
		buildRichTextSubMessage("hello"),
	}, nil, []byte(`{"ok":true}`))

	rich := msg.GetBotForwardedMessage().GetMessage().GetRichResponseMessage()
	if rich == nil {
		t.Fatal("missing rich response")
	}
	if len(rich.GetSubmessages()) != 1 {
		t.Fatalf("submessages = %d, want 1", len(rich.GetSubmessages()))
	}
	if len(rich.GetUnifiedResponse().GetData()) == 0 {
		t.Fatal("missing unified response data")
	}
	if !rich.GetContextInfo().GetIsForwarded() {
		t.Fatal("rich response should be marked forwarded")
	}
}

func TestBuildMarkdownUnifiedResponse(t *testing.T) {
	data, err := buildMarkdownUnifiedResponse("**hello**")
	if err != nil {
		t.Fatal(err)
	}

	var payload struct {
		ResponseID string `json:"response_id"`
		Sections   []any  `json:"sections"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.ResponseID == "" || len(payload.Sections) != 1 {
		t.Fatalf("bad markdown payload: %+v", payload)
	}
}
