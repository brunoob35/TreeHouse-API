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
                       id_professora INT UNSIGNED
);