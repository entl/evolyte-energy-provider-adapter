-- name: CreateInverter :one
INSERT INTO inverters (
    user_id,
    vendor,
    model,
    serial_number,
    total_lifetime_production_kwh,
    installation_date,
    created_at,
    updated_at
)
VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8
)
RETURNING
    id, user_id, vendor, model, serial_number,
    total_lifetime_production_kwh, installation_date, created_at, updated_at;

-- name: GetInverterById :one
SELECT * FROM inverters WHERE id = $1;

-- name: GetInvertersByUserId :many
SELECT * FROM inverters WHERE user_id = $1;

-- name: GetInverters :many
SELECT * FROM inverters LIMIT $1 OFFSET $2;

-- name: UpdateInverter :exec
UPDATE inverters
SET
    vendor = COALESCE(sqlc.narg(vendor), vendor),
    model = COALESCE(sqlc.narg(model), model),
    serial_number = COALESCE(sqlc.narg(serial_number), serial_number),
    total_lifetime_production_kwh = COALESCE(sqlc.narg(total_lifetime_production_kwh), total_lifetime_production_kwh),
    installation_date = COALESCE(sqlc.narg(installation_date), installation_date),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteInverter :exec
DELETE FROM inverters WHERE id = $1;