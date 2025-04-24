#!/bin/bash


export PGPASSWORD="postgres"

pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "task" --exclu.de-table=task_runs --exclude-table=task_config_secrets > task.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "integration" --exclude-table=credentials --exclude-table=integrations --exclude-table=integration_type_setups > integration.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "integration_types" --exclude-table=integration_plugin_binaries --exclude-table=task_binaries > integration_types.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "auth" --exclude-table=api_keys > auth.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "compliance" --exclude-table=benchmark_assignments --exclude-table=framework_compliance_summaries > compliance.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "dex" > dex.bak
pg_dump --host "localhost" --port "5432" --username "postgres" --no-password --format=t --blobs --no-owner --no-privileges --no-comments --no-subscriptions --verbose "core" --exclude-table=dashboard_widgets --exclude-table=dashboards --exclude-table=run_named_query_run_caches --exclude-table=sessions --exclude-table=widgets --exclude-table=chat_clarifications --exclude-table=chat_suggestions --exclude-table=chatbot_secrets --exclude-table=chats --exclude-table=platform_configurations > core.bak

