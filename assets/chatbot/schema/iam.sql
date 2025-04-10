---------------------------------------------------------
-- Create Schemas
---------------------------------------------------------
CREATE SCHEMA entraid;
CREATE SCHEMA azure;
CREATE SCHEMA aws;


---------------------------------------------------------
-- Entra ID / HR Tables (Schema: entraid)
---------------------------------------------------------

-- HR Table: Admin Consent
CREATE TABLE entraid.entraid_admin_consent_request_policy (
    id   text PRIMARY KEY,                           -- Placeholder primary key.
    data json                                        -- Placeholder data for 8 columns (original count: 8).
);

-- HR Table: App Registration
CREATE TABLE entraid.entraid_app_registration (
    display_name                         text,         -- The display name for the application.
    id                                   text PRIMARY KEY,  -- The unique identifier for the application.
    app_id                               text,         -- The unique identifier assigned to the application by Azure AD.
    created_date_time                    timestamptz,  -- Date/time the application was registered (UTC).
    description                          text,         -- Free text describing the application.
    is_authorization_service_enabled     boolean,      -- Indicates if authorization service is enabled.
    oauth2_require_post_response         boolean,      -- If POST requests are allowed for OAuth 2.0 token requests.
    publisher_domain                     text,         -- Verified publisher domain for the application.
    sign_in_audience                     text,         -- Supported Microsoft accounts for sign-in.
    api                                  json,         -- Web API settings.
    identifier_uris                      json,         -- URIs that identify the application.
    info                                 json,         -- Basic profile information (marketing, support, etc.).
    key_credentials                      json,         -- Collection of key credentials.
    owner_ids                            json,         -- IDs of the application owners.
    parental_control_settings            json,         -- Parental control settings.
    password_credentials                 json,         -- Collection of password credentials.
    spa                                  json,         -- Single-page application settings (e.g., redirect URIs).
    tags_src                             json,         -- Custom strings for categorization.
    web                                  json,         -- Web application settings.
    tags                                 json,         -- Tags (ColumnDescriptionTags).
    title                                text,         -- Title (ColumnDescriptionTitle).
    tenant_id                            text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier (ColumnDescriptionTenant).
);

-- HR Table: Application
CREATE TABLE entraid.entraid_application (
    display_name                   text,              -- Display name for the application.
    id                             text PRIMARY KEY,  -- Unique identifier for the application.
    app_id                         text,              -- Unique application identifier assigned by Azure AD.
    created_date_time              timestamptz,       -- Registration date/time (UTC).
    description                    text,              -- Description for end users.
    is_authorization_service_enabled boolean,           -- Indicates if authorization service is enabled.
    oauth2_require_post_response   boolean,           -- If OAuth 2.0 POST requests are allowed.
    publisher_domain               text,              -- Verified publisher domain.
    sign_in_audience               text,              -- Supported Microsoft accounts for sign-in.
    api                            json,              -- Web API configuration.
    identifier_uris                json,              -- URIs identifying the application.
    info                           json,              -- Basic profile information.
    key_credentials                json,              -- Application key credentials.
    owner_ids                      json,              -- IDs of the owners.
    parental_control_settings      json,              -- Parental control settings.
    password_credentials           json,              -- Password credentials.
    spa                            json,              -- Single-page application settings.
    tags_src                       json,              -- Custom categorization strings.
    web                            json,              -- Web application settings.
    tags                           json,              -- Tags (ColumnDescriptionTags).
    title                          text,              -- Title (ColumnDescriptionTitle).
    tenant_id                      text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier (correlated to entraid_tenant).
);

-- HR Table: Auth Policy
CREATE TABLE entraid.entraid_authorization_policy (
    id   text PRIMARY KEY,                           -- Placeholder primary key.
    data json                                        -- Placeholder data for 12 columns.
);

-- HR Table: Conditional Access Policy
CREATE TABLE entraid.entraid_conditional_access_policy (
    id   text PRIMARY KEY,                           -- Placeholder primary key.
    data json                                        -- Placeholder data for 22 columns.
);

-- HR Table: Device
CREATE TABLE entraid.entraid_device (
    id text PRIMARY KEY,                              -- Unique device identifier.
    display_name text,                                -- Display name for the device.
    account_enabled boolean,                          -- Indicates if the device account is enabled.
    device_id text,                                   -- Unique identifier from Azure Device Registration Service.
    approximate_last_sign_in_date_time timestamptz,   -- Last sign-in timestamp (UTC).
    filter text,                                      -- OData query string for device searches.
    is_compliant boolean,                             -- Device compliance status.
    is_managed boolean,                               -- Indicates if the device is managed.
    mdm_app_id text,                                  -- MDM application identifier.
    operating_system text,                            -- Device operating system.
    operating_system_version text,                    -- Operating system version.
    profile_type text,                                -- Classification of the device.
    trust_type text,                                  -- Type of trust (Workplace, AzureAd, or ServerAd).
    extension_attributes json,                        -- Extension attributes (JSON).
    member_of json,                                   -- JSON array of group/role memberships.
    title text,                                       -- Custom title or label.
    tenant_id text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);

-- HR Table: Audit Report
CREATE TABLE entraid.entraid_directory_audit_report (
    id   text PRIMARY KEY,                           -- Placeholder primary key.
    data json                                        -- Placeholder data for 14 columns.
);

