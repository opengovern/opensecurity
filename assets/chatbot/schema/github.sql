-- Table: github_repository_vulnerability_alert
CREATE TABLE github_repository_vulnerability_alert (
    number INTEGER NOT NULL,                          -- The unique identifier for the vulnerability alert within the repository.
    repository_id INTEGER NOT NULL REFERENCES github.github_repository(id) ON DELETE CASCADE, -- Repository ID (foreign key)
    state TEXT NOT NULL,                              -- Current state of the alert (e.g., "OPEN", "DISMISSED", "FIXED").
    dependency_scope TEXT,                            -- Scope of the dependency affected (e.g., "RUNTIME").
    dismiss_comment TEXT,                             -- Comment provided when dismissing the alert (if applicable).
    dismiss_reason TEXT,                              -- Reason for dismissing the alert (e.g., "INACCURATE", "WONT_FIX").
    dismissed_at TIMESTAMPTZ,                         -- Timestamp when the alert was dismissed (if applicable).
    dismisser JSON,                                   -- Information about the user who dismissed the alert (if applicable).
                                                      --   created_at: Timestamp when the user was created.
                                                      --   email: Email address of the user.
                                                      --   login: Login name of the user.
                                                      --   updated_at: Timestamp when the user was last updated.
                                                      --   url: URL of the user's profile.
    fixed_at TIMESTAMPTZ,                             -- Timestamp when the vulnerability was fixed (if applicable).
    security_advisory JSON,                            -- Detailed information about the security advisory associated with the alert.
                                                      --   classification: Classification of the advisory.
                                                      --   cvss: CVSS score and vector string.
                                                      --   description: Description of the advisory.
                                                      --   ghsa_id: GitHub Security Advisory ID.
                                                      --   id: Unique ID of the advisory.
                                                      --   identifiers: Array of identifiers.
                                                      --   node_id: Node ID of the advisory.
                                                      --   notifications_permalink: Permalink to notifications.
                                                      --   origin: Origin of the advisory.
                                                      --   permalink: Permalink to the advisory.
                                                      --   published_at: Timestamp when the advisory was published.
                                                      --   references: Array of references.
                                                      --   severity: Severity of the advisory.
                                                      --   summary: Summary of the advisory.
                                                      --   updated_at: Timestamp when the advisory was last updated.
                                                      --   withdrawn_at: Timestamp when the advisory was withdrawn.
    security_vulnerability JSON,                      -- Specific details about the security vulnerability.
                                                      --   advisory: Details of the security advisory.
                                                      --   first_patched_version: First patched version.
                                                      --   package: Details of the vulnerable package.
                                                      --   severity: Severity of the vulnerability.
                                                      --   updated_at: Timestamp when the vulnerability was last updated.
                                                      --   vulnerable_version_range: Vulnerable version range.
    vulnerable_manifest_filename TEXT,                -- Name of the manifest file containing the vulnerable dependency (e.g., "go.mod").
    vulnerable_manifest_path TEXT,                    -- Path to the vulnerable manifest file within the repository.
    vulnerable_requirements TEXT,                      -- Version requirements for the vulnerable dependency.
    severity TEXT NOT NULL,                           -- Severity of the vulnerability (e.g., "HIGH", "CRITICAL").
    cvss_score NUMERIC,                               -- CVSS score representing the severity of the vulnerability.
    PRIMARY KEY (number, repository_id)              -- Composite primary key.
);

