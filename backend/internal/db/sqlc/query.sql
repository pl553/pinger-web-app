-- name: UpsertContainerStatus :exec
INSERT INTO containers (container_ip, ping_time_ms, last_successful_ping, status)
VALUES ($1, $2, $3, $4)
ON CONFLICT (container_ip)
DO UPDATE SET
  ping_time_ms = EXCLUDED.ping_time_ms,
  last_successful_ping = EXCLUDED.last_successful_ping,
  status = EXCLUDED.status;

-- name: ListContainerStatuses :many
SELECT container_ip, ping_time_ms, last_successful_ping, status
FROM containers
ORDER BY last_successful_ping DESC;

-- name: GetContainerStatusByIp :one
SELECT *
FROM containers
WHERE container_ip = $1;
