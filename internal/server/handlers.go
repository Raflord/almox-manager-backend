package server

import (
	"almox-manager-backend/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RecordInputBody struct {
	Material      string `json:"material"`
	AverageWeight int    `json:"average_weight"`
	Unit          string `json:"unit"`
	Operator      string `json:"operator"`
	Shift         string `json:"shift"`
}

type FilteredInputBody struct {
	Material    string     `json:"material"`
	FirstDate   *time.Time `json:"first_date"`
	SeccondDate *time.Time `json:"seccond_date"`
}
type DayData struct {
	Material    string `json:"material"`
	TotalWeight int    `json:"total_weight"`
}

func (s *FiberServer) HandleGetLatestCell(c *fiber.Ctx) error {
	records, err := s.db.QueryLatestRecords()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(records)
}

func (s *FiberServer) HandleGetDayCell(c *fiber.Ctx) error {
	records, err := s.db.QueryDayRecords()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	sum := make(map[string]int)

	for _, item := range records {
		sum[item.Material] += item.AverageWeight
	}

	var daySum []DayData
	for material, total := range sum {
		daySum = append(daySum, DayData{
			Material:    material,
			TotalWeight: total,
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(daySum)
}

func (s *FiberServer) HandleGetFiltered(c *fiber.Ctx) error {
	var body FilteredInputBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	records, err := s.db.QueryFilteredRecords(database.FilteredInput{
		Material:    body.Material,
		FirstDate:   body.FirstDate,
		SeccondDate: body.SeccondDate,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(records)
}

func (s *FiberServer) HandlePostCellulose(c *fiber.Ctx) error {
	var body RecordInputBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	err = s.db.CreateNewRecord(database.RecordInput{
		Material:      body.Material,
		AverageWeight: body.AverageWeight,
		Unit:          body.Unit,
		Operator:      body.Operator,
		Shift:         body.Shift,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}
func (s *FiberServer) HandlePutCellulose(c *fiber.Ctx) error {
	return c.JSON("message: PutCellulose")
}
func (s *FiberServer) HandleDeleteCellulose(c *fiber.Ctx) error {
	return c.JSON("message: DeleteCellulose")
}
