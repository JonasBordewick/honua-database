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

CREATE TABLE IF NOT EXISTS hass_services (
    id SERIAL PRIMARY KEY,
    identity TEXT NOT NULL,
    CONSTRAINT fk_identity FOREIGN KEY(identity) REFERENCES identities(identifier) ON DELETE CASCADE,
    domain TEXT NOT NULL,
    name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT uc_domain UNIQUE (identity, domain)
);

CREATE TABLE IF NOT EXISTS allowed_services (
    entity_id INTEGER NOT NULL,
    CONSTRAINT fk_entity_id FOREIGN KEY(entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    service_id INTEGER NOT NULL,
    CONSTRAINT fk_service_id FOREIGN KEY(service_id) REFERENCES hass_services(id) ON DELETE CASCADE,
    PRIMARY KEY(entity_id, service_id)
);

CREATE TABLE IF NOT EXISTS allowed_sensors (
    device_id INTEGER NOT NULL,
    CONSTRAINT fk_device_id FOREIGN KEY(device_id) REFERENCES entities(id) ON DELETE CASCADE,
    sensor_id INTEGER NOT NULL,
    CONSTRAINT fk_sensor_id FOREIGN KEY(sensor_id) REFERENCES entities(id) ON DELETE CASCADE,
    PRIMARY KEY(device_id, sensor_id)
);

CREATE TABLE IF NOT EXISTS conditions (
    id SERIAL PRIMARY KEY,
    type INTEGER NOT NULL,
    sensor_id INTEGER,
    CONSTRAINT fk_sensor_id FOREIGN KEY(sensor_id) REFERENCES entities(id) ON DELETE CASCADE,
    before TEXT,
    after TEXT,
    below INTEGER,
    above INTEGER,
    comparison_state TEXT,
    parent_id INTEGER,
    CONSTRAINT fk_parent_id FOREIGN KEY(parent_id) REFERENCES conditions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    entitiy_id INTEGER NOT NULL,
    CONSTRAINT fk_entitiy_id FOREIGN KEY(entitiy_id) REFERENCES entities(id) ON DELETE CASCADE,
    event_based_evaluation BOOLEAN NOT NULL,
    periodic_trigger_type INTEGER,
    description TEXT NOT NULL,
    condition_id INTEGER NOT NULL,
    CONSTRAINT fk_condition_id FOREIGN KEY(condition_id) REFERENCES conditions(id) ON DELETE CASCADE,
    is_enabled BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER NOT NULL,
    CONSTRAINT fk_rule_id FOREIGN KEY(rule_id) REFERENCES rules(id) ON DELETE CASCADE,
    is_then_action BOOLEAN NOT NULL DEFAULT true,
    service_id INTEGER NOT NULL,
    CONSTRAINT fk_service_id FOREIGN KEY(service_id) REFERENCES hass_services(id) ON DELETE CASCADE
);