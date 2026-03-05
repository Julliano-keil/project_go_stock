-- Roda DEPOIS de 20250211_new_schema (ordem alfabética).
-- Remove coluna salt da tabela usuario.
ALTER TABLE usuario DROP COLUMN salt;

-- Remove id_unidade_estoque de item_estoque (FK primeiro, depois coluna).
ALTER TABLE item_estoque
    DROP FOREIGN KEY fk_item_unidade,
    DROP COLUMN id_unidade_estoque;
