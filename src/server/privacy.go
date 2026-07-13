package server

import (
	"context"
	"slices"
	"time"

	__ "github.com/devlikeapro/gows/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validPrivacySettingNames = []types.PrivacySettingType{
	types.PrivacySettingTypeGroupAdd,
	types.PrivacySettingTypeLastSeen,
	types.PrivacySettingTypeStatus,
	types.PrivacySettingTypeProfile,
	types.PrivacySettingTypeReadReceipts,
	types.PrivacySettingTypeOnline,
	types.PrivacySettingTypeCallAdd,
	types.PrivacySettingTypeMessages,
	types.PrivacySettingTypeDefense,
	types.PrivacySettingTypeStickers,
}

var validPrivacySettingValues = []types.PrivacySetting{
	types.PrivacySettingAll,
	types.PrivacySettingContacts,
	types.PrivacySettingContactAllowlist,
	types.PrivacySettingContactBlacklist,
	types.PrivacySettingMatchLastSeen,
	types.PrivacySettingKnown,
	types.PrivacySettingNone,
	types.PrivacySettingOnStandard,
	types.PrivacySettingOff,
}

func (s *Server) GetPrivacySettings(ctx context.Context, req *__.Session) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	settings := cli.GetPrivacySettings(ctx)
	return toJson(settings)
}

func (s *Server) SetPrivacySetting(ctx context.Context, req *__.PrivacySettingRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	name := types.PrivacySettingType(req.GetName())
	if !slices.Contains(validPrivacySettingNames, name) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid privacy setting name: %s", req.GetName())
	}
	value := types.PrivacySetting(req.GetValue())
	if !slices.Contains(validPrivacySettingValues, value) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid privacy setting value: %s", req.GetValue())
	}
	settings, err := cli.SetPrivacySetting(ctx, name, value)
	if err != nil {
		return nil, err
	}
	return toJson(settings)
}

func (s *Server) SetDefaultDisappearingTimer(ctx context.Context, req *__.DefaultDisappearingTimerRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	if err := cli.SetDefaultDisappearingTimer(ctx, time.Duration(req.GetSeconds())*time.Second); err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}
