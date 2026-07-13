package gows

import (
	"context"
	"encoding/json"

	"go.mau.fi/whatsmeow"
	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

func (gows *GoWS) SendIQ(ctx context.Context, namespace, iqType string, to types.JID, content []waBinary.Node) (*waBinary.Node, error) {
	if iqType == "set" {
		return gows.int.SendIQ(ctx, whatsmeow.DangerousInfoQuery{
			Namespace: namespace,
			Type:      "set",
			To:        to,
			Content:   content,
		})
	}
	return gows.int.SendIQ(ctx, whatsmeow.DangerousInfoQuery{
		Namespace: namespace,
		Type:      "get",
		To:        to,
		Content:   content,
	})
}

func DecodeIQContent(data string) ([]waBinary.Node, error) {
	if data == "" {
		return nil, nil
	}
	var content []waBinary.Node
	if err := json.Unmarshal([]byte(data), &content); err != nil {
		return nil, err
	}
	return content, nil
}
