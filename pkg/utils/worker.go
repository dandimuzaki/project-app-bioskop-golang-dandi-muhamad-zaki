package utils

import (
	"fmt"
	"sync"

	"github.com/project-app-bioskop-golang/internal/dto"
	"go.uber.org/zap"
)
 
type EmailJob struct {
  EmailContent dto.EmailRequest
	Config Configuration
	Log *zap.Logger
}

type TicketJob struct {
	Config Configuration
	Log *zap.Logger
  Data dto.TicketEmail
}
 
// Worker pool
func StartEmailWorkers(workerCount int, jobs <-chan EmailJob, stop <-chan struct{}, metrics *Metrics, wg *sync.WaitGroup) {
  wg.Add(workerCount)

  for i := 1; i <= workerCount; i++ {
    go func(id int) {
      defer wg.Done()

      for {
        select {
        case job, ok := <-jobs:
          if !ok {
            job.Log.Info(fmt.Sprintf("worker %v jobs channel closed", id))
            return
          }

          // send email
          if err := SendEmail(job.EmailContent, job.Config, job.Log); err != nil {
            job.Log.Error("Email failed to send", zap.Error(err))
            metrics.Failed()
            continue
          }

          metrics.Sent()

        case <-stop:
            fmt.Println("worker ", id, " received stop signal")
            return
        }
      }
    }(i)
  }
}

func StartTicketWorkers(workerCount int, jobs <-chan TicketJob, stop <-chan struct{}, metrics *Metrics, wg *sync.WaitGroup) {
  wg.Add(workerCount)

  for i := 1; i <= workerCount; i++ {
    go func(id int) {
      defer wg.Done()

      for {
        select {
        case job, ok := <-jobs:
          if !ok {
              job.Log.Info(fmt.Sprintf("worker %v jobs channel closed", id))
              return
          }

          var attachments []dto.Attachment
          for _, t := range job.Data.Tickets {
            png, err := GenerateQR(t.QRToken, job.Config, job.Log)
            if err != nil {
                job.Log.Error("Failed to generate QR", zap.Error(err))
                continue
            }

            attachment := dto.Attachment{
              FileName: fmt.Sprintf("ticket-%v-%v.png", i, job.Data.BookingDate),
              FileByte: png,
              ContentType: "image/png",
            }

            attachments = append(attachments, attachment)
          }

          fmt.Println(len(attachments))

          // Format email content
          body := SendTicket(job.Data)
          to := job.Data.Profile.Email
          subject := "Your Ticket Is Ready"
          content := dto.EmailRequest{
            To: to,
            Subject: subject,
            Body: body,
            Attachments: attachments,
          }

          // Send email
          if err := SendEmail(content, job.Config, job.Log); err != nil {
            job.Log.Error("Email failed to send", zap.Error(err))
            metrics.Failed()
            continue
          }

          metrics.Sent()

        case <-stop:
          fmt.Println("worker ", id, " received stop signal")
          return
        }
      }
    }(i)
  }
}