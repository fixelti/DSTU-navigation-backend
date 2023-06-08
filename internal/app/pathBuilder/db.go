package pathBuilder

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"

	"navigation/internal/appError"
	"navigation/internal/database/client/postgresql"
	"navigation/internal/logging"
	"navigation/internal/models"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

// Получение ссылок секторов
func (r *repository) GetSectorLink() ([]models.SectorLink, appError.AppError) {
	txError    := appError.NewAppError("can't start transaction")
	queryError := appError.NewAppError("failed to complete the request")
	scanError  := appError.NewAppError("can't scan database response")

	var sectorLink []models.SectorLink
	req := `SELECT number_sector, link FROM sector_link;`
	tx, err := r.client.Begin(context.Background())
	if err != nil {
		_ = tx.Rollback(context.Background())
		txError.Wrap("GetSectorLink")
		txError.Err = err
		return nil, *txError
	}

	rows, err := tx.Query(context.Background(), req)
	if err != nil {
		_ = tx.Rollback(context.Background())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			queryError.Wrap("GetSectorLink")
			queryError.Err = pgErr
			return nil, *queryError
		}
		queryError.Wrap("GetSectorLink")
		queryError.Err = err
		return nil, *queryError
	}

	for rows.Next() {
		var sl models.SectorLink
		err := rows.Scan(&sl.NumberSector, &sl.NumberLink)
		if err != nil {
			scanError.Wrap("GetSectorLink")
			scanError.Err = err
			return nil, *scanError
		}
		sectorLink = append(sectorLink, sl)
	}

	_ = tx.Commit(context.Background())
	return sectorLink, appError.AppError{}
}

// получение номера сеткора
func (r *repository) GetSector(number string, building uint) (int, appError.AppError) {
	var sector models.Sector
	txError    := appError.NewAppError("can't start transaction")
	queryError := appError.NewAppError("failed to complete the request")
	// // scanError  := appError.NewAppError("can't scan database response")

	req :=
		`SELECT 
	sector.number
	FROM auditorium
	JOIN sector ON auditorium.id_sector=sector.id
	JOIN floor ON sector.id_floor=floor.id
	JOIN building ON floor.id_building=building.id
	WHERE auditorium.number = $1 
	AND building.number = $2;`

	tx, err := r.client.Begin(context.Background())
	// fmt.Println("Work: ", tx, err, r.client)
	
	if err != nil {
		_ = tx.Rollback(context.Background())
		txError.Wrap("db GetSector")
		txError.Err = err
		return 0, *txError
	}
	err = tx.QueryRow(
		context.Background(),
		req,
		number,
		building).Scan(&sector.Number)

	if err != nil {
		_ = tx.Rollback(context.Background())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			queryError.Wrap("db GetSector")
			queryError.Err = pgErr
			return 0, *queryError
		} else {
			queryError.Wrap("db GetSector")
			queryError.Err = err
			return 0, *queryError
		}
	}
	_ = tx.Commit(context.Background())
	return sector.Number, appError.AppError{}
}

// получение номера переходного сектора
func (r *repository) GetTransitionSector(sectorNumber, type_transtion_sector int) (int, appError.AppError) {
	txError    := appError.NewAppError("can't start transaction")
	queryError := appError.NewAppError("failed to complete the request")
	// scanError  := appError.NewAppError("can't scan database response")
	var sector models.Sector 
	type_transtion_sector = 1 // вообще это из базы можно будет убрать
	req :=
		`SELECT transition.number 
	FROM transition
	JOIN sector ON sector.id_transition = transition.id
	WHERE sector.number = $1 
	AND type_transition = $2`

	tx, err := r.client.Begin(context.Background())
	if err != nil {
		_ = tx.Rollback(context.Background())
		txError.Wrap("GetTransitionSector")
		txError.Err = err
		return 0, *txError
	}

	err = tx.QueryRow(
		context.Background(),
		req,
		sectorNumber,
		type_transtion_sector).Scan(&sector.Number)

	if err != nil {
		_ = tx.Rollback(context.Background())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			queryError.Wrap("GetTransitionSector")
			queryError.Err = pgErr
			return 0, *queryError
		}
		queryError.Wrap("GetTransitionSector")
		queryError.Err = err
		return 0, *queryError
	}
	_ = tx.Commit(context.Background())
	return sector.Number, appError.AppError{}
}

// получение номера переходного сектора для лифта
func (r *repository) GetTransitionSector2(sectorNumber, type_transtion_sector int) (int, appError.AppError) {
	txError    := appError.NewAppError("can't start transaction")
	queryError := appError.NewAppError("failed to complete the request")
	// scanError  := appError.NewAppError("can't scan database response")
	var sector models.Sector
	req :=
		`SELECT sector.number
		FROM transition
		JOIN sector ON sector.id_transition = transition.id
		WHERE transition.number = $1
		AND type_transition = $2;`

	tx, err := r.client.Begin(context.Background())
	if err != nil {
		_ = tx.Rollback(context.Background())
		txError.Wrap("GetTransitionSector")
		txError.Err = err
		return 0, *txError
	}

	err = tx.QueryRow(
		context.Background(),
		req,
		sectorNumber,
		type_transtion_sector).Scan(&sector.Number)

	if err != nil {
		_ = tx.Rollback(context.Background())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			queryError.Wrap("GetTransitionSector")
			queryError.Err = pgErr
			return 0, *queryError
		}
		queryError.Wrap("GetTransitionSector")
		queryError.Err = err
		return 0, *queryError
	}
	_ = tx.Commit(context.Background())
	return sector.Number, appError.AppError{}
}