-------------------------------------------------------------------------------
-- 1) HubSpot Owners
-------------------------------------------------------------------------------
CREATE TABLE hubspot_owner (
    id               TEXT PRIMARY KEY,  -- e.g. "owner_001"
    portal_id        TEXT NOT NULL,
    first_name       TEXT,
    last_name        TEXT,
    email            TEXT,
    user_id          BIGINT,           -- numeric user ID if needed
    teams            JSONB,            -- store team information in JSON (or JSON if preferred)
    title            TEXT,             -- owner title/job role
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_owner IS E'
Stores owner (user) records from HubSpot, including name, email, and user/team info.
';

COMMENT ON COLUMN hubspot_owner.id               IS 'Primary key: unique owner ID (e.g., "owner_001").';
COMMENT ON COLUMN hubspot_owner.portal_id        IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_owner.teams            IS 'JSON array or object containing team info.';
COMMENT ON COLUMN hubspot_owner.title            IS 'Job title or role of the owner.';
COMMENT ON COLUMN hubspot_owner.archived         IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_owner.archived_at      IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_owner.created_at       IS 'Timestamp when this owner record was created in HubSpot.';
COMMENT ON COLUMN hubspot_owner.updated_at       IS 'Timestamp when this owner record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 2) HubSpot Companies
-------------------------------------------------------------------------------
CREATE TABLE hubspot_company (
    id               TEXT PRIMARY KEY,  -- e.g. "company_001"
    portal_id        TEXT NOT NULL,
    title            TEXT,             -- e.g. "Acme Corporation"
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_company IS E'
Stores company (business/organization) records from HubSpot.
';

COMMENT ON COLUMN hubspot_company.id          IS 'Primary key: unique company ID (e.g., "company_001").';
COMMENT ON COLUMN hubspot_company.portal_id   IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_company.title       IS 'Company name (e.g., "Acme Corporation").';
COMMENT ON COLUMN hubspot_company.archived    IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_company.archived_at IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_company.created_at  IS 'When this company record was created in HubSpot.';
COMMENT ON COLUMN hubspot_company.updated_at  IS 'When this company record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 3) HubSpot Contacts
-------------------------------------------------------------------------------
CREATE TABLE hubspot_contact (
    id               TEXT PRIMARY KEY,   -- e.g. "contact_001"
    portal_id        TEXT NOT NULL,
    first_name       TEXT,
    last_name        TEXT,
    email            TEXT,
    title            TEXT,               -- optional job title or contact "title"
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_contact IS E'
Stores individual person records (contacts) from HubSpot.
';

COMMENT ON COLUMN hubspot_contact.id          IS 'Primary key: unique contact ID (e.g. "contact_001").';
COMMENT ON COLUMN hubspot_contact.portal_id   IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_contact.email       IS 'Primary email address for the contact.';
COMMENT ON COLUMN hubspot_contact.title       IS 'Optional job title or other descriptive title.';
COMMENT ON COLUMN hubspot_contact.archived    IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_contact.archived_at IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_contact.created_at  IS 'When this contact record was created in HubSpot.';
COMMENT ON COLUMN hubspot_contact.updated_at  IS 'When this contact record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 4) HubSpot Campaigns
-------------------------------------------------------------------------------
CREATE TABLE hubspot_campaign (
    id               TEXT PRIMARY KEY,  -- e.g. "campaign_001"
    portal_id        TEXT NOT NULL,
    title            TEXT,
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_campaign IS E'
Stores campaign records from HubSpot (e.g., marketing campaigns).
';

COMMENT ON COLUMN hubspot_campaign.id          IS 'Primary key: unique campaign ID (e.g., "campaign_001").';
COMMENT ON COLUMN hubspot_campaign.portal_id   IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_campaign.title       IS 'Campaign name or label.';
COMMENT ON COLUMN hubspot_campaign.archived    IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_campaign.archived_at IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_campaign.created_at  IS 'When this campaign record was created in HubSpot.';
COMMENT ON COLUMN hubspot_campaign.updated_at  IS 'When this campaign record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 5) HubSpot HubDB
-------------------------------------------------------------------------------
CREATE TABLE hubspot_hub_db (
    id               TEXT PRIMARY KEY,   -- e.g. "hubdb_001"
    portal_id        TEXT NOT NULL,
    title            TEXT,               -- name of the HubDB table
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_hub_db IS E'
Stores metadata about HubDB tables from HubSpot.
';

COMMENT ON COLUMN hubspot_hub_db.id          IS 'Primary key: unique HubDB table ID (e.g., "hubdb_001").';
COMMENT ON COLUMN hubspot_hub_db.portal_id   IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_hub_db.title       IS 'Human-friendly name of the HubDB table.';
COMMENT ON COLUMN hubspot_hub_db.archived    IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_hub_db.archived_at IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_hub_db.created_at  IS 'When this HubDB record was created in HubSpot.';
COMMENT ON COLUMN hubspot_hub_db.updated_at  IS 'When this HubDB record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 6) HubSpot Communications
-------------------------------------------------------------------------------
CREATE TABLE hubspot_communications (
    id                             TEXT PRIMARY KEY,  -- e.g. "communication_001"
    hs_communication_channel_type  TEXT NOT NULL,     -- e.g. SMS, WHATS_APP, LINKEDIN_MESSAGE, FACEBOOK
    hs_communication_logged_from   TEXT NOT NULL,     -- e.g. "CRM"
    hs_communication_body          TEXT,              -- text/content of the message
    hs_timestamp                   TIMESTAMPTZ NOT NULL,
    hubspot_owner_id               TEXT REFERENCES hubspot_owner(id) ON DELETE SET NULL,
    archived                       BOOLEAN DEFAULT FALSE,
    archived_at                    TIMESTAMPTZ,
    created_at                     TIMESTAMPTZ NOT NULL,
    updated_at                     TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_communications IS E'
Stores records of external messages (SMS, WhatsApp, LinkedIn, Facebook, etc.) known as "Communications" in HubSpot.
';

COMMENT ON COLUMN hubspot_communications.id                            IS 'Primary key: unique communications ID (e.g., "communication_001").';
COMMENT ON COLUMN hubspot_communications.hs_communication_channel_type IS 'Type of message channel (SMS, WHATS_APP, etc.).';
COMMENT ON COLUMN hubspot_communications.hs_communication_logged_from  IS 'Always set to "CRM" for these logs.';
COMMENT ON COLUMN hubspot_communications.hs_communication_body         IS 'The body/content of the message.';
COMMENT ON COLUMN hubspot_communications.hs_timestamp                  IS 'Timestamp of when the message was sent/received.';
COMMENT ON COLUMN hubspot_communications.hubspot_owner_id              IS 'References the HubSpot owner record. ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_communications.archived                      IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_communications.archived_at                   IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_communications.created_at                    IS 'When this communication record was created in HubSpot.';
COMMENT ON COLUMN hubspot_communications.updated_at                    IS 'When this communication record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 7) HubSpot Deal
-------------------------------------------------------------------------------

-- HUBSPOT DEALS --
CREATE TABLE IF NOT EXISTS hubspot_deals (
                                             portal_id TEXT,
                                             deal_id TEXT PRIMARY KEY,
                                             deal_name TEXT,
                                             deal_stage TEXT,
                                             pipeline TEXT,
                                             close_date DATETIME,
                                             amount NUMERIC,
                                             deal_type TEXT,
                                             description TEXT,
                                             owner_id TEXT,
                                             created_at DATETIME,
                                             updated_at DATETIME,
                                             archived BOOLEAN,
                                             closed_amount NUMERIC,
                                             closed_amount_home_currency NUMERIC,
                                             deal_stage_probability NUMERIC,
                                             forecast_amount NUMERIC,
                                             forecast_amount_home_currency NUMERIC,
                                             last_sales_activity_date DATETIME,
                                             object_id TEXT,
                                             priority TEXT,
                                             time_before_close NUMERIC,
                                             notes_last_activity_date DATETIME,
                                             notes_last_contacted DATETIME,
                                             associations jsonb,
                                             dealstages jsonb
);
COMMENT ON TABLE hubspot_deals IS E'
Represents a HubSpot deal, typically a sales opportunity or transaction.
';

-- Table to store pipeline information
CREATE TABLE IF NOT EXISTS hubspot_pipelines (
                                                 portal_id TEXT NOT NULL,       -- The HubSpot portal ID associated with the pipeline. Cannot be null.
                                                 pipeline_id TEXT NOT NULL,     -- Unique identifier for the pipeline within the HubSpot portal. Cannot be null.
                                                 label TEXT,                    -- The descriptive label or name of the pipeline. Can be null.
                                                 display_order INTEGER,         -- The order in which the pipeline appears in the HubSpot interface. Can be null.
                                                 archived BOOLEAN,              -- Indicates whether the pipeline has been archived (true) or not (false). Can be null.
                                                 PRIMARY KEY (portal_id, pipeline_id) -- Composite primary key: ensures uniqueness of pipelines within a portal.
    );

-- Table to store stages within each pipeline
CREATE TABLE IF NOT EXISTS hubspot_pipeline_stages (
                                                       portal_id TEXT NOT NULL,       -- The HubSpot portal ID associated with the pipeline stage. Cannot be null.
                                                       stage_id TEXT NOT NULL,         -- Unique identifier for the pipeline stage within the HubSpot portal. Cannot be null.
                                                       pipeline_id TEXT NOT NULL,     -- The identifier of the pipeline to which this stage belongs. Cannot be null.
                                                       label TEXT,                    -- The descriptive label or name of the pipeline stage. Can be null.
                                                       display_order INTEGER NOT NULL, -- The order in which the stage appears within the pipeline. Cannot be null.
                                                       PRIMARY KEY (portal_id, stage_id), -- Composite primary key: ensures uniqueness of pipeline stages within a portal.
    FOREIGN KEY (portal_id, pipeline_id) REFERENCES hubspot_pipelines(portal_id, pipeline_id) -- Establishes a relationship with the hubspot_pipelines table, ensuring referential integrity and that each stage belongs to a valid pipeline within the correct portal.
    );

-- Insert the default pipeline
INSERT INTO hubspot_pipelines (pipeline_id, name, label) VALUES
('default', 'defaultpipeline', 'Default Sales Pipeline');

-- Insert the stages for the default pipeline
INSERT INTO hubspot_pipeline_stages (stage_id, pipeline_id, label, value, stage_order) VALUES
('appointmentscheduled', 'default', 'Appointment Scheduled', 'appointmentscheduled', 1),
('qualifiedtobuy', 'default', 'Qualified to Buy', 'qualifiedtobuy', 2),
('presentationscheduled', 'default', 'Presentation Scheduled', 'presentationscheduled', 3),
('decisionmakerboughtin', 'default', 'Decision Maker Bought-In', 'decisionmakerboughtin', 4),
('contractsent', 'default', 'Contract Sent', 'contractsent', 5),
('closedwon', 'default', 'Closed Won', 'closedwon', 6),
('closedlost', 'default', 'Closed Lost', 'closedlost', 7);
-------------------------------------------------------------------------------
-- 8) HubSpot Leads
-------------------------------------------------------------------------------
CREATE TABLE hubspot_leads (
    id               TEXT PRIMARY KEY,   -- e.g. "lead_001"
    contact_id       TEXT REFERENCES hubspot_contact(id) ON DELETE SET NULL,
    owner_id         TEXT REFERENCES hubspot_owner(id)   ON DELETE SET NULL,
    portal_id        TEXT NOT NULL,
    title            TEXT,
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_leads IS E'
Represents leads in HubSpot (potential customers/contacts with interest in products/services).
';

COMMENT ON COLUMN hubspot_leads.id          IS 'Primary key: unique lead ID (e.g., "lead_001").';
COMMENT ON COLUMN hubspot_leads.contact_id  IS 'Optional link to a contact, ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_leads.owner_id    IS 'Owner ID from the hubspot_owner table, ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_leads.portal_id   IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_leads.title       IS 'Short title or description of the lead.';
COMMENT ON COLUMN hubspot_leads.archived    IS 'TRUE if archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_leads.archived_at IS 'When it was archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_leads.created_at  IS 'When this lead was created in HubSpot.';
COMMENT ON COLUMN hubspot_leads.updated_at  IS 'When this lead was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 9) HubSpot Marketing Email
-------------------------------------------------------------------------------
CREATE TABLE hubspot_marketing_email (
    id                  TEXT NOT NULL,
    contact_id          TEXT REFERENCES hubspot_contact(id) ON DELETE SET NULL,
    company_id          TEXT REFERENCES hubspot_company(id) ON DELETE SET NULL,
    portal_id           TEXT,
    title               TEXT,
    subject             TEXT,
    feedback_survey_id  TEXT,
    publish_date        TIMESTAMPTZ,
    archived            BOOLEAN,
    archived_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ,
    CONSTRAINT pk_hubspot_marketing_email PRIMARY KEY (id)
);

COMMENT ON TABLE hubspot_marketing_email IS E'
Stores information about Marketing Emails from HubSpot (e.g. "GET /marketing/v3/emails/").
Columns:
• contact_id, company_id: optional references with ON DELETE SET NULL.
• publish_date: when the email was sent or scheduled.
• archived, archived_at: for soft-deletions in HubSpot.
• created_at, updated_at: creation/modification timestamps in HubSpot.
';

COMMENT ON COLUMN hubspot_marketing_email.id                 IS 'Primary key: unique ID assigned to the marketing email by HubSpot.';
COMMENT ON COLUMN hubspot_marketing_email.contact_id         IS 'Optional reference to a single contact (ON DELETE SET NULL).';
COMMENT ON COLUMN hubspot_marketing_email.company_id         IS 'Optional reference to a single company (ON DELETE SET NULL).';
COMMENT ON COLUMN hubspot_marketing_email.portal_id          IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_marketing_email.title              IS 'Internal label or title of the email.';
COMMENT ON COLUMN hubspot_marketing_email.subject            IS 'Subject line for the email.';
COMMENT ON COLUMN hubspot_marketing_email.feedback_survey_id IS 'If linked to a feedback survey, the ID.';
COMMENT ON COLUMN hubspot_marketing_email.publish_date       IS 'When the email was published/sent.';
COMMENT ON COLUMN hubspot_marketing_email.archived           IS 'TRUE if the email is archived/deleted, FALSE if active.';
COMMENT ON COLUMN hubspot_marketing_email.archived_at        IS 'Timestamp of archival if archived=TRUE.';
COMMENT ON COLUMN hubspot_marketing_email.created_at         IS 'When this email record was created in HubSpot.';
COMMENT ON COLUMN hubspot_marketing_email.updated_at         IS 'When this email record was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 10) HubSpot Ticket
-------------------------------------------------------------------------------
CREATE TABLE hubspot_ticket (
    id               TEXT PRIMARY KEY,   -- e.g. "ticket_001"
    company_id       TEXT REFERENCES hubspot_company(id) ON DELETE SET NULL,
    contact_id       TEXT REFERENCES hubspot_contact(id) ON DELETE SET NULL,
    owner_id         TEXT REFERENCES hubspot_owner(id)   ON DELETE SET NULL,
    portal_id        TEXT,
    title            TEXT,
    archived         BOOLEAN DEFAULT FALSE,
    archived_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_ticket IS E'
Represents a support/help request ticket in HubSpot.
';

COMMENT ON COLUMN hubspot_ticket.id          IS 'Primary key: unique ticket ID (e.g., "ticket_001").';
COMMENT ON COLUMN hubspot_ticket.company_id  IS 'Optional reference to company, ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_ticket.contact_id  IS 'Optional reference to contact, ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_ticket.owner_id    IS 'Optional owner reference, ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_ticket.portal_id   IS 'HubSpot portal/account ID.';
COMMENT ON COLUMN hubspot_ticket.title       IS 'Short title of the support ticket.';
COMMENT ON COLUMN hubspot_ticket.archived    IS 'TRUE if ticket is archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_ticket.archived_at IS 'When archived if archived=TRUE.';
COMMENT ON COLUMN hubspot_ticket.created_at  IS 'When this ticket was created in HubSpot.';
COMMENT ON COLUMN hubspot_ticket.updated_at  IS 'When this ticket was last updated in HubSpot.';



-------------------------------------------------------------------------------
-- 11) Association Types (Reference Table)
-------------------------------------------------------------------------------
CREATE TABLE hubspot_association_types (
    type_id    INT PRIMARY KEY,     -- e.g. 1, 5, 211
    category   TEXT NOT NULL,       -- "HUBSPOT_DEFINED" or "USER_DEFINED"
    label      TEXT,                -- optional label (e.g. "Primary", "Billing contact")
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE hubspot_association_types IS E'
Stores metadata for association types used in HubSpot. 
Each row represents a distinct associationTypeId with a category (HUBSPOT_DEFINED or USER_DEFINED) 
and an optional label (e.g. "Decision maker").
';

COMMENT ON COLUMN hubspot_association_types.type_id    IS 'Numeric code for the association type (e.g. 1, 5, 30).';
COMMENT ON COLUMN hubspot_association_types.category   IS 'HUBSPOT_DEFINED or USER_DEFINED.';
COMMENT ON COLUMN hubspot_association_types.label      IS 'Optional descriptive label for the association type.';
COMMENT ON COLUMN hubspot_association_types.created_at IS 'When this type record was inserted here.';
COMMENT ON COLUMN hubspot_association_types.updated_at IS 'When this type record was last updated.';



-------------------------------------------------------------------------------
-- 12) Associations (Bridging Table)
-------------------------------------------------------------------------------
CREATE TABLE hubspot_associations (
    association_id         BIGSERIAL PRIMARY KEY,
    from_object_type_id    TEXT NOT NULL,  -- e.g. "0-1" for contacts, "0-3" for deals
    from_object_id         TEXT NOT NULL,  -- record ID in HubSpot (e.g. "23125848331")
    to_object_type_id      TEXT NOT NULL,
    to_object_id           TEXT NOT NULL,
    association_category   TEXT NOT NULL,  -- "HUBSPOT_DEFINED" or "USER_DEFINED"
    association_type_id    INT REFERENCES hubspot_association_types(type_id) ON DELETE SET NULL,
    created_at             TIMESTAMPTZ NOT NULL,
    updated_at             TIMESTAMPTZ NOT NULL,
    archived               BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE hubspot_associations IS E'
Links two HubSpot records (e.g. contact to deal, communication to contact) with a specified association type.
';

COMMENT ON COLUMN hubspot_associations.association_id       IS 'Primary key for the association row.';
COMMENT ON COLUMN hubspot_associations.from_object_type_id  IS 'Object type ID of the "source" record (e.g. "0-1" for contacts).';
COMMENT ON COLUMN hubspot_associations.from_object_id       IS 'The actual record ID in HubSpot for the "source".';
COMMENT ON COLUMN hubspot_associations.to_object_type_id    IS 'Object type ID of the "target" record (e.g. "0-2" for companies).';
COMMENT ON COLUMN hubspot_associations.to_object_id         IS 'The actual record ID in HubSpot for the "target".';
COMMENT ON COLUMN hubspot_associations.association_category IS 'Category of the association (HUBSPOT_DEFINED or USER_DEFINED).';
COMMENT ON COLUMN hubspot_associations.association_type_id  IS 'References hubspot_association_types.type_id; ON DELETE SET NULL.';
COMMENT ON COLUMN hubspot_associations.archived             IS 'TRUE if the association is archived/deleted in HubSpot.';
COMMENT ON COLUMN hubspot_associations.created_at           IS 'Timestamp when this association row was created locally.';
COMMENT ON COLUMN hubspot_associations.updated_at           IS 'Timestamp when this association row was last updated locally.';

