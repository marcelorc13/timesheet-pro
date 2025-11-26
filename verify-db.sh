#!/bin/sh
# Database verification script for Render deployment
# Run this via Render Shell to verify migrations were applied

echo "=== Checking Database Connection ==="
if [ -z "$POSTGRES_URL" ]; then
    echo "ERROR: POSTGRES_URL not set!"
    exit 1
fi
echo "✓ POSTGRES_URL is set"

echo ""
echo "=== Checking Goose Version Table ==="
psql "$POSTGRES_URL" -c "SELECT version_id, is_applied, tstamp FROM goose_db_version ORDER BY id;"

echo ""
echo "=== Checking if tables exist ==="
psql "$POSTGRES_URL" -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;"

echo ""
echo "=== Checking users table structure ==="
psql "$POSTGRES_URL" -c "\d users"

echo ""
echo "=== Checking organizations table structure ==="
psql "$POSTGRES_URL" -c "\d organizations"

echo ""
echo "=== Checking user_organizations table structure ==="
psql "$POSTGRES_URL" -c "\d user_organizations"

echo ""
echo "=== Sample data counts ==="
psql "$POSTGRES_URL" -c "SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'organizations', COUNT(*) FROM organizations
UNION ALL  
SELECT 'user_organizations', COUNT(*) FROM user_organizations;"
