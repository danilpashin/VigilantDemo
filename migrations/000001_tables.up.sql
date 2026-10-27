CREATE TABLE customers (
    id uuid PRIMARY KEY,
    status text NOT NULL
        CHECK (status IN ('ACTIVE', 'RESTRICTED', 'CLOSED')),
    home_country char(2),
    home_city text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE cards (
    id uuid PRIMARY KEY,
    customer_id uuid NOT NULL REFERENCES customers(id),
    token text NOT NULL UNIQUE,
    last4 char(4) NOT NULL,
    status text NOT NULL
        CHECK (status IN ('ACTIVE', 'BLOCKED', 'EXPIRED', 'CLOSED')),
    expires_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_cards_customer_id ON cards (customer_id);

CREATE TABLE merchants (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    mcc char(4) NOT NULL,
    country char(2) NOT NULL,
    city text NOT NULL,
    risk_status text NOT NULL
        CHECK (risk_status IN ('KNOWN', 'VERIFIED', 'HIGH_RISK', 'BLOCKED')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE rulesets (
    version integer PRIMARY KEY CHECK (version > 0),
    content_checksum text NOT NULL,
    config_snapshot jsonb NOT NULL,
    activated_at timestamptz,
    deactivated_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (
        deactivated_at IS NULL
        OR activated_at IS NULL
        OR deactivated_at >= activated_at
    )
);

CREATE TABLE authorization_attempts (
    id uuid PRIMARY KEY,
    customer_id uuid NOT NULL REFERENCES customers(id),
    card_id uuid NOT NULL REFERENCES cards(id),
    merchant_id uuid REFERENCES merchants(id),
    occurred_at timestamptz NOT NULL,
    amount_minor bigint NOT NULL CHECK (amount_minor > 0),
    currency char(3) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_authorization_attempts_customer_id
    ON authorization_attempts (customer_id);
CREATE INDEX idx_authorization_attempts_card_id
    ON authorization_attempts (card_id);
CREATE INDEX idx_authorization_attempts_merchant_id
    ON authorization_attempts (merchant_id);
CREATE INDEX idx_authorization_attempts_occurred_at
    ON authorization_attempts (occurred_at);

CREATE TABLE risk_decisions (
    id uuid PRIMARY KEY,
    authorization_id uuid NOT NULL UNIQUE
        REFERENCES authorization_attempts(id),
    verdict text NOT NULL
        CHECK (verdict IN ('ALLOW', 'REVIEW', 'BLOCK')),
    risk_score integer NOT NULL CHECK (risk_score BETWEEN 0 AND 1000),
    ruleset_version integer NOT NULL REFERENCES rulesets(version),
    processing_time_ms integer NOT NULL CHECK (processing_time_ms >= 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE decision_rule_hits (
    decision_id uuid NOT NULL REFERENCES risk_decisions(id) ON DELETE CASCADE,
    rule_code text NOT NULL,
    score integer NOT NULL DEFAULT 0,
    PRIMARY KEY (decision_id, rule_code)
);

CREATE TABLE roles (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    login text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    role_id integer NOT NULL REFERENCES roles(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE audit_logs (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer REFERENCES users(id) ON DELETE SET NULL,
    action text NOT NULL,
    entity_type text,
    entity_id uuid,
    details jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
