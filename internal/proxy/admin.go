package proxy

import (
	"encoding/json"
	"net/http"
	"os"
)

func (p *ProxyHandler) AdminHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/clients", p.adminClients)
	mux.HandleFunc("/admin/stats", p.adminStats)
	mux.HandleFunc("/admin/cache/purge", p.adminPurge)
	mux.HandleFunc("/admin/cache/purge-all", p.adminPurgeAll)
	return p.withAdminAuth(mux)
}

func (p *ProxyHandler) withAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Admin-Key")
		if key == "" || key != os.Getenv("ADMIN_KEY") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (p *ProxyHandler) adminClients(w http.ResponseWriter, r *http.Request) {
	list := p.registry.All()
	json.NewEncoder(w).Encode(list)
}

func (p *ProxyHandler) adminStats(w http.ResponseWriter, r *http.Request) {
	// placeholder pour des stats simplifiées (dans la vraie version : agrégées depuis Prometheus)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","metrics":"check /metrics for prometheus"}`))
}

func (p *ProxyHandler) adminPurge(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "client required", 400)
		return
	}

	if c, ok := p.caches[clientID]; ok {
		c.Clear()
		w.WriteHeader(204)
		return
	}
	http.Error(w, "client cache not found", 404)
}

func (p *ProxyHandler) adminPurgeAll(w http.ResponseWriter, r *http.Request) {
	for _, c := range p.caches {
		c.Clear()
	}
	w.WriteHeader(204)
}
