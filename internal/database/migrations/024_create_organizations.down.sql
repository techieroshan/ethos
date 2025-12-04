-- Down migration for organizations
-- +migrate Down

DROP TRIGGER IF EXISTS update_organization_members_updated_at ON organization_members;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
DROP INDEX IF EXISTS idx_organization_domains_verified;
DROP INDEX IF EXISTS idx_organization_domains_org_id;
DROP INDEX IF EXISTS idx_organization_members_role;
DROP INDEX IF EXISTS idx_organization_members_user_id;
DROP INDEX IF EXISTS idx_organization_members_org_id;
DROP INDEX IF EXISTS idx_organizations_subscription_status;
DROP INDEX IF EXISTS idx_organizations_created_by;
DROP INDEX IF EXISTS idx_organizations_domain;
DROP INDEX IF EXISTS idx_organizations_slug;
DROP TABLE IF EXISTS organization_domains;
DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS organizations;
