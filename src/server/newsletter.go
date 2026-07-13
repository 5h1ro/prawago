package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/devlikeapro/gows/gows"
	"github.com/devlikeapro/gows/media"
	__ "github.com/devlikeapro/gows/proto"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

const (
	newsletterMutationUpdate      = "7150902998257522"
	newsletterMutationDelete      = "8316537688363079"
	newsletterMutationChangeOwner = "7341777602580933"
	newsletterMutationDemote      = "6551828931592903"
)

func toNewsletter(n *types.NewsletterMetadata) *__.Newsletter {
	var picture string
	if n.ThreadMeta.Picture != nil {
		picture = n.ThreadMeta.Picture.URL
		if picture == "" {
			picture = n.ThreadMeta.Picture.DirectPath
		}
	}

	var preview string
	preview = n.ThreadMeta.Preview.URL
	if preview == "" {
		preview = n.ThreadMeta.Preview.DirectPath
	}
	var role string
	if n.ViewerMeta != nil {
		role = string(n.ViewerMeta.Role)
	}
	return &__.Newsletter{
		Id:              n.ID.String(),
		Name:            n.ThreadMeta.Name.Text,
		Description:     n.ThreadMeta.Description.Text,
		Invite:          n.ThreadMeta.InviteCode,
		Picture:         picture,
		Preview:         preview,
		Verified:        n.ThreadMeta.VerificationState == types.NewsletterVerificationStateVerified,
		Role:            role,
		SubscriberCount: int64(n.ThreadMeta.SubscriberCount),
	}
}

func (s *Server) GetSubscribedNewsletters(ctx context.Context, req *__.NewsletterListRequest) (*__.NewsletterList, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	resp, err := cli.GetSubscribedNewsletters(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*__.Newsletter, len(resp))
	for i, n := range resp {
		list[i] = toNewsletter(n)
	}
	return &__.NewsletterList{Newsletters: list}, nil
}

func (s *Server) GetNewsletterInfo(ctx context.Context, req *__.NewsletterInfoRequest) (result *__.Newsletter, err error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	id := req.GetId()
	var resp *types.NewsletterMetadata
	if gows.HasNewsletterSuffix(id) {
		jid, err := types.ParseJID(id)
		if err != nil {
			return nil, err
		}
		resp, err = cli.GetNewsletterInfo(ctx, jid)
	} else {
		resp, err = cli.GetNewsletterInfoWithInvite(ctx, id)
	}
	if err != nil {
		return nil, err
	}
	return toNewsletter(resp), nil
}

func (s *Server) CreateNewsletter(ctx context.Context, req *__.CreateNewsletterRequest) (*__.Newsletter, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	params := whatsmeow.CreateNewsletterParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Picture:     req.GetPicture(),
	}
	resp, err := cli.CreateNewsletter(ctx, params)
	if err != nil {
		return nil, err
	}
	return toNewsletter(resp), nil

}

