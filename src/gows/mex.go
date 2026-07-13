package gows

import (
	"context"
	"encoding/json"
)

func (gows *GoWS) SendMexIQ(ctx context.Context, queryID string, variables any, dataPath string) (json.RawMessage, error) {
	data, err := gows.int.SendMexIQ(ctx, queryID, variables)
	if err != nil || dataPath == "" {
		return data, err
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	if child, ok := root[dataPath]; ok {
		return child, nil
	}
	return data, nil
}
