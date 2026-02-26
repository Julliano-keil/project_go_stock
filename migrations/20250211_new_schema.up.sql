
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS equipments;
DROP TABLE IF EXISTS subcategories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS stock_units;


CREATE TABLE categoria (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL
);

CREATE TABLE sub_categoria (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_categoria INT NOT NULL,
    nome VARCHAR(255) NOT NULL,
    CONSTRAINT fk_subcategoria_categoria
        FOREIGN KEY (id_categoria)
        REFERENCES categoria(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE unidades_de_estoque (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL
);

CREATE TABLE equipamento (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    id_subCategoria INT NOT NULL,
    id_unidade_estoque INT NOT NULL,
    CONSTRAINT fk_equipamento_subcategoria
        FOREIGN KEY (id_subCategoria)
        REFERENCES sub_categoria(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_equipamento_unidade
        FOREIGN KEY (id_unidade_estoque)
        REFERENCES unidades_de_estoque(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE item_estoque (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_equipamento INT NOT NULL,
    id_unidade_estoque INT NOT NULL,
    status_code VARCHAR(50) NOT NULL,
    codigo VARCHAR(100) NOT NULL,
    CONSTRAINT fk_item_equipamento
        FOREIGN KEY (id_equipamento)
        REFERENCES equipamento(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_item_unidade
        FOREIGN KEY (id_unidade_estoque)
        REFERENCES unidades_de_estoque(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE usuario (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL
);

CREATE TABLE movimentacao (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tipo_movimentacao VARCHAR(50) NOT NULL,
    data DATETIME NOT NULL,
    id_user INT NOT NULL,
    CONSTRAINT fk_movimentacao_usuario
        FOREIGN KEY (id_user)
        REFERENCES usuario(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE movimentacao_historico (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_movimentacao INT NOT NULL,
    id_item_estoque INT NOT NULL,
    CONSTRAINT fk_hist_movimentacao
        FOREIGN KEY (id_movimentacao)
        REFERENCES movimentacao(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_hist_item
        FOREIGN KEY (id_item_estoque)
        REFERENCES item_estoque(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
