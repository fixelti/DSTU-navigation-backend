package getPathPoints

import (
	"strconv"

	"navigation/internal/appError"
	"navigation/internal/database/client/postgresql"
	"navigation/internal/logging"
	"navigation/internal/models"
)

var (
	stairs = 4
)

var (
	TransitionError = appError.NewAppError("not info for transition")
)

type data struct {
	audPoints          models.Coordinates // координаты аудитории.
	audBorderPoints    models.Coordinates // координаты места отрисовки пути (одна из границ аудитории).
	sectorBorderPoints models.Coordinates // координаты одной из границ сектора.

	sectorNumber     int    // номер сектора.
	nextSectorNumber int    // номер следующего сектора.
	audNumber        string // номер аудитории.
	endPoints string

	logger     *logging.Logger // логирования.
	client postgresql.Client      // для обращения к базе данных.

	points []models.Coordinates // массив координат. Для построения пути.

	transition       int
	transitionNumber int

	sectorType int
}

func newData(audNumber string, endPoints string, 
	sectorEntry, sectorExit, nextSectorNumber int, 
	logger *logging.Logger, 
	client postgresql.Client,
	transition, transitionNumber int) (*data, appError.AppError) {
	var err appError.AppError
	data := &data{
		audNumber:        audNumber,
		endPoints: endPoints,
		sectorNumber: sectorExit, //TODO: тут может быть ошибка. Может пердаваться не верный сектор. 
		nextSectorNumber: nextSectorNumber,
		logger:           logger,
		client:       client,
		transition: transition,
		transitionNumber: transitionNumber,
	}

	if data.transition == stair {
		sectorEntry, sectorExit = sectorExit, sectorEntry
	}

	err = data.getPoints(sectorEntry, sectorExit)
	if err.Err != nil {
		err.Wrap("newData")
		return nil, err
	}

	return data, err
}

// получение audPoints, audBorderPoints, sectorBorderPoints
func (d *data) getPoints(entry, exit int) appError.AppError {
	var err appError.AppError
	repository := NewRepository(d.client, d.logger)
	// получаем координаты аудитории по ее номеру.
	if d.transition == stair { // Если лестнциа
		d.audBorderPoints, err = repository.getAudBorderPoint(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.audPoints, err = repository.getTransitionPoints(d.transitionNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

	} else if d.transition == noTransition { // Если перехода между этажами нет
		d.audPoints, err = repository.getAudPoints(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.audBorderPoints, err = repository.getAudBorderPoint(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.sectorBorderPoints, err = repository.getSectorBorderPoint(entry, exit)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}
 
	} else if d.transition == elevator { // Если лифт
		d.audPoints, err = repository.getAudPoints(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.audBorderPoints, err = repository.getAudBorderPoint(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		if len(strconv.Itoa(exit)) == 4 {
			d.audBorderPoints, err = repository.getTransitionSectorBorderPoint(d.transitionNumber)
			if err.Err != nil {
				err.Wrap("getPoints")
				return err
			}
		} else {
			d.sectorBorderPoints, err = repository.getSectorBorderPoint(entry, exit)
			if err.Err != nil {
				err.Wrap("getPoints")
				return err
			}
		}
		
	} else if d.transition == transitionToAud { // Если от перехоного сектора до конечной аудитории
		d.audPoints, err = repository.getTransitionPoints(d.transitionNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.audBorderPoints, err = repository.getTransitionSectorBorderPoint(d.transitionNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.sectorBorderPoints, err = repository.getSectorBorderPoint(entry, exit)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

	} else if d.transition == aud2Aud { // Если аудитории находятся в одном секторе

		d.audPoints, err = repository.getAudPoints(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.audBorderPoints, err = repository.getAudBorderPoint(d.audNumber)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}

		d.sectorBorderPoints, err = repository.getAudBorderPoint(d.endPoints)
		if err.Err != nil {
			err.Wrap("getPoints")
			return err
		}
	} else {
		err.Err = TransitionError
		err.Wrap("getPoints")
	}

	return err
}
