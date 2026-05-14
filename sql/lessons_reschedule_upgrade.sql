ALTER TABLE treehousedb.aulas
    ADD COLUMN data_aula_original DATETIME NULL AFTER data_aula,
    ADD COLUMN data_aula_solicitada DATETIME NULL AFTER data_aula_original;
