USE treehousedb;

SET NAMES utf8mb4;
SET @app_usuario_id = NULL;

-- =========================================================
-- CARGA DE TESTE - UNIVERSO SIMPSONS
-- =========================================================
-- Objetivo:
-- Popular clientes, alunos, turmas, aulas e contratos sem criar usuarios.
--
-- Caracteristicas:
-- - idempotente para os registros desta carga
-- - usa nomes de personagens dos Simpsons
-- - usa CPFs falsos, mas validos
-- - cria turmas de grupo e turmas VIP vinculadas a contratos
--
-- Execucao sugerida:
-- SOURCE /caminho/para/sql/seed_simpsons_test_data.sql;
-- =========================================================

-- ---------------------------------------------------------
-- GARANTIAS MINIMAS DE DOMINIO
-- ---------------------------------------------------------

INSERT INTO aulas_status (id, nome_status) VALUES
    (1, 'pendente'),
    (2, 'realizada'),
    (3, 'cancelada'),
    (4, 'remarcada'),
    (5, 'pendente reagendamento'),
    (6, 'indenizada')
ON DUPLICATE KEY UPDATE
    nome_status = VALUES(nome_status);

INSERT INTO contratos_status (id, nome_status) VALUES
    (1, 'Ativo'),
    (2, 'Pendente'),
    (3, 'Vencido'),
    (4, 'Inativo')
ON DUPLICATE KEY UPDATE
    nome_status = VALUES(nome_status);

INSERT INTO contratos_tipos (id, nome_tipo) VALUES
    (1, 'Anual'),
    (2, 'Semestral'),
    (3, 'Trimestral'),
    (4, 'Mensal'),
    (5, 'Temporário')
ON DUPLICATE KEY UPDATE
    nome_tipo = VALUES(nome_tipo);

-- ---------------------------------------------------------
-- LIMPEZA DA CARGA ANTERIOR
-- ---------------------------------------------------------

DELETE aa
FROM alunos_aulas aa
INNER JOIN aulas a
    ON a.id = aa.id_aula
INNER JOIN turmas t
    ON t.id = a.id_turma
WHERE t.nome LIKE 'Teste - Simpsons%';

DELETE c
FROM contratos c
INNER JOIN alunos a
    ON a.id = c.id_aluno
WHERE a.nome IN (
    'Bart Simpson',
    'Lisa Simpson',
    'Maggie Simpson',
    'Milhouse Van Houten',
    'Martin Prince',
    'Ralph Wiggum',
    'Rod Flanders',
    'Todd Flanders'
);

DELETE a
FROM aulas a
INNER JOIN turmas t
    ON t.id = a.id_turma
WHERE t.nome LIKE 'Teste - Simpsons%';

DELETE at
FROM alunos_turmas at
INNER JOIN turmas t
    ON t.id = at.id_turma
WHERE t.nome LIKE 'Teste - Simpsons%';

DELETE t
FROM turmas t
WHERE t.nome LIKE 'Teste - Simpsons%';

DELETE ca
FROM clientes_alunos ca
INNER JOIN clientes c
    ON c.id = ca.id_cliente
WHERE c.email IN (
    'homer.simpson.seed@example.com',
    'marge.simpson.seed@example.com',
    'kirk.vanhouten.seed@example.com',
    'luann.vanhouten.seed@example.com',
    'ned.flanders.seed@example.com',
    'maude.flanders.seed@example.com',
    'clancy.wiggum.seed@example.com',
    'seymour.skinner.seed@example.com',
    'edna.krabappel.seed@example.com',
    'apu.nahasapeemapetilon.seed@example.com'
);

DELETE ec
FROM enderecos_clientes ec
INNER JOIN clientes c
    ON c.id = ec.id_cliente
WHERE c.email IN (
    'homer.simpson.seed@example.com',
    'marge.simpson.seed@example.com',
    'kirk.vanhouten.seed@example.com',
    'luann.vanhouten.seed@example.com',
    'ned.flanders.seed@example.com',
    'maude.flanders.seed@example.com',
    'clancy.wiggum.seed@example.com',
    'seymour.skinner.seed@example.com',
    'edna.krabappel.seed@example.com',
    'apu.nahasapeemapetilon.seed@example.com'
);

DELETE FROM alunos
WHERE nome IN (
    'Bart Simpson',
    'Lisa Simpson',
    'Maggie Simpson',
    'Milhouse Van Houten',
    'Martin Prince',
    'Ralph Wiggum',
    'Rod Flanders',
    'Todd Flanders'
);

DELETE FROM clientes
WHERE email IN (
    'homer.simpson.seed@example.com',
    'marge.simpson.seed@example.com',
    'kirk.vanhouten.seed@example.com',
    'luann.vanhouten.seed@example.com',
    'ned.flanders.seed@example.com',
    'maude.flanders.seed@example.com',
    'clancy.wiggum.seed@example.com',
    'seymour.skinner.seed@example.com',
    'edna.krabappel.seed@example.com',
    'apu.nahasapeemapetilon.seed@example.com'
);

DELETE FROM enderecos
WHERE rua LIKE 'Rua Seed - %'
   OR rua LIKE 'Avenida Seed - %';

-- ---------------------------------------------------------
-- CLIENTES
-- ---------------------------------------------------------

INSERT INTO clientes (
    nome,
    email,
    cpf,
    rg,
    telefone,
    ativo,
    nascimento,
    lgpd_aceito,
    lgpd_aceito_em,
    lgpd_finalidade
) VALUES
    (
        'Homer Simpson',
        'homer.simpson.seed@example.com',
        '12345678909',
        'HS74201',
        '11987650001',
        TRUE,
        '1980-05-12',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Marge Simpson',
        'marge.simpson.seed@example.com',
        '98765432100',
        'MS74202',
        '11987650002',
        TRUE,
        '1982-03-19',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Kirk Van Houten',
        'kirk.vanhouten.seed@example.com',
        '31415926590',
        'KVH1001',
        '11987650003',
        TRUE,
        '1979-08-09',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Luann Van Houten',
        'luann.vanhouten.seed@example.com',
        '27182818205',
        'LVH1002',
        '11987650004',
        TRUE,
        '1981-10-17',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Ned Flanders',
        'ned.flanders.seed@example.com',
        '24681357928',
        'NF01001',
        '11987650005',
        TRUE,
        '1978-01-10',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Maude Flanders',
        'maude.flanders.seed@example.com',
        '13579246828',
        'MF01002',
        '11987650006',
        TRUE,
        '1980-06-21',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Clancy Wiggum',
        'clancy.wiggum.seed@example.com',
        '16180339805',
        'CW22001',
        '11987650007',
        TRUE,
        '1977-11-04',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Seymour Skinner',
        'seymour.skinner.seed@example.com',
        '45678912364',
        'SS33001',
        '11987650008',
        TRUE,
        '1975-09-30',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Edna Krabappel',
        'edna.krabappel.seed@example.com',
        '11235813207',
        'EK33002',
        '11987650009',
        TRUE,
        '1979-04-14',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    ),
    (
        'Apu Nahasapeemapetilon',
        'apu.nahasapeemapetilon.seed@example.com',
        '78912345664',
        'AN44001',
        '11987650010',
        TRUE,
        '1981-12-02',
        TRUE,
        CURRENT_TIMESTAMP,
        'Carga de teste para validacao de fluxos de clientes, contratos, aulas e turmas.'
    );

-- ---------------------------------------------------------
-- ENDERECOS DOS CLIENTES
-- ---------------------------------------------------------