-- HR Table: Directory Role
CREATE TABLE entraid.entraid_directory_role (
    id               text PRIMARY KEY,              -- Unique directory role identifier.
    description      text,                          -- Role description.
    display_name     text,                          -- Role display name.
    role_template_id text,                          -- Template ID on which the role is based.
    member_ids       json,                          -- JSON array of member IDs.
    title            text,                          -- Optional title or label.
    tenant_id        text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);

-- HR Table: Directory Setting
CREATE TABLE entraid.entraid_directory_setting (
    id           text PRIMARY KEY,                  -- Unique settings identifier.
    display_name text,                              -- Display name derived from template.
    template_id  text,                              -- Template identifier.
    name         text,                              -- Setting name.
    value        text,                              -- Setting value.
    title        text,                              -- Optional title or label.
    tenant_id    text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);

-- HR Table: Domain
CREATE TABLE entraid.entraid_domain (
    id                   text PRIMARY KEY,          -- Fully qualified domain name.
    authentication_type  text,                      -- Authentication type (Managed or Federated).
    is_default           boolean,                   -- Default domain for user creation.
    is_admin_managed     boolean,                   -- If DNS record management is delegated to Microsoft 365.
    is_initial           boolean,                   -- If this is the initial domain created.
    is_root              boolean,                   -- If the domain is a verified root domain.
    is_verified          boolean,                   -- Ownership verification status.
    supported_services   text,                      -- Capabilities assigned to the domain.
    title                text,                      -- Title (ColumnDescriptionTitle).
    tenant_id            text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);

-- HR Table: Entra ID Tenant
CREATE TABLE entraid.entraid_tenant (
    tenant_id text PRIMARY KEY,                     -- Unique tenant identifier.
    display_name text,                              -- Tenant display name.
    tenant_type text,                               -- Tenant type (e.g., commercial, government, educational).
    created_date_time timestamptz,                  -- Tenant creation timestamp (UTC).
    verified_domains json,                          -- JSON list of verified domains.
    on_premises_sync_enabled boolean,               -- Indicates if on-premises sync is enabled.
    metadata text,                                  -- Metadata associated with the tenant.
    platform_account_id text,                       -- Platform account ID.
    platform_resource_id text                       -- Unique resource ID in opengovernance.
);

-- HR Table: Sign-In Report
CREATE TABLE entraid.entraid_sign_in_report (
    id                                  text PRIMARY KEY,  -- Unique sign-in activity ID.
    created_date_time                   timestamptz,       -- Timestamp when sign-in was initiated (UTC).
    user_display_name                   text,              -- Display name of the user.
    user_principal_name                 text,              -- User principal name.
    user_id                             text REFERENCES entraid.entraid_user(id),  -- User ID (foreign key to entraid_user).
    app_id                              text,              -- Azure AD app ID (service principal).
    app_display_name                    text,              -- Application display name.
    ip_address                          text,              -- Client IP address.
    client_app_used                     text,              -- Legacy client used for sign-in.
    correlation_id                      text,              -- Request correlation ID.
    conditional_access_status           text,              -- Status of triggered conditional access policy.
    is_interactive                      boolean,           -- Indicates if sign-in was interactive.
    risk_detail                         text,              -- Details on risk state.
    risk_level_aggregated               text,              -- Aggregated risk level.
    risk_level_during_sign_in           text,              -- Risk level during sign-in.
    risk_state                          text,              -- Overall risk state.
    resource_display_name               text,              -- Display name of the resource.
    resource_id                         text,              -- Identifier for the resource.
    risk_event_types                    json,              -- JSON array of risk event types.
    status                              json,              -- JSON object with sign-in status details.
    device_detail                       json,              -- JSON object with device details.
    location                            json,              -- JSON object with sign-in location details.
    applied_conditional_access_policies json,              -- JSON array of triggered policies.
    title                               text,              -- Title (ColumnDescriptionTitle).
    tenant_id                           text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);

-- HR Table: Enterprise App
CREATE TABLE entraid.entraid_enterprise_application (
    id                          text PRIMARY KEY,  -- Unique identifier for the service principal.
    display_name                text,              -- Display name for the service principal.
    app_id                      text,              -- Associated application identifier.
    account_enabled             boolean,           -- Indicates if the service principal is enabled.
    app_display_name            text,              -- Display name exposed by the application.
    app_owner_organization_id   text,              -- Tenant ID where the application is registered.
    app_role_assignment_required boolean,         -- Indicates if an app role assignment is required.
    service_principal_type      text,              -- Type of service principal (application, managed identity, or legacy).
    sign_in_audience            text,              -- Supported Microsoft account types.
    app_description             text,              -- Description exposed by the application.
    description                 text,              -- Internal description.
    login_url                   text,              -- URL to redirect the user for sign-in.
    logout_url                  text,              -- URL used to log out a user.
    add_ins                     json,              -- Custom behavior settings.
    alternative_names           json,              -- Alternative names for identity details.
    app_roles                   json,              -- Roles exposed by the application.
    info                        json,              -- Basic profile information.
    key_credentials             json,              -- Key credentials.
    notification_email_addresses json,             -- Notification email addresses.
    owner_ids                   json,              -- IDs of the owners.
    password_credentials        json,              -- Password credentials.
    oauth2_permission_scopes    json,              -- OAuth2 permission scopes.
    reply_urls                  json,              -- Token redirect URLs.
    service_principal_names     json,              -- JSON array of identifier URIs.
    tags_src                    json,              -- Custom categorization strings.
    tags                        json,              -- Tags (ColumnDescriptionTags).
    title                       text,              -- Title (ColumnDescriptionTitle).
    tenant_id                   text REFERENCES entraid.entraid_tenant(tenant_id),  -- Tenant identifier.
    metadata                    text,              -- Azure resource metadata.
    platform_account_id         text,              -- Platform account ID.
    platform_resource_id        text               -- Unique resource ID in opengovernance.
);

