package server

import (
	"context"
	"encoding/json"

	__ "github.com/devlikeapro/gows/proto"
)

const (
	usernameQueryCheck           = "26124072630599520"
	usernameQueryCheckMulti      = "27134626522840290"
	usernameQuerySet             = "27108705368767936"
	usernameQueryGet             = "32618050064506056"
	usernameQueryRecommendations = "26077456248616956"
	usernameQueryPinSet          = "25529696019976770"
)

func (s *Server) CheckUsername(ctx context.Context, req *__.CheckUsernameRequest) (*__.Json, error) {
	return s.sendUsernameMex(ctx, req.GetSession().GetId(), usernameQueryCheck, map[string]any{
		"username":            req.GetUsername(),
		"include_suggestions": req.GetIncludeSuggestions(),
	}, "xwa2_username_check")
}

func (s *Server) CheckUsernameMulti(ctx context.Context, req *__.CheckUsernameMultiRequest) (*__.Json, error) {
	return s.sendUsernameMex(ctx, req.GetSession().GetId(), usernameQueryCheckMulti, map[string]any{
		"usernames": req.GetUsernames(),
	}, "xwa2_username_check_multi")
}

func (s *Server) SetUsername(ctx context.Context, req *__.SetUsernameRequest) (*__.Json, error) {
	source := req.GetSource()
	if source == "" {
		source = "USER_INPUT"
	}
	variables := map[string]any{
		"username": req.GetUsername(),
		"reserved": false,
		"source":   source,
	}
	if req.GetSessionId() != "" {
		variables["session_id"] = req.GetSessionId()
	}
	if req.GetPin() != "" {
		variables["pin"] = req.GetPin()
	}
	return s.sendUsernameMex(ctx, req.GetSession().GetId(), usernameQuerySet, variables, "xwa2_username_set")
}

func (s *Server) DeleteUsername(ctx context.Context, req *__.Session) (*__.Json, error) {
	return s.sendUsernameMex(ctx, req.GetId(), usernameQuerySet, map[string]any{
		"username": nil,
	}, "xwa2_username_delete")
}

func (s *Server) GetMyUsername(ctx context.Context, req *__.Session) (*__.OptionalString, error) {
	data, err := s.sendRawMex(ctx, req.GetId(), usernameQueryGet, map[string]any{}, "xwa2_username_get")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &__.OptionalString{Value: resp.Username}, nil
}

func (s *Server) SetUsernamePin(ctx context.Context, req *__.UsernamePinRequest) (*__.Json, error) {
	return s.sendUsernameMex(ctx, req.GetSession().GetId(), usernameQueryPinSet, map[string]any{
		"pin": req.GetPin(),
	}, "xwa2_username_pin_set")
}

func (s *Server) GetUsernameRecommendations(ctx context.Context, req *__.UsernameRecommendationsRequest) (*__.Json, error) {
	variables := map[string]any{}
	if req.GetSource() != "" {
		variables["source"] = req.GetSource()
	}
	return s.sendUsernameMex(ctx, req.GetSession().GetId(), usernameQueryRecommendations, variables, "xwa2_username_get_recommendations")
}

func (s *Server) SendMexIQ(ctx context.Context, req *__.MexIQRequest) (*__.Json, error) {
	var variables any = map[string]any{}
	if req.GetVariablesJson() != "" {
		if err := json.Unmarshal([]byte(req.GetVariablesJson()), &variables); err != nil {
			return nil, err
		}
	}
	data, err := s.sendRawMex(ctx, req.GetSession().GetId(), req.GetQueryId(), variables, req.GetDataPath())
	if err != nil {
		return nil, err
	}
	return &__.Json{Data: string(data)}, nil
}

func (s *Server) sendUsernameMex(ctx context.Context, sessionID string, queryID string, variables any, dataPath string) (*__.Json, error) {
	data, err := s.sendRawMex(ctx, sessionID, queryID, variables, dataPath)
	if err != nil {
		return nil, err
	}
	return &__.Json{Data: string(data)}, nil
}

func (s *Server) sendRawMex(ctx context.Context, sessionID string, queryID string, variables any, dataPath string) (json.RawMessage, error) {
	cli, err := s.Sm.Get(sessionID)
	if err != nil {
		return nil, err
	}
	return cli.SendMexIQ(ctx, queryID, variables, dataPath)
}
