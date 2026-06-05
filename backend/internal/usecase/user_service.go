package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api-vpn/internal/model"

	"github.com/gofrs/uuid/v5"
)

var ErrInvalidPlan = errors.New("invalid plan")
var ErrTrialAlreadyUsed = errors.New("trial already used")
var ErrTrialActiveSubscription = errors.New("active subscription exists")

const trialDuration = 24 * time.Hour
const trialMaxIPs = 2

type TrialActivationsRepo interface {
	HasUsed(ctx context.Context, telegramID int64) (bool, error)
	Record(ctx context.Context, telegramID int64) error
}

// PlanDays maps plan codes to duration in days.
var PlanDays = map[string]int{
	"1m":  30,
	"3m":  90,
	"6m":  180,
	"12m": 365,
}

var PlanLabels = map[string]string{
	"1m":  "1 месяц",
	"3m":  "3 месяца",
	"6m":  "6 месяцев",
	"12m": "12 месяцев",
}

type UsersRepo interface {
	Upsert(ctx context.Context, p model.UpsertUserParams) (model.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (model.User, error)
}

type SubscriptionsRepo interface {
	Create(ctx context.Context, p model.CreateSubscriptionParams) (model.Subscription, error)
	DeactivateActiveForUser(ctx context.Context, telegramUserID int64) error
	GetActiveForUser(ctx context.Context, telegramUserID int64, now time.Time) (model.Subscription, error)
	UpdateActiveClientUUID(ctx context.Context, telegramUserID int64, clientUUID string, now time.Time) error
}

type UserProfile struct {
	User         model.User
	Subscription *model.Subscription
	Client       *model.VPNClient
	Access       *model.XUIAccess
}

type UserService struct {
	users         UsersRepo
	subs          SubscriptionsRepo
	trials        TrialActivationsRepo
	clients       *VPNClients
	servers       *VPNServers
	xui           *XUIAccess
	defaultServer string
	defaultMaxIPs int
	now           func() time.Time
}

func NewUserService(
	users UsersRepo,
	subs SubscriptionsRepo,
	trials TrialActivationsRepo,
	clients *VPNClients,
	servers *VPNServers,
	xui *XUIAccess,
	defaultServer string,
	defaultMaxIPs int,
) *UserService {
	if defaultMaxIPs <= 0 {
		defaultMaxIPs = 2
	}
	if strings.TrimSpace(defaultServer) == "" {
		defaultServer = "default"
	}
	return &UserService{
		users:         users,
		subs:          subs,
		trials:        trials,
		clients:       clients,
		servers:       servers,
		xui:           xui,
		defaultServer: defaultServer,
		defaultMaxIPs: defaultMaxIPs,
		now:           time.Now,
	}
}

func (s *UserService) UpsertUser(ctx context.Context, p model.UpsertUserParams) (model.User, error) {
	return s.users.Upsert(ctx, p)
}

func (s *UserService) GetProfile(ctx context.Context, telegramID int64) (UserProfile, error) {
	out := UserProfile{}
	u, err := s.users.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return UserProfile{}, err
		}
		out.User = model.User{TelegramID: telegramID}
	} else {
		out.User = u
	}

	sub, err := s.subs.GetActiveForUser(ctx, telegramID, s.now())
	if err == nil {
		out.Subscription = &sub
	}

	client, err := s.clients.GetActiveByTelegramUserID(ctx, telegramID)
	if err == nil {
		out.Client = &client
		if s.xui != nil {
			if acc, err := s.xui.Get(ctx, client.ClientUUID); err == nil {
				out.Access = &acc
			}
		}
	}

	return out, nil
}

