package server

import (
	"almox-manager-backend/internal/database"

	"github.com/gofiber/fiber/v2"
)

type RecordInputBody struct {
	Material      string `json:"material"`
	AverageWeight int    `json:"average_weight"`
	Unit          string `json:"unit"`
	Operator      string `json:"operator"`
	Shift         string `json:"shift"`
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
