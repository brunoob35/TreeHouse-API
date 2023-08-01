USE (database);

DROP database IF EXISTS treehousedb;
create database treehousedb;
CREATE TABLE usuarios (
                          id_usuario INT NOT NULL AUTO_INCREMENT,
                          nome_usuario VARCHAR(50) NOT NULL,
                          email_usuario VARCHAR(100) NOT NULL unique,
                          senha VARCHAR(255) NOT NULL,
                          id_acesso INT NOT NULL DEFAULT 0,
                          id_funcao INT NOT NULL DEFAULT 0,
                          cpf VARCHAR(14),
                          rg VARCHAR(20),
                          celular VARCHAR(20),
                          data_nascimento DATE,
                          ativo TINYINT DEFAULT 1,
                          data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          PRIMARY KEY (id_usuario)
);

-- Tabela Ususarios
USE treehousedb;

-- Tabela Aluno
CREATE TABLE alunos (
                        id_aluno INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        nome VARCHAR(255),
                        ativo TINYINT DEFAULT 1,
                        data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

);

-- Tabela Profesora
CREATE TABLE profesoras (
                            id_professora INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                            nome VARCHAR(255),
                            ativo TINYINT DEFAULT 1,
                            data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela Livros
CREATE TABLE livros (
                        id_livros INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        nome_livro VARCHAR(255)
);

-- Tabela Turma
CREATE TABLE turmas (
                        id_turma INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        nome_turma VARCHAR(50),
                        id_professor INT UNSIGNED,
                        ativo TINYINT DEFAULT 1,
                        data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela Aulas
CREATE TABLE alunos_aulas (
                              id_aluno INT,
                              id_aula INT
);

-- Tabela professores nas aulas
CREATE TABLE professores_aulas (
                                   id_professor INT,
                                   id_aula INT
);

-- Tabela aulas
CREATE TABLE aulas (
                       id_aula INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       datahora_aula DATE,
                       id_turma INT,
                       id_status_aula INT,
                       data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

);

-- Tabela alunos nas turmas
CREATE TABLE alunos_turmas (
                               id_aluno INT,
                               id_turma INT
);

-- Tabela Status Aulas
CREATE TABLE status_aulas (
                              id_status INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                              status_nome VARCHAR(255)
);

-- Tabela Anos Letivos
CREATE TABLE anos_letivos (
                              id_ano_letivo INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                              ano_letivo_nome VARCHAR(255)
);

-- Tabela Salarios
CREATE TABLE salarios (
                          id_salario INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                          valor VARCHAR(255),
                          data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela Pais e Alunos
CREATE TABLE pais_e_alunos (
                               id_pai INT,
                               id_aluno INT
);

CREATE TABLE acessos_ids (
                             id_acesso INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                             acesso_nome VARCHAR(255)
);