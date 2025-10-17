-- name: InsertEvent :one
INSERT INTO events (feature_id, user_id, variant, event_type)
VALUES ($1, $2, $3, $4)
RETURNING id, feature_id, user_id, variant, event_type, created_at;

-- name: ListEventsByFeatureID :many
SELECT
  id,
  feature_id,
  user_id,
  variant,
  event_type,
  created_at
FROM events
WHERE feature_id = $1
ORDER BY created_at DESC;
