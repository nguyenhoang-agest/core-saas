package service

import (
	"context"
	"strings"
	"time"

	tenantapi "github.com/example/core-backend/modules/tenant/api"
	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/repository"
	"github.com/example/core-backend/pkg/auth"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tenantRecord struct {
	ID          string     `gorm:"column:id;primaryKey"`
	Name        string     `gorm:"column:name"`
	Description string     `gorm:"column:description"`
	Status      string     `gorm:"column:status"`
	IsolateMode string     `gorm:"column:isolate_mode"`
	IsHost      bool       `gorm:"column:is_host"`
	CreateAt    time.Time  `gorm:"column:create_at"`
	UpdateAt    time.Time  `gorm:"column:update_at"`
	DeleteAt    *time.Time `gorm:"column:delete_at"`
}

func (tenantRecord) TableName() string {
	return "tenants"
}

type IAuthService interface {
	RegisterOrg(ctx context.Context, req *dto.RegisterOrgRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	GetCurrentUser(ctx context.Context, tenantID, userID string, isHost bool) (*dto.CurrentUserResponse, error)
}

type authService struct {
	db         *gorm.DB
	userRepo   repository.IUserRepository
	bus        *eventbus.Bus
	resolver   *CachedPermissionResolver
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	db *gorm.DB,
	userRepo repository.IUserRepository,
	bus *eventbus.Bus,
	resolver *CachedPermissionResolver,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) IAuthService {
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	return &authService{
		db:         db,
		userRepo:   userRepo,
		bus:        bus,
		resolver:   resolver,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) RegisterOrg(ctx context.Context, req *dto.RegisterOrgRequest) (*dto.AuthResponse, error) {
	if req == nil {
		return nil, exceptions.NewBadRequestError(response_code.InvalidPayload, nil)
	}

	orgName := strings.TrimSpace(req.OrgName)
	if len(orgName) < 3 {
		return nil, exceptions.NewBadRequestError(errcode.OrgNameRequired, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, exceptions.NewBadRequestError(errcode.UserNameRequired, nil)
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		return nil, exceptions.NewBadRequestError(errcode.UserEmailRequired, nil)
	}

	if len(req.Password) < 8 {
		return nil, exceptions.NewBadRequestError(errcode.PasswordTooShort, nil)
	}

	// 1. Check if tenant organization name already exists
	tenantID := uuid.NewString()
	now := time.Now().UTC()

	if s.db != nil {
		var existingTenant tenantRecord
		err := s.db.WithContext(ctx).Table("tenants").
			Where("name = ? AND delete_at IS NULL", orgName).
			First(&existingTenant).Error
		if err == nil {
			return nil, exceptions.NewBadRequestError(errcode.OrgAlreadyExists, nil)
		}

		newTenant := tenantRecord{
			ID:          tenantID,
			Name:        orgName,
			Status:      "active",
			IsolateMode: "shared",
			IsHost:      false,
			CreateAt:    now,
			UpdateAt:    now,
		}

		if err := s.db.WithContext(ctx).Table("tenants").Create(&newTenant).Error; err != nil {
			return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
		}

		// 4. Seed isolated system roles for this new tenant
		if err := SyncSystemRoles(ctx, s.db, tenantID, false); err != nil {
			return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
		}
	}

	// 2. Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
	}

	// 5. Create Owner user in users table and user_roles join table
	ownerUser := &model.User{
		BaseModel: model.BaseModel{
			Id:       uuid.NewString(),
			CreateAt: now,
			UpdateAt: now,
		},
		TenantID:     tenantID,
		Email:        email,
		Name:         name,
		Role:         string(rbac.RoleOwner),
		Status:       model.StatusActive,
		PasswordHash: passwordHash,
		LastLoginAt:  &now,
	}

	if err := s.userRepo.Create(ctx, ownerUser); err != nil {
		return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
	}

	// 6. Emit event on eventbus if available
	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, tenantapi.TenantCreatedTopic, tenantapi.TenantCreatedEvent{
			ID:          tenantID,
			Name:        orgName,
			IsolateMode: "shared",
			IsHost:      false,
			OwnerEmail:  email,
			OwnerName:   name,
		})
	}

	// 7. Issue JWT Token Pair
	tokens, err := auth.GenerateTokenPair(
		s.jwtSecret,
		ownerUser.Id,
		tenantID,
		ownerUser.Email,
		ownerUser.Role,
		false,
		s.accessTTL,
		s.refreshTTL,
	)
	if err != nil {
		return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
	}

	return &dto.AuthResponse{
		User:         ownerUser,
		TenantID:     tenantID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	if req == nil {
		return nil, exceptions.NewUnauthorizeError(response_code.WrongCredentials, nil)
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		return nil, exceptions.NewUnauthorizeError(response_code.WrongCredentials, nil)
	}

	var user *model.User

	if req.TenantID != "" {
		u, err := s.userRepo.GetByEmail(ctx, req.TenantID, email)
		if err != nil || u == nil {
			return nil, exceptions.NewUnauthorizeError(response_code.WrongCredentials, nil)
		}
		user = u
	} else {
		users, err := s.userRepo.GetByEmailAnyTenant(ctx, email)
		if err != nil || len(users) == 0 {
			return nil, exceptions.NewUnauthorizeError(response_code.WrongCredentials, nil)
		}

		if len(users) == 1 {
			user = &users[0]
		} else {
			// Multiple workspace accounts: pick the active one matching the password
			var matched *model.User
			for i := range users {
				if users[i].Status == model.StatusActive && auth.CheckPassword(users[i].PasswordHash, req.Password) {
					matched = &users[i]
					break
				}
			}
			if matched == nil {
				return nil, exceptions.NewUnauthorizeError(response_code.WrongCredentials, nil)
			}
			user = matched
		}
	}

	// 1. Password verification
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		return nil, exceptions.NewUnauthorizeError(response_code.WrongCredentials, nil)
	}

	// 2. User status verification
	if user.Status != model.StatusActive {
		return nil, exceptions.NewForbiddenError(response_code.ForbiddenGeneral)
	}

	// 3. Tenant verification
	isHost := false
	if s.db != nil {
		var tenant tenantRecord
		if err := s.db.WithContext(ctx).Table("tenants").
			Where("id = ? AND delete_at IS NULL", user.TenantID).
			First(&tenant).Error; err != nil {
			return nil, exceptions.NewForbiddenError(response_code.ForbiddenGeneral)
		}

		if tenant.Status != "active" {
			return nil, exceptions.NewForbiddenError(response_code.ForbiddenGeneral)
		}
		isHost = tenant.IsHost
	}

	// 4. Update last login timestamp
	now := time.Now().UTC()
	_ = s.userRepo.UpdateLastLogin(ctx, user.TenantID, user.Id, now)
	user.LastLoginAt = &now

	// 5. Generate token pair
	tokens, err := auth.GenerateTokenPair(
		s.jwtSecret,
		user.Id,
		user.TenantID,
		user.Email,
		user.Role,
		isHost,
		s.accessTTL,
		s.refreshTTL,
	)
	if err != nil {
		return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
	}

	return &dto.AuthResponse{
		User:         user,
		TenantID:     user.TenantID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	if req == nil || strings.TrimSpace(req.RefreshToken) == "" {
		return nil, exceptions.NewUnauthorizeError(response_code.RefreshTokenExpired, nil)
	}

	claims, err := auth.ValidateToken(s.jwtSecret, strings.TrimSpace(req.RefreshToken))
	if err != nil || claims.TokenType != auth.TokenTypeRefresh {
		return nil, exceptions.NewUnauthorizeError(response_code.RefreshTokenExpired, nil)
	}

	user, err := s.userRepo.GetByID(ctx, claims.TenantID, claims.Subject)
	if err != nil || user == nil || user.Status != model.StatusActive {
		return nil, exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
	}

	isHost := false
	if s.db != nil {
		var tenant tenantRecord
		if err := s.db.WithContext(ctx).Table("tenants").
			Where("id = ? AND delete_at IS NULL", user.TenantID).
			First(&tenant).Error; err != nil || tenant.Status != "active" {
			return nil, exceptions.NewForbiddenError(response_code.ForbiddenGeneral)
		}
		isHost = tenant.IsHost
	}

	tokens, err := auth.GenerateTokenPair(
		s.jwtSecret,
		user.Id,
		user.TenantID,
		user.Email,
		user.Role,
		isHost,
		s.accessTTL,
		s.refreshTTL,
	)
	if err != nil {
		return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
	}

	return &dto.AuthResponse{
		User:         user,
		TenantID:     user.TenantID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

func (s *authService) GetCurrentUser(ctx context.Context, tenantID, userID string, isHost bool) (*dto.CurrentUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, tenantID, userID)
	if err != nil || user == nil {
		return nil, exceptions.NewNotFoundError(errcode.NotFoundUser)
	}

	tenantName := tenantID
	if s.db != nil {
		var tenant tenantRecord
		if err := s.db.WithContext(ctx).Table("tenants").
			Where("id = ? AND delete_at IS NULL", tenantID).
			First(&tenant).Error; err == nil && tenant.Name != "" {
			tenantName = tenant.Name
		}
	}

	var permissions []string
	if s.resolver != nil {
		perms, err := s.resolver.ResolvePermissions(ctx, tenantID, userID, isHost)
		if err == nil {
			permissions = perms
		}
	}
	if permissions == nil {
		permissions = user.EffectivePermissions()
	}

	return &dto.CurrentUserResponse{
		User: user,
		Tenant: &dto.TenantSummary{
			ID:     tenantID,
			Name:   tenantName,
			IsHost: isHost,
		},
		Permissions: permissions,
	}, nil
}
