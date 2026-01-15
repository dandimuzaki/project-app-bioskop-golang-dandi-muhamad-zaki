package utils

// import (
// 	"fmt"
// 	"sync"
// )

// type PaymentJob struct {
//     ID int
// }

// // Worker pool
// func StartPaymentWorkers(workerCount int, jobs <-chan PaymentJob, stop <-chan struct{}, metrics *Metrics, wg *sync.WaitGroup) {
//     wg.Add(workerCount)

//     for i := 1; i <= workerCount; i++ {
//         go func(id int) {
//             defer wg.Done()

//             for {
//                 select {
//                 case job, ok := <-jobs:
//                     if !ok {
//                         fmt.Println("worker ", id, " jobs channel closed")
//                         return
//                     }

//                     // update payment

//                     if err := SendEmail(job.Email); err != nil {
//                         fmt.Println("email failed:", err)
//                         metrics.Failed()
//                         continue
//                     }

//                     metrics.Sent()

//                 case <-stop:
//                     fmt.Println("worker ", id, " received stop signal")
//                     return
//                 }
//             }
//         }(i)
//     }
// }
 
 