COMMENT ON TABLE github_repository_vulnerability_alert IS 'General vulnerability alerts in repositories. This table contains a broader range of vulnerability alerts, including those identified by Dependabot (which are also found in `github_repository_dependabot_alert`) as well as other vulnerabilities. Use this table to see all potential vulnerabilities, investigate non-dependency alerts, and track your overall security posture. Regularly monitor for new alerts and prioritize fixes based on severity and risk.';
-- Table: github_repository_dependabot_alert
CREATE TABLE github_repository_dependabot_alert (
    alert_number INTEGER NOT NULL,                  -- The unique identifier for the Dependabot alert within the repository.
    repository_id INTEGER NOT NULL REFERENCES github.github_repository(id) ON DELETE CASCADE, -- Repository ID (foreign key referencing github.github_repository), ON DELETE CASCADE ensures related alerts are deleted.
    PRIMARY KEY (alert_number, repository_id),      -- Composite primary key: alert number and repository ID, ensures unique alerts per repository.
    state TEXT NOT NULL,                            -- The current state of the alert (e.g., "open", "dismissed", "fixed").
    dependency_package_ecosystem TEXT,              -- The package ecosystem for the vulnerable dependency (e.g., "npm", "go").
    dependency_package_name TEXT,                   -- The name of the vulnerable package.
    dependency_manifest_path TEXT,                 -- The path to the manifest file where the vulnerable dependency is declared.
    dependency_scope TEXT,                          -- The scope of the dependency (e.g., "runtime", "development").
    security_advisory_ghsa_id TEXT,                 -- The GitHub Security Advisory ID (GHSA) associated with the alert.
    security_advisory_cve_id TEXT,                  -- The Common Vulnerabilities and Exposures ID (CVE) associated with the alert.
    security_advisory_summary TEXT,                 -- A brief summary of the security advisory.
    security_advisory_description TEXT,             -- A detailed description of the security advisory.
    security_advisory_severity TEXT,                -- The severity of the vulnerability (e.g., "high", "medium", "low").
    security_advisory_cvss_score NUMERIC,           -- The CVSS score of the vulnerability.
    security_advisory_cvss_vector_string TEXT,      -- The CVSS vector string representing the vulnerability's characteristics.
    security_advisory_cwes JSON,                    -- An array of CWEs (Common Weakness Enumerations) associated with the vulnerability.
                                                    --   cwe_id: CWE ID
                                                    --   cwe_name: CWE name
    security_advisory_published_at TIMESTAMPTZ,     -- The timestamp when the security advisory was published.
    security_advisory_updated_at TIMESTAMPTZ,       -- The timestamp when the security advisory was last updated.
    created_at TIMESTAMPTZ NOT NULL,                -- The timestamp when the Dependabot alert was created.
    updated_at TIMESTAMPTZ NOT NULL,                -- The timestamp when the Dependabot alert was last updated.
    dismissed_at TIMESTAMPTZ,                       -- The timestamp when the Dependabot alert was dismissed (if applicable).
    dismissed_reason TEXT,                          -- The reason why the Dependabot alert was dismissed (if applicable).
    dismissed_comment TEXT,                         -- Any comments associated with dismissing the Dependabot alert.
    fixed_at TIMESTAMPTZ                            -- The timestamp when the vulnerability addressed by the alert was fixed.
);

COMMENT ON TABLE github_repository_dependabot_alert IS 'Dependabot alerts for vulnerable dependencies in repositories. This table is a subset of `github_repository_vulnerability_alert`, focusing specifically on vulnerabilities identified by Dependabot. Use this table to find vulnerable dependencies, track automated security updates, and prioritize fixes.';dependa---------------------------------------------------------
-- GitHub Tables (Schema: github)
---------------------------------------------------------


-- Table: github_user
-- Stores information about individual GitHub user accounts, including profile details, activity metrics, and settings.
CREATE TABLE github_user (
    login TEXT PRIMARY KEY,                             -- The login name of the user, used for identification and @-mentions.
    id INTEGER UNIQUE NOT NULL,                        -- The unique numeric ID assigned to the user by GitHub.
    node_id TEXT UNIQUE,                               -- The node ID of the user in GitHub's GraphQL API, used for API interactions.
    name TEXT,                                         -- The user's display name, which can be different from their login.
    email TEXT,                                        -- The user's publicly visible email address (if provided).
    company TEXT,                                      -- The company affiliated with the user, as displayed on their profile.
    location TEXT,                                     -- The user's declared location, which can be a city, region, or country.
    url TEXT,                                          -- The URL of the user's GitHub profile page.
    created_at TIMESTAMP WITH TIME ZONE,                -- The timestamp when the user account was created (UTC).
    updated_at TIMESTAMP WITH TIME ZONE,                -- The timestamp when the user's profile was last updated (UTC).
    has_sponsors_listing BOOLEAN,                       -- Indicates whether the user has a GitHub Sponsors listing enabled.
    interaction_ability JSON,                          -- Contains settings related to the user's interaction abilities (e.g., if they have limited interactions).
    sponsors_listing JSON,                             -- Contains details about the user's GitHub Sponsors listing (if enabled).
    status JSON,                                       -- Contains information about the user's current status (e.g., emoji, message).
);


