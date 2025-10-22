package proxy

import (
	"io"
	"net/http"
	"time"

	"cachedproxy/internal/cache"
	"cachedproxy/internal/utils"
)

type ProxyHandler struct {
	logger   *utils.Logger
	registry *Registry
	caches   map[string]*cache.MemoryCache
	rl       *RateLimiter
}

func NewProxyHandler(logger *utils.Logger, reg *Registry) http.Handler {
	return &ProxyHandler{
		logger:   logger,
		registry: reg,
		caches:   make(map[string]*cache.MemoryCache),
		rl:       NewRateLimiter(50),
	}
}

var backendTransport = &http.Transport{
	MaxIdleConns:          200,
	MaxIdleConnsPerHost:   50,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var backendClient = &http.Client{
	Transport: backendTransport,
	Timeout:   15 * time.Second,
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = "acme" // default client to be accessible on a web browser
	}

	client, ok := p.registry.Get(clientID)
	if !ok {
		http.Error(w, "unknown client", http.StatusForbidden)
		return
	}

	ip := GetIP(r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
	if !p.rl.Allow(client.ID, ip, client.RateRPS) {
		w.Header().Set("Retry-After", "1")
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		p.logger.Infof("___________________")
		p.logger.Infof("WARNING ! -> Attack suspected. Blocked after too many requests. Sended alert on admin panel.")
		p.logger.Infof("[RL] block %s ip=%s %s %s", clientID, ip, r.Method, r.URL.Path)
		p.logger.Infof("___________________")
		return
	}

	// select cache for this client
	cacheClient, ok := p.caches[clientID]
	if !ok {
		cacheClient = cache.NewMemoryCache(client.TTL)
		p.caches[clientID] = cacheClient
	}

	if r.Method != http.MethodGet {
		resp := p.forwardRequest(client, r)
		writeResponse(w, resp)
		p.logger.Infof("[PASS] %s %s %s %d (%.2fms)", clientID, r.Method, r.URL.Path, resp.StatusCode, msSince(start))
		return
	}

	key := clientID + ":" + r.URL.String()

	if item, found := cacheClient.Get(key); found {
		writeHeaders(w, item.Header)
		w.WriteHeader(item.StatusCode)
		w.Write(item.Body)
		p.logger.Infof("[HIT] %s %s (%.2fms)", clientID, r.URL.Path, msSince(start))
		recordMetrics(clientID, item.StatusCode, true, time.Since(start))
		return
	}


	resp := p.forwardRequest(client, r)
	cacheClient.Set(key, resp.StatusCode, resp.Header, resp.Body)
	writeResponse(w, resp)
	p.logger.Infof("[MISS] %s %s (%.2fms)", clientID, r.URL.Path, time.Since(start).Seconds()*1000)
	recordMetrics(clientID, resp.StatusCode, false, time.Since(start))
}

type ResponseRecorder struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (p *ProxyHandler) forwardRequest(client *Client, r *http.Request) *ResponseRecorder {
	req, err := http.NewRequest(r.Method, client.BackendURL+r.URL.RequestURI(), r.Body)
	if err != nil {
		return &ResponseRecorder{StatusCode: http.StatusInternalServerError, Header: make(http.Header), Body: []byte(err.Error())}
	}

	// Copy headers except hop-by-hop
	for k, vals := range r.Header {
		switch k {
		case "Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization", "TE", "Trailers", "Transfer-Encoding", "Upgrade":
			continue
		}
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	resp, err := backendClient.Do(req)
	if err != nil {
		return &ResponseRecorder{StatusCode: http.StatusBadGateway, Header: make(http.Header), Body: []byte("backend unreachable")}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return &ResponseRecorder{StatusCode: resp.StatusCode, Header: resp.Header.Clone(), Body: body}
}

func writeResponse(w http.ResponseWriter, rec *ResponseRecorder) {
	writeHeaders(w, rec.Header)
	w.WriteHeader(rec.StatusCode)
	w.Write(rec.Body)
}

func writeHeaders(w http.ResponseWriter, headers map[string][]string) {
	for k, vals := range headers {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
}

func msSince(t time.Time) float64 { return float64(time.Since(t).Milliseconds()) }
