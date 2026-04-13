package auth

import (
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"

	isql "github.com/masudur-rahman/styx/sql"
)

type sqlAuthRepo struct {
	db     isql.Engine
	logger logr.Logger
}

// NewSQLAuthRepository creates a new styx-based auth repository.
func NewSQLAuthRepository(db isql.Engine, logger logr.Logger) repos.AuthRepository {
	return &sqlAuthRepo{
		db:     db.Table(models.RefreshToken{}.TableName()),
		logger: logger,
	}
}

func (r *sqlAuthRepo) CreateRefreshToken(token *models.RefreshToken) error {
	r.logger.Infow("creating refresh token", "userID", token.UserID)
	_, err := r.db.InsertOne(token)
	return err
}

func (r *sqlAuthRepo) GetRefreshTokenByUUID(uuid string) (*models.RefreshToken, error) {
	r.logger.Infow("finding refresh token by UUID", "uuid", uuid)
	var token models.RefreshToken
	found, err := r.db.FindOne(&token, models.RefreshToken{TokenUUID: uuid})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrRefreshTokenNotFound{UUID: uuid}
	}
	return &token, nil
}

func (r *sqlAuthRepo) RevokeRefreshToken(uuid string) error {
	r.logger.Infow("revoking refresh token", "uuid", uuid)
	token, err := r.GetRefreshTokenByUUID(uuid)
	if err != nil {
		return err
	}
	token.Revoked = 1
	return r.db.ID(token.ID).UpdateOne(token)
}

func (r *sqlAuthRepo) RevokeAllUserTokens(userID int64) error {
	r.logger.Infow("revoking all tokens for user", "userID", userID)
	var tokens []models.RefreshToken
	filter := models.RefreshToken{UserID: userID}
	if err := r.db.FindMany(&tokens, filter); err != nil {
		return err
	}
	for i := range tokens {
		if tokens[i].Revoked == 1 {
			continue
		}
		tokens[i].Revoked = 1
		if err := r.db.ID(tokens[i].ID).UpdateOne(&tokens[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *sqlAuthRepo) DeleteExpiredTokens() error {
	r.logger.Infow("deleting expired refresh tokens")
	now := time.Now().Unix()
	var tokens []models.RefreshToken
	if err := r.db.FindMany(&tokens); err != nil {
		return err
	}
	for _, t := range tokens {
		if t.ExpiresAt < now {
			if err := r.db.DeleteOne(models.RefreshToken{ID: t.ID}); err != nil {
				return err
			}
		}
	}
	return nil
}
