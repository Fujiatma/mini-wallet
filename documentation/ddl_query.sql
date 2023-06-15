-- ddl
CREATE TABLE customers
(
    id           VARCHAR(36) PRIMARY KEY,
    name         VARCHAR(75) NOT NULL,
    customer_xid VARCHAR(75) UNIQUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallets
(
    id          VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) UNIQUE,
    status      VARCHAR(75),
    enabled_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    balance     INT       DEFAULT 0,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions
(
    id               VARCHAR(36) PRIMARY KEY,
    wallet_id        VARCHAR(36),
    transaction_type ENUM('deposit', 'withdrawal'),
    amount           INT          NOT NULL,
    reference_id     VARCHAR(75) NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE deposits
(
    id           VARCHAR(36) PRIMARY KEY,
    wallet_id    VARCHAR(36),
    deposited_by VARCHAR(75),
    status       VARCHAR(75),
    deposited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    amount       INT          NOT NULL,
    reference_id VARCHAR(75) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE withdrawals
(
    id           VARCHAR(36) PRIMARY KEY,
    wallet_id    VARCHAR(36),
    withdrawn_by VARCHAR(75),
    status       VARCHAR(75),
    withdrawn_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    amount       INT          NOT NULL,
    reference_id VARCHAR(75) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE wallets
    ADD CONSTRAINT fk_wallets_customers
        FOREIGN KEY (customer_id)
            REFERENCES customers (id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE;

ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_wallets
        FOREIGN KEY (wallet_id)
            REFERENCES wallets (id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE;

ALTER TABLE deposits
    ADD CONSTRAINT fk_deposits_wallets
        FOREIGN KEY (wallet_id)
            REFERENCES wallets (id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE;

ALTER TABLE withdrawals
    ADD CONSTRAINT fk_withdrawals_wallets
        FOREIGN KEY (wallet_id)
            REFERENCES wallets (id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE;