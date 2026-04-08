-- =========================================================
-- INSERTS INICIAIS
-- =========================================================

# INSERT INTO treehousedb.usuarios (senha, nome, email, ativo, cpf, rg)
# VALUES (
#            '$2a$10$6iGqzERlawUn1AnEOe49XOWHqYWe2M.4mES1h.dWraAv.KveQw1oy',
#            'Gestor Perfil',
#            'gestor@gestio.com',
#            TRUE,
#            19262973004,
#            273284101
#        );
#
# SET @gestor_id = LAST_INSERT_ID();
#
# INSERT INTO treehousedb.usuarios_permissoes (id_usuario, id_permissao)
# VALUES (@gestor_id, 1);
#
# INSERT INTO treehousedb.usuarios (senha, nome, email, ativo, cpf, rg)
# VALUES (
#            '$2a$10$2lT9uGk6h7dD4y3nE9bOAuq9dF8kL1aP2xC4fM9sN5vR7wT8yZ1Qa',
#            'Professor Demonstração',
#            'professor@gestio.com',
#            TRUE,
#            60396492088,
#            192078860
#        );
#
# SET @professor_id = LAST_INSERT_ID();

INSERT INTO treehousedb.usuarios_permissoes (id_usuario, id_permissao)
VALUES (@professor_id, 2);


-- IMPORTANTE:
-- Os IDs abaixo sao bit flags e devem continuar sendo potencias de 2.
INSERT INTO treehousedb.permissoes (id, nome) VALUES
                                      (1, 'gestao'),
                                      (2, 'professor'),
                                      (4, 'gestao master');

INSERT INTO treehousedb.aulas_status (id, nome_status) VALUES
                                               (1, 'pendente'),
                                               (2, 'realizada'),
                                               (3, 'cancelada'),
                                               (4, 'remarcada'),
                                               (5, 'pendente reagendamento'),
                                               (6, 'indenizada');



INSERT INTO treehousedb.alunos_turmas (id_aluno, id_turma) VALUES
                                                               (1, 1),
                                                               (2, 2),
                                                               (3, 3),
                                                               (4, 4),
                                                               (5, 5),
                                                               (6, 6),
                                                               (7, 7),
                                                               (8, 8),
                                                               (9, 9),
                                                               (10, 10),
                                                               (11, 1),
                                                               (12, 2),
                                                               (13, 3),
                                                               (14, 4),
                                                               (15, 5),
                                                               (16, 6),
                                                               (17, 7),
                                                               (18, 8),
                                                               (19, 9),
                                                               (20, 10),
                                                               (21, 1),
                                                               (22, 2),
                                                               (23, 3),
                                                               (24, 4),
                                                               (25, 5),
                                                               (26, 6),
                                                               (27, 7),
                                                               (28, 8),
                                                               (29, 9),
                                                               (30, 10),
                                                               (31, 1),
                                                               (32, 2),
                                                               (33, 3),
                                                               (34, 4),
                                                               (35, 5),
                                                               (36, 6),
                                                               (37, 7),
                                                               (38, 8),
                                                               (39, 9),
                                                               (40, 10),
                                                               (41, 1),
                                                               (42, 2),
                                                               (43, 3),
                                                               (44, 4),
                                                               (45, 5),
                                                               (46, 6),
                                                               (47, 7),
                                                               (48, 8),
                                                               (49, 9),
                                                               (50, 10),
                                                               (51, 1),
                                                               (52, 2),
                                                               (53, 3),
                                                               (54, 4),
                                                               (55, 5),
                                                               (56, 6),
                                                               (57, 7),
                                                               (58, 8),
                                                               (59, 9),
                                                               (60, 10),
                                                               (61, 1),
                                                               (62, 2),
                                                               (63, 3),
                                                               (64, 4),
                                                               (65, 5),
                                                               (66, 6),
                                                               (67, 7),
                                                               (68, 8),
                                                               (69, 9),
                                                               (70, 10);


-- Inserting 50 random entries with specific time slots
SET @start_date = '2023-12-01';
SET @end_date = '2023-12-10';

INSERT INTO treehousedb.aulas (datahora_aula, datahora_fim_aula, id_turma, id_status_aula)
VALUES
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 10:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 11:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 13:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5)),
    (CONCAT(DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 10) DAY), ' 14:00:00'), NULL, FLOOR(1 + RAND() * 10), FLOOR(1 + RAND() * 5));


INSERT INTO treehousedb.alunos_aulas (id_aluno, id_aula, id_professor)
VALUES
    (6,6,1),
    (56,6,1),
    (36,6,1),
    (1,1,1),
    (11,1,1),
    (21,1,1),
    (31,1,1),
    (41,1,1);

