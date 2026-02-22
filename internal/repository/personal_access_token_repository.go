package repository

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "golang-auth/internal/domain"

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
    
    query := `INSERT INTO personal_access_tokens (id, token_hash, user_id, token_name, last_used_at, expires_at, created_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
    
    idUserBytes, _ := token.UserID.MarshalBinary()
    idTokenBytes, _ := token.ID.MarshalBinary()

    _, err := repo.db.ExecContext(ctx, query, 
        idTokenBytes,
        token.TokenHash,
        idUserBytes,
        token.TokenName,
        token.LastUsedAt,
        token.ExpiresAt,
        token.CreatedAt,
    )
    
    return err
}

func (repo *tokenRepository) FindByToken(ctx context.Context, token string) (*domain.PersonalAccessToken, error) {
    
    query := `SELECT id, token_hash, user_id, token_name, last_used_at, expires_at, created_at
              FROM personal_access_tokens WHERE token_hash = ?`
    
    t := &domain.PersonalAccessToken{}
    var idBin []byte 
    var userBin []byte
    
    // Buat penampung sementara untuk kolom yang bisa NULL
    var lastUsedAt sql.NullTime
    var expiresAt sql.NullTime

    err := repo.db.QueryRowContext(ctx, query, token).Scan(
        &idBin,
        &t.TokenHash,
        &userBin,
        &t.TokenName,
        &lastUsedAt,
        &expiresAt, 
        &t.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, sql.ErrNoRows 
        }
        return nil, err
    }

    // Convert byte ke UUID
    t.ID, _ = uuid.FromBytes(idBin)
    t.UserID, _ = uuid.FromBytes(userBin)

    // Pindahkan data dari sql.NullTime ke pointer struct jika datanya valid (tidak NULL)
    if lastUsedAt.Valid {
        t.LastUsedAt = &lastUsedAt.Time
    }
    if expiresAt.Valid {
        t.ExpiresAt = &expiresAt.Time
    }

    return t, nil
}

func (repo *tokenRepository) UpdateLastUsed(ctx context.Context, token string) error {
    
    query := `UPDATE personal_access_tokens SET last_used_at = ? WHERE token_hash = ?`

    now := time.Now()
    res, err := repo.db.ExecContext(ctx, query, now, token) 
    if err != nil {
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        return errors.New("token not found")
    }
    return nil
}

func (repo*tokenRepository) Delete(ctx context.Context, token string) error {
    query := `DELETE FROM personal_access_tokens WHERE token_hash = ?`
    _, err := repo.db.ExecContext(ctx, query, token)
    
    return err
}

func (repo *tokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
    query := `DELETE FROM personal_access_tokens WHERE user_id = ?`
    
    userBin, _ := userID.MarshalBinary()
    _, err := repo.db.ExecContext(ctx, query, userBin)

    return err
}