INSERT INTO enderecos (
    cep,
    rua,
    numero,
    bairro,
    cidade,
    estado,
    pais,
    complemento
) VALUES
    ('01001001', 'Rua Seed - Evergreen Terrace Casa Homer', '742', 'Springfield Residencial', 'Springfield', 'SP', 'Brasil', 'Casa amarela'),
    ('01001002', 'Rua Seed - Evergreen Terrace Casa Marge', '742', 'Springfield Residencial', 'Springfield', 'SP', 'Brasil', 'Entrada lateral'),
    ('01001003', 'Rua Seed - Van Houten Residence', '53', 'Springfield Norte', 'Springfield', 'SP', 'Brasil', 'Portao azul'),
    ('01001004', 'Rua Seed - Luann Residence', '53', 'Springfield Norte', 'Springfield', 'SP', 'Brasil', 'Sobrado'),
    ('01001005', 'Rua Seed - Casa Flanders', '744', 'Evergreen District', 'Springfield', 'SP', 'Brasil', 'Ao lado dos Simpsons'),
    ('01001006', 'Rua Seed - Casa Maude', '744', 'Evergreen District', 'Springfield', 'SP', 'Brasil', 'Quintal amplo'),
    ('01001007', 'Rua Seed - Delegacia Wiggum', '12', 'Centro', 'Springfield', 'SP', 'Brasil', 'Fundos da delegacia'),
    ('01001008', 'Rua Seed - Escola Skinner', '19', 'Centro Escolar', 'Springfield', 'SP', 'Brasil', 'Sala da diretoria'),
    ('01001009', 'Rua Seed - Sala Krabappel', '19', 'Centro Escolar', 'Springfield', 'SP', 'Brasil', 'Bloco B'),
    ('01001010', 'Avenida Seed - Kwik-E-Mart', '24', 'Springfield Sul', 'Springfield', 'SP', 'Brasil', 'Loja principal');

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Evergreen Terrace Casa Homer'
WHERE c.email = 'homer.simpson.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Evergreen Terrace Casa Marge'
WHERE c.email = 'marge.simpson.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Van Houten Residence'
WHERE c.email = 'kirk.vanhouten.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Luann Residence'
WHERE c.email = 'luann.vanhouten.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Casa Flanders'
WHERE c.email = 'ned.flanders.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Casa Maude'
WHERE c.email = 'maude.flanders.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Delegacia Wiggum'
WHERE c.email = 'clancy.wiggum.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Escola Skinner'
WHERE c.email = 'seymour.skinner.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Rua Seed - Sala Krabappel'
WHERE c.email = 'edna.krabappel.seed@example.com';

INSERT INTO enderecos_clientes (id_cliente, id_endereco)
SELECT c.id, e.id
FROM clientes c
INNER JOIN enderecos e
    ON e.rua = 'Avenida Seed - Kwik-E-Mart'
WHERE c.email = 'apu.nahasapeemapetilon.seed@example.com';

-- ---------------------------------------------------------
-- ALUNOS
-- ---------------------------------------------------------

INSERT INTO alunos (
    nome,
    livro,
    alfabetizacao,
    nascimento,
    ativo
) VALUES
    ('Bart Simpson', 'Colecao Fonica - Nivel 2', 'Em desenvolvimento', '2014-04-01', TRUE),
    ('Lisa Simpson', 'Leitura Avancada - Nivel 4', 'Avancada', '2016-05-09', TRUE),
    ('Maggie Simpson', 'Primeiros Sons', 'Inicial', '2021-01-12', TRUE),
    ('Milhouse Van Houten', 'Colecao Fonica - Nivel 2', 'Em desenvolvimento', '2014-07-01', TRUE),
    ('Martin Prince', 'Interpretacao de Texto - Nivel 3', 'Avancada', '2014-02-14', TRUE),
    ('Ralph Wiggum', 'Alfabetizacao Lúdica', 'Inicial', '2015-09-03', TRUE),
    ('Rod Flanders', 'Leituras em Familia', 'Intermediaria', '2016-08-08', TRUE),
    ('Todd Flanders', 'Leituras em Familia', 'Inicial', '2018-10-10', TRUE);

-- ---------------------------------------------------------
-- VINCULOS CLIENTE/ALUNO
-- ---------------------------------------------------------

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome IN ('Bart Simpson', 'Lisa Simpson', 'Maggie Simpson')
WHERE c.email = 'homer.simpson.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome IN ('Bart Simpson', 'Lisa Simpson', 'Maggie Simpson')
WHERE c.email = 'marge.simpson.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome = 'Milhouse Van Houten'
WHERE c.email = 'kirk.vanhouten.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome = 'Milhouse Van Houten'
WHERE c.email = 'luann.vanhouten.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome = 'Martin Prince'
WHERE c.email = 'seymour.skinner.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome = 'Martin Prince'
WHERE c.email = 'edna.krabappel.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome = 'Ralph Wiggum'
WHERE c.email = 'clancy.wiggum.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome IN ('Rod Flanders', 'Todd Flanders')
WHERE c.email = 'ned.flanders.seed@example.com';

INSERT INTO clientes_alunos (id_cliente, id_aluno)
SELECT c.id, a.id
FROM clientes c
INNER JOIN alunos a
    ON a.nome IN ('Rod Flanders', 'Todd Flanders')
WHERE c.email = 'maude.flanders.seed@example.com';