-- HR Table: Service Principal
CREATE TABLE entraid.entraid_service_principal (
    id TEXT PRIMARY KEY,                                    -- Unique identifier for the service principal.
    display_name TEXT,                                       -- Human-readable name of the service principal.
    app_id TEXT,                                           -- Unique identifier for the associated application (Client ID).
    account_enabled BOOLEAN,                                -- Indicates whether the service principal's account is enabled.
    app_display_name TEXT,                                  -- Display name of the associated application.
    app_owner_organization_id TEXT,                         -- Tenant ID of the organization that owns the associated application.
    app_role_assignment_required BOOLEAN,                   -- Indicates whether role assignment is required for the service principal.
    service_principal_type TEXT,                            -- Type of service principal (e.g., Application, Managed Identity).
    sign_in_audience TEXT,                                  -- Supported account types for sign-in.
    app_description TEXT,                                   -- Description of the associated application.
    description TEXT,                                       -- Description of the service principal.
    login_url TEXT,                                         -- URL for signing in to the service principal.
    logout_url TEXT,                                        -- URL for signing out of the service principal.
    add_ins JSONB,                                           -- Custom behavior settings for the service principal.
    alternative_names JSONB,                                 -- Alternative names or aliases for the service principal.
    app_roles JSONB,                                         -- List of application roles exposed by the service principal.
                                                            -- Example: [{"id": "xxx", "allowedMemberTypes": ["User"], "displayName": "Read", "isEnabled": true, "description": "Allows reading resources", "value": "reader"}]
    info JSONB,                                              -- Basic information about the service principal.
    key_credentials JSONB,                                   -- List of key credentials associated with the service principal.
                                                            -- Example: [{
                                                            --     "CustomKeyIdentifier": "E9FEE4B06988140E4E4ADE2BBB9F00CBBAFEAB33",  -- Custom identifier for the key credential.
                                                            --     "DisplayName": "testing-pem",                                   -- Display name of the key credential.
                                                            --     "EndDateTime": "2034-11-04T18:28:55Z",                         -- Date and time when the key credential expires.
                                                            --     "Key": null,                                                   -- Key value (usually not populated in this context).
                                                            --     "KeyId": "a30d42d8-8b69-4544-8fb9-85573d8acd5c",              -- Unique identifier for the key credential.
                                                            --     "StartDateTime": "2024-11-06T18:28:55Z",                       -- Date and time when the key credential becomes valid.
                                                            --     "TypeEscaped": "AsymmetricX509Cert",                                 -- Type of key credential (e.g., Asymmetric X509 Certificate).
                                                            --     "Usage": "Verify"                                              -- Usage of the key credential (e.g., Verify, Sign).
                                                            -- }]
    notification_email_addresses JSONB,                       -- List of email addresses for receiving notifications related to the service principal.
    owner_ids JSONB,                                         -- List of identifiers of the owners of the service principal.
    password_credentials JSONB,                              -- List of password credentials associated with the service principal.
                                                            -- Example: [{
                                                            --     "CustomKeyIdentifier": null,                               -- Custom identifier for the password credential.
                                                            --     "EndDateTime": "2025-08-01T20:43:08.766Z",               -- Date and time when the password credential expires.
                                                            --     "Hint": "dQl",                                             -- Hint for the password credential.
                                                            --     "KeyId": "667b6187-d733-48fa-8d48-b8896be7b19b",          -- Unique identifier for the password credential.
                                                            --     "SecretText": null,                                       -- Password value (usually not populated in this context).
                                                            --     "StartDateTime": "2025-02-02T21:43:08.766Z"             -- Date and time when the password credential becomes valid.
                                                            -- }]
    oauth2_permission_scopes JSONB,                          -- OAuth2 permission scopes granted to the service principal.
    reply_urls JSONB,                                        -- Reply URLs for the service principal.
    service_principal_names JSONB,                           -- Service principal names (SPNs) associated with the service principal.
    tags_src JSONB,                                          -- Source tags associated with the service principal.
    tags JSONB,                                              -- Additional tags associated with the service principal.
    title TEXT,                                             -- Title or label for the service principal.
    tenant_id TEXT REFERENCES entraid.entraid_tenant(tenant_id), -- Identifier of the tenant that the service principal belongs to.
    metadata TEXT,                                           -- Additional metadata associated with the service principal.
    cloud_environment TEXT                                -- Cloud environment where the service principal is located.
);

-- HR Table: User
CREATE TABLE entraid.entraid_user (
    id text PRIMARY KEY,                        -- Unique user identifier.
    display_name text,                          -- Display name for the user.
    user_principal_name text,                   -- User principal name (email).
    account_enabled boolean,                    -- Indicates if the user account is enabled.
    created_date_time timestamp with time zone, -- User creation timestamp (UTC).
    last_sign_in_date_time timestamp with time zone, -- Last sign-in timestamp (UTC).
    user_type text,                             -- Classification of user type.
    mail text,                                  -- SMTP email address.
    job_title text,                             -- User job title.
    identities json,                            -- User identities (JSON).
    password_policies text,                     -- Password policies.
    sign_in_sessions_valid_from_date_time timestamp with time zone, -- Tokens valid from this timestamp (UTC).
    usage_location text,                        -- Two-letter country code (ISO 3166).
    im_addresses json,                          -- Instant messaging addresses (JSON).
    other_mails json,                           -- Additional email addresses (JSON).
    tenant_id text REFERENCES entraid.entraid_tenant(tenant_id),  -- Tenant identifier.
    metadata text,                              -- Metadata of the Azure resource.
    platform_account_id text,                   -- Platform account ID.
    platform_resource_id text                   -- Unique resource ID in opengovernance.
);