-- Table: GitHub Organization
CREATE TABLE github_organization (
    login text UNIQUE NOT NULL,                         -- The login name used to identify the organization on GitHub. This is unique across all organizations.
    id integer UNIQUE NOT NULL PRIMARY KEY,             -- The unique integer ID assigned to the organization by GitHub. This is the primary key for the table.
    node_id text UNIQUE,                               -- A unique identifier used for accessing the organization through GitHub's GraphQL API.
    name text,                                          -- The display name of the organization, which may be different from the login.
    created_at timestamptz,                            -- The date and time (in UTC) when the organization was created on GitHub.
    updated_at timestamptz,                            -- The date and time (in UTC) when the organization's information was last updated.
    description text,                                   -- A text description of the organization, provided by the organization owners.
    email text,                                        -- The primary email address associated with the organization, used for contact purposes.
    url text,                                          -- The URL of the organization's main page on GitHub.
    interaction_ability json,                          -- A JSON object containing settings related to how users can interact with the organization.
    is_verified boolean,                               -- A boolean value indicating whether the organization has been verified by GitHub. Verified organizations have confirmed their profile information.
    location text,                                     -- The geographical location of the organization, if provided in their profile.
    saml_identity_provider json,                       -- A JSON object containing details about the SAML identity provider used by the organization for single sign-on.
                                                       /* Example:
                                                        {
                                                            "digest_method": "http://www.w3.org/2001/04/xmlenc#sha256",
                                                            "issuer": "https://sts.windows.net/e15efe5a-7f1c-4bd4-a8ee-9943d6ab620b/",
                                                            "signature_method": "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
                                                            "sso_url": "https://login.microsoftonline.com/e15efe5a-7f1c-4bd4-a8ee-9943d6ab620b/saml2"
                                                        }
                                                       */
    website_url text,                                  -- The URL of the organization's external website, if provided.
    billing_email text,                                -- The email address used for billing purposes.
    two_factor_requirement_enabled boolean,            -- A boolean value indicating whether two-factor authentication is required for all members of the organization.
    default_repo_permission text,                      -- The default permission level for newly created repositories within the organization.
    plan_filled_seats integer,                         -- The number of seats currently occupied by members in the organization's GitHub plan.
    plan_name text,                                    -- The name of the GitHub plan that the organization is subscribed to.
    plan_seats integer,                                -- The total number of seats available in the organization's GitHub plan.
    web_commit_signoff_required boolean               -- A boolean value indicating whether web commit signoff is required for the organization.
);


