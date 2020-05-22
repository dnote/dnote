-- remove-billing-columns.sql drops billing related columns that are now obsolete.

-- +migrate Up

ALTER TABLE users DROP COLUMN IF EXISTS stripe_customer_id;
ALTER TABLE users DROP COLUMN IF EXISTS billing_country;

-- +migrate Down

