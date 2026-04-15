package stock_movement

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lince/datastore"
	"lince/entities"
)

type stockMovementRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

// NewStockMovementRepository cria um novo repositório de movimento de estoque.
func NewStockMovementRepository(settings datastore.SettingsRepository) datastore.StockMovementRepository {
	return stockMovementRepository{
		conn: settings.Connection,
	}
}

func (r stockMovementRepository) ListStockMovements(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.StockMovement, error) {
	db := r.conn(company)

	query := `
		SELECT
			m.id,
			m.tipo_movimentacao,
			m.data,
			m.id_user,
			(SELECT COUNT(*) FROM movimentacao_historico mh WHERE mh.id_movimentacao = m.id) AS quantidade
		FROM movimentacao m
		ORDER BY m.data DESC, m.id DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []entities.StockMovement
	for rows.Next() {
		var sm entities.StockMovement
		if err := rows.Scan(&sm.ID, &sm.TipoMovimentacao, &sm.Data, &sm.IDUser, &sm.Quantidade); err != nil {
			return nil, err
		}
		list = append(list, sm)
	}
	return list, rows.Err()
}

func (r stockMovementRepository) GetStockMovementByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.StockMovement, error) {
	db := r.conn(company)

	query := `
		SELECT
			m.id,
			m.tipo_movimentacao,
			m.data,
			m.id_user,
			(SELECT COUNT(*) FROM movimentacao_historico mh WHERE mh.id_movimentacao = m.id) AS quantidade
		FROM movimentacao m
		WHERE m.id = ?
	`

	var sm entities.StockMovement
	err := db.QueryRowContext(ctx, query, id).Scan(
		&sm.ID,
		&sm.TipoMovimentacao,
		&sm.Data,
		&sm.IDUser,
		&sm.Quantidade,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	idsItemEstoque, err := r.getMovementItemIDs(ctx, db, id)
	if err != nil {
		return nil, err
	}

	sm.IDsItemEstoque = idsItemEstoque
	return &sm, nil
}

func (r stockMovementRepository) CreateStockMovement(
	ctx context.Context,
	company entities.CompanyDatabaseConfig,
	tipoMovimentacao string,
	idUser int64,
	idsItemEstoque []int64,
) (*entities.StockMovement, error) {
	db := r.conn(company)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}

	rollback := func(innerErr error) (*entities.StockMovement, error) {
		_ = tx.Rollback()
		return nil, innerErr
	}

	for _, idItem := range idsItemEstoque {
		existe, err := itemExists(ctx, tx, idItem)
		if err != nil {
			return rollback(err)
		}
		if !existe {
			return rollback(fmt.Errorf("item de estoque %d não encontrado", idItem))
		}

		ultimoTipo, existeMovimentoAnterior, err := getLastMovementTypeByItem(ctx, tx, idItem)
		if err != nil {
			return rollback(err)
		}

		switch tipoMovimentacao {
		case "SAIDA":
			if !existeMovimentoAnterior || ultimoTipo != "ENTRADA" {
				return rollback(fmt.Errorf("saldo insuficiente para saída do item %d", idItem))
			}
		case "ENTRADA":
			if existeMovimentoAnterior && ultimoTipo == "ENTRADA" {
				return rollback(fmt.Errorf("item %d já está em estoque", idItem))
			}
		default:
			return rollback(errors.New("tipo_movimentacao inválido"))
		}
	}

	query := `
	INSERT INTO movimentacao (tipo_movimentacao, data, id_user) VALUES (?, ?, ?)
	`

	dataMovimento := time.Now()
	res, err := tx.ExecContext(
		ctx,
		query,
		tipoMovimentacao,
		dataMovimento,
		idUser,
	)
	if err != nil {
		return rollback(err)
	}

	idMovimentacao, err := res.LastInsertId()
	if err != nil {
		return rollback(err)
	}

	queryHistoric := `
	INSERT INTO movimentacao_historico (id_movimentacao, id_item_estoque) VALUES (?, ?)
	`

	for _, idItem := range idsItemEstoque {
		if _, err := tx.ExecContext(
			ctx,
			queryHistoric,
			idMovimentacao,
			idItem,
		); err != nil {
			return rollback(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entities.StockMovement{
		ID:               idMovimentacao,
		TipoMovimentacao: tipoMovimentacao,
		Data:             dataMovimento,
		IDUser:           idUser,
		Quantidade:       len(idsItemEstoque),
		IDsItemEstoque:   idsItemEstoque,
	}, nil
}

func (r stockMovementRepository) CreateItemEstoque(
	ctx context.Context,
	company entities.CompanyDatabaseConfig,
	idEquipamento int64,
	statusCode string,
	codigo string,
) (*entities.ItemEstoque, error) {
	db := r.conn(company)

	query := `SELECT 1 FROM equipamento WHERE id = ? LIMIT 1`

	var existe int
	err := db.QueryRowContext(
		ctx,
		query,
		idEquipamento,
	).Scan(&existe)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("equipamento %d não encontrado", idEquipamento)
		}
		return nil, err
	}

	queryItemEstoque := `INSERT INTO item_estoque (id_equipamento, status_code, codigo) VALUES (?, ?, ?)`

	res, err := db.ExecContext(
		ctx,
		queryItemEstoque,
		idEquipamento,
		statusCode,
		codigo,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entities.ItemEstoque{
		ID:            id,
		IDEquipamento: idEquipamento,
		StatusCode:    statusCode,
		Codigo:        codigo,
	}, nil
}

func (r stockMovementRepository) RankCategoriesByMovement(
	ctx context.Context,
	company entities.CompanyDatabaseConfig,
) ([]entities.CategoriaMovimentacaoRanking, error) {
	db := r.conn(company)

	query := `
		SELECT
			c.id,
			c.nome,
			COUNT(mh.id) AS total_movimentacao
		FROM movimentacao_historico mh
		INNER JOIN item_estoque ie ON ie.id = mh.id_item_estoque
		INNER JOIN equipamento eq ON eq.id = ie.id_equipamento
		INNER JOIN sub_categoria sc ON sc.id = eq.id_subCategoria
		INNER JOIN categoria c ON c.id = sc.id_categoria
		GROUP BY c.id, c.nome
		ORDER BY total_movimentacao DESC, c.nome ASC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.CategoriaMovimentacaoRanking
	pos := 1
	for rows.Next() {
		var row entities.CategoriaMovimentacaoRanking
		if err := rows.Scan(&row.ID, &row.Nome, &row.TotalMovimentacao); err != nil {
			return nil, err
		}
		row.Posicao = pos
		pos++
		list = append(list, row)
	}

	return list, rows.Err()
}

