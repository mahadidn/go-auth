package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang-auth/internal/domain"
	"time"

	"github.com/google/uuid"
)

type tokenRepository struct {
	db *sql.DB
}

func NewPersonalAccessTokenRepository(db *sql.DB) domain.PersonalAccessTokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (repo *tokenRepository) Create(ctx context.Context, token *domain.PersonalAccessToken) error {
	
	query := `INSERT INTO personal_access_tokens (token, user_id, token_name, last_used_at, expires_at, created_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	idBytes, err := token.UserID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, 
		token.Token,
		idBytes,
		token.TokenName,
		token.LastUsedAt,
		token.ExpiresAt,
		token.CreatedAt,
	)
	
	return err
}

func (repo *tokenRepository) FindByToken(ctx context.Context, token string) (*domain.PersonalAccessToken, error) {
    query := `SELECT token, user_id, token_name, last_used_at, expires_at, created_at
              FROM personal_access_tokens WHERE token = ?`
    
    t := &domain.PersonalAccessToken{}
    var userBin []byte
    
    // 1. Buat penampung sementara untuk kolom yang bisa NULL
    var lastUsedAt sql.NullTime
    var expiresAt sql.NullTime

    err := repo.db.QueryRowContext(ctx, query, token).Scan(
        &t.Token,
        &userBin,
        &t.TokenName,
        &lastUsedAt, // Scan ke sql.NullTime
        &expiresAt,  // Scan ke sql.NullTime
        &t.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("token not found")
        }
        return nil, err
    }

    t.UserID, _ = uuid.FromBytes(userBin)

    // 2. Pindahkan data dari sql.NullTime ke pointer struct jika datanya valid (tidak NULL)
    if lastUsedAt.Valid {
        t.LastUsedAt = &lastUsedAt.Time
    }
    if expiresAt.Valid {
        t.ExpiresAt = &expiresAt.Time
    }

    return t, nil
}

func (repo *tokenRepository) UpdateLastUsed(ctx context.Context, token string) error {
	
	query := `UPDATE personal_access_tokens SET last_used_at = ? WHERE token = ?`

	now := time.Now()
	res, err := repo.db.ExecContext(ctx, query, &now, token)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("token not found")
	}
	return nil
}


func (repo *tokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM personal_access_tokens WHERE token = ?`
	_, err := repo.db.ExecContext(ctx, query, token)
	
	return err
}

func (repo *tokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM personal_access_tokens WHERE user_id = ?`
	
	userBin, _ := userID.MarshalBinary()
	_, err := repo.db.ExecContext(ctx, query, userBin)

	return err
}

