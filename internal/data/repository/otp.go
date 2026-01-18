package repository

import (
	"context"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type OTPRepository interface {
	Create(o entity.OTP) error
	Find(otp string) (*string, error)
	Update(otp string) error
}

type otpRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewOTPRepository(db database.PgxIface, log *zap.Logger) OTPRepository {
	return &otpRepository{
		db: db,
		Logger: log,
	}
}

func (r *otpRepository) Create(o entity.OTP) error {
	query := `
		INSERT INTO otp_codes (email, otp_hash, expired_at, created_at)
		VALUES ($1, $2, NOW() + interval '5 minute', NOW())
		RETURNING email, otp_hash, expired_at
	`
	err := r.db.QueryRow(context.Background(), query, o.Email, o.OTPHash).Scan(&o.Email, &o.OTPHash, &o.ExpiredAt)
	if err != nil {
		r.Logger.Error("Error query create otp: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *otpRepository) Find(email string) (*string, error) {
	query := `
		SELECT otp_hash
		FROM otp_codes
		WHERE email = $1 AND expired_at > NOW()
	`
	var otp string
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&otp,
	)

	if err != nil {
		return nil, err
	}

	return &otp, err
}

func (r *otpRepository) Update(otp string) error {
	query := `
		UPDATE otp_codes
		SET used_at = NOW()
		WHERE otp_hash = $1
	`
	_, err := r.db.Exec(context.Background(), query, otp)
	if err != nil {
		r.Logger.Error("Error query create otp: ", zap.Error(err))
		return err
	}

	return nil
}