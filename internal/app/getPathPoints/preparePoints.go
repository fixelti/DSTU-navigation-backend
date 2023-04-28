package getPathPoints

import (
	"strconv"

	"navigation/internal/models"
)

// подготовка данных
func (d *data) preparePoints(pointsType, axis int, borderPoint, points models.Coordinates) models.Coordinates {
	switch pointsType {
	// начальных путь от границ аудитории.
	case audStartPoints:
		var coordinates models.Coordinates

		XX := borderPoint.X + 1
		YX := (borderPoint.Y + (borderPoint.Height + borderPoint.Y)) / 2
		// YX := borderPoint.Y + 1

		XY := (borderPoint.X + (borderPoint.Widht + borderPoint.X)) / 2
		// XY := borderPoint.X + 1
		YY := borderPoint.Y + 1

		if axis == AxisX {
			coordinates.X = XX
			coordinates.Y = YX
			coordinates.Widht = WidhtX
			coordinates.Height = HeightX

		} else if axis == AxisY {
			coordinates.X = XY
			coordinates.Y = YY
			coordinates.Widht = WidhtY
			coordinates.Height = HeightY
		}

		return coordinates

	// от начального пути до вхождения в сектор
	case auditory2Sector:
		//TODO: возможно, надо будет и тут поменять, как это сделано в блоке else.
		if axis == AxisX {
			if borderPoint.X < points.X {
				return models.Coordinates{
					X:      points.X + points.Widht,
					Y:      d.points[0].Y + d.points[0].Height,
					Widht:  -WidhtX,
					Height: HeightX}
			} else {
				if d.points[0].Height == HeightY || d.points[0].Widht == WidhtX {
					return models.Coordinates{
						X:      points.X + points.Widht,
						Y:      d.points[0].Y + d.points[0].Height - HeightX,
						Widht:  WidhtX,
						Height: HeightX}
				} else {
					// fmt.Println("This work")
					return models.Coordinates{
						X:      points.X + points.Widht,
						Y:      d.points[0].Y + d.points[0].Height,
						Widht:  WidhtX,
						Height: HeightX}
				}
			}
		} else {
			if len(strconv.Itoa(d.sectorNumber)) == stairs {
				return models.Coordinates{
					X:      d.points[0].X,
					Y:      points.Y + points.Height,
					Widht:  WidhtY,
					Height: HeightY}
			} else {
				return models.Coordinates{
					X:      points.X + points.Widht, // менял, чтобы 330 не съезжала
					Y:      points.Y + points.Height,
					Widht:  WidhtY,
					Height: borderPoint.Y - (points.Y + points.Height)}
			}
		}

	// соеденить путь и сектор
	case path2Sector:
		if axis == AxisX {
			if len(strconv.Itoa(d.sectorNumber)) == stairs {
				return models.Coordinates{
					X:      d.points[0].X,
					Y:      points.Y + points.Height,
					Widht:  borderPoint.X - points.X,
					Height: HeightX}
			} else {
				// d.logger.Logger.Infoln("path to sector: ")
				if borderPoint.X < points.X {
					// d.logger.Logger.Infoln("border point < points: ")
					return models.Coordinates{
						X:      points.X + points.Widht,
						Y:      points.Y + points.Height,
						// Y: d.points[0].Y + d.points[0].Height,
						Widht:  borderPoint.X - points.X,
						Height: HeightX}
				} else {
					return models.Coordinates{
						X:      points.X + points.Widht,
						// Y:      d.points[0].Y + d.points[0].Height,
						Y:      points.Y + points.Height,
						Widht:  borderPoint.X - points.X,
						Height: HeightX}
				}
			}
		} else {
			if borderPoint.Y > points.Y {
				return models.Coordinates{
					X:      points.X + points.Widht,
					Y:      points.Y + points.Height - HeightX,
					Widht:  WidhtY,
					Height: borderPoint.Y - points.Y}
			} else {
				return models.Coordinates{
					X:      points.X + points.Widht,
					Y:      points.Y + points.Height,
					Widht:  WidhtY,
					Height: borderPoint.Y - (points.Y + points.Height)}
			}
		}

	// путь, который прокладывается между секторами
	case sector2Sector:
		if axis == AxisX {
			return models.Coordinates{
				X: points.X + points.Widht,
				Y: points.Y + points.Height,
				Widht: borderPoint.X - (points.X + points.Widht),
				Height: HeightX,
			}

		} else {
			return models.Coordinates{
				X: points.X,
				Y: points.Y + points.Height,
				Widht: WidhtY,
				Height: borderPoint.Y - (points.Y + points.Height),
			}
		}
		
	default:
		return models.Coordinates{}
	}

}