-- ---------------------------------------------------------
-- ENDERECOS DAS TURMAS
-- ---------------------------------------------------------

INSERT INTO enderecos (
    cep,
    rua,
    numero,
    bairro,
    cidade,
    estado,
    pais,
    complemento
) VALUES
    ('02002001', 'Rua Seed - Escola Grupo Springfield', '100', 'Centro Escolar', 'Springfield', 'SP', 'Brasil', 'Sala 01'),
    ('02002002', 'Rua Seed - Escola Grupo Lisa Ralph', '101', 'Centro Escolar', 'Springfield', 'SP', 'Brasil', 'Sala 02'),
    ('02002003', 'Rua Seed - Estudio Flanders', '744', 'Evergreen District', 'Springfield', 'SP', 'Brasil', 'Sala de musica'),
    ('02002004', 'Rua Seed - VIP Bart', '742', 'Springfield Residencial', 'Springfield', 'SP', 'Brasil', 'Biblioteca'),
    ('02002005', 'Rua Seed - VIP Lisa', '742', 'Springfield Residencial', 'Springfield', 'SP', 'Brasil', 'Escritorio'),
    ('02002006', 'Rua Seed - VIP Milhouse', '53', 'Springfield Norte', 'Springfield', 'SP', 'Brasil', 'Quarto de estudos'),
    ('02002007', 'Rua Seed - VIP Ralph', '12', 'Centro', 'Springfield', 'SP', 'Brasil', 'Sala reservada');

-- ---------------------------------------------------------
-- TURMAS
-- ---------------------------------------------------------

INSERT INTO turmas (
    id_professor,
    id_endereco,
    nome,
    descricao_recorrencia,
    recorrencia_json
) VALUES
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - Escola Grupo Springfield' LIMIT 1),
        'Teste - Simpsons Grupo Springfield',
        'Segundas e quartas, 09:00, 8 aulas',
        '{"weekdays":["segunda","quarta"],"lesson_count":8,"start_date":"2026-06-15","start_time":"09:00"}'
    ),
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - Escola Grupo Lisa Ralph' LIMIT 1),
        'Teste - Simpsons Leitura Avancada',
        'Tercas e quintas, 14:00, 8 aulas',
        '{"weekdays":["terca","quinta"],"lesson_count":8,"start_date":"2026-06-16","start_time":"14:00"}'
    ),
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - Estudio Flanders' LIMIT 1),
        'Teste - Simpsons Flanders em Casa',
        'Sabados, 10:00, 6 aulas',
        '{"weekdays":["sabado"],"lesson_count":6,"start_date":"2026-06-20","start_time":"10:00"}'
    ),
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - VIP Bart' LIMIT 1),
        'Teste - Simpsons VIP Bart',
        'Tercas, 16:00, 4 aulas',
        '{"weekdays":["terca"],"lesson_count":4,"start_date":"2026-06-16","start_time":"16:00"}'
    ),
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - VIP Lisa' LIMIT 1),
        'Teste - Simpsons VIP Lisa',
        'Quartas, 17:00, 4 aulas',
        '{"weekdays":["quarta"],"lesson_count":4,"start_date":"2026-06-17","start_time":"17:00"}'
    ),
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - VIP Milhouse' LIMIT 1),
        'Teste - Simpsons VIP Milhouse',
        'Quintas, 18:00, 4 aulas',
        '{"weekdays":["quinta"],"lesson_count":4,"start_date":"2026-06-18","start_time":"18:00"}'
    ),
    (
        NULL,
        (SELECT id FROM enderecos WHERE rua = 'Rua Seed - VIP Ralph' LIMIT 1),
        'Teste - Simpsons VIP Ralph',
        'Sextas, 15:30, 4 aulas',
        '{"weekdays":["sexta"],"lesson_count":4,"start_date":"2026-05-15","start_time":"15:30"}'
    );

-- ---------------------------------------------------------
-- VINCULOS ALUNO/TURMA
-- ---------------------------------------------------------

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons Grupo Springfield'
WHERE a.nome IN ('Bart Simpson', 'Milhouse Van Houten', 'Martin Prince');

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons Leitura Avancada'
WHERE a.nome IN ('Lisa Simpson', 'Ralph Wiggum');

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons Flanders em Casa'
WHERE a.nome IN ('Rod Flanders', 'Todd Flanders');

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons VIP Bart'
WHERE a.nome = 'Bart Simpson';

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons VIP Lisa'
WHERE a.nome = 'Lisa Simpson';

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons VIP Milhouse'
WHERE a.nome = 'Milhouse Van Houten';

