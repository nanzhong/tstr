api: API_PG_DSN=$DATABASE_URL go run -mod=vendor ./cmd/tstr api --bootstrap-token dev
runner: go run -mod=vendor ./cmd/tstr run --insecure --access-token dev
ui: go run -mod=vendor ./cmd/tstr ui --access-token dev
db: ./scripts/start_db.sh
