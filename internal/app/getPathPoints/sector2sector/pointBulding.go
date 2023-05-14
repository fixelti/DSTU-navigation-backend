package sectorToSector

import (
	axes "navigation/internal/app/getPathPoints/axis"
	"navigation/internal/appError"
	"navigation/internal/models"
)

func (s *sectorToSectorController) building(iterator int, borderSector models.Coordinates) appError.AppError {
	boolean := true
	repository := NewRepository(s.client, s.logger)
	axis := axes.DefenitionAxis(borderSector.Widht, borderSector.Height, s.constData.axisX, s.constData.axisY)
	lastPathSector := false
	for boolean {
		if s.checkOccurrence(s.Points[iterator], axis, borderSector) {
			var points models.Coordinates
			s.pathAlignment(borderSector, axis)

			axis = axes.ChangeAxis(axis, s.constData.axisX, s.constData.axisY)

			if lastPathSector == true {
				points = s.preparation(axis, borderSector, s.Points[iterator], true)
			} else {
				points = s.preparation(axis, borderSector, s.Points[iterator], false)
			}

			s.Points = append(s.Points, points)
	
			s.OldAxis = axis
			boolean = false
		} else {
			lastPathSector = true
		
			points := s.preparation(axis, borderSector, s.Points[iterator], false)

			ok, err := repository.checkBorderAud(points)
			if err.Err != nil {
				err.Wrap("otherPathPoints")
				return err
			}
			ok2, err := repository.checkBorderSector(points)
			if err.Err != nil {
				err.Wrap("otherPathPoints")
				return err
			}

			if !ok && !ok2 {
				//TODO написать изменения направления или типо что-то такого
			}

			s.Points = append(s.Points, points)
		}

		iterator += 1
	}

	return appError.AppError{}
}

// точки от начала пути до вхождение в пределы сектора


// проверка на вхождение точек пути в пределы сектора.
func (s *sectorToSectorController) checkOccurrence(points models.Coordinates, axis int, borderSector models.Coordinates) bool {
	switch axis {
	case s.constData.axisX:
		ph := points.X + points.Widht
		x1 := borderSector.X
		x2 := borderSector.X + borderSector.Widht
		if x1 <= ph && ph <= x2 {
			return true
		} else {
			return false
		}
	case s.constData.axisY:
		ph := points.Y + points.Height
		y1 := borderSector.Y
		y2 := borderSector.Y + borderSector.Height
		if y1 <= ph && ph <= y2 {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

// выравнивание пути
func (s *sectorToSectorController) pathAlignment(sectorBorderPoint models.Coordinates, axis int) {
	lenght := len(s.Points)
	path := s.Points[lenght-1]
	switch axis {
	case s.constData.axisX:
		points := models.Coordinates{
			X: (path.X),
			Y: (path.Y)}
		sectorPoints := (sectorBorderPoint.X + (sectorBorderPoint.Widht + sectorBorderPoint.X)) / 2
		if sectorPoints > path.X {
			points.Widht = sectorPoints - path.X
			points.Height = s.constData.heightX
			s.Points[lenght-1].Widht = points.Widht
		} else if sectorPoints < path.X {
			points.Widht = (sectorBorderPoint.X + (sectorBorderPoint.Widht / 2)) - (path.X)
			points.Height = s.constData.heightX
			s.Points[lenght-1].Widht = points.Widht
		}
	case s.constData.axisY:
		points := models.Coordinates{
			X: (path.X),
			Y: (path.Y)}
		sectorPoints := (sectorBorderPoint.Y + (sectorBorderPoint.Height + sectorBorderPoint.Y)) / 2
		if sectorPoints > path.Y {
			points.Widht = s.constData.widhtY
			points.Height = sectorPoints - path.Y
			s.Points[lenght-1].Height = points.Height
		} else if sectorPoints < path.Y {
			points.Widht = s.constData.widhtY
			points.Height = sectorPoints - path.Y
			s.Points[lenght-1].Height = points.Height
		}
	default:
		s.logger.Errorln("Path Alignment default")
	}
}