-- HR Table: Group
CREATE TABLE entraid.entraid_group (
    id                             text PRIMARY KEY,  -- Unique group identifier.
    display_name                   text,              -- Group display name.
    description                    text,              -- Optional group description.
    classification                 text,              -- Group classification (e.g., low, medium, high impact).
    created_date_time              timestamptz,       -- Creation timestamp (UTC).
    expiration_date_time           timestamptz,       -- Expiration timestamp (UTC).
    is_assignable_to_role          boolean,           -- If the group can be assigned to an Azure AD role.
    is_subscribed_by_mail          boolean,           -- If the group is subscribed to email conversations.
    mail                           text,              -- Group SMTP address.
    mail_enabled                   boolean,           -- Indicates if the group is mail-enabled.
    mail_nickname                  text,              -- Mail alias for the group.
    membership_rule                text,              -- Dynamic membership rule.
    membership_rule_processing_state text,           -- Processing state of the membership rule.
    on_premises_domain_name        text,              -- On-premises domain name (if synced).
    on_premises_last_sync_date_time timestamptz,      -- Last on-premises sync timestamp (UTC).
    on_premises_net_bios_name      text,              -- On-premises NetBIOS name.
    on_premises_sam_account_name   text,              -- On-premises SAM account name.
    on_premises_security_identifier text,             -- On-premises security identifier.
    on_premises_sync_enabled       boolean,           -- If on-premises sync is enabled.
    renewed_date_time              timestamptz,       -- Last renewal timestamp (UTC).
    security_enabled               boolean,           -- If the group is a security group.
    security_identifier            text,              -- Group security identifier.
    visibility                     text,              -- Visibility setting (e.g., Private, Public).
    assigned_labels                json,              -- Sensitivity labels (JSON).
    group_types                    json,              -- Types of groups (JSON).
    member_ids                     json,              -- JSON array of member IDs.
    owner_ids                      json,              -- JSON array of owner IDs.
    proxy_addresses                json,              -- JSON array of proxy addresses.
    resource_behavior_options      json,              -- JSON array of behavior options.
    resource_provisioning_options  json,              -- JSON array of provisioned resources.
    nested_groups                  json,              -- JSON array of nested group memberships.
    tags                           text,              -- Tags associated with the group.
    title                          text,              -- Additional title or label.
    tenant_id                      text REFERENCES entraid.entraid_tenant(tenant_id),  -- Tenant identifier.
    metadata                       text,              -- Metadata for the Azure resource.
    platform_account_id            text,              -- Platform account ID.
    platform_resource_id           text               -- Unique resource ID in opengovernance.
);

-- HR Table: Group Membership
CREATE TABLE entraid.entraid_group_membership (
    id                  text PRIMARY KEY,              -- Unique membership record identifier.
    display_name        text,                          -- Display name of the member.
    group_id            text REFERENCES entraid.entraid_group(id),  -- Group ID (foreign key).
    account_enabled     boolean,                       -- Member's account enabled status.
    user_principal_name text,                          -- Member's principal name.
    user_type           text,                          -- Type of user (e.g., Member, Guest).
    state               text,                          -- Membership state (e.g., Active, Inactive).
    security_identifier text,                          -- Security identifier (SID) for the user.
    proxy_addresses     text,                          -- Alternate emails or UPNs.
    mail                text,                          -- Primary email address.
    title               text,                          -- Title (ColumnDescriptionTitle).
    tenant_id           text REFERENCES entraid.entraid_tenant(tenant_id),  -- Tenant identifier.
    metadata            text,                          -- Optional membership metadata.
    platform_account_id text,                          -- Platform account ID.
    platform_resource_id text                           -- Unique resource ID in opengovernance.
);

-- HR Table: Identity Provider
CREATE TABLE entraid.entraid_identity_provider (
    id            text PRIMARY KEY,                  -- Unique identity provider identifier.
    name          text,                              -- Provider display name.
    type          text,                              -- Provider type (e.g., Google, Facebook).
    client_id     text,                              -- Client ID provided by the identity provider.
    client_secret text,                              -- Client secret (write-only).
    title         text,                              -- Title (ColumnDescriptionTitle).
    tenant_id     text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);

-- HR Table: Registration Details of User
CREATE TABLE entraid.entraid_user_registration_details (
    user_object_id text PRIMARY KEY REFERENCES entraid.entraid_user(id),  -- Unique user object ID (FK to entraid_user).
    registration_system_preferred_authentication_methods json,         -- Preferred authentication methods (JSON).
    registration_is_system_preferred_authentication_method_enabled boolean,  -- If system-preferred method is enabled.
    registration_has_admin_privileges_assigned boolean,                -- If admin privileges are assigned.
    registration_is_mfa_capable boolean,                               -- If the user is capable of MFA.
    registration_is_mfa_registered boolean,                            -- If the user has registered for MFA.
    registration_is_sspr_capable boolean,                              -- If the user is capable of SSPR.
    registration_is_sspr_registered boolean,                           -- If the user has registered for SSPR.
    registration_is_sspr_enabled boolean,                              -- If SSPR is enabled.
    registration_is_passwordless_capable boolean,                      -- If the user is capable of passwordless authentication.
    registration_last_updated_date_time timestamptz,                   -- Last update timestamp (UTC).
    registration_methods_registered json,                              -- Authentication methods registered (JSON).
    registration_user_preferred_method_for_secondary_authentication text,  -- Preferred secondary authentication method.
    tenant_id text REFERENCES entraid.entraid_tenant(tenant_id)         -- Tenant identifier.
);

