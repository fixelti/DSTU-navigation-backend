package drawPath

import (
	"fmt"
	"navigation/internal/appError"
	"navigation/internal/logging"
	"navigation/internal/models"
)

const (
	AxisX = 1
	AxisY = 2

	WidhtX  = 130
	HeightX = 30

	WidhtY  = 30
	HeightY = 130

	plus  = 0
	minus = 1
)

type Path struct {
	AudienceCoordinates models.Coordinates
	AudienceBorderPoint models.Coordinates
	SectorBorderPoint   models.Coordinates
	SectorNumber        int
	AudienceNumber      string
	Path                []models.Coordinates
	Repository          Repository
}

func NewPath(
	audienceCoordinates,
	audienceBorderPoint,
	sectorBorderPoint models.Coordinates,
	sectorNumber int,
	audienceNumber string,
	repository Repository) *Path {
	return &Path{
		AudienceCoordinates: audienceCoordinates,
		AudienceBorderPoint: audienceBorderPoint,
		SectorBorderPoint:   sectorBorderPoint,
		SectorNumber:        sectorNumber,
		AudienceNumber:      audienceNumber,
		Repository:          repository,
	}
}

var (
	User000004 = appError.NewError("drawPath", "GetSelector", "Input does not match desired length", "-", "US-000004")
)

func (d *Path) DrawInitPath() error {

	err := d.drawPathAuditory()
	if err != nil {
		return err
	}

	err = d.drawPathSector()
	if err != nil {
		return err
	}

	return nil
}

func (d *Path) drawPathAuditory() error {
	var err error
	axis := d.defenitionAxis(d.AudienceBorderPoint.Widht, d.AudienceBorderPoint.Height)

	switch axis {

	case AxisX:
		fmt.Println("Work AxisX")
		err := d.drawAudX()
		if err != nil {
			logging.GetLogger().Errorln("DrawPathAuditory case AxisX. Error - ", err)
			return err
		}

	case AxisY:
		fmt.Println("Work AxisY")
		err := d.drawAudY()
		if err != nil {
			logging.GetLogger().Errorln("DrawPathAuditory case AxisY. Error - ", err.Error())
			return err
		}

	default:
		logging.GetLogger().Errorln("DrawPathAuditory case default. Error - ", err)
		err = User000004
	}

	return err
}

func (d *Path) drawPathSector() error {
	iterator := 0
	axis := d.defenitionAxis(d.SectorBorderPoint.Widht, d.SectorBorderPoint.Height)
	boolean := true

	for boolean {
		if d.checkPath2Sector(d.Path[iterator], axis) {
			points := d.getDrawPoints2Sector(d.Path[iterator], axis)

			d.Path = append(d.Path, points)
			boolean = false
		} else {
			// определяем в каком направлении рисовать
			points := d.getDrawPoints(d.Path[iterator], axis)
			if points == (models.Coordinates{}) {
				return User000004
			}

			ok, err := d.Repository.checkBorderAud(points)
			if err != nil {
				return User000004
			}

			ok2, err := d.Repository.checkBorderSector(points)
			if err != nil {
				return User000004
			}

			if !ok && !ok2 {
				//TODO написать изменения направления или типо что-то такого
			}

			d.Path = append(d.Path, points)
		}

		iterator += 1
	}

	return nil
}

func (d *Path) getDrawPoints(path models.Coordinates, axis int) models.Coordinates {

	switch axis {
	case AxisX:
		points := models.Coordinates{
			X: (d.Path[0].X + d.Path[0].Widht),
			Y: (path.Y + path.Height)}
		sectorPoints := (d.SectorBorderPoint.Y + (d.SectorBorderPoint.Height + d.SectorBorderPoint.Y)) / 2
		if sectorPoints > path.X {
			points.Widht = WidhtY
			points.Height = HeightY
			return points
		} else {
			points.Widht = -WidhtY
			points.Height = -HeightY
			return points
		}
	case AxisY:
		points := models.Coordinates{
			X: (path.X + path.Widht),
			Y: (d.Path[0].Y + d.Path[0].Height)}
		sectorPoints := (d.SectorBorderPoint.X + (d.SectorBorderPoint.Widht + d.SectorBorderPoint.X)) / 2
		if sectorPoints > path.X {
			points.Widht = WidhtX
			points.Height = HeightX
			return points
		} else {
			points.Widht = -WidhtX
			points.Height = -HeightX
			return points
		}
	default:
		return models.Coordinates{}
	}
}

func (d *Path) getDrawPoints2Sector(path models.Coordinates, axis int) models.Coordinates {
	points := models.Coordinates{
		X: (path.X + path.Widht),
		Y: (path.Y + path.Height)}

	switch axis {
	case AxisX:
		fmt.Println("Work 21")
		sectorPoints := (d.SectorBorderPoint.X + (d.SectorBorderPoint.Widht + d.SectorBorderPoint.X)) / 2
		if sectorPoints > path.X {
			points.Widht = d.SectorBorderPoint.X - (path.X + path.Widht)
			points.Height = HeightX
			return points
		} else {
			points.Widht = -d.SectorBorderPoint.X - (path.X + path.Widht)
			points.Height = -HeightX
			return points
		}
	case AxisY:
		//TODO тут разобраться
		fmt.Println("Work 22")
		sectorPoints := (d.SectorBorderPoint.Y + (d.SectorBorderPoint.Height + d.SectorBorderPoint.Y)) / 2
		if sectorPoints < path.Y {
			points.Widht = WidhtY
			points.Height = d.SectorBorderPoint.Y - (path.Y + path.Height)
			return points
		} else {
			points.Widht = -WidhtY
			points.Height = -d.SectorBorderPoint.Y - (path.Y + path.Height)
			return points
		}
	default:
		return models.Coordinates{}
	}
}

func (d *Path) checkPath2Sector(path models.Coordinates, axis int) bool {
	switch axis {
	case AxisX:
		ph := path.Y + path.Height
		y1 := d.SectorBorderPoint.Y
		y2 := d.SectorBorderPoint.Y + d.SectorBorderPoint.Height
		if y1 <= ph && ph <= y2 {
			return true
		} else {
			return false
		}
	case AxisY:
		ph := path.X + path.Widht
		x1 := d.SectorBorderPoint.X
		x2 := d.SectorBorderPoint.X + d.SectorBorderPoint.Widht
		fmt.Println(ph, x1, x2)
		if x1 <= ph && ph <= x2 {
			fmt.Println("Work 1")
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func (d *Path) defenitionAxis(width, height int) int {
	if width == 1 {
		return AxisX
	} else if height == 1 {
		return AxisY
	} else {
		return 0
	}
}
