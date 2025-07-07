#!/bin/bash

set -e

# Database migration script for OnCall Scheduler
# Usage: ./migrate.sh [environment]
# environment: docker (default) or prod

ENVIRONMENT=${1:-docker}

# Database connection parameters
if [ "$ENVIRONMENT" = "docker" ]; then
    DB_HOST=${DB_HOST:-db}
    DB_PORT=${DB_PORT:-5432}
    DB_NAME=${DB_NAME:-oncall}
    DB_USER=${DB_USER:-oncall}
    DB_PASSWORD=${DB_PASSWORD:-oncall123}
elif [ "$ENVIRONMENT" = "prod" ]; then
    if [ -z "$DATABASE_URL" ]; then
        echo "Error: DATABASE_URL environment variable is required for production"
        exit 1
    fi
    # Parse DATABASE_URL for individual components if needed
    # For now, we'll use DATABASE_URL directly with psql
else
    echo "Error: Invalid environment. Use 'docker' or 'prod'"
    exit 1
fi

echo "Starting database migration for environment: $ENVIRONMENT"

# Function to run SQL file
run_migration() {
    local file=$1
    echo "Running migration: $file"
    
    if [ "$ENVIRONMENT" = "docker" ]; then
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$file"
    else
        psql "$DATABASE_URL" -f "$file"
    fi
    
    if [ $? -eq 0 ]; then
        echo "✓ Migration completed: $file"
    else
        echo "✗ Migration failed: $file"
        exit 1
    fi
}

# Wait for database to be ready (for docker environment)
if [ "$ENVIRONMENT" = "docker" ]; then
    echo "Waiting for database to be ready..."
    for i in {1..30}; do
        if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
            echo "Database is ready!"
            break
        fi
        echo "Waiting for database... ($i/30)"
        sleep 2
    done
    
    if [ $i -eq 30 ]; then
        echo "Error: Database connection timeout"
        exit 1
    fi
fi

# Run migrations in order
MIGRATIONS_DIR="./migrations"
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Error: Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

# Check if there are any migration files
if [ -z "$(ls -A $MIGRATIONS_DIR/*.sql 2>/dev/null)" ]; then
    echo "No migration files found in $MIGRATIONS_DIR"
    exit 0
fi

# Run all migration files in order
for migration_file in $MIGRATIONS_DIR/*.sql; do
    if [ -f "$migration_file" ]; then
        run_migration "$migration_file"
    fi
done

echo "All migrations completed successfully!"

# Optional: Create a migrations log table to track applied migrations
if [ "$ENVIRONMENT" = "docker" ]; then
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
    CREATE TABLE IF NOT EXISTS migration_log (
        id SERIAL PRIMARY KEY,
        migration_name VARCHAR(255) NOT NULL,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    INSERT INTO migration_log (migration_name) 
    SELECT '001_init.sql' 
    WHERE NOT EXISTS (SELECT 1 FROM migration_log WHERE migration_name = '001_init.sql');
    " > /dev/null 2>&1
else
    psql "$DATABASE_URL" -c "
    CREATE TABLE IF NOT EXISTS migration_log (
        id SERIAL PRIMARY KEY,
        migration_name VARCHAR(255) NOT NULL,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    INSERT INTO migration_log (migration_name) 
    SELECT '001_init.sql' 
    WHERE NOT EXISTS (SELECT 1 FROM migration_log WHERE migration_name = '001_init.sql');
    " > /dev/null 2>&1
fi

echo "Migration log updated."