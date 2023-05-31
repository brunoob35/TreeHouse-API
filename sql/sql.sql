USE (database);

CREATE TABLE usuarios (
                          id_usuario INT NOT NULL AUTO_INCREMENT,
                          nome_usuario VARCHAR(50) NOT NULL,
                          email_usuario VARCHAR(100) NOT NULL unique,
                          senha VARCHAR(255) NOT NULL,
                          id_acesso INT NOT NULL,
                          cpf VARCHAR(14),
                          rg VARCHAR(20),
                          celular VARCHAR(20),
                          data_nascimento DATE,
                          data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          PRIMARY KEY (id_usuario)
);



CREATE TABLE turma (
                       id_turma INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       nome_turma VARCHAR(50),
                       data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       id_professora INT UNSIGNED

);

use treehousedb;

-- Tabela Aluno
CREATE TABLE alunos (
                        id_aluno INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        nome VARCHAR(255),
                        idade INT,
                        datahora_aniversario DATETIME,
                        livro INT,
                        professora INT,
                        id_pais INT,
                        id_ano_letivo INT,
                        alfabetizacao INT,
                        id_turma INT,
                        data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

);

-- Tabela Profesora
CREATE TABLE profesoras (
                            id_professora INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                            id_usuario_professora INT,
                            nome VARCHAR(255),
                            aniversario_professora DATE,
                            id_turma INT,
                            id_salario INT,
                            data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

);

-- Tabela Livros
CREATE TABLE livros (
                        id_livros INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        nome_livro VARCHAR(255)
);

-- Tabela Turma
CREATE TABLE turmas (
                        id_turma INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        nome_turma VARCHAR(255) NOT NULL,
                        id_professor INT,
                        data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

-- Tabela Aulas
CREATE TABLE aulas (
                       id_aula INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       datahora_aula DATE,
                       id_turma INT,
                       id_status_aula INT,
                       data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

);

-- Tabela Alunos das turmas
CREATE TABLE alunos_turmas (
                               id_aluno INT,
                               id_turma INT
);

-- Tabela Alunos nas aulas
CREATE TABLE alunos_aulas (
                              id_aluno INT,
                              id_aula INT
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
                          data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
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