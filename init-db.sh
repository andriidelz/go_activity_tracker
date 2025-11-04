#!/bin/bash
set -e

psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'activity_test'" | grep -q 1 || \
  psql -U postgres -c 'CREATE DATABASE activity_test;'

psql -U postgres activity_test -c "
  CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
  CREATE EXTENSION IF NOT EXISTS \"pg_trgm\";
"
