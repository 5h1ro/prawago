package server

import (
	"context"
	"github.com/devlikeapro/gows/media"
	__ "github.com/devlikeapro/gows/proto"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
)

func (s *Server) SetProfileName(ctx context.Context, req *__.ProfileNameRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	patch := appstate.BuildSettingPushName(req.Name)
	err = cli.SendAppState(ctx, patch)
	if err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

func (s *Server) SetProfileStatus(ctx context.Context, req *__.ProfileStatusRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	err = cli.SetStatusMessage(ctx, req.GetStatus())
	if err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

func (s *Server) SetProfilePicture(ctx context.Context, req *__.SetProfilePictureRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	content := req.Picture
	var picture []byte
	if len(content) != 0 {
		picture, err = media.ProfilePicture(req.Picture)
		if err != nil {
			return nil, err
		}
	}
	_, err = cli.SetGroupPhoto(ctx, types.EmptyJID, picture)
	if err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

func (s *Server) GetBusinessProfile(ctx context.Context, req *__.JidRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	profile, err := cli.GetBusinessProfile(ctx, jid)
	if err != nil {
		return nil, err
	}
	return toJson(profile)
}

func (s *Server) GetBotList(ctx context.Context, req *__.Session) (*__.JsonList, error) {
	cli, err := s.Sm.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	list, err := cli.GetBotListV2(ctx)
	if err != nil {
		return nil, err
	}
	return toJsonList(list)
}

func (s *Server) GetBotProfiles(ctx context.Context, req *__.Session) (*__.JsonList, error) {
	cli, err := s.Sm.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	list, err := cli.GetBotListV2(ctx)
	if err != nil {
		return nil, err
	}
	profiles, err := cli.GetBotProfiles(ctx, list)
	if err != nil {
		return nil, err
	}
	return toJsonList(profiles)
}
