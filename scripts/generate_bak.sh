#!/bin/bash


export PGPASSWORD="postgres"

pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "integration" --exclude-table=integrations --exclude-table=integration_type_setups > integration.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "integration_types"  > integration_types.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "auth" > auth.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "compliance" --exclude-table=framework_compliance_summaries > compliance.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "dex" > dex.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "core" --exclude-table=platform_configurations > core.bak