-- HR Table: Managed Identity
CREATE TABLE entraid.entraid_managed_identity (
    id                             text PRIMARY KEY,          -- Unique managed identity identifier.
    display_name                   text,                      -- Display name (e.g., resource name).
    app_id                         text,                      -- Associated application's unique identifier.
    account_enabled                boolean,                   -- If the account is enabled.
    app_display_name               text,                      -- Display name exposed by the application.
    app_owner_organization_id      text,                      -- Tenant ID where the application is registered.
    app_role_assignment_required   boolean,                   -- If an app role assignment is required.
    identity_type                  text,                      -- Identity type: SystemAssigned or UserAssigned.
    sign_in_audience               text,                      -- Supported Microsoft account types.
    app_description                text,                      -- Description exposed by the application.
    description                    text,                      -- Free text description.
    login_url                      text,                      -- URL for authentication redirection.
    logout_url                     text,                      -- URL for logging out.
    add_ins                        json,                      -- Custom behavior definitions.
    alternative_names              json,                      -- Alternative names (JSON).
    app_roles                      json,                      -- Exposed roles (JSON).
    info                           json,                      -- Basic profile information (JSON).
    key_credentials                json,                      -- Collection of key credentials (JSON).
    notification_email_addresses   json,                      -- Notification email addresses (JSON).
    owner_ids                      json,                      -- Owner IDs (JSON).
    password_credentials           json,                      -- Password credentials (JSON).
    oauth2_permission_scopes       json,                      -- OAuth2 permission scopes (JSON).
    reply_urls                     json,                      -- Redirect URIs (JSON).
    service_principal_names        json,                      -- Identifier URIs (JSON).
    tags_src                       json,                      -- Custom categorization strings (JSON).
    tags                           json,                      -- Tags (ColumnDescriptionTags) (JSON).
    title                          text,                      -- Title (ColumnDescriptionTitle).
    tenant_id                      text REFERENCES entraid.entraid_tenant(tenant_id),  -- Tenant identifier.
    metadata                       text,                      -- Azure resource metadata.
    platform_account_id            text,                      -- Platform account ID.
    platform_resource_id           text                       -- Unique resource ID in opengovernance.
);

-- HR Table: Security Defaults Policy
CREATE TABLE entraid.entraid_security_defaults_policy (
    id          text PRIMARY KEY,                  -- Policy identifier.
    display_name text,                             -- Display name for the policy.
    is_enabled  boolean,                           -- If security defaults is enabled.
    description text,                             -- Policy description.
    title       text,                             -- Title (ColumnDescriptionTitle).
    tenant_id   text REFERENCES entraid.entraid_tenant(tenant_id)  -- Tenant identifier.
);



---------------------------------------------------------
-- Azure Tables (Schema: azure)
---------------------------------------------------------

-- Table: Azure Subscription
-- Table: Azure Subscription

CREATE TABLE azure.azure_subscription (

    qualified_subscription_id text,  -- Fully qualified subscription ID. Example: "/subscriptions/4a645ffb-85be-487f-8c21-8dd368b088fd"
    id text PRIMARY KEY,                        -- Subscription ID. Example: "4a645ffb-85be-487f-8c21-8dd368b088fd"
    display_name text,                          -- Friendly name for the subscription. Example: "Azure subscription 2"
    tenant_id text REFERENCES entraid.entraid_tenant(tenant_id), -- Tenant identifier (from entraid schema).
    state text,                                 -- Subscription state (e.g., 'StateEnabled', 'StateWarned', etc.). Example: "Enabled"
    authorization_source text,                  -- Authorization source (e.g., 'Legacy, RoleBased'). Example: "RoleBased"
    managed_by_tenants json,                    -- JSON array of managing tenants.
    subscription_policies json,                 -- JSON object representing subscription policies. Example: {"locationPlacementId": "Public_2014-09-01", "quotaId": "PayAsYouGo_2014-09-01", "spendingLimit": "Off"}
    tags json                                   -- Additional subscription tags (JSON).

);

-- Table: Azure Role Definition
CREATE TABLE azure.azure_role_definition (
    short_id TEXT,             -- The friendly ID/name that identifies the role definition. Example: "80dcbedb-47ef-405d-95bd-188a1b4ac406"
    id TEXT PRIMARY KEY,        -- Unique role definition identifier. Example: "/subscriptions/4a645ffb-85be-487f-8c21-8dd368b088fd/providers/Microsoft.Authorization/roleDefinitions/80dcbedb-47ef-405d-95bd-188a1b4ac406"
    resource_type TEXT,        -- Contains the resource type. Example: "Microsoft.Authorization/roleDefinitions"
    name TEXT,                -- Name of the Role. Example: "Elastic SAN Owner"
    role_type TEXT,            -- Type of the role definition. Example: "BuiltInRole"
    description TEXT,         -- Description of the role definition. Example: "Allows for full access to all resources under Azure Elastic SAN including changing network security policies to unblock data path access"
    assignable_scopes JSON,  -- JSON array of valid scopes. Example: ["/"] 
    permissions JSON,         -- JSON block with allowed/denied actions. Example: [{"actions": ["Microsoft.Authorization/*/read", "Microsoft.ResourceHealth/availabilityStatuses/read", "Microsoft.Resources/deployments/*", "Microsoft.Resources/subscriptions/resourceGroups/read", "Microsoft.ElasticSan/elasticSans/*", "Microsoft.ElasticSan/locations/*"], "dataActions":, "notActions":, "notDataActions":}]
    title TEXT                -- Optional display label (ColumnDescriptionTitle). Example: "elastic san owner"
);

