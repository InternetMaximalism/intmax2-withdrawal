-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE withdrawals (
    id                    uuid not null default uuid_generate_v4(),
    status                int not null,
    transfer_data         json not null,
    transfer_merkle_proof json not null,
    transaction           json not null,
    tx_merkle_proof       json not null,
    enough_balance_proof  json not null,
    transfer_hash         varchar(255) not null UNIQUE,
    block_number          int not null,
    block_hash            varchar(255) not null,
    created_at            timestamptz not null default now(),
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE withdrawals;
