USE treehousedb;

SET @schema_name := DATABASE();

DROP PROCEDURE IF EXISTS apply_lgpd_patch;

DELIMITER $$

CREATE PROCEDURE apply_lgpd_patch()
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'nome_protegido'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN nome_protegido TEXT NULL AFTER nome;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'cpf_protegido'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN cpf_protegido TEXT NULL AFTER rg;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'rg_protegido'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN rg_protegido TEXT NULL AFTER cpf_protegido;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'cpf_hash'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN cpf_hash VARCHAR(64) NULL AFTER rg_protegido;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'rg_hash'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN rg_hash VARCHAR(64) NULL AFTER cpf_hash;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'lgpd_aceito'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN lgpd_aceito BOOLEAN NOT NULL DEFAULT FALSE AFTER nascimento;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'lgpd_aceito_em'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN lgpd_aceito_em DATETIME NULL AFTER lgpd_aceito;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND COLUMN_NAME = 'lgpd_finalidade'
  ) THEN
    ALTER TABLE usuarios
      ADD COLUMN lgpd_finalidade TEXT NULL AFTER lgpd_aceito_em;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND INDEX_NAME = 'uq_usuarios_cpf_hash'
  ) THEN
    ALTER TABLE usuarios
      ADD CONSTRAINT uq_usuarios_cpf_hash UNIQUE (cpf_hash);
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'usuarios'
      AND INDEX_NAME = 'uq_usuarios_rg_hash'
  ) THEN
    ALTER TABLE usuarios
      ADD CONSTRAINT uq_usuarios_rg_hash UNIQUE (rg_hash);
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'nome_protegido'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN nome_protegido TEXT NULL AFTER nome;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'cpf_protegido'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN cpf_protegido TEXT NULL AFTER rg;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'rg_protegido'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN rg_protegido TEXT NULL AFTER cpf_protegido;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'cpf_hash'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN cpf_hash VARCHAR(64) NULL AFTER rg_protegido;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'rg_hash'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN rg_hash VARCHAR(64) NULL AFTER cpf_hash;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'lgpd_aceito'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN lgpd_aceito BOOLEAN NOT NULL DEFAULT FALSE AFTER nascimento;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'lgpd_aceito_em'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN lgpd_aceito_em DATETIME NULL AFTER lgpd_aceito;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND COLUMN_NAME = 'lgpd_finalidade'
  ) THEN
    ALTER TABLE clientes
      ADD COLUMN lgpd_finalidade TEXT NULL AFTER lgpd_aceito_em;
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND INDEX_NAME = 'uq_clientes_cpf_hash'
  ) THEN
    ALTER TABLE clientes
      ADD CONSTRAINT uq_clientes_cpf_hash UNIQUE (cpf_hash);
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'clientes'
      AND INDEX_NAME = 'uq_clientes_rg_hash'
  ) THEN
    ALTER TABLE clientes
      ADD CONSTRAINT uq_clientes_rg_hash UNIQUE (rg_hash);
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @schema_name
      AND TABLE_NAME = 'alunos'
      AND COLUMN_NAME = 'nome_protegido'
  ) THEN
    ALTER TABLE alunos
      ADD COLUMN nome_protegido TEXT NULL AFTER nome;
  END IF;
END$$

DELIMITER ;

CALL apply_lgpd_patch();

DROP PROCEDURE IF EXISTS apply_lgpd_patch;
