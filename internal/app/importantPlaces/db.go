package importantPlaces

import (
	"context"
	"errors"
	"fmt"
	"navigation/internal/appError"
	"navigation/internal/database/client/postgresql"
	"navigation/internal/logging"
	"navigation/internal/models"

	"github.com/jackc/pgconn"
)

const (
	file = "db.go"

	createFunction 	= "create"
	readFunction 	= "read"
	updateFunction 	= "update"
	deleteFunction 	= "delete"
	listFunction 	= "list"
)

var (
	txError 	= appError.NewAppError("can't start transaction")
	queryError 	= appError.NewAppError("failed to complite the request")
	scanError  	= appError.NewAppError("can't scan database response")
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func newRepository(
	client postgresql.Client, 
	logger *logging.Logger,
) Repository {
	return &repository {
		client: client,
		logger: logger,
	}
}

func (r *repository) Create(places models.ImportantPlaces) (models.ImportantPlaces, appError.AppError) {
	var newImportantPlaces models.ImportantPlaces
	req := `INSERT INTO important_places(name, id_auditorium) VALUES ($1, $2) RETURNING id;`

	tx, err := r.client.Begin(context.Background())
	if err != nil {
		_ = tx.Rollback(context.Background())
		txError.Wrap(fmt.Sprintf("file: %s, function: %s", file, createFunction))
		txError.Err = err
		return models.ImportantPlaces{}, *txError
	}

	err = tx.QueryRow(
		context.Background(),
		req,
		places.Name,
		places.AuditoryID).Scan(&newImportantPlaces.ID)

	if err != nil {
		_ = tx.Rollback(context.Background())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			queryError.Wrap(fmt.Sprintf("file: %s, function: %s", file, createFunction))
			queryError.Err = pgErr
			return models.ImportantPlaces{}, *queryError
		}
		queryError.Wrap(fmt.Sprintf("file: %s, function: %s", file, createFunction))
		queryError.Err = err
		return models.ImportantPlaces{}, *queryError
	}
	_ = tx.Commit(context.Background())
	return newImportantPlaces, appError.AppError{}	
}

func (r *repository) Read(id int) (models.ImportantPlaces, appError.AppError) {
	var importantPlaces models.ImportantPlaces
	request :=
	`SELECT *
	FROM important_places 
	JOIN auditorium 
	WHERE id = $1;`

	tx, err := r.client.Begin(context.Background())
	if err != nil {
		_ = tx.Rollback(context.Background())
		txError.Wrap(fmt.Sprintf("file: %s, function: %s", file, readFunction))
		txError.Err = err
		return models.ImportantPlaces{}, *txError
	}

	err = tx.QueryRow(
		context.Background(),
		request,
		id).Scan(
		&importantPlaces.ID,
		&importantPlaces.Name,
		&importantPlaces.AuditoryID)

	if err != nil {
		_ = tx.Rollback(context.Background())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			queryError.Wrap(fmt.Sprintf("file: %s, function: %s", file, readFunction))
			queryError.Err = pgErr
			return models.ImportantPlaces{}, *queryError
		}
		queryError.Wrap(fmt.Sprintf("file: %s, function: %s", file, readFunction))
		queryError.Err = err
		return models.ImportantPlaces{}, *queryError
	}
	_ = tx.Commit(context.Background())
	return importantPlaces, appError.AppError{}
}

func (r *repository) Update(oldpPlaces models.ImportantPlaces, newPlaces models.ImportantPlaces) (models.ImportantPlaces, appError.AppError) {}

func (r *repository) Delete(id int) (appError.AppError) {}

func (r *repository) List(numberBuild models.ImportantPlaces) ([]models.ImportantPlaces, appError.AppError) {}