-- Table: GitHub Organization External Identity
CREATE TABLE github_organization_external_identity (
    guid TEXT NOT NULL,                        -- Unique GUID for the external identity.
    organization_id INTEGER NOT NULL REFERENCES github.github_organization(id) ON DELETE CASCADE, -- Organization ID (foreign key).
    external_id TEXT,                          -- The ID of the identity as it exists in the external identity provider (e.g., the user's ID in an external directory). 
    external_provider_id TEXT,                -- The directory/tenant ID of the external identity provider (e.g., the tenant ID of the external directory). This ID shows where the identity comes from.
    user_id INTEGER REFERENCES github.github_user(id) ON DELETE CASCADE, -- The ID of the associated GitHub user (foreign key).
    saml_identity JSON,                        -- SAML identity information (JSON).
                                                --     Contains attributes and properties exposed by the identity provider.
    scim_identity JSON,                        -- SCIM identity information (JSON).
                                                --     Contains attributes and properties exposed by the identity provider.
    organization_invitation JSON,              -- Organization invitation details (JSON).
                                                --     created_at: Invitation creation timestamp.
                                                --     email: Invited email address.
                                                --     invitation_type: Type of invitation.
                                                --     invitee: Details of the invited user.
                                                --     inviter: Details of the inviting user.
                                                --     organization: Details of the organization.
                                                --     role: Invited role.
    metadata JSON,                              -- Additional metadata (JSON).
    PRIMARY KEY (guid, organization_id)        -- Composite primary key using GUID and organization ID.
);

-- Github Team
CREATE TABLE github_team (
    organization_id INTEGER REFERENCES github.github_organization(id) ON DELETE CASCADE, -- Organization ID (foreign key).
    slug TEXT NOT NULL,                          -- Team slug (URL-friendly identifier).
    name TEXT NOT NULL,                          -- Team name.
    id INTEGER UNIQUE NOT NULL PRIMARY KEY,     -- Unique team ID assigned by GitHub (primary key).
    node_id TEXT UNIQUE,                        -- Node ID for GitHub GraphQL API.
    description TEXT,                            -- Team description.
    privacy TEXT,                                -- Team privacy setting (e.g., VISIBLE, SECRET).
    permission TEXT,                            -- The permission level of the team (e.g., 'pull', 'push', 'admin').
    url TEXT,                                  -- URL for the team page.
    html_url TEXT,                             -- The URL of the team's page on GitHub.
    members_count INTEGER,                      -- Total team members.
    repos_count INTEGER,                         -- Count of accessible repositories.
    parent_team_id INTEGER,                    -- The ID of the parent team, if any (null if no parent).
    team_sync JSON,                             -- A JSON object containing information about team synchronization with an identity provider (null if not externally managed or error).
    notification_setting TEXT,                  -- The notification setting for the team.
);

CREATE TABLE github_team_member (
    team_id INTEGER NOT NULL REFERENCES github.github_team(id) ON DELETE CASCADE,  -- Team ID (foreign key). This column stores the ID of the team that the member (user or organization) belongs to. It references the 'id' column in the 'github_team' table, ensuring that only valid team IDs can be used. The ON DELETE CASCADE clause ensures that if a team is deleted, the corresponding memberships in this table are also deleted.
    member_principal_type TEXT NOT NULL,            -- The type of principal (e.g., 'User', 'Organization'). This column indicates whether the member is a user or an organization.
    member_principal_id INTEGER NOT NULL,              -- User or Organization ID (foreign key). This column stores the ID of the member, which can be a user or an organization. It references either the 'id' column in the 'github_user' table or the 'id' column in the 'github_organization' table, depending on the value of 'member_principal_type'. The ON DELETE CASCADE clause ensures that if the corresponding user or organization is deleted, the membership in this table is also deleted.
    PRIMARY KEY (team_id, member_principal_id)      -- Composite primary key to ensure unique team/member combinations. This composite primary key enforces that a user or organization can only be a member of a specific team once, preventing duplicate entries.
);

-- This function is a trigger function that enforces data integrity by ensuring that every user added to a team is a valid GitHub user.
CREATE OR REPLACE FUNCTION github.check_team_member_user_id()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the member_principal_id exists in the github_user table.
    IF NEW.member_principal_type = 'User' THEN
        IF NOT EXISTS (SELECT 1 FROM github.github_user WHERE id = NEW.member_principal_id) THEN
            RAISE EXCEPTION 'Invalid user ID for team member.'; -- If the user ID is not found, raise an exception to prevent the operation and signal the invalid data.
        END IF;
    END IF;

    RETURN NEW; -- If the user ID is valid, allow the operation to proceed.
END;
$$ LANGUAGE plpgsql;

-- This trigger ensures that the `check_team_member_user_id` function is executed before an insert or update operation on the `github_team_member` table.
CREATE TRIGGER check_team_member_user_id_trigger
BEFORE INSERT OR UPDATE ON github.github_team_member -- Specify that this trigger should be activated before INSERT or UPDATE operations on the 'github_team_member' table.
FOR EACH ROW EXECUTE FUNCTION github.check_team_member_user_id(); -- Execute the 'check_team_member_user_id' function for each row that is being inserted or updated.


CREATE TABLE github_repository (
    id INTEGER NOT NULL PRIMARY KEY,  -- Unique repository ID from GitHub.
    organization_id INTEGER NOT NULL,  -- Organization ID (foreign key).
    name TEXT NOT NULL,  -- Repository name.
    repository_full_name TEXT UNIQUE NOT NULL,  -- Full repository name (owner/repo).
    node_id TEXT UNIQUE,  -- Node ID for GitHub GraphQL API.
    name_with_owner TEXT,  -- Repository name with owner.
    description TEXT,  -- Repository description.
    created_at TIMESTAMP WITH TIME ZONE,  -- Creation timestamp (UTC).
    updated_at TIMESTAMP WITH TIME ZONE,  -- Last update timestamp (UTC).
    pushed_at TIMESTAMP WITH TIME ZONE,  -- Last push timestamp (UTC).
    is_active BOOLEAN,  -- Active status.
    is_empty BOOLEAN,  -- If the repository is empty.
    is_fork BOOLEAN,  -- If the repository is a fork.
    is_security_policy_enabled BOOLEAN,  -- If a security policy is enabled.
    owner JSON,  -- Repository owner details (JSON).
    homepage_url TEXT,  -- Repository homepage URL.
    license_info JSON,  -- License information (JSON).
    topics JSON,  -- List of topics (JSON).
    visibility TEXT,  -- Repository visibility (public, private, etc.).
    default_branch_ref JSON,  -- Default branch details (JSON).
    permissions JSON,  -- Permissions details (JSON).
    parent JSON,  -- Parent repository details (JSON).
    source JSON,  -- Source repository details (JSON).
    primary_language TEXT,  -- Primary programming language.
    languages JSON,  -- Languages used (JSON).
    repo_settings JSON,  -- Repository settings (JSON), including:
                         --     allow_auto_merge: If auto-merge is allowed
                         --     allow_rebase_merge: If rebase merge is allowed
                         --     allow_update_branch: If update branch is allowed
                         --     archived: If the repository is archived
                         --     delete_branch_on_merge: If branches are deleted on merge
                         --     forking_allowed: If forking is allowed
                         --     has_discussions_enabled: If discussions are enabled
                         --     has_downloads: If downloads are enabled
                         --     has_issues_enabled: If issues are enabled
                         --     has_pages: If pages are enabled
                         --     has_projects_enabled: If projects are enabled
                         --     has_wiki_enabled: If wiki is enabled
                         --     is_in_organization: If the repository belongs to an organization
                         --     is_mirror: If the repository is a mirror
                         --     is_private: If the repository is private
                         --     is_template: If the repository is a template
                         --     locked: If the repository is locked
                         --     merge_commit_allowed: If merge commits are allowed
                         --     merge_commit_message: Merge commit message
                         --     merge_commit_title: Merge commit title
                         --     squash_merge_allowed: If squash merging is allowed
                         --     squash_merge_commit_message: Squash merge commit message
                         --     squash_merge_commit_title: Squash merge commit title
                         --     web_commit_signoff_required: If web commit signoff is required
    security_settings JSON,  -- Security settings (JSON), including:
                             --     dependabot_security_updates_enabled: If Dependabot security updates are enabled
                             --     private_vulnerability_reporting_enabled: If private vulnerability reporting is enabled
                             --     secret_scanning_enabled: If secret scanning is enabled
                             --     secret_scanning_push_protection_enabled: If secret scanning push protection is enabled
                             --     vulnerability_alerts_enabled: If vulnerability alerts are enabled
    repo_urls JSON,  -- Repository URLs (JSON).
                     --     clone_url: The URL to clone the repository over HTTPS.
                     --     git_url: The URL to clone the repository over Git.
                     --     html_url: The URL of the repository's website.
                     --     ssh_url: The URL to clone the repository over SSH.
                     --     svn_url: The URL to clone the repository over SVN.
    metrics JSON,  -- Metrics and statistics (JSON), including:
                    --     branches: The number of branches in the repository.
                    --     commits: The number of commits in the repository.
                    --     forks: The number of times the repository has been forked.
                    --     issues: The number of open issues in the repository.
                    --     open_issues: The number of open issues in the repository.
                    --     pull_requests: The number of open pull requests in the repository.
                    --     releases: The number of releases in the repository.
                    --     size: The size of the repository in kilobytes.
                    --     stargazers: The number of users who have starred the repository.
                    --     subscribers: The number of users who are subscribed to the repository.
                    --     tags: The number of tags in the repository.
    metadata JSON  -- Additional metadata (JSON), including:
                   --     contact_links: Contact information for the repository.
                   --     possible_commit_emails: An array of possible email addresses for the repository's commits.
                   --     pull_request_templates: An array of pull request templates for the repository.
                   --     archived_at: The timestamp when the repository was archived, if applicable.
                   --     code_of_conduct: The code of conduct for the repository.
                   --     issue_templates: An array of issue templates for the repository.
                   --     lock_reason: The reason why the repository is locked, if applicable.
                   --     security_policy_url: The URL of the repository's security policy.
);


-- Table: GitHub Repository Ruleset
CREATE TABLE github_repository_ruleset (
    repository_full_name TEXT NOT NULL REFERENCES github.github_repository(repository_full_name) ON DELETE CASCADE, -- Repository full name (foreign key).
    organization text REFERENCES github.github_organization(organization) ON DELETE CASCADE, -- Owning organization (foreign key).
    repository_name TEXT,                          -- Repository name (redundant, for clarity).
    name TEXT NOT NULL,                            -- Name of the ruleset.
    id TEXT NOT NULL,                              -- Unique ruleset identifier.
    created_at TIMESTAMP WITH TIME ZONE,           -- Creation timestamp (UTC) for the ruleset.
    database_id INTEGER,                           -- Database identifier for the ruleset.
    enforcement TEXT,                              -- Enforcement level (e.g., "enabled", "disabled").
    rules JSON,                                     -- JSON array of rules.
    bypass_actors JSON,                            -- JSON array of actors allowed to bypass the ruleset.
    conditions JSON,                               -- JSON object with conditions for the ruleset.
    metadata JSON,                                  -- Additional metadata (JSON).
    PRIMARY KEY (repository_full_name, id)          -- Composite primary key.
);



-- Table: Contains definitions of organization-specific (or repository-specific) roles within GitHub,
-- including associated permissions and metadata. Used by OpenSecurity
-- (formerly OpenComply) for security and compliance reporting.

CREATE TABLE github_organization_role_definition (
    id              INTEGER        NOT NULL,  -- Unique numeric ID of the role definition in OpenSecurity
    name            TEXT           NOT NULL,  -- Name of the role (e.g., "Security Admin", "Compliance Manager")
    description     TEXT,                   -- Brief description of the role, outlining its purpose/responsibilities
    permissions     JSON,                   -- JSON array of permissions granted to this role (e.g. "read_audit_logs")
    organization_id INTEGER        NOT NULL,  -- ID of the GitHub organization to which this role belongs
    created_at      TIMESTAMPTZ    NOT NULL,  -- Timestamp when this role definition was created
    updated_at      TIMESTAMPTZ    NOT NULL,  -- Timestamp when this role definition was last updated
    source          TEXT,                   -- Indicates where the definition originates (e.g., "Organization", manual, etc.)
    base_role       TEXT,                   -- (Optional) Parent/base role from which this definition inherits permissions
    type            TEXT                    -- Classifies the role: could be "organization-roles", "custom-repository-roles", or other categories
);







-- Function: github.is_github_configured_for_sso_with_entraid(integer)
CREATE OR REPLACE FUNCTION github.is_github_configured_for_sso_with_entraid(org_id INTEGER)
RETURNS BOOLEAN
AS $$
BEGIN
  -- Purpose: This function checks if a GitHub organization is configured for SSO with Entra ID (formerly Azure AD).
  -- It also handles cases where the `saml_identity_provider` field is NULL (no SSO configuration)
  -- and cases where there is no corresponding entry in the `entraid.entraid_tenant` table.
  RETURN EXISTS (
    SELECT 1
    FROM github.github_organization go
    LEFT JOIN entraid.entraid_tenant et ON go.saml_identity_provider ->> 'issuer' LIKE CONCAT('https://sts.windows.net/', et.tenant_id, '/%')
    WHERE go.id = org_id  -- Filter by the provided organization ID
      AND go.saml_identity_provider IS NOT NULL -- Ensure SAML details are present
      AND et.tenant_id IS NOT NULL             -- Ensure corresponding Entra ID tenant exists
  );
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION github.is_github_configured_for_sso_with_entraid(INTEGER) IS '
Checks if a GitHub organization (by ID) is configured for SSO with Entra ID.

Returns TRUE if the organization''s SAML issuer in `github_organization` matches an existing Entra ID tenant, otherwise FALSE (including cases with no SSO or missing Entra ID tenant).

Useful for determining SSO status, even when SSO is not configured or Entra ID details are missing.

Example: SELECT github.is_github_configured_for_sso_with_entraid(12345);
';

-- Table-Level Comment
-- Table-Level Comments for GitHub Tables (Schema: github)

COMMENT ON TABLE github_user IS 'Stores information about individual GitHub user accounts, including profile details, activity metrics, and settings.';
COMMENT ON TABLE github_organization IS 'Stores comprehensive information about GitHub organizations, including profile details, membership, activity, and configuration. Used for auditing, tracking, and analysis of organizations.';
COMMENT ON TABLE github_team IS 'Stores detailed information about teams within GitHub organizations, including profile, membership, activity, and hierarchical relationships. Used to audit settings, track membership, analyze activity, and manage access control.';
COMMENT ON TABLE github_team_member IS 'This table stores information about team memberships in GitHub. Each row represents a user belonging to a specific team. It helps track which users are members of which teams, facilitating analysis of team composition and access control.';

COMMENT ON FUNCTION github.check_team_member_user_id IS 'This function is a trigger function that enforces data integrity by ensuring that every user added to a team is a valid GitHub user. It does this by checking if the provided ID exists in the respective table (github_user or github_organization). If the ID is not found, it prevents the operation (insert or update) and raises an exception to signal the invalid data.';

COMMENT ON TRIGGER check_team_member_user_id_trigger ON github.github_team_member IS 'This trigger ensures that the `check_team_member_user_id` function is executed before an insert or update operation on the `github_team_member` table. This proactive check helps maintain the integrity of the team membership data by validating user IDs before any changes are applied to the table.';


COMMENT ON TABLE github_organization_external_identity IS 'Stores information about external identities connected to a GitHub organization, tracking the link between external identity providers (SAML or SCIM) and GitHub users. Crucial for managing external user access and ensuring proper synchronization.';
COMMENT ON TABLE github_organization_member IS 'Provides details about members within a GitHub organization, including profile information, membership details, activity indicators, and connections. Valuable for understanding team composition, access control, engagement, and contributions.';
COMMENT ON TABLE github_repository IS 'Stores comprehensive information about GitHub repositories, including details about the repository itself, activity, ownership, and related resources. Crucial for tracking activity, managing settings, understanding codebase characteristics, and analyzing usage.';
COMMENT ON TABLE github_repository_ruleset IS 'Stores information about rulesets configured for GitHub repositories. This includes details about the ruleset itself (name, enforcement level, conditions), the rules it contains, and actors who can bypass it. This table is crucial for understanding and managing the enforcement of rules within repositories, ensuring code quality, security, and compliance with organizational policies.';
COMMENT ON TABLE github_organization_role_definition IS 'Contains definitions of organization-specific (or repository-specific) roles within GitHub, including associated permissions and metadata. Used by OpenSecurity (formerly OpenComply) for security and compliance reporting.';
