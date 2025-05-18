# 🧩 System Design: Go URL Shortener

## 1. 🧠 About the URL Shortener System

A URL shortener is a service that converts long URLs into shorter, fixed-length aliases that redirect to the original link.
This system is built in Go using Gin for routing, MongoDB for persistence, Redis for caching, and supports user-based link generation with optional analytics.

---

## 2. ✅ Functional Requirements

* Create short URLs for long URLs
* Redirect short URLs to original URLs
* Store URL mappings per user
* Cache short URL resolutions for fast redirection

---

## 3. 🛡️ Non-Functional Requirements

* High availability and low latency
* Scalability to handle millions of URLs and redirections
* Fault tolerance (cache fallback to DB)
* Consistent 7-character short code uniqueness
* CI for automated testing , build and ensure code compatibility 

---

## 4. 📈 System Capacity (Target: 1 Million URLs)

* **Shortened URLs stored**: 1,000,000
* **Daily clicks (read-heavy)**: Assuming 100:1 read/write → \~100M clicks/day
* **Peak QPS**: \~1,200 reads/sec, 12 writes/sec
* **URL lifetime**: Persistent (can be extended to support expiry)

---

## 5. 🗄️ Database Capacity Estimate

### urls collection:

* 1M records @ \~0.5KB each → \~500MB

### click\_events collection:

* 100M events/month @ \~0.5KB each → \~50GB/month

Redis cache (frequently accessed short URLs): \~10,000 keys → \~5MB

---

## 6. 📦 User Schema

```json
{
  "_id": "uuid",
  "user_name": "john_doe",
  "email": "john@example.com"
}
```

---

## 7. 🔗 URL Schema

```json
{
  "_id": "uuid",
  "user_id": "uuid",
  "short_url": "cxv8ui0",
  "long_url": "https://example.com",
  "created_at": "timestamp"
}
```

---

## 8. ⚖️ Trade-offs

| Decision                | Trade-off                                                        |
| ----------------------- | ---------------------------------------------------------------- |
| Random 7-char Base62    | Simpler but may risk collision (mitigated with retries)          |
| Redis read cache        | Improves speed, but adds complexity & memory use                 |
| MongoDB                 | Flexible schema, but eventually consistency vs. relational joins |


---

## 9. ⚙️ Scalability Additions

### Load Balancer Usage

To handle high traffic and ensure fault tolerance, the system can be fronted by a load balancer (e.g., NGINX, AWS ELB). The load balancer will:

* Distribute requests evenly across multiple Go backend instances
* Provide health checks and reroute around failed nodes
* Terminate TLS and forward HTTP traffic to services

### Database Sharding

To horizontally scale MongoDB beyond single-instance storage limits:

* Shard `urls` collection by short URL prefix or user ID hash
* Shard `click_events` by short URL hash or timestamp for write distribution
* Each shard is managed by a MongoDB shard cluster (config servers + routers)

This approach reduces hotspots, improves parallelism, and accommodates future growth to 100M+ URLs and billions of redirections.

---

## 10. 🚀 Improvements

* Add rate limiting and abuse protection
* Support link expiration or one-time use
* Add custom aliases (e.g., /my-brand)
* Export analytics dashboard for users
* Add backup & archiving strategies for old click logs
* Auto-prune stale links

