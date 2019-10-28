-- create-weekly-repetition.sql creates the default repetition rules for the users
-- that used to have the weekly email digest on Friday 20:00 UTC

-- +migrate Up

WITH next_friday AS (
  SELECT * FROM generate_series(
    CURRENT_DATE + INTERVAL '1 day',
    CURRENT_DATE + INTERVAL '7 days',
    INTERVAL '1 day'
  ) AS day
)
INSERT INTO repetition_rules
(
  user_id,
  title,
  enabled,
  hour,
  minute,
  frequency,
  last_active,
  book_domain,
  note_count,
  next_active
) SELECT
  t1.id,
  'Default weekly repetition',
  true,
  20,
  0,
  604800000, -- 7 days
  0,
  'all',
  20,
  extract(epoch FROM date_trunc('day', t1.day) + INTERVAL '20 hours') * 1000
FROM (
  SELECT * FROM users
  INNER JOIN next_friday ON EXTRACT(ISODOW FROM day) = '5' -- next friday
  WHERE users.cloud = true
) as t1;


-- +migrate Down

DELETE FROM repetition_rules WHERE title = 'Default weekly repetition' AND enabled AND hour = 8 AND minute = 0;
