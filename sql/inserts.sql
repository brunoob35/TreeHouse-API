INSERT INTO
    treehousedb.acessos_ids (acesso_nome)
VALUES
    ('Gerente Sistema'), ('Gerente Organização'), ('Professor'), ('Assistente Sistema');

insert into
    treehousedb.usuarios (nome_usuario, email_usuario, senha, id_acesso, cpf, rg, celular, data_nascimento)
VALUES
    ('Bruno Quoos', 'bruno@treehouse.com', '123456789', 1, '0000000000', '00000000000', '43 996630496', '2023-12-06');