-- Table: Azure Role Assignment
CREATE TABLE azure.azure_role_assignment (
    cloud_environment TEXT,            -- Azure Cloud Environment (e.g., AzurePublicCloud). Example: "AzurePublicCloud"
    name TEXT,                        -- Friendly name for the role assignment. Example: "0149e7b3-0e31-489e-ae3d-231d5e6776cb"
    id TEXT PRIMARY KEY,              -- Unique role assignment identifier. Example: "/subscriptions/4a645ffb-85be-487f-8c21-8dd368b088fd/providers/Microsoft.Authorization/roleAssignments/0149e7b3-0e31-489e-ae3d-231d5e6776cb"
    scope TEXT,                        -- Scope (management group, subscription, etc.). Can be:
                                        --  * Root Management Group: "/"
                                        --  * Management Group: "/providers/Microsoft.Management/managementGroups/{groupId}"
                                        --  * Subscription: "/subscriptions/{subscriptionId}"
                                        --  * Resource Group: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}" 
                                        --  * Resource: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
    resource_type TEXT,                        -- Usually "Microsoft.Authorization/roleAssignments". Example: "Microsoft.Authorization/roleAssignments"
    principal_id TEXT,                -- Assigned principal's object ID. Example: "62a06e30-7435-4797-9c9b-be41bdc5dd1e"
    principal_type TEXT CHECK (principal_type IN ('User', 'Group', 'ServicePrincipal')), -- Principal type with CHECK constraint 
    created_on TIMESTAMPTZ,            -- Creation timestamp (UTC). Example: "2024-12-25T01:36:45.75172Z"
    updated_on TIMESTAMPTZ,            -- Last update timestamp (UTC). Example: "2024-12-25T01:36:45.75172Z"
    role_definition_id TEXT REFERENCES azure.azure_role_definition(id), -- Foreign key to role definition, referencing the id column.
    subscription_id TEXT REFERENCES azure.azure_subscription(id),     -- Foreign key to subscription, referencing the id column.
    title TEXT                        -- Title (ColumnDescriptionTitle). Example: "0149e7b3-0e31-489e-ae3d-231d5e6776cb"
);


---------------------------------------------------------
-- AWS Tables (Schema: aws)
---------------------------------------------------------

-- Table: AWS Account
CREATE TABLE aws.aws_account (
    name text,                              -- Account name.
    arn text,                               -- Account ARN.
    account_id text,                        -- AWS account unique identifier.
    organization_id text,                   -- Organization identifier (if applicable).
    organization_arn text,                  -- Organization ARN.
    organization_feature_set text,          -- Feature set available (e.g., ALL or CONSOLIDATED_BILLING).
    organization_master_account_arn text,   -- Master account ARN.
    organization_master_account_email text, -- Master account email.
    organization_master_account_id text,    -- Master account ID.
    organization_available_policy_types json,  -- JSON structure for region opt-in status.
    account_email text,                     -- Account email.
    account_status text,                    -- Account status.
    title text,                             -- Custom title or label.
    account_aliases json                    -- JSON array of account aliases.
);

-- Table: AWS IAM Access Key
CREATE TABLE aws.aws_iam_access_key (
    access_key_id text PRIMARY KEY,          -- Unique access key identifier.
    user_name text,                          -- IAM user associated with the key.
    status text,                             -- Access key status (e.g., Active, Inactive).
    create_date timestamptz,                 -- Creation timestamp (UTC).
    access_key_last_used_date timestamptz,   -- Last used timestamp (UTC).
    access_key_last_used_service text,       -- AWS service last accessed.
    access_key_last_used_region text,        -- AWS region where the key was last used.
    title text                               -- Custom title for the access key.
);

-- Table: AWS IAM Account Password Policy
CREATE TABLE aws.aws_iam_account_password_policy (
    allow_users_to_change_password boolean,  -- If IAM users can change their password.
    expire_passwords boolean,                -- If passwords expire.
    hard_expiry boolean,                     -- If password changes are prevented after expiry.
    max_password_age integer,                -- Maximum password age (days).
    minimum_password_length integer,         -- Minimum password length.
    password_reuse_prevention integer,       -- Number of previous passwords disallowed.
    require_lowercase_characters boolean,    -- If lowercase characters are required.
    require_numbers boolean,                 -- If numbers are required.
    require_symbols boolean,                 -- If symbols are required.
    require_uppercase_characters boolean       -- If uppercase characters are required.
);

-- Table: AWS IAM Policy
CREATE TABLE aws.aws_iam_policy (
    name text,                              -- Friendly name for the policy.
    policy_id text,                         -- Unique policy identifier.
    path text,                              -- Policy path.
    arn text,                               -- Policy ARN.
    is_aws_managed boolean,                 -- True if the policy is AWS managed.
    is_attachable boolean,                  -- If the policy can be attached.
    create_date timestamptz,                -- Creation timestamp (UTC).
    update_date timestamptz,                -- Last update timestamp (UTC).
    attachment_count integer,               -- Number of entities the policy is attached to.
    is_attached boolean,                    -- True if attached to any entity.
    default_version_id text,                -- Identifier for the default policy version.
    permissions_boundary_usage_count integer, -- Count for permissions boundary usage.
    policy json,                            -- Policy document (JSON).
    policy_std json,                        -- Canonical JSON form for searching.
    tags_src json,                          -- Source tags (JSON).
    tags json,                              -- Additional policy tags (JSON).
    title text                              -- Custom title for the policy.
);

