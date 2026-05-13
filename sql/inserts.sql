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