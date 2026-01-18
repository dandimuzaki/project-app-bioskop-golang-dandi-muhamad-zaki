package repository

import (
	"context"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type TicketRepository interface{
	GetByBookingID(id int) ([]dto.Ticket, error)
}

type ticketRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewTicketRepository(db database.PgxIface, log *zap.Logger) TicketRepository {
	return &ticketRepository{
		db: db,
		Logger: log,
	}
}

func (r *ticketRepository) GetByBookingID(id int) ([]dto.Ticket, error) {
	query := `SELECT s.seat_code, t.qr_token
	FROM tickets t
	LEFT JOIN seats s ON s.id = t.seat_id
	WHERE t.booking_id = $1
	`
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("Error query get tickets: ", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var tickets []dto.Ticket
	for rows.Next() {
		var t dto.Ticket
		err := rows.Scan(&t.SeatCode, &t.QRToken)
		if err != nil {
			r.Logger.Error("Error query get tickets: ", zap.Error(err))
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}