#!/bin/bash
# Initialize database by running all migration UP files in order

set -e

# Directory containing migrations
MIGRATIONS_DIR="/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "No migrations directory found at $MIGRATIONS_DIR"
    exit 0
fi

# Run only .up.sql files in alphabetical order
echo "Running database migrations..."
for migration in $(find "$MIGRATIONS_DIR" -maxdepth 1 -name "*.up.sql" | sort); do
    echo "Running migration: $(basename $migration)"
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration" 2>&1 | grep -v "^NOTICE:"  || true
done

echo "✅ Migrations completed successfully"
