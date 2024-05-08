CREATE TABLE IF NOT EXISTS "client_info"
(
    id                  BIGSERIAL PRIMARY KEY,
    contract_id         VARCHAR(255) NOT NULL,
    phone_number        VARCHAR(255) NOT NULL,
    address             VARCHAR(255) NOT NULL,
    payment_sum         VARCHAR(255) NOT NULL,
    comment             VARCHAR(255) NOT NULL,
    location_latitude   VARCHAR(255) NOT NULL,
    location_longitude  VARCHAR(255) NOT NULL,
    address_foto_path   VARCHAR(255) NOT NULL,
    payment_foto_path   VARCHAR(255) NOT NULL
);
