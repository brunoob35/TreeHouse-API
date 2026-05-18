ALTER TABLE treehousedb.turmas
    ADD COLUMN id_endereco INT UNSIGNED NULL AFTER id_professor,
    ADD CONSTRAINT fk_turmas_endereco
        FOREIGN KEY (id_endereco) REFERENCES treehousedb.enderecos (id)
            ON UPDATE CASCADE
            ON DELETE SET NULL,
    ADD INDEX idx_turmas_id_endereco (id_endereco);
