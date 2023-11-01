DROP TABLE IF EXISTS vehicle_data;
DROP TABLE IF EXISTS account;

CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE TABLE account (
    account_id VARCHAR(32) DEFAULT 'UID' || nextval('user_id_seq') || to_char(current_timestamp, 'YYYYMMDDHH24MISS') || nextval('user_id_seq'),
    username VARCHAR(20) NOT NULL,
    email VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    is_subscribe BOOLEAN,
    PRIMARY KEY (account_id)
);

CREATE TABLE vehicle_data (
    v_id VARCHAR(32),
    account_id VARCHAR(32),
    plate_number VARCHAR(20), 
    PRIMARY KEY (v_id),
    FOREIGN KEY (account_id) REFERENCES account(account_id)
);

CREATE INDEX IF NOT EXISTS acc_id_index ON account(account_id);
CREATE INDEX IF NOT EXISTS v_id_index ON vehicle_data(v_id);