func itemExists(ctx context.Context, tx *sql.Tx, idItem int64) (bool, error) {
	var existe int

	query := `SELECT 1 FROM item_estoque WHERE id = ? LIMIT 1`

	err := tx.QueryRowContext(
		ctx,
		query,
		idItem,
	).Scan(&existe)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getLastMovementTypeByItem(ctx context.Context, tx *sql.Tx, idItem int64) (string, bool, error) {
	var tipo string
	err := tx.QueryRowContext(
		ctx,
		`
		SELECT m.tipo_movimentacao
		FROM movimentacao_historico mh
		INNER JOIN movimentacao m ON m.id = mh.id_movimentacao
		WHERE mh.id_item_estoque = ?
		ORDER BY m.data DESC, m.id DESC
		LIMIT 1
		`,
		idItem,
	).Scan(&tipo)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return tipo, true, nil
}

func (r stockMovementRepository) getMovementItemIDs(ctx context.Context, db *sql.DB, idMovimentacao int64) ([]int64, error) {
	rows, err := db.QueryContext(
		ctx,
		`
		SELECT id_item_estoque
		FROM movimentacao_historico
		WHERE id_movimentacao = ?
		ORDER BY id
		`,
		idMovimentacao,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var idItem int64
		if err := rows.Scan(&idItem); err != nil {
			return nil, err
		}
		ids = append(ids, idItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
