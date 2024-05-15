CREATE TABLE IF NOT EXISTS "client_info"
(
    id                  BIGSERIAL PRIMARY KEY,
    contract_id         VARCHAR(255) DEFAULT '',
    phone_number        VARCHAR(255) DEFAULT '',
    address             VARCHAR(255) DEFAULT '',
    payment_sum         VARCHAR(255) DEFAULT '',
    comment             VARCHAR(255) DEFAULT '',
    location            VARCHAR(255) DEFAULT '',
    address_foto_path   VARCHAR(255) DEFAULT '',
    payment_foto_path   VARCHAR(255) DEFAULT '',
    user_name           VARCHAR(255) DEFAULT '',
    chat_id             BIGINT NOT NULL,
    step                FLOAT DEFAULT 0
);
