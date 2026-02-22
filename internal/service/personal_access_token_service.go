package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang-auth/internal/domain"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type personalAccessTokenService struct {
	personalAccessTokenRepository domain.PersonalAccessTokenRepository
	db *sql.DB
	validate *validator.Validate
}

func NewPersonalAccessTokenService(personalAccessTokenRepository domain.PersonalAccessTokenRepository, db *sql.DB, validate *validator.Validate) domain.PersonalAccessTokenService {
	return &personalAccessTokenService{
		personalAccessTokenRepository: personalAccessTokenRepository,
		db: db,
		validate: validate,
	}
}

func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (service *personalAccessTokenService) Create(ctx context.Context, req domain.PersonalAccessTokenRequest) (string, time.Time, error) {
	// validasi struct
	if err := service.validate.Struct(req); err != nil {
		return "", time.Time{}, err
	}

	// generate raw token
	// 32 bytes akan menghasilkan 64 karakter string hex
	randomStr, err := generateSecureToken(32)
	if err != nil {
		return "", time.Time{}, err
	}

	// tambahkan prefix agar token mudah diidentifikasi
	rawToken := "auth_pat_" + randomStr

	// hash token pakai SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	// siapkan data untuk disimpan ke database
	tokenID, _ := uuid.NewV7()
	now := time.Now()

	// expires at
	expiresAt := time.Now().Add(time.Hour * 48) // set token expire 2 hari

	entity := &domain.PersonalAccessToken{
		ID: tokenID,
		TokenHash: hashedToken,
		UserID: req.UserID,
		TokenName: req.TokenName,
		CreatedAt: now,
		ExpiresAt: &expiresAt,
	}

	// simpan ke database
	err = service.personalAccessTokenRepository.Create(ctx, entity)
	if err != nil {
		return "", time.Time{}, err
	}
	
	return rawToken, expiresAt, nil
}

func (service *personalAccessTokenService) FindByToken(ctx context.Context, token string) (*domain.PersonalAccessToken, error) {
	
	// hash ulang token asli yg dikirim oleh user
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))


	res, err := service.personalAccessTokenRepository.FindByToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	// cek token sudah expired atau belum
	if res.ExpiresAt != nil && time.Now().After(*res.ExpiresAt){
		// hapus token dari DB karna sudah basi
		_ = service.personalAccessTokenRepository.Delete(ctx, hashedToken)
		return nil, errors.New("token expired")
	}
	
	return res, nil
}

func (service *personalAccessTokenService) Delete(ctx context.Context, token string) error {
	
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	err := service.personalAccessTokenRepository.Delete(ctx, hashedToken)
	return err
}

func (service *personalAccessTokenService) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	
	err := service.personalAccessTokenRepository.DeleteByUserID(ctx, userID)
	return err
}

func (service *personalAccessTokenService) UpdateLastUsed(ctx context.Context, token string) error {
	
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	err := service.personalAccessTokenRepository.UpdateLastUsed(ctx, hashedToken)
	return err
}
