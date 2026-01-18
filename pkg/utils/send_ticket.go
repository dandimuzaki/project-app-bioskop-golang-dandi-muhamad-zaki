package utils

import (
	"fmt"
	"strings"

	"github.com/project-app-bioskop-golang/internal/dto"
)

func SendTicket(ticket dto.TicketEmail) string {
	var seats []string
	for _, t := range ticket.Tickets {
		seats = append(seats, t.SeatCode)
	}

	return fmt.Sprintf(`
	<!DOCTYPE html>
		<html>
			<body style="margin:0; padding:0; font-family: Arial, Helvetica, sans-serif; background-color:#f6f6f6;">
				<table width="%v" cellpadding="0" cellspacing="0" style="padding:20px;">
					<tr>
						<td align="center">
							<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; padding:24px;">
						<!-- Header -->
						<tr>
							<td align="center" style="padding-bottom:16px;">
								<h2 style="margin:0; color:#222;">🎬 Your Ticket Is Ready</h2>
							</td>
						</tr>

						<!-- Greeting -->
						<tr>
							<td style="padding-bottom:12px; color:#333;">
								<p style="margin:0;">Hi <strong>%s</strong>,</p>
							</td>
						</tr>

						<!-- Confirmation -->
						<tr>
							<td style="padding-bottom:16px; color:#333;">
								<p style="margin:0;">
									Your payment was successful. Below is your official cinema ticket.
									Please present the QR code at the entrance gate.
								</p>
							</td>
						</tr>

						<!-- Ticket Details -->
						<tr>
							<td style="padding:16px; background:#f9f9f9; border-radius:6px; color:#333;">
								<p style="margin:4px 0;"><strong>Movie:</strong> %s</p>
								<p style="margin:4px 0;"><strong>Cinema:</strong> %s</p>
								<p style="margin:4px 0;"><strong>Studio:</strong> %s</p>
								<p style="margin:4px 0;"><strong>Date:</strong> %s</p>
								<p style="margin:4px 0;"><strong>Start Time:</strong> %s</p>
								<p style="margin:4px 0;"><strong>Seat:</strong> %s</p>
							</td>
						</tr>

						<!-- Notes -->
						<tr>
							<td style="padding-top:12px; color:#555; font-size:14px;">
								<ul style="padding-left:18px; margin:0;">
									<li>Please arrive at least <strong>15 minutes</strong> before the show.</li>
									<li>This ticket is valid for <strong>one-time entry only</strong>.</li>
									<li>Do not share your QR code with others.</li>
								</ul>
							</td>
						</tr>

						<!-- Footer -->
						<tr>
							<td style="padding-top:24px; color:#777; font-size:13px;">
								<p style="margin:0;">
									Enjoy the movie 🍿<br/>
									<strong>Cinema Booking Team</strong>
								</p>
								<p style="margin-top:8px; font-size:12px;">
									If you have any issues, please contact our support team.
								</p>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
			</body>
		</html>
	`, "100%", ticket.Profile.Name, ticket.Movie.Title, ticket.Cinema.Name, 
	ticket.Studio.Name, ticket.BookingDate, ticket.Screening.StartTime, 
	strings.Join(seats, ", "))
}