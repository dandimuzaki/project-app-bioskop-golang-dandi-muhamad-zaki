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
	tx, err := r.db.Begin(context.Background()); 
	if err != nil {
		r.Logger.Error("Error start db transaction: ", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	// Validate booking
	var b entity.Booking
	query := `SELECT id, status FROM bookings WHERE id = $1`
	err = tx.QueryRow(context.Background(), query, p.BookingID).Scan(&b.ID, &b.Status)
	if b.Status == "cancelled" {
		return nil, errors.New("booking is already cancelled or failed")
	}
	if b.Status == "paid" {
		return nil, nil // continue without creating new payment
	}
	if err != nil {
		r.Logger.Error("Error query get booking by id: ", zap.Error(err))
		return nil, err
	}

	// Create payment
	var payment entity.Payment
	query = `INSERT INTO payments (booking_id, payment_method, amount, status, created_at, updatedt_at)
	VALUES ($1, $2, $3, 'pending', NOW(), NOW())
	RETURNING id`
	err = tx.QueryRow(context.Background(), query, p.BookingID, p.PaymentMethod, p.Amount).Scan(&payment.ID)
	if err != nil {
		r.Logger.Error("Error query create payments: ", zap.Error(err))
		return nil, err
	}

	return &payment.ID, nil
}

func (r *paymentRepository) GetPaymentByID(id int) (*entity.Payment, error) {
	// Get payment detail
	query := `SELECT p.id, booking_id, pm.id, pm.name, amount, status, created_at, updatedt_at
	FROM payments p
	LEFT JOIN payment_methods pm ON p.payment_method = pm.id
	WHERE p.id = $1
	`
	var payment entity.Payment
	var paymentMethod entity.PaymentMethod
	err := r.db.QueryRow(context.Background(), query, id).Scan(&payment.ID, &payment.BookingID, &paymentMethod.ID, &paymentMethod.Name, &payment.Amount, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		r.Logger.Error("Error query get payment by id: ", zap.Error(err))
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(id int, p dto.UpdatePayment) error {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background()); 
	if err != nil {
		r.Logger.Error("Error start db transaction: ", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	// Update payment
	query := `UPDATE payments
	SET status = $1,
	transaction_id = $2,
	updated_at = NOW()
	WHERE id = $3`
	_, err = tx.Exec(context.Background(), query, p.Status, p.TransactionID, id)
	if err != nil {
		r.Logger.Error("Error query create payments: ", zap.Error(err))
		return err
	}

	return nil
}