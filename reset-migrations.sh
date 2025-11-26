#!/bin/sh
# Force reset migrations in Render database
# WARNING: This will drop all tables and re-run migrations from scratch
echo "⚠️  WARNING: This script will DROP ALL TABLES and re-run migrations!"
echo ""
echo "Press Ctrl+C now to cancel, or wait 5 seconds to continue..."
sleep 5

if [ -z "$POSTGRES_URL" ]; then
    echo "ERROR: POSTGRES_URL not set!"
    exit 1
fi

echo ""
echo "=== Dropping goose version table ==="
psql "$POSTGRES_URL" -c "DROP TABLE IF EXISTS goose_db_version CASCADE;"

echo ""
echo "=== Dropping all application tables ==="
psql "$POSTGRES_URL" -c "
DROP TABLE IF EXISTS timesheet_entries CASCADE;
DROP TABLE IF EXISTS timesheets CASCADE;
DROP TABLE IF EXISTS addresses CASCADE;
DROP TABLE IF EXISTS user_organizations CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;
DROP TABLE IF EXISTS users CASCADE;
"

echo ""
echo "=== Verifying tables are dropped ==="
psql "$POSTGRES_URL" -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';"

echo ""
echo "=== Running fresh migrations ==="
cd /app/migrations
goose postgres "$POSTGRES_URL" up

echo ""
echo "=== Verifying migrations were applied ==="
psql "$POSTGRES_URL" -c "SELECT version_id, is_applied, tstamp FROM goose_db_version ORDER BY id;"

echo ""
echo "=== Listing tables ==="
psql "$POSTGRES_URL" -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;"

echo ""
echo "✅ Migration reset complete!"
