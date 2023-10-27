DROP TABLE IF EXISTS vehicle_data;
DROP TABLE IF EXISTS account;

CREATE SEQUENCE user_id_seq;
CREATE TABLE account (
    account_id VARCHAR(255) DEFAULT 'UID' || nextval('user_id_seq') || to_char(current_timestamp, 'YYYYMMDDHH24MISS') || nextval('user_id_seq'),
    email VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    is_subscribe BOOLEAN,
    PRIMARY KEY (user_id)
);

CREATE TABLE vehicle_data (
    v_id VARCHAR(255) DEFAULT 'VID' || plate_number,
    account_id VARCHAR(255),
    plate_number VARCHAR(20), 
    PRIMARY KEY (v_id),
    FOREIGN KEY (account_id) REFERENCES account(account_id)
);

CREATE INDEX acc_id_index ON account(account_id);
CREATE INDEX v_id_index ON vehicle_data(v_id);