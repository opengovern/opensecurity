#!/bin/bash
set -e

dt=$(date '+%d/%m/%Y %H:%M:%S');
echo "$dt - Running init script the 1st time Primary PostgreSql container is created...";

authDatabaseName="auth"
authUserName="auth_service"

informationDatabaseName="information"
informationUserName="information_service"

dexDatabaseName="dex"
dexUserName="dex_service"

exporterUserName="postgres_exporter"
steampipeUserName="steampipe_user"

migratorDatabaseName="migrator"
migratorUserName="migrator_worker"

describeDatabaseName="describe"
describeUserName="describe_scheduler"

assistantDatabaseName="assistant"
assistantUserName="assistant_service"

policyDatabaseName="policy"
policyUserName="policy_service"


complianceDatabaseName="compliance"
complianceUserName="compliance_service"

coreDatabaseName="core"
coreUserName="core_service"

reporterDatabaseName="reporter"
reporterUserName="reporter_service"

alertingDatabaseName="alerting"
alertingUserName="alerting_service"

integrationDatabaseName="integration"
integrationUserName="integration_service"

taskDatabaseName="task"
taskUserName="task_service"

integrationTypesDatabaseName="integration_types"

echo "$dt - Running: psql -v ON_ERROR_STOP=1 --username postgres --dbname postgres ...";

PGPASSWORD="postgres" psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL

SELECT 'CREATE DATABASE $informationDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$informationDatabaseName')\gexec
SELECT 'ALTER ROLE $informationUserName WITH PASSWORD ''$POSTGRES_INFORMATION_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$informationUserName')\gexec
SELECT 'CREATE USER $informationUserName WITH PASSWORD ''$POSTGRES_INFORMATION_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$informationUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$informationDatabaseName" to $informationUserName;

\c "$informationDatabaseName"
GRANT ALL ON SCHEMA public TO $informationUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

SELECT 'CREATE DATABASE $integrationTypesDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$integrationTypesDatabaseName')\gexec


SELECT 'CREATE DATABASE $dexDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$dexDatabaseName')\gexec
SELECT 'ALTER ROLE $dexUserName WITH PASSWORD ''$DEX_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$dexUserName')\gexec
SELECT 'CREATE USER $dexUserName WITH PASSWORD ''$DEX_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$dexUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$dexDatabaseName" to $dexUserName;

\c "$dexDatabaseName"
GRANT ALL ON SCHEMA public TO $dexUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

SELECT 'CREATE DATABASE $describeDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$describeDatabaseName')\gexec
SELECT 'ALTER ROLE $describeUserName WITH PASSWORD ''$POSTGRES_DESCRIBE_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$describeUserName')\gexec
SELECT 'CREATE USER $describeUserName WITH PASSWORD ''$POSTGRES_DESCRIBE_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$describeUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$describeDatabaseName" to $describeUserName;

\c $describeDatabaseName
GRANT ALL ON SCHEMA public TO $describeUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;

SELECT 'CREATE DATABASE $policyDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$policyDatabaseName')\gexec
SELECT 'ALTER ROLE $policyUserName WITH PASSWORD ''$POSTGRES_POLICY_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$policyUserName')\gexec
SELECT 'CREATE USER $policyUserName WITH PASSWORD ''$POSTGRES_POLICY_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$policyUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$policyDatabaseName" to $policyUserName;

\c $policyDatabaseName
GRANT ALL ON SCHEMA public TO $policyUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;


SELECT 'CREATE DATABASE $complianceDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$complianceDatabaseName')\gexec
SELECT 'ALTER ROLE $complianceUserName WITH PASSWORD ''$POSTGRES_COMPLIANCE_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$complianceUserName')\gexec
SELECT 'CREATE USER $complianceUserName WITH PASSWORD ''$POSTGRES_COMPLIANCE_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$complianceUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$complianceDatabaseName" to $complianceUserName;

\c $complianceDatabaseName
GRANT ALL ON SCHEMA public TO $complianceUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;

SELECT 'CREATE DATABASE $authDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$authDatabaseName')\gexec
SELECT 'ALTER ROLE $authUserName WITH PASSWORD ''$POSTGRES_AUTH_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$authUserName')\gexec
SELECT 'CREATE USER $authUserName WITH PASSWORD ''$POSTGRES_AUTH_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$authUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$authDatabaseName" to $authUserName;

\c "$authDatabaseName"
GRANT ALL ON SCHEMA public TO $authUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;

SELECT 'CREATE DATABASE $coreDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$coreDatabaseName')\gexec
SELECT 'ALTER ROLE $coreUserName WITH PASSWORD ''$POSTGRES_CORE_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$coreUserName')\gexec
SELECT 'CREATE USER $coreUserName WITH PASSWORD ''$POSTGRES_CORE_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$coreUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$coreDatabaseName" to $coreUserName;

\c "$coreDatabaseName"
GRANT ALL ON SCHEMA public TO $coreUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;

SELECT 'CREATE DATABASE $integrationDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$integrationDatabaseName')\gexec
SELECT 'ALTER ROLE $integrationUserName WITH PASSWORD ''$POSTGRES_INTEGRATION_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$integrationUserName')\gexec
SELECT 'CREATE USER $integrationUserName WITH PASSWORD ''$POSTGRES_INTEGRATION_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$integrationUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$integrationDatabaseName" to $integrationUserName;

