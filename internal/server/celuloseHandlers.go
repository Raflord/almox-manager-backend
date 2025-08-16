package server

import (
	"context"
	"time"

	"almox-manager-backend/internal/database"
	"almox-manager-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
)

type LoadBody struct {
	Id            string `json:"id"`
	Material      string `json:"material"`
	AverageWeight int32  `json:"averageWeight"`
	Unit          string `json:"unit"`
	CreatedAt     string `json:"createdAt,omitempty"`
	Timezone      string `json:"timezone"`
	Operator      string `json:"operator"`
	Shift         string `json:"shift"`
}

type LoadFilteredBody struct {
	Material    string `json:"material"`
	FirstDate   string `json:"firstDate"`
	SeccondDate string `json:"seccondDate"`
}

func (s *FiberServer) HandleGetLatest(c *fiber.Ctx) error {
	loads, err := s.db.GetLatest(context.Background())
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(loads)
}

func (s *FiberServer) HandleGetSummary(c *fiber.Ctx) error {
	now := pgtype.Timestamp{Time: time.Now(), Valid: true}

	loads, err := s.db.GetSummary(c.Context(), now)
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

	records, err := s.db.GetFiltered(c.Context(), database.GetFilteredParams{
		Material: body.Material,
		Column2:  body.FirstDate,
		Column3:  body.SeccondDate,
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

	err = s.db.CreateLoad(c.Context(), database.CreateLoadParams{
		ID:            ulid.Make().String(),
		Material:      body.Material,
		AverageWeight: body.AverageWeight,
		Unit:          body.Unit,
		CreatedAt:     pgtype.Timestamp{Time: parsedDateTime, Valid: true},
		Timezone:      pgtype.Text{String: body.Timezone, Valid: true},
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

	err = s.db.UpdateLoad(c.Context(), database.UpdateLoadParams{
		ID:        body.Id,
		Material:  body.Material,
		CreatedAt: pgtype.Timestamp{Time: parsedDateTime, Valid: true},
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
	err := s.db.DeleteLoad(c.Context(), c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