INSERT INTO alunos_turmas (id_aluno, id_turma)
SELECT a.id, t.id
FROM alunos a
INNER JOIN turmas t
    ON t.nome = 'Teste - Simpsons VIP Ralph'
WHERE a.nome = 'Ralph Wiggum';

-- ---------------------------------------------------------
-- AULAS
-- ---------------------------------------------------------

INSERT INTO aulas (
    id_status,
    id_professor,
    id_turma,
    assunto,
    vocabulario,
    saldo,
    observacoes,
    data_aula,
    data_aula_original,
    data_aula_solicitada
) VALUES
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Grupo Springfield' LIMIT 1),
        'Rimas e fonemas',
        'rima, som inicial, aliteracao',
        'Boa participacao do grupo',
        'Aula concluida com dinamica em equipe.',
        '2026-06-15 09:00:00',
        NULL,
        NULL
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Grupo Springfield' LIMIT 1),
        'Leitura guiada',
        'paragrafo, pausa, entonacao',
        'Bart e Milhouse precisaram de reforco em pontuacao.',
        'Martin liderou a leitura coletiva.',
        '2026-06-17 09:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Grupo Springfield' LIMIT 1),
        'Producao de frases',
        'sujeito, verbo, objeto',
        '',
        'Planejada para uso de cartoes ilustrados.',
        '2026-06-22 09:00:00',
        NULL,
        NULL
    ),
    (
        4,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Grupo Springfield' LIMIT 1),
        'Recontagem de historia',
        'sequencia, personagem, narrador',
        '',
        'Aula remarcada por conflito escolar.',
        '2026-06-24 09:00:00',
        '2026-06-23 09:00:00',
        '2026-06-24 09:00:00'
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Leitura Avancada' LIMIT 1),
        'Interpretacao de conto curto',
        'contexto, inferencia, protagonista',
        'Lisa avançou bem; Ralph precisou de apoio.',
        'Atividade com mapa mental.',
        '2026-06-16 14:00:00',
        NULL,
        NULL
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Leitura Avancada' LIMIT 1),
        'Leitura em voz alta',
        'diccao, ritmo, pausa',
        'Boa evolucao de fluencia.',
        'Uso de cronometro e gravacao.',
        '2026-06-18 14:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Leitura Avancada' LIMIT 1),
        'Resumo e reescrita',
        'resumo, coesao, conectivos',
        '',
        'Aula futura para consolidacao.',
        '2026-06-23 14:00:00',
        NULL,
        NULL
    ),
    (
        5,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Leitura Avancada' LIMIT 1),
        'Compreensao textual',
        'titulo, ideia central, evidencia',
        '',
        'Aguardando nova data conforme disponibilidade da familia.',
        '2026-06-25 14:00:00',
        '2026-06-25 14:00:00',
        '2026-06-27 10:00:00'
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Flanders em Casa' LIMIT 1),
        'Leitura biblica infantil',
        'capitulo, verso, moral',
        '',
        'Primeira aula em ambiente domiciliar.',
        '2026-06-20 10:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons Flanders em Casa' LIMIT 1),
        'Escrita de pequenas frases',
        'frase, espaco, letra maiuscula',
        '',
        'Aula programada para o fim de semana seguinte.',
        '2026-06-27 10:00:00',
        NULL,
        NULL
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Bart' LIMIT 1),
        'Ortografia contextual',
        'ss, ç, ch',
        'Bart respondeu bem a exercicios competitivos.',
        'Sessao individual concluida.',
        '2026-06-16 16:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Bart' LIMIT 1),
        'Leitura com tirinhas',
        'quadrinho, fala, balão',
        '',
        'Aula futura com material tematico.',
        '2026-06-23 16:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Bart' LIMIT 1),
        'Redacao curta',
        'inicio, meio, fim',
        '',
        'Atividade de historia curta.',
        '2026-06-30 16:00:00',
        NULL,
        NULL
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Lisa' LIMIT 1),
        'Interpretacao de texto poetico',
        'estrofe, verso, metafora',
        'Lisa demonstrou autonomia elevada.',
        'Aula individual com excelente desempenho.',
        '2026-06-17 17:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Lisa' LIMIT 1),
        'Analise de argumentos',
        'opiniao, tese, evidencia',
        '',
        'Aula planejada para debate orientado.',
        '2026-06-24 17:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Lisa' LIMIT 1),
        'Escrita autoral',
        'coerencia, repertorio, revisao',
        '',
        'Produção de pequeno ensaio.',
        '2026-07-01 17:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Milhouse' LIMIT 1),
        'Fluencia de leitura',
        'ritmo, entonacao, pausa',
        '',
        'Primeira aula individual de reforco.',
        '2026-06-18 18:00:00',
        NULL,
        NULL
    ),
    (
        1,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Milhouse' LIMIT 1),
        'Vocabulário aplicado',
        'sinonimo, antonimo, contexto',
        '',
        'Aula futura com cartas de vocabulario.',
        '2026-06-25 18:00:00',
        NULL,
        NULL
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Ralph' LIMIT 1),
        'Consciência fonologica',
        'som final, rima, silaba',
        'Ralph manteve atencao parcial, mas completou a atividade.',
        'Aula realizada com apoio visual.',
        '2026-05-15 15:30:00',
        NULL,
        NULL
    ),
    (
        2,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Ralph' LIMIT 1),
        'Associacao imagem-palavra',
        'imagem, palavra, som',
        'Evolucao discreta porém consistente.',
        'Sessao concluida com jogos.',
        '2026-05-22 15:30:00',
        NULL,
        NULL
    ),
    (
        3,
        NULL,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Ralph' LIMIT 1),
        'Sequencia alfabetica',
        'ordem, letra inicial, reconhecimento',
        '',
        'Aula cancelada por indisponibilidade da familia.',
        '2026-05-29 15:30:00',
        NULL,
        NULL
    );