-- Table: AWS IAM Policy Attachment
CREATE TABLE aws.aws_iam_policy_attachment (
    policy_arn text,                        -- ARN of the attached policy.
    is_attached boolean,                    -- Attachment status.
    policy_groups json,                     -- JSON array of groups attached.
    policy_roles json,                      -- JSON array of roles attached.
    policy_users json                       -- JSON array of users attached.
);

-- Table: AWS IAM Virtual MFA Device
CREATE TABLE aws.aws_iam_virtual_mfa_device (
    serial_number text PRIMARY KEY,         -- Serial number for the virtual MFA device.
    enable_date timestamptz,                 -- Enablement timestamp (UTC).
    assignment_status text,                  -- Assignment status (e.g., Assigned/Unassigned).
    user_id text,                           -- IAM user ID associated with the device.
    user_name text,                         -- IAM user name.
    user json,                              -- IAM user details (JSON).
    tags_src json,                          -- Source tags (JSON).
    tags json,                              -- Additional tags (JSON).
    title text                              -- Custom title for the MFA device.
);

-- Table: AWS IAM Group
CREATE TABLE aws.aws_iam_group (
    name text,                              -- Friendly name for the IAM group.
    group_id text PRIMARY KEY,              -- Unique group identifier.
    path text,                              -- Group path.
    arn text,                               -- Group ARN.
    create_date timestamptz,                -- Creation timestamp (UTC).
    inline_policies json,                   -- JSON array of inline policy documents.
    inline_policies_std json,               -- Canonical inline policies (JSON).
    attached_policy_arns json,              -- JSON array of attached managed policy ARNs.
    users json,                             -- JSON array of IAM users in the group.
    title text                              -- Custom title for the group.
);

-- Table: AWS IAM Role
CREATE TABLE aws.aws_iam_role (
    name text,                              -- Friendly name for the IAM role.
    arn text,                               -- Role ARN.
    role_id text PRIMARY KEY,               -- Unique role identifier.
    create_date timestamptz,                -- Creation timestamp (UTC).
    description text,                       -- Role description.
    instance_profile_arns json,             -- JSON array of associated instance profiles.
    max_session_duration integer,           -- Maximum session duration (seconds).
    path text,                              -- Role path.
    permissions_boundary_arn text,          -- ARN for the permissions boundary.
    permissions_boundary_type text,         -- Type of permissions boundary.
    role_last_used_date timestamptz,          -- Last used timestamp (UTC).
    role_last_used_region text,              -- AWS region where last used.
    tags_src json,                          -- Source tags (JSON).
    inline_policies json,                   -- JSON array of inline policies.
    inline_policies_std json,               -- Canonical inline policies (JSON).
    attached_policy_arns json,              -- JSON array of attached managed policy ARNs.
    assume_role_policy json,                -- Assume role policy document (JSON).
    assume_role_policy_std json,            -- Canonical assume role policy (JSON).
    title text,                             -- Custom title for the role.
    tags json                               -- Additional role tags (JSON).
);

---------------------------------------------------------
-- Table-Level Comments
---------------------------------------------------------
COMMENT ON TABLE aws.aws_iam_access_key IS 'The aws_iam_access_key table stores details about AWS IAM access keys, including creation dates, last used information, and status. Primary Goal: To enable monitoring and auditing of API credentials and ensure the security and proper usage of access keys. Use Cases: Credential rotation, detecting inactive or compromised keys, and enforcing security compliance.';
COMMENT ON TABLE aws.aws_iam_account_password_policy IS 'The aws_iam_account_password_policy table stores the password policy settings for an AWS account, including rules for password complexity, expiration, and reuse. Primary Goal: To enforce and audit password security requirements across IAM users. Use Cases: Ensuring compliance with security policies, mitigating password-based attacks, and supporting internal audits.';
COMMENT ON TABLE aws.aws_iam_policy IS 'The aws_iam_policy table contains detailed information about IAM policies, including metadata, versioning, and attachment statistics. Primary Goal: To manage and audit the permissions and access control policies within an AWS environment. Use Cases: Policy compliance, security reviews, and troubleshooting permission issues.';
COMMENT ON TABLE aws.aws_iam_policy_attachment IS 'The aws_iam_policy_attachment table maps IAM policies to the AWS IAM entities (users, groups, roles) to which they are attached. Primary Goal: To provide visibility into policy enforcement and the distribution of permissions across the AWS environment. Use Cases: Auditing policy assignments, verifying least privilege configurations, and troubleshooting access anomalies.';
-- (Note: The aws_iam_user table was not included above; if present, its comment should be added accordingly.)
COMMENT ON TABLE aws.aws_iam_virtual_mfa_device IS 'The aws_iam_virtual_mfa_device table records details about virtual MFA devices, including enablement dates, assignment status, and associated user information. Primary Goal: To track and manage MFA devices for enhanced account security. Use Cases: Enforcing multi-factor authentication, monitoring MFA usage, and auditing security measures for user accounts.';
COMMENT ON TABLE aws.aws_iam_group IS 'The aws_iam_group table stores detailed information about AWS IAM groups, including group identifiers, membership, and attached policies. Primary Goal: To facilitate group-based access control and streamline permission management. Use Cases: Managing group memberships, auditing group policies, and implementing role-based access controls.';
COMMENT ON TABLE aws.aws_iam_role IS 'The aws_iam_role table captures detailed information about AWS IAM roles, including creation details, usage metrics, policy attachments, and permissions boundaries. Primary Goal: To manage roles that govern cross-account access, application permissions, and temporary security credentials. Use Cases: Role assumption tracking, access control enforcement, and auditing cross-service permissions within AWS. Note: All IAM entities in these tables are AWS‑exclusive, and their identifiers are non‑SSO; they cannot be correlated to external identity providers such as Okta or EntraID.';
COMMENT ON TABLE aws.aws_account IS 'This table stores comprehensive information about AWS accounts, including identifiers, ARNs, organizational context, and policy opt‑in statuses. It helps administrators audit and manage AWS accounts, monitor account statuses, and maintain proper governance. This table is also used to retrieve a list of AWS Accounts.';

