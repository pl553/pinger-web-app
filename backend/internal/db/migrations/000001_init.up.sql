CREATE TYPE container_status AS ENUM ('UP', 'DOWN');

CREATE TABLE containers (
    container_ip VARCHAR(39) PRIMARY KEY,
    ping_time_ms INT NOT NULL,
    last_successful_ping TIMESTAMPTZ,
    status container_status NOT NULL
);