-- ---------------------------------------------------------
-- CONTRATOS
-- ---------------------------------------------------------

INSERT INTO contratos (
    id_cliente_representante,
    id_cliente_responsavel,
    id_aluno,
    id_tipo_contrato,
    id_status,
    id_turma,
    valor,
    email_representante,
    cpf_representante,
    rg,
    telefone_representante,
    est_civil_representante,
    desconto_porcentagem,
    valor_final,
    parcelas,
    parcelas_descricao,
    numero_aulas,
    periodicidade,
    tempo_aula,
    tempo_contrato,
    inicio_contrato,
    vencimento_contrato,
    primeira_aula
) VALUES
    (
        (SELECT id FROM clientes WHERE email = 'marge.simpson.seed@example.com' LIMIT 1),
        (SELECT id FROM clientes WHERE email = 'homer.simpson.seed@example.com' LIMIT 1),
        (SELECT id FROM alunos WHERE nome = 'Bart Simpson' LIMIT 1),
        1,
        1,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Bart' LIMIT 1),
        4800.00,
        'marge.simpson.seed@example.com',
        '98765432100',
        'MS74202',
        '11987650002',
        'Casada',
        6.25,
        4500.00,
        12,
        '12x de R$ 375,00',
        24,
        '1x por semana',
        '50 min',
        '12 meses',
        '2026-06-16 00:00:00',
        '2026-06-22 23:59:59',
        '2026-06-16 16:00:00'
    ),
    (
        (SELECT id FROM clientes WHERE email = 'homer.simpson.seed@example.com' LIMIT 1),
        (SELECT id FROM clientes WHERE email = 'marge.simpson.seed@example.com' LIMIT 1),
        (SELECT id FROM alunos WHERE nome = 'Lisa Simpson' LIMIT 1),
        2,
        1,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Lisa' LIMIT 1),
        3600.00,
        'homer.simpson.seed@example.com',
        '12345678909',
        'HS74201',
        '11987650001',
        'Casado',
        5.00,
        3420.00,
        6,
        '6x de R$ 570,00',
        16,
        '1x por semana',
        '50 min',
        '6 meses',
        '2026-06-17 00:00:00',
        '2026-12-20 23:59:59',
        '2026-06-17 17:00:00'
    ),
    (
        (SELECT id FROM clientes WHERE email = 'luann.vanhouten.seed@example.com' LIMIT 1),
        (SELECT id FROM clientes WHERE email = 'kirk.vanhouten.seed@example.com' LIMIT 1),
        (SELECT id FROM alunos WHERE nome = 'Milhouse Van Houten' LIMIT 1),
        4,
        2,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Milhouse' LIMIT 1),
        1200.00,
        'luann.vanhouten.seed@example.com',
        '27182818205',
        'LVH1002',
        '11987650004',
        'Divorciada',
        0.00,
        1200.00,
        4,
        '4x de R$ 300,00',
        8,
        '1x por semana',
        '50 min',
        '1 mes',
        '2026-06-18 00:00:00',
        '2026-07-18 23:59:59',
        '2026-06-18 18:00:00'
    ),
    (
        (SELECT id FROM clientes WHERE email = 'clancy.wiggum.seed@example.com' LIMIT 1),
        (SELECT id FROM clientes WHERE email = 'clancy.wiggum.seed@example.com' LIMIT 1),
        (SELECT id FROM alunos WHERE nome = 'Ralph Wiggum' LIMIT 1),
        3,
        3,
        (SELECT id FROM turmas WHERE nome = 'Teste - Simpsons VIP Ralph' LIMIT 1),
        1800.00,
        'clancy.wiggum.seed@example.com',
        '16180339805',
        'CW22001',
        '11987650007',
        'Casado',
        0.00,
        1800.00,
        3,
        '3x de R$ 600,00',
        12,
        '1x por semana',
        '50 min',
        '3 meses',
        '2026-03-01 00:00:00',
        '2026-05-30 23:59:59',
        '2026-05-15 15:30:00'
    ),
    (
        (SELECT id FROM clientes WHERE email = 'maude.flanders.seed@example.com' LIMIT 1),
        (SELECT id FROM clientes WHERE email = 'ned.flanders.seed@example.com' LIMIT 1),
        (SELECT id FROM alunos WHERE nome = 'Rod Flanders' LIMIT 1),
        4,
        2,
        NULL,
        950.00,
        'maude.flanders.seed@example.com',
        '13579246828',
        'MF01002',
        '11987650006',
        'Casada',
        0.00,
        950.00,
        1,
        'Parcela unica',
        4,
        '1x por semana',
        '50 min',
        '1 mes',
        '2026-06-20 00:00:00',
        '2026-07-20 23:59:59',
        '2026-06-20 10:00:00'
    ),
    (
        (SELECT id FROM clientes WHERE email = 'ned.flanders.seed@example.com' LIMIT 1),
        (SELECT id FROM clientes WHERE email = 'maude.flanders.seed@example.com' LIMIT 1),
        (SELECT id FROM alunos WHERE nome = 'Todd Flanders' LIMIT 1),
        5,
        4,
        NULL,
        800.00,
        'ned.flanders.seed@example.com',
        '24681357928',
        'NF01001',
        '11987650005',
        'Viúvo',
        0.00,
        800.00,
        2,
        '2x de R$ 400,00',
        4,
        '1x por semana',
        '50 min',
        '2 meses',
        '2026-02-01 00:00:00',
        '2026-03-31 23:59:59',
        '2026-02-07 10:00:00'
    );

