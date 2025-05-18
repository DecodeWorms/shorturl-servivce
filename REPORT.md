
📚 URL Shortener Project Report

🚑 Challenges Faced

1. CI Mongo readiness: MongoDB wasn't available immediately in GitHub Actions, resolved using nc wait loop

2. Short code collision: Ensured uniqueness by checking existing entries in DB with retries

3. Test flakiness: Fixed with mocked interfaces and consistent random seeding

🔍 Insights Gained

1. Leveraged multi-stage Docker builds for clean images

2. Introduced Gin middleware for request tracing

3. Effective use of mocks (GoMock) for service-layer testing


🚀 Future Improvements

1. User dashboard for click analytics

2. Expiry + one-time links

3. Real-time metrics for admin (Prometheus/Grafana)

🔎 How to Read This Project

main.go → Entry point

services/ → Business logic (URL, User)

mocks/ → Auto-generated GoMock files

models/ → All entity types

handlers/ → HTTP handlers

storage/ → Mongo/Redis implementations