func (s *UserService) ActivateTrial(ctx context.Context, telegramID int64, profile model.UpsertUserParams) (UserProfile, error) {
	used, err := s.trials.HasUsed(ctx, telegramID)
	if err != nil {
		return UserProfile{}, err
	}
	if used {
		return UserProfile{}, ErrTrialAlreadyUsed
	}
	// Use activeErr, not err: err still holds HasUsed result (often nil) and must not be reused here.
	_, activeErr := s.clients.GetActiveByTelegramUserID(ctx, telegramID)
	if activeErr == nil {
		return UserProfile{}, ErrTrialActiveSubscription
	}
	if !errors.Is(activeErr, ErrNotFound) {
		return UserProfile{}, activeErr
	}

	profile.TelegramID = telegramID
	if _, err := s.users.Upsert(ctx, profile); err != nil {
		return UserProfile{}, err
	}
	if _, err := s.servers.GetActiveByID(ctx, s.defaultServer); err != nil {
		return UserProfile{}, ErrInvalidServer
	}

	_ = s.clients.DeactivateActiveByTelegramUserID(ctx, telegramID)

	now := s.now().UTC()
	expires := now.Add(trialDuration)

	note := telegramNote(profile)
	uid, err := uuid.NewV4()
	if err != nil {
		return UserProfile{}, err
	}
	tgID := telegramID
	client, err := s.clients.Create(ctx, model.CreateVPNClientParams{
		ClientUUID:     uid,
		ServerID:       s.defaultServer,
		TelegramUserID: &tgID,
		MaxIPs:         trialMaxIPs,
		KeyExpiresAt:   expires,
		Note:           &note,
	})
	if err != nil {
		return UserProfile{}, err
	}

	if s.xui == nil {
		return UserProfile{}, fmt.Errorf("xui access not configured")
	}
	access, err := s.xui.Provision(ctx, client.ClientUUID)
	if err != nil {
		return UserProfile{}, err
	}

	if err := s.trials.Record(ctx, telegramID); err != nil {
		return UserProfile{}, err
	}

	_ = s.subs.DeactivateActiveForUser(ctx, telegramID)
	uuidStr := client.ClientUUID.String()
	_, err = s.subs.Create(ctx, model.CreateSubscriptionParams{
		TelegramUserID: telegramID,
		PlanCode:       "trial",
		PlanLabel:      "Бесплатно 24 часа",
		StartsAt:       now,
		EndsAt:         client.KeyExpiresAt,
		ClientUUID:     &uuidStr,
		IsMock:         true,
	})
	if err != nil {
		return UserProfile{}, err
	}

	u, _ := s.users.GetByTelegramID(ctx, telegramID)
	sub, _ := s.subs.GetActiveForUser(ctx, telegramID, s.now())
	return UserProfile{User: u, Subscription: &sub, Client: &client, Access: &access}, nil
}

