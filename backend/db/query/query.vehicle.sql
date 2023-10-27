-- name: GetVehicleByID :one
SELECT * FROM vehicle_data WHERE v_id = $1;

-- name: InsertVehicle :one
INSERT INTO vehicle_data (v_id, account_id, plate_number) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateVehicle :one
UPDATE vehicle_data SET account_id = $2, plate_number = $3 WHERE v_id = $1 RETURNING *;