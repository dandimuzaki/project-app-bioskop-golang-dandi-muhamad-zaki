package repository

import (
	"context"
	"errors"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type PaymentRepository interface{
	Create(b dto.PaymentRequest) (*int, error)
	GetPaymentByID(id int) (*entity.Payment, error)
	Update(p dto.UpdatePayment) error
}

type paymentRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewPaymentRepository(db database.PgxIface, log *zap.Logger) PaymentRepository {
	return &paymentRepository{
		db: db,
		Logger: log,
	}
}

func (r *paymentRepository) Create(p dto.PaymentRequest) (*int, error) {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()

	// Validate booking
	var bookingStatus string
	query := `SELECT status
		FROM bookings
		WHERE id = $1
		FOR UPDATE`
	err = tx.QueryRow(context.Background(), query, p.BookingID).Scan(&bookingStatus)

	if err != nil {
		r.Logger.Error("Error query get booking: ", zap.Error(err))
		return nil, err
	}

	if bookingStatus == "cancelled" {
		return nil, errors.New("booking is already cancelled")
	}

	// Idempotency: reuse existing pending payment
	var existingPaymentID int
	query = `SELECT id
		FROM payments
		WHERE booking_id = $1 AND status = 'pending'`
	err = tx.QueryRow(context.Background(), query, p.BookingID).Scan(&existingPaymentID)

	if err == nil {
		_ = tx.Commit(context.Background())
		return &existingPaymentID, nil
	}

	// Create payment
	var payment entity.Payment
	query = `INSERT INTO payments (booking_id, payment_method, amount, status, created_at, updated_at)
	VALUES ($1, $2, $3, 'pending', NOW(), NOW())
	RETURNING id`
	err = tx.QueryRow(context.Background(), query, p.BookingID, p.PaymentMethod, p.Amount).Scan(&payment.ID)
	if err != nil {
		r.Logger.Error("Error query create payment: ", zap.Error(err))
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &payment.ID, nil
}

func (r *paymentRepository) GetPaymentByID(id int) (*entity.Payment, error) {
	// Get payment detail
	query := `SELECT p.id, booking_id, pm.name, amount, status, created_at, updatedt_at
	FROM payments p
	LEFT JOIN payment_methods pm ON p.payment_method = pm.id
	WHERE p.id = $1
	`
	var payment entity.Payment
	err := r.db.QueryRow(context.Background(), query, id).Scan(&payment.ID, &payment.BookingID, &payment.PaymentMethod, &payment.Amount, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		r.Logger.Error("Error query get payment by id: ", zap.Error(err))
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(p dto.UpdatePayment) error {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()

	// Update payment
	query := `UPDATE payments
	SET status = $1,
	transaction_id = $2,
	updated_at = NOW()
	WHERE id = $3`
	_, err = tx.Exec(context.Background(), query, p.Status, p.TransactionID, p.PaymentID)
	if err != nil {
		r.Logger.Error("Error query update payment: ", zap.Error(err))
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}