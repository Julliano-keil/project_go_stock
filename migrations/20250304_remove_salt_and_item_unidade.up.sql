
ALTER TABLE usuario DROP COLUMN salt;

ALTER TABLE item_estoque
    DROP FOREIGN KEY fk_item_unidade,
    DROP COLUMN id_unidade_estoque;