\c $integrationDatabaseName
GRANT ALL ON SCHEMA public TO $integrationUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;

SELECT 'CREATE DATABASE $taskDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$taskDatabaseName')\gexec
SELECT 'ALTER ROLE $taskUserName WITH PASSWORD ''$POSTGRES_TASK_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$taskUserName')\gexec
SELECT 'CREATE USER $taskUserName WITH PASSWORD ''$POSTGRES_TASK_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$taskUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$taskDatabaseName" to $taskUserName;

\c $taskDatabaseName
GRANT ALL ON SCHEMA public TO $taskUserName;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;

SELECT 'CREATE DATABASE $migratorDatabaseName'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$migratorDatabaseName')\gexec
SELECT 'ALTER ROLE $migratorUserName WITH PASSWORD ''$POSTGRES_MIGRATOR_DB_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$migratorUserName')\gexec
SELECT 'CREATE USER $migratorUserName WITH PASSWORD ''$POSTGRES_MIGRATOR_DB_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$migratorUserName')\gexec
GRANT ALL PRIVILEGES ON DATABASE "$migratorDatabaseName" to $migratorUserName;

\c "$migratorDatabaseName"
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS citext;
GRANT ALL ON SCHEMA public TO $migratorUserName;
GRANT pg_read_all_data TO $coreUserName;
GRANT pg_write_all_data TO $coreUserName;
GRANT ALL ON SCHEMA public TO $coreUserName;
GRANT pg_read_all_data TO $complianceUserName;
GRANT pg_write_all_data TO $complianceUserName;
GRANT ALL ON SCHEMA public TO $complianceUserName;

SELECT 'ALTER ROLE $steampipeUserName WITH PASSWORD ''$POSTGRES_STEAMPIPE_USER_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$steampipeUserName')\gexec
SELECT 'CREATE USER $steampipeUserName WITH PASSWORD ''$POSTGRES_STEAMPIPE_USER_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$steampipeUserName')\gexec

\connect "$complianceDatabaseName";
GRANT pg_read_all_data TO $migratorUserName;
GRANT pg_write_all_data TO $migratorUserName;
GRANT ALL ON SCHEMA public TO $migratorUserName;
\connect "$coreDatabaseName";
GRANT pg_read_all_data TO $migratorUserName;
GRANT pg_write_all_data TO $migratorUserName;
GRANT ALL ON SCHEMA public TO $migratorUserName;
\connect "$authDatabaseName";
GRANT pg_read_all_data TO $migratorUserName;
GRANT pg_write_all_data TO $migratorUserName;
GRANT ALL ON SCHEMA public TO $migratorUserName; 
\connect "$onboardDatabaseName";
GRANT pg_read_all_data TO $migratorUserName;
GRANT pg_write_all_data TO $migratorUserName;
GRANT ALL ON SCHEMA public TO $migratorUserName;
GRANT pg_read_all_data TO $steampipeUserName;
\connect "$coreDatabaseName";
GRANT pg_read_all_data TO $migratorUserName;
GRANT pg_write_all_data TO $migratorUserName;
GRANT ALL ON SCHEMA public TO $migratorUserName;
GRANT pg_read_all_data TO $steampipeUserName;
\connect "$integrationDatabaseName";
GRANT pg_read_all_data TO $migratorUserName;
GRANT pg_write_all_data TO $migratorUserName;
GRANT ALL ON SCHEMA public TO $migratorUserName;
GRANT pg_read_all_data TO $steampipeUserName;
\connect "$integrationTypesDatabaseName";
GRANT pg_read_all_data TO $migrationUserName;
GRANT pg_write_all_data TO $migrationUserName;
GRANT pg_read_all_data TO $steampipeUserName;
GRANT pg_read_all_data TO $integrationUserName;

\connect "postgres";
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
SELECT 'ALTER ROLE $exporterUserName WITH PASSWORD ''$POSTGRES_EXPORTER_PASSWORD'''
WHERE EXISTS (select from pg_catalog.pg_roles where rolname = '$exporterUserName')\gexec
SELECT 'CREATE USER $exporterUserName WITH PASSWORD ''$POSTGRES_EXPORTER_PASSWORD'''
WHERE NOT EXISTS (select from pg_catalog.pg_roles where rolname = '$exporterUserName')\gexec
GRANT pg_monitor TO $exporterUserName;

ALTER USER $migratorUserName WITH SUPERUSER;

EOSQL

echo "$dt - Init script is completed";
export PGPASSWORD='postgres'
set "PGPASSWORD=postgres"
PGPASSWORD="postgres"
pg_restore -h localhost -p 5432 -U postgres -d $authDatabaseName -v "$authDatabaseName.bak";
pg_restore -h localhost -p 5432 -U postgres -d $integrationDatabaseName -v "$integrationDatabaseName.bak";
pg_restore -h localhost -p 5432 -U postgres -d $complianceDatabaseName -v "$complianceDatabaseName.bak";
pg_restore -h localhost -p 5432 -U postgres -d $dexDatabaseName -v "$dexDatabaseName.bak";
pg_restore -h localhost -p 5432 -U postgres -d $coreDatabaseName -v "$coreDatabaseName.bak";


echo "$dt - Restore is completed";  
