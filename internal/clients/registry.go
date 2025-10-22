package clients

import (
	"sync"
	"time"
)

type Client struct {
	ID         string
	APIKey     string
	BackendURL string
	DefaultTTL time.Duration
	RateRPS    int // plafond global client
}

type Registry struct {
	byID map[string]*Client
	mu   sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{byID: make(map[string]*Client)}
}
func (r *Registry) Get(id string) (*Client, bool) {
	r.mu.RLock(); defer r.mu.RUnlock()
	c, ok := r.byID[id]; return c, ok
}
func (r *Registry) Set(c *Client) {
	r.mu.Lock(); defer r.mu.Unlock()
	r.byID[c.ID] = c
}
func (r *Registry) All() []*Client {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*Client, 0, len(r.byID))
	for _, c := range r.byID { out = append(out, c) }
	return out
}
