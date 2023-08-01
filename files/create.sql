CREATE TABLE IF NOT EXISTS identities (
    identifier TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY,
    identity TEXT NOT NULL,
    CONSTRAINT fk_identity FOREIGN KEY(identity) REFERENCES identities(identifier) ON DELETE CASCADE,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    is_device BOOLEAN NOT NULL,
    allow_rules BOOLEAN NOT NULL,
    has_attribute BOOLEAN NOT NULL,
    attribute TEXT,
    is_victron_sensor BOOLEAN NOT NULL,
    has_numeric_state BOOLEAN NOT NULL,
    CONSTRAINT uc_entity UNIQUE (identity, entity_id)
);

CREATE TABLE IF NOT EXISTS states (
    id SERIAL PRIMARY KEY,
    entity_id INTEGER NOT NULL,
    CONSTRAINT fk_entity_id FOREIGN KEY(entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    state TEXT NOT NULL,
    record_time TIMESTAMPTZ DEFAULT timezone('Europe/Berlin', now()) NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata (
    id SERIAL PRIMARY KEY,
    filepath TEXT NOT NULL,
    executed_at TIMESTAMPTZ DEFAULT timezone('Europe/Berlin', now()) NOT NULL
);