

CREATE TABLE stock_units (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name        VARCHAR(255)    NOT NULL,
    location    VARCHAR(255)    NOT NULL DEFAULT '',
    created_at  DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_stock_units_name ON stock_units (name);



CREATE TABLE categories (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name        VARCHAR(255)    NOT NULL,
    created_at  DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_categories_name ON categories (name);


CREATE TABLE subcategories (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name        VARCHAR(255)    NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    created_at  DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    CONSTRAINT fk_subcategories_category
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_subcategories_category ON subcategories (category_id);



CREATE TABLE equipments (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name            VARCHAR(255)    NOT NULL,
    serial_number   VARCHAR(255)    NOT NULL,
    category_id     BIGINT UNSIGNED NOT NULL,
    subcategory_id  BIGINT UNSIGNED NOT NULL,
    stock_unit_id   BIGINT UNSIGNED NOT NULL,
    created_at      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    CONSTRAINT fk_equipments_category
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_equipments_subcategory
        FOREIGN KEY (subcategory_id) REFERENCES subcategories (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_equipments_stock_unit
        FOREIGN KEY (stock_unit_id) REFERENCES stock_units (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_equipments_category    ON equipments (category_id);
CREATE INDEX idx_equipments_subcategory ON equipments (subcategory_id);
CREATE INDEX idx_equipments_stock_unit  ON equipments (stock_unit_id);
CREATE UNIQUE INDEX uk_equipments_serial ON equipments (serial_number);



CREATE TABLE stock_movements (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    equipment_id   BIGINT UNSIGNED NOT NULL,
    stock_unit_id  BIGINT UNSIGNED NOT NULL,
    movement_type  VARCHAR(10)     NOT NULL COMMENT 'ENTRADA ou SAIDA',
    quantity       INT             NOT NULL,
    movement_date  DATE            NOT NULL,
    created_at     DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    CONSTRAINT chk_movement_type CHECK (movement_type IN ('ENTRADA', 'SAIDA')),
    CONSTRAINT fk_movements_equipment
        FOREIGN KEY (equipment_id) REFERENCES equipments (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_movements_stock_unit
        FOREIGN KEY (stock_unit_id) REFERENCES stock_units (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_movements_equipment     ON stock_movements (equipment_id);
CREATE INDEX idx_movements_stock_unit   ON stock_movements (stock_unit_id);
CREATE INDEX idx_movements_equipment_date ON stock_movements (equipment_id, movement_date);
