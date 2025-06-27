-- +goose Up
-- Add soft delete column to debates table
ALTER TABLE debates ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- +goose Down
-- Remove soft delete column from debates table
ALTER TABLE debates DROP COLUMN IF EXISTS deleted_at; 