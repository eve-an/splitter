-- name: ListFeatures :many
SELECT
  f.id AS feature_id,
  f.name AS feature_name,
  f.description AS feature_description,
  f.active AS feature_active,
  f.created_at AS feature_created_at,
  v.id AS variant_id,
  v.name AS variant_name,
  v.weight AS variant_weight
FROM features f
LEFT JOIN variants v ON f.id = v.feature_id
ORDER BY f.id;

-- name: GetFeature :many
SELECT
  f.id AS feature_id,
  f.name AS feature_name,
  f.description AS feature_description,
  f.active AS feature_active,
  f.created_at AS feature_created_at,
  v.id AS variant_id,
  v.name AS variant_name,
  v.weight AS variant_weight
FROM features f
LEFT JOIN variants v ON f.id = v.feature_id
WHERE f.id = $1;

-- name: InsertFeature :one
INSERT INTO features (name, description, active)
VALUES ($1, $2, $3)
RETURNING id;

-- name: InsertVariant :one
INSERT INTO variants (feature_id, name, weight)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateFeature :exec
UPDATE features
SET name = $1,
    description = $2,
    active = $3
WHERE id = $4;

-- name: DeleteVariantsByFeature :exec
DELETE FROM variants WHERE feature_id = $1;

-- name: DeleteFeature :exec
DELETE FROM features WHERE id = $1;
