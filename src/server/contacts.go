package server

import (
	"context"
	"fmt"

	"github.com/devlikeapro/gows/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (s *Server) UpdateContact(ctx context.Context, req *__.UpdateContactRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.Session.Id)
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.Jid)
	if err != nil {
		return nil, fmt.Errorf("error parsing jid: %w", err)
	}
	err = cli.UpdateContact(ctx, jid, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

func (s *Server) GetContactById(ctx context.Context, req *__.EntityByIdRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	user, err := types.ParseJID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("error parsing jid %v: %w", req.Id, err)
	}

	contact, err := cli.Storage.Contacts.GetContact(user)
	if err != nil {
		return nil, fmt.Errorf("error getting contact %v: %w", user, err)
	}
	response, err := toJson(contact)
	if err != nil {
		return nil, fmt.Errorf("error marshaling contact %v: %w", user, err)
	}
	return response, nil
}

func (s *Server) GetContacts(ctx context.Context, req *__.GetContactsRequest) (*__.JsonList, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	pagination := toPagination(req.Pagination)
	sort := toStorageSort(req.SortBy)
	contacts, err := cli.Storage.Contacts.GetAllContacts(sort, pagination)
	if err != nil {
		return nil, err
	}
	response, err := toJsonList(contacts)
	if err != nil {
		return nil, fmt.Errorf("error marshaling contacts: %w", err)
	}
	return response, nil
}

func (s *Server) ResolveBusinessMessageLink(ctx context.Context, req *__.GroupCodeRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	target, err := cli.ResolveBusinessMessageLink(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}
	return toJson(target)
}

func (s *Server) ResolveContactQRLink(ctx context.Context, req *__.GroupCodeRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	target, err := cli.ResolveContactQRLink(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}
	return toJson(target)
}

func (s *Server) GetContactQRLink(ctx context.Context, req *__.ContactQRLinkRequest) (*__.OptionalString, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	link, err := cli.GetContactQRLink(ctx, req.GetRevoke())
	if err != nil {
		return nil, err
	}
	return &__.OptionalString{Value: link}, nil
}

func (s *Server) GetUserInfo(ctx context.Context, req *__.JidsRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jids, err := parseJIDs(req.GetJids())
	if err != nil {
		return nil, err
	}
	info, err := cli.GetUserInfo(ctx, jids)
	if err != nil {
		return nil, err
	}
	byJID := make(map[string]types.UserInfo, len(info))
	for jid, userInfo := range info {
		byJID[jid.String()] = userInfo
	}
	return toJson(byJID)
}

func (s *Server) GetUserDevices(ctx context.Context, req *__.JidsRequest) (*__.JsonList, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jids, err := parseJIDs(req.GetJids())
	if err != nil {
		return nil, err
	}
	devices, err := cli.GetUserDevices(ctx, jids)
	if err != nil {
		return nil, err
	}
	return toJsonList(jidsToStrings(devices))
}

func (s *Server) GetBlocklist(ctx context.Context, req *__.Session) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	blocklist, err := cli.GetBlocklist(ctx)
	if err != nil {
		return nil, err
	}
	return toJson(map[string]interface{}{
		"dhash": blocklist.DHash,
		"jids":  jidsToStrings(blocklist.JIDs),
	})
}

func (s *Server) UpdateBlocklist(ctx context.Context, req *__.BlocklistRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	var action events.BlocklistChangeAction
	switch req.GetAction() {
	case string(events.BlocklistChangeActionBlock):
		action = events.BlocklistChangeActionBlock
	case string(events.BlocklistChangeActionUnblock):
		action = events.BlocklistChangeActionUnblock
	default:
		return nil, fmt.Errorf("unknown blocklist action: %s", req.GetAction())
	}
	blocklist, err := cli.UpdateBlocklist(ctx, jid, action)
	if err != nil {
		return nil, err
	}
	return toJson(map[string]interface{}{
		"dhash": blocklist.DHash,
		"jids":  jidsToStrings(blocklist.JIDs),
	})
}

func parseJIDs(values []string) ([]types.JID, error) {
	jids := make([]types.JID, 0, len(values))
	for ind, value := range values {
		jid, err := types.ParseJID(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JID at index %d (%s): %w", ind, value, err)
		}
		jids = append(jids, jid)
	}
	return jids, nil
}

func jidsToStrings(jids []types.JID) []string {
	values := make([]string, len(jids))
	for i, jid := range jids {
		values[i] = jid.String()
	}
	return values
}
