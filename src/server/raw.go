package server

import (
	"context"
	"encoding/json"

	"github.com/devlikeapro/gows/gows"
	__ "github.com/devlikeapro/gows/proto"
	"go.mau.fi/whatsmeow/types"
)

func (s *Server) SendIQ(ctx context.Context, req *__.IQRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	to := types.ServerJID
	if req.GetTo() != "" {
		to, err = types.ParseJID(req.GetTo())
		if err != nil {
			return nil, err
		}
	}
	content, err := gows.DecodeIQContent(req.GetContentJson())
	if err != nil {
		return nil, err
	}
	resp, err := cli.SendIQ(ctx, req.GetNamespace(), req.GetType(), to, content)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return &__.Json{Data: string(data)}, nil
}