-- ---------------------------------------------------------
-- RESUMO RAPIDO
-- ---------------------------------------------------------
SELECT
    (SELECT COUNT(*) FROM clientes WHERE email LIKE '%.seed@example.com') AS clientes_seed,
    (SELECT COUNT(*) FROM alunos WHERE nome IN (
        'Bart Simpson',
        'Lisa Simpson',
        'Maggie Simpson',
        'Milhouse Van Houten',
        'Martin Prince',
        'Ralph Wiggum',
        'Rod Flanders',
        'Todd Flanders'
    )) AS alunos_seed,
    (SELECT COUNT(*) FROM turmas WHERE nome LIKE 'Teste - Simpsons%') AS turmas_seed,
    (SELECT COUNT(*) FROM aulas WHERE id_turma IN (
        SELECT id FROM turmas WHERE nome LIKE 'Teste - Simpsons%'
    )) AS aulas_seed,
    (SELECT COUNT(*) FROM contratos WHERE id_aluno IN (
        SELECT id FROM alunos WHERE nome IN (
            'Bart Simpson',
            'Lisa Simpson',
            'Maggie Simpson',
            'Milhouse Van Houten',
            'Martin Prince',
            'Ralph Wiggum',
            'Rod Flanders',
            'Todd Flanders'
        )
    )) AS contratos_seed;