func (s *UserService) MockActivate(ctx context.Context, telegramID int64, planCode string, profile model.UpsertUserParams) (UserProfile, error) {
	planCode = strings.TrimSpace(strings.ToLower(planCode))
	days, ok := PlanDays[planCode]
	if !ok {
		return UserProfile{}, ErrInvalidPlan
	}
	label := PlanLabels[planCode]
	if label == "" {
		label = planCode
	}

	profile.TelegramID = telegramID
	if _, err := s.users.Upsert(ctx, profile); err != nil {
		return UserProfile{}, err
	}

	if _, err := s.servers.GetActiveByID(ctx, s.defaultServer); err != nil {
		return UserProfile{}, ErrInvalidServer
	}

	now := s.now().UTC()
	expires := now.Add(time.Duration(days) * 24 * time.Hour)

	var client model.VPNClient
	existing, err := s.clients.GetActiveByTelegramUserID(ctx, telegramID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return UserProfile{}, err
	}

	if errors.Is(err, ErrNotFound) {
		latest, lerr := s.clients.GetLatestByTelegramUserID(ctx, telegramID)
		if lerr == nil {
			client, err = s.clients.Extend(ctx, latest.ClientUUID, days)
			if err != nil {
				return UserProfile{}, err
			}
		} else if errors.Is(lerr, ErrNotFound) {
			note := telegramNote(profile)
			uid, err := uuid.NewV4()
			if err != nil {
				return UserProfile{}, err
			}
			tgID := telegramID
			client, err = s.clients.Create(ctx, model.CreateVPNClientParams{
				ClientUUID:     uid,
				ServerID:       s.defaultServer,
				TelegramUserID: &tgID,
				MaxIPs:         s.defaultMaxIPs,
				KeyExpiresAt:   expires,
				Note:           &note,
			})
			if err != nil {
				return UserProfile{}, err
			}
		} else {
			return UserProfile{}, lerr
		}
	} else {
		client, err = s.clients.Extend(ctx, existing.ClientUUID, days)
		if err != nil {
			return UserProfile{}, err
		}
	}

	if s.xui == nil {
		return UserProfile{}, fmt.Errorf("xui access not configured")
	}
	access, err := s.xui.Provision(ctx, client.ClientUUID)
	if err != nil {
		return UserProfile{}, err
	}

	_ = s.subs.DeactivateActiveForUser(ctx, telegramID)
	uuidStr := client.ClientUUID.String()
	_, err = s.subs.Create(ctx, model.CreateSubscriptionParams{
		TelegramUserID: telegramID,
		PlanCode:       planCode,
		PlanLabel:      label,
		StartsAt:       now,
		EndsAt:         client.KeyExpiresAt,
		ClientUUID:     &uuidStr,
		IsMock:         true,
	})
	if err != nil {
		return UserProfile{}, err
	}

	u, _ := s.users.GetByTelegramID(ctx, telegramID)
	sub, _ := s.subs.GetActiveForUser(ctx, telegramID, s.now())
	return UserProfile{User: u, Subscription: &sub, Client: &client, Access: &access}, nil
}

func (s *UserService) GetConfig(ctx context.Context, telegramID int64) (model.XUIAccess, error) {
	client, err := s.clients.GetActiveByTelegramUserID(ctx, telegramID)
	if err != nil {
		return model.XUIAccess{}, err
	}
	if s.xui == nil {
		return model.XUIAccess{}, fmt.Errorf("xui access not configured")
	}
	acc, err := s.xui.Get(ctx, client.ClientUUID)
	if err == nil {
		return acc, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return model.XUIAccess{}, err
	}
	return s.xui.Provision(ctx, client.ClientUUID)
}

func (s *UserService) RefreshConfig(ctx context.Context, telegramID int64) (model.XUIAccess, error) {
	client, err := s.clients.GetActiveByTelegramUserID(ctx, telegramID)
	if err != nil {
		return model.XUIAccess{}, err
	}
	if s.xui == nil {
		return model.XUIAccess{}, fmt.Errorf("xui access not configured")
	}

	if err := s.xui.RemoveFromPanel(ctx, client); err != nil {
		return model.XUIAccess{}, err
	}
	if err := s.clients.Deactivate(ctx, client.ClientUUID); err != nil && !errors.Is(err, ErrNotFound) {
		return model.XUIAccess{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return model.XUIAccess{}, err
	}
	tgID := telegramID
	newClient, err := s.clients.Create(ctx, model.CreateVPNClientParams{
		ClientUUID:     uid,
		ServerID:       client.ServerID,
		TelegramUserID: &tgID,
		MaxIPs:         client.MaxIPs,
		KeyExpiresAt:   client.KeyExpiresAt,
		Note:           client.Note,
	})
	if err != nil {
		return model.XUIAccess{}, err
	}

	access, err := s.xui.Provision(ctx, newClient.ClientUUID)
	if err != nil {
		return model.XUIAccess{}, err
	}

	_ = s.subs.UpdateActiveClientUUID(ctx, telegramID, newClient.ClientUUID.String(), s.now())

	return access, nil
}

func telegramNote(p model.UpsertUserParams) string {
	if p.Username != nil && strings.TrimSpace(*p.Username) != "" {
		return strings.TrimSpace(*p.Username)
	}
	return "tg_" + strconv.FormatInt(p.TelegramID, 10)
}
