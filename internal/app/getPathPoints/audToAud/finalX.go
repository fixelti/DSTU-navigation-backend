package audToAud

import (
	"navigation/internal/models"
)

func (m *middleController) upX(borderPoint, points models.Coordinates) models.Coordinates {
	var result models.Coordinates
	if borderPoint.X < points.X {
		result = models.Coordinates{
			X: points.X + points.Widht,
            Y: points.Y + points.Height,
            Widht: borderPoint.X - (points.X + points.Widht),
            Height: 5,
		}
	} else if borderPoint.X > points.X {
		result = models.Coordinates{
			X: points.X,
            Y: points.Y + points.Height,
            Widht: borderPoint.X - points.X,
            Height: 5,
		}
	}

	return result
}

func (m *middleController) downX(borderPoint, points models.Coordinates) models.Coordinates {
	var result models.Coordinates
	if borderPoint.X < points.X {
		result = models.Coordinates{
			X: points.X + points.Widht,
            Y: points.Y + points.Height,
            Widht: borderPoint.X - (points.X + points.Widht),
            Height: -5,
		}
	} else if borderPoint.X > points.X {
		result = models.Coordinates{
			X: points.X,
            Y: points.Y + points.Height,
            Widht: borderPoint.X - points.X,
            Height: -5,
		}
	}

	return result
}

func (m *middleController) leftAndRightX(borderPoint, points models.Coordinates) models.Coordinates {
	var result models.Coordinates
	// var lenght = len(m.Points)

	// m.Points[lenght-1].Y = borderPoint.Y + 10
	// points.Y = borderPoint.Y + 10
	
	result = models.Coordinates{
		X: points.X + points.Widht,
		Y: points.Y,
		Widht: borderPoint.X - (points.X + points.Widht),
		Height: 5,
	}

	return result
}