#!/bin/bash

echo "Setting up project..."
go mod tidy
docker-compose up -d
# Run migrations (assume using go-migrate or manual)
psql -U user -d github.com/dhanarrizky/Golang-template -f migrations/001_create_users_table.sql
psql -U user -d github.com/dhanarrizky/Golang-template -f migrations/002_add_index.sql