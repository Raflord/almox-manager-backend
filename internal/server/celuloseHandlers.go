package server

import (
	"almox-manager-backend/internal/types"
	"almox-manager-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type LoadBody struct {
	Id            string `json:"id"`
	Material      string `json:"material"`
	AverageWeight int    `json:"average_weight"`
	Unit          string `json:"unit"`
	CreatedAt     string `json:"created_at,omitempty"`
	Timezone      string `json:"timezone"`
	Operator      string `json:"operator"`
	Shift         string `json:"shift"`
}

type LoadFilteredBody struct {
	Material    string `json:"material"`
	FirstDate   string `json:"first_date"`
	SeccondDate string `json:"seccond_date"`
}

func (s *FiberServer) HandleGetLatest(c *fiber.Ctx) error {
	loads, err := s.db.QueryLatest()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(loads)
}

func (s *FiberServer) HandleGetDay(c *fiber.Ctx) error {
	loads, err := s.db.QueryDay()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(loads)
}

func (s *FiberServer) HandleFiltered(c *fiber.Ctx) error {
	var body LoadFilteredBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	records, err := s.db.QueryFiltered(types.LoadFiltered{
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

func (s *FiberServer) HandleCreate(c *fiber.Ctx) error {
	var body LoadBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	parsedDateTime, err := utils.ParseDateTime(body.CreatedAt)
	if err != nil {
		return err
	}

	err = s.db.CreateLoad(types.Load{
		Material:      body.Material,
		AverageWeight: body.AverageWeight,
		Unit:          body.Unit,
		CreatedAt:     parsedDateTime,
		Timezone:      body.Timezone,
		Operator:      body.Operator,
		Shift:         body.Shift,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (s *FiberServer) HandleUpdate(c *fiber.Ctx) error {
	var body LoadBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	parsedDateTime, err := utils.ParseDateTime(body.CreatedAt)
	if err != nil {
		return err
	}

	err = s.db.UpdateLoad(types.Load{
		Id:        body.Id,
		Material:  body.Material,
		CreatedAt: parsedDateTime,
		Operator:  body.Operator,
		Shift:     body.Shift,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *FiberServer) HandleDelete(c *fiber.Ctx) error {
	err := s.db.DeleteLoad(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
