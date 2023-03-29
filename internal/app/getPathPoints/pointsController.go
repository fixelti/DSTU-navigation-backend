package getPathPoints

import (
	"log"
	"navigation/internal/logging"
	"navigation/internal/models"
)

const (
	AxisX = 1 // указывает на ось x.
	AxisY = 2 // указывает на ось y.

	WidhtX  = 40 // ширина на оси x.
	HeightX = 20 // высота на оси x.

	WidhtY  = 20 // ширина на оси y.
	HeightY = 40 // высота на оси y.

	plus  = 0 // значение будет положительным.
	minus = 1 // значение будет отрицательным.
)

// эти константы будут использовать для рассчета данных. Они буду передаваться в switch.
const (
	audStartPoints = 1 // для начального пути от границ аудитории.
	auditory2Sector = 2
	path2Sector     = 3
	sector2Sector   = 4
)

type pointsController struct {
	logger     *logging.Logger // логирования.
	repository Repository      // для обращения к базе данных.

	audStart string // номер начальной аудитории.
	audEnd   string // номер конечной аудитории.
	sectors  []int  // массив номеров секторов
}

func NewPointsController(audStart, audEnd string, sectors[]int, logger *logging.Logger, repository Repository) *pointsController {
	return &pointsController{
		logger:     logger,
		repository: repository,
		audStart: audStart,
		audEnd: audEnd,
		sectors: sectors,
	}
}

func (p *pointsController) getPathPoints() ([]models.Coordinates, error) {
	/* находим минимальное значение между номерами двух секторов.
	   необходимо для внутренней логики.
	*/
	entry, exit, err := min(p.sectors[0], p.sectors[1])
	if err != nil {
		return nil, err
	}

	// получаем новый объекта типа 'data'. С данными этого типа будет происходить вся работа.
	data, err := newData(p.audStart, entry, exit, p.sectors[1], p.logger, p.repository)
	if err != nil {
		return nil, err
	}

	// построение начального пути. От границы аудитории.
	err = data.setAudStartPoints()
	if err != nil {
		return nil, err
	}

	// borderSector, err := p.repository.getSectorBorderPoint(entry, exit)
	// if err != nil {
	// 	return nil, err
	// }

	// // построение пути вплоть до вхождение в область точек сектора.
	// err = data.middlePoints(borderSector)
	// if err != nil {
	// 	return nil, err
	// }

	// for i := 1; i < len(p.sectors)-1; i++ {

	// 	entry, exit, err := min(p.sectors[i], p.sectors[i+1])
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	borderSector, err := p.repository.getSectorBorderPoint(entry, exit)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	err = data.sector2Sector(borderSector)
	// 	if err != nil {
	// 		return nil,err
	// 	}
	// }

	// entry, exit, err = min(p.sectors[len(p.sectors)-1], p.sectors[len(p.sectors)-2])
	// if err != nil {
	// 	return nil,  err
	// }

	// // получаем новый объекта типа 'data'. С данными этого типа будет происходить вся работа.
	// dataEnd, err := newData(p.audStart, entry, exit, p.sectors[1], p.logger, p.repository)
	// if err != nil {
	// 	return nil, err
	// }

	// err = dataEnd.setAudStartPoints()
	// if err != nil {
	// 	return nil, err
	// }

	// borderSector, err = p.repository.getSectorBorderPoint(entry, exit)
	// if err != nil {
	// 	return nil, err
	// }

	// // построение пути вплоть до вхождение в область точек сектора.
	// err = dataEnd.middlePoints(borderSector)
	// if err != nil {
	// 	return nil, err
	// }

	// data.points = append(data.points, dataEnd.points...)

	log.Println("data - ", data.points)
	return data.points, nil
}

// TODO: добавить ошибки, если переменные одинаковые.
func min(a, b int) (int, int, error) {
	if a < b {
		return a, b, nil
	} else if a > b {
		return b, a, nil
	}
	return a, b, nil
}
