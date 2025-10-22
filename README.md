# ⚡ LiteShield Proxy

**LiteShield** is a lightweight, high-performance **reverse proxy** built in Go with **caching, DDoS mitigation, and client isolation** features — designed for developers, agencies, and SaaS businesses who need a fast and secure edge layer for their web apps.

---

## 🚀 Overview

LiteShield acts as an intelligent reverse proxy between your users and your backend services.  
It automatically caches responses, protects from excessive requests, and isolates traffic per client — offering a simple, modular way to scale and secure your web infrastructure.

```
[User] → [LiteShield Proxy] → [Backend Server / API]
```

---

## 🧠 Key Features

| Feature | Description |
|----------|-------------|
| ⚡ **Smart Caching** | Automatically caches GET responses with TTL-based invalidation. Reduces backend load and improves latency. |
| 🧱 **Client Isolation** | Each client has a separate cache namespace — no data leaks, perfect for multi-tenant systems. |
| 🛡 **DDoS Mitigation** | Lightweight rate-limiting & connection control at the proxy layer to block flooding or abusive traffic. |
| 🔍 **Transparent Logging** | Detailed request metrics for monitoring and analytics. |
| ⚙️ **Easy Integration** | Drop-in reverse proxy — no complex config, no dependency on nginx or envoy. |
| 🧩 **Extensible Design** | Clean Go codebase — easy to extend for authentication, load balancing, or CDN-like features. |

---

## 🧱 Architecture

```
┌─────────────────────────────┐
│         LiteShield          │
│ ┌─────────────────────────┐ │
│ │       HTTP Proxy        │ │
│ │  - Forward requests     │ │
│ │  - Cache responses      │ │
│ └─────────────────────────┘ │
│ ┌─────────────────────────┐ │
│ │      Cache Layer        │ │
│ │  - Memory-based TTL     │ │
│ │  - Per-client storage   │ │
│ └─────────────────────────┘ │
│ ┌─────────────────────────┐ │
│ │   Admin & Metrics API   │ │
│ │  - Clients overview     │ │
│ │  - Purge cache          │ │
│ │  - Prometheus metrics   │ │
│ └─────────────────────────┘ │
└─────────────────────────────┘
```

---

## 🧰 Tech Stack

- **Go 1.22+** — proxy, caching, admin API  
- **React + Chakra UI** — modern dashboard UI  
- **Prometheus (optional)** — metrics monitoring  
- **Docker / Swarm ready** — easy deployment

---

## ⚙️ Usage

### Run the proxy
```bash
go run cmd/proxy/main.go
```

### Environment variables
| Variable | Description | Example |
|-----------|-------------|----------|
| `BACKEND_URL` | Target backend to proxy | `http://localhost:8080` |
| `ADMIN_KEY` | API key for admin endpoints | `supersecret` |
| `CACHE_TTL` | Default cache lifetime | `60s` |

---

## 🧩 Admin API

| Endpoint | Method | Description |
|-----------|---------|-------------|
| `/admin/clients` | `GET` | List all registered clients |
| `/admin/cache/purge?client={id}` | `POST` | Purge cache for a specific client |
| `/metrics` | `GET` | Prometheus metrics endpoint |

All admin routes require `X-Admin-Key` header.

---

## 🌐 Dashboard (LiteShield UI)

A lightweight React dashboard lets you:
- View **cache statistics** in real-time  
- Manage **clients**  
- Monitor **logs & uptime**

Run the dashboard locally:
```bash
cd web-panel
npm install
npm run dev
```

---

## 🧩 Why LiteShield?

✅ **For Developers:** Instant setup, low code footprint, and Go-level performance.  
✅ **For Agencies:** Cache & proxy multiple client sites securely.  
✅ **For SaaS Teams:** Unified edge layer that scales without nginx/traefik headaches.  
✅ **For Ops:** Integrates with Docker Swarm, metrics ready, configurable TTL & keys.

---

## 🛡 Future Roadmap

- [ ] Redis distributed cache support  
- [ ] WebSocket real-time dashboard  
- [ ] IP-based rate limiting  
- [ ] Persistent caching (disk / S3)  
- [ ] Multi-node sync  

---

## 🧑‍💻 Author

**LiteShield** is an open project by **Tom Caillaud**,  
built for scalability, speed, and simplicity.
