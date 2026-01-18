package utils

import "sync"
 
type Metrics struct {
    mu          sync.Mutex
    EmailSent   int
    EmailFailed int
}
 
func (m *Metrics) Sent() {
    m.mu.Lock()
    m.EmailSent++
    m.mu.Unlock()
}
 
func (m *Metrics) Failed() {
    m.mu.Lock()
    m.EmailFailed++
    m.mu.Unlock()
}
 
 