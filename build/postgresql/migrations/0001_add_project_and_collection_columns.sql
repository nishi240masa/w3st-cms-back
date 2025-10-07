-- Migration: add project_id and collection_ids to existing tables (idempotent)
-- Run this against the Postgres DB for existing deployments

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'api_collections' AND column_name = 'project_id'
  ) THEN
    ALTER TABLE api_collections ADD COLUMN project_id INT NOT NULL DEFAULT 1;
  END IF;
END
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'content_entries' AND column_name = 'project_id'
  ) THEN
    ALTER TABLE content_entries ADD COLUMN project_id INT NOT NULL DEFAULT 1;
  END IF;
END
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'api_keys' AND column_name = 'project_id'
  ) THEN
    ALTER TABLE api_keys ADD COLUMN project_id INT NOT NULL DEFAULT 1;
  END IF;
END
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'api_keys' AND column_name = 'collection_ids'
  ) THEN
    ALTER TABLE api_keys ADD COLUMN collection_ids INT[] NOT NULL DEFAULT '{}';
  END IF;
END
$$;