func (s *Server) NewsletterToggleMute(ctx context.Context, req *__.NewsletterToggleMuteRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	if !gows.IsNewsletter(jid) {
		return nil, errors.New("invalid jid, not a newsletter")
	}
	err = cli.NewsletterToggleMute(ctx, jid, req.GetMute())
	if err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

// NewsletterToggleFollow
func (s *Server) NewsletterToggleFollow(ctx context.Context, req *__.NewsletterToggleFollowRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	if !gows.IsNewsletter(jid) {
		return nil, errors.New("invalid jid, not a newsletter")
	}
	if req.Follow {
		err = cli.FollowNewsletter(ctx, jid)
	} else {
		err = cli.UnfollowNewsletter(ctx, jid)
	}
	if err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

func (s *Server) NewsletterUpdateName(ctx context.Context, req *__.JidStringRequest) (*__.Json, error) {
	return s.newsletterUpdate(ctx, req.GetSession().GetId(), req.GetJid(), map[string]any{
		"name":     req.GetValue(),
		"settings": nil,
	})
}

func (s *Server) NewsletterUpdateDescription(ctx context.Context, req *__.JidStringRequest) (*__.Json, error) {
	return s.newsletterUpdate(ctx, req.GetSession().GetId(), req.GetJid(), map[string]any{
		"description": req.GetValue(),
		"settings":    nil,
	})
}

func (s *Server) NewsletterReactionMode(ctx context.Context, req *__.JidStringRequest) (*__.Json, error) {
	return s.newsletterUpdate(ctx, req.GetSession().GetId(), req.GetJid(), map[string]any{
		"settings": map[string]any{
			"reaction_codes": map[string]any{
				"value": req.GetValue(),
			},
		},
	})
}

func (s *Server) NewsletterUpdatePicture(ctx context.Context, req *__.NewsletterPictureRequest) (*__.Json, error) {
	picture, err := media.ProfilePicture(req.GetPicture())
	if err != nil {
		return nil, err
	}
	return s.newsletterUpdate(ctx, req.GetSession().GetId(), req.GetJid(), map[string]any{
		"picture":  base64.StdEncoding.EncodeToString(picture),
		"settings": nil,
	})
}

func (s *Server) NewsletterRemovePicture(ctx context.Context, req *__.JidRequest) (*__.Json, error) {
	return s.newsletterUpdate(ctx, req.GetSession().GetId(), req.GetJid(), map[string]any{
		"picture":  "",
		"settings": nil,
	})
}

func (s *Server) NewsletterDelete(ctx context.Context, req *__.JidRequest) (*__.Json, error) {
	return s.newsletterMex(ctx, req.GetSession().GetId(), newsletterMutationDelete, req.GetJid(), map[string]any{}, "")
}

func (s *Server) NewsletterChangeOwner(ctx context.Context, req *__.NewsletterUserRequest) (*__.Json, error) {
	return s.newsletterMex(ctx, req.GetSession().GetId(), newsletterMutationChangeOwner, req.GetJid(), map[string]any{
		"user_id": req.GetUserLid(),
	}, "")
}

func (s *Server) NewsletterDemote(ctx context.Context, req *__.NewsletterUserRequest) (*__.Json, error) {
	return s.newsletterMex(ctx, req.GetSession().GetId(), newsletterMutationDemote, req.GetJid(), map[string]any{
		"user_id": req.GetUserLid(),
	}, "")
}

func (s *Server) NewsletterSubscribeLiveUpdates(ctx context.Context, req *__.JidRequest) (*__.OptionalUInt64, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	duration, err := cli.NewsletterSubscribeLiveUpdates(ctx, jid)
	if err != nil {
		return nil, err
	}
	return &__.OptionalUInt64{Value: uint64(duration.Seconds())}, nil
}

func (s *Server) NewsletterMarkViewed(ctx context.Context, req *__.NewsletterMessageIdsRequest) (*__.Empty, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	serverIDs := make([]types.MessageServerID, len(req.GetServerIds()))
	for i, id := range req.GetServerIds() {
		serverIDs[i] = types.MessageServerID(id)
	}
	if err := cli.NewsletterMarkViewed(ctx, jid, serverIDs); err != nil {
		return nil, err
	}
	return &__.Empty{}, nil
}

func (s *Server) newsletterUpdate(ctx context.Context, sessionID, jid string, updates map[string]any) (*__.Json, error) {
	return s.newsletterMex(ctx, sessionID, newsletterMutationUpdate, jid, map[string]any{
		"updates": updates,
	}, "xwa2_newsletter_update")
}

func (s *Server) newsletterMex(ctx context.Context, sessionID, queryID, jid string, variables map[string]any, dataPath string) (*__.Json, error) {
	cli, err := s.Sm.Get(sessionID)
	if err != nil {
		return nil, err
	}
	parsed, err := types.ParseJID(jid)
	if err != nil {
		return nil, err
	}
	if !gows.IsNewsletter(parsed) {
		return nil, errors.New("invalid jid, not a newsletter")
	}
	variables["newsletter_id"] = parsed.String()
	data, err := cli.SendMexIQ(ctx, queryID, variables, dataPath)
	if err != nil {
		return nil, err
	}
	return &__.Json{Data: string(data)}, nil
}

func (s *Server) GetNewsletterMessagesByInvite(ctx context.Context, req *__.GetNewsletterMessagesByInviteRequest) (*__.Json, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	code := req.Invite
	params := &gows.GetNewsletterMessagesByInviteParams{
		Count: int(req.Limit),
	}
	newsletterMessages, err := cli.GetNewsletterMessagesByInvite(code, params)
	if err != nil {
		return nil, fmt.Errorf("error getting newsletter messages: %w", err)
	}
	response, err := toJson(newsletterMessages)
	if err != nil {
		return nil, fmt.Errorf("error marshaling messages: %w", err)
	}
	return response, nil
}

func (s *Server) GetNewsletterMessages(ctx context.Context, req *__.GetNewsletterMessagesRequest) (*__.JsonList, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	if !gows.IsNewsletter(jid) {
		return nil, errors.New("invalid jid, not a newsletter")
	}
	messages, err := cli.GetNewsletterMessages(ctx, jid, &whatsmeow.GetNewsletterMessagesParams{
		Count:  int(req.GetLimit()),
		Before: types.MessageServerID(req.GetBefore()),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting newsletter messages: %w", err)
	}
	return toJsonList(messages)
}

func (s *Server) GetNewsletterMessageUpdates(ctx context.Context, req *__.GetNewsletterMessageUpdatesRequest) (*__.JsonList, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}
	jid, err := types.ParseJID(req.GetJid())
	if err != nil {
		return nil, err
	}
	if !gows.IsNewsletter(jid) {
		return nil, errors.New("invalid jid, not a newsletter")
	}
	var since time.Time
	if req.GetSince() != 0 {
		since = time.Unix(int64(req.GetSince()), 0)
	}
	messages, err := cli.GetNewsletterMessageUpdates(ctx, jid, &whatsmeow.GetNewsletterUpdatesParams{
		Count: int(req.GetLimit()),
		Since: since,
		After: types.MessageServerID(req.GetAfter()),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting newsletter message updates: %w", err)
	}
	return toJsonList(messages)
}

func (s *Server) SearchNewslettersByView(ctx context.Context, req *__.SearchNewslettersByViewRequest) (*__.NewsletterSearchPageResult, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}

	params := gows.SearchNewsletterByViewParams{
		Page: gows.SearchPageParams{
			Count:       int(req.Page.Limit),
			StartCursor: req.Page.StartCursor,
		},
		View:       req.View,
		Categories: req.Categories,
		Countries:  req.Countries,
	}
	resp, err := cli.SearchNewsletterByView(params)
	if err != nil {
		return nil, fmt.Errorf("error searching newsletters: %w", err)
	}
	result := searchResultToGrpcResponse(resp)
	return result, nil
}

func searchResultToGrpcResponse(resp *gows.SearchNewsletterResult) *__.NewsletterSearchPageResult {
	newsletters := make([]*__.Newsletter, len(resp.Newsletters))
	for i, n := range resp.Newsletters {
		newsletters[i] = toNewsletter(n)
	}
	page := &__.SearchPageResult{
		StartCursor:     resp.Page.StartCursor,
		EndCursor:       resp.Page.EndCursor,
		HasNextPage:     resp.Page.HasNextPage,
		HasPreviousPage: resp.Page.HasPreviousPage,
	}
	result := &__.NewsletterSearchPageResult{
		Newsletters: &__.NewsletterList{Newsletters: newsletters},
		Page:        page,
	}
	return result
}

func (s *Server) SearchNewslettersByText(ctx context.Context, req *__.SearchNewslettersByTextRequest) (*__.NewsletterSearchPageResult, error) {
	cli, err := s.Sm.Get(req.GetSession().GetId())
	if err != nil {
		return nil, err
	}

	params := gows.SearchNewsletterByTextParams{
		Page: gows.SearchPageParams{
			Count:       int(req.Page.Limit),
			StartCursor: req.Page.StartCursor,
		},
		Text:       req.Text,
		Categories: req.Categories,
	}
	resp, err := cli.SearchNewsletterByText(params)
	if err != nil {
		return nil, fmt.Errorf("error searching newsletters: %w", err)
	}
	result := searchResultToGrpcResponse(resp)
	return result, nil
}