COMMENT ON TABLE entraid.entraid_directory_role IS 'A directory role in Entra ID represents a set of permissions that can be assigned to identities. Members of this role (identified in the member_ids column) can be user identities (internal or external via B2B) or application identities (service principals, managed identities). Roles simplify permission management and enforce least‑privilege access.';
COMMENT ON TABLE entraid.entraid_directory_setting IS 'Stores configuration settings used by Entra ID. These settings can apply to various scenarios, including user identities, device identities, or application configurations.';
COMMENT ON TABLE entraid.entraid_application IS 'Stores Azure AD application details: IDs, configuration (account types, redirect URIs), and OAuth 2.0 permission scopes (e.g., "Access Slack"). Used for auditing application settings, managing access, and ensuring compliance.';
COMMENT ON TABLE entraid.entraid_user IS 'Stores information about user accounts in Entra ID, including internal users, external (B2B) guests, and synced users. Contains attributes like UPNs, display names, and contact details. Used to verify user information, audit events, troubleshoot logins, and generate reports.';
COMMENT ON TABLE entraid.entraid_user_registration_details IS 'This table stores secondary registration information for users within Microsoft Entra ID, focusing on authentication methods and capabilities like MFA, SSPR, and passwordless login. It is used for auditing and security reporting, and does not include primary user attributes like name or type.';
COMMENT ON TABLE entraid.entraid_group IS 'Stores information about Entra ID groups, including display names, descriptions, membership rules, on‑premises sync details, and security settings. Group types include M365, Security, Dynamic, and Cloud groups. Groups can contain internal/external (B2B) users and application identities (service principals). Used to verify group configurations, audit lifecycle events, and troubleshoot group assignments and security policies. Provides a complete view of group structures and membership for identity management and compliance audits when used with entraid_group_membership.';
COMMENT ON TABLE entraid.entraid_group_membership IS 'Records which identities belong to which groups, including user identities (internal or external B2B) and application identities (service principals), along with membership states and tenant IDs. Used to verify group membership, check user enablement, and audit membership across tenants for compliance.';
COMMENT ON TABLE entraid.entraid_service_principal IS 'Stores comprehensive information about service principals within an Azure Active Directory tenant. This includes details about the service principal''s identity (display name, app ID, type), associated application (display name, description, roles), authentication and authorization configurations (sign‑in audience, key credentials, password credentials, OAuth2 permissions), and other relevant properties (owners, URLs, tags). This table is essential for managing, auditing, and analyzing service principals, which are crucial for enabling applications, services, and automated processes to access resources securely within the Azure AD environment.';
COMMENT ON TABLE entraid.entraid_tenant IS 'Stores information about Entra ID tenants, each representing a dedicated and isolated instance of Azure Active Directory (Azure AD). Each tenant has its own set of users, groups, applications, and other directory objects. This table includes tenant identifier, display name, tenant type, creation time, verified domains, and on‑premises synchronization status. Used to audit tenant configurations, manage tenant‑wide settings, and understand the relationship between tenants and their associated Azure resources. Furthermore, this table is the definitive source for identifying all Azure AD/Entra ID directories.';
COMMENT ON TABLE entraid.entraid_device IS 'Stores information about devices registered in Entra ID (Azure Active Directory), including device identifiers, operating system details, registration status, compliance status, and other relevant properties. Used to inventory devices, monitor device health and compliance, troubleshoot issues, and manage device access to organizational resources.';

COMMENT ON TABLE azure.azure_role_definition IS 'Stores Role Definitions within Azure, including scope, permissions, and classification (BuiltIn or Custom). Use it to manage and audit the actions permitted (and not permitted) by each role and the scopes (management groups, subscriptions, resource groups, etc.) to which each role can be assigned.';
COMMENT ON TABLE azure.azure_role_assignment IS 'Details which principal (internal or external B2B user, group, or application identity such as a service principal) is assigned which role at which scope in Azure. It references azure_role_definition.id to link each assignment to its role definition. This table is crucial for auditing and ensuring least privilege, as roles assigned at higher-level scopes apply to child resources. Note: This table exclusively tracks Azure Role Assignments and is not intended for general identity inventory purposes.';
COMMENT ON TABLE azure.azure_subscription IS 'This table stores comprehensive information about Azure subscriptions, including identifiers, display names, tenant associations, state, and authorization details. It helps administrators monitor subscription status, manage policies, and ensure proper governance. This table is also used to retrieve a list of Azure subscriptions.';
