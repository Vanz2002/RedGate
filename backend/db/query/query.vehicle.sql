-- name: GetVehicleByID :one
SELECT * FROM vehicle_data WHERE v_id = $1;

-- name: InsertVehicle :one
INSERT INTO vehicle_data (v_id, account_id, plate_number) VALUES ($1, $2, $3) RETURNING *;

-- name: GetPlateID :one
SELECT vd.v_id, a.is_subscribe
FROM vehicle_data AS vd
INNER JOIN account AS a ON vd.account_id = a.account_id
WHERE vd.account_id = $1;

-- name: VerifyVehicle :one
SELECT vd.plate_number
FROM vehicle_data AS vd
INNER JOIN account AS a ON vd.account_id = a.account_id
WHERE vd.v_id = $1 AND a.is_subscribe = true;