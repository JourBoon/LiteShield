package proxy

import (
	"net"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter limite le débit par client et par IP (par client)
type RateLimiter struct {
	mu         sync.Mutex
	perClient  map[string]*rate.Limiter            // clientID -> limiter global
	perIP      map[string]map[string]*rate.Limiter // clientID -> ip -> limiter
	defaultRPS int
	ipRPS      int
	ipBurst    int
}

func NewRateLimiter(defaultRPS int) *RateLimiter {
	if defaultRPS <= 0 {
		defaultRPS = 50
	}
	return &RateLimiter{
		perClient:  make(map[string]*rate.Limiter),
		perIP:      make(map[string]map[string]*rate.Limiter),
		defaultRPS: defaultRPS,
		ipRPS:      10,  // défaut: 10 req/s par IP
		ipBurst:    20,  // burst par IP
	}
}

func (r *RateLimiter) clientLimiter(clientID string, rps int) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rps <= 0 {
		rps = r.defaultRPS
	}
	lim, ok := r.perClient[clientID]
	if !ok {
		lim = rate.NewLimiter(rate.Limit(rps), rps*2) // burst = 2x
		r.perClient[clientID] = lim
	}
	return lim
}

func (r *RateLimiter) ipLimiter(clientID, ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.perIP[clientID]
	if !ok {
		m = make(map[string]*rate.Limiter)
		r.perIP[clientID] = m
	}
	lim, ok := m[ip]
	if !ok {
		lim = rate.NewLimiter(rate.Limit(r.ipRPS), r.ipBurst)
		m[ip] = lim
	}
	return lim
}

// Allow vérifie les limites client + IP
func (r *RateLimiter) Allow(clientID, ip string, clientRPS int) bool {
	cl := r.clientLimiter(clientID, clientRPS)
	il := r.ipLimiter(clientID, ip)
	return cl.Allow() && il.Allow()
}

// GetIP récupère l’IP du client (X-Forwarded-For si présent)
func GetIP(remoteAddr string, xff string) string {
	// Priorité au premier IP du X-Forwarded-For
	if xff != "" {
		// XFF peut contenir "ip1, ip2, ip3"
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return trimSpace(xff[:i])
			}
		}
		return trimSpace(xff)
	}
	// Sinon, parse RemoteAddr "ip:port"
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func trimSpace(s string) string {
	// petite trim sans strings import pour rester minimal
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s) - 1
	for end >= 0 && (s[end] == ' ' || s[end] == '\t') {
		end--
	}
	if end < start {
		return ""
	}
	return s[start : end+1]
}

// Optionnel: nettoyage périodique des map IP mais on peut rajouter un TTL d’inactivité
// (Non nécessaire pour le MVP)
