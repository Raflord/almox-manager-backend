package server

import (
	"context"
	"fmt"
	"time"

	"almox-manager-backend/internal/database"
	"almox-manager-backend/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
)

type LoadBody struct {
	Material      string `json:"material" validate:"required"`
	AverageWeight int32  `json:"averageWeight" validate:"required"`
	Unit          string `json:"unit" validate:"required"`
	CreatedAt     string `json:"createdAt" validate:"required"`
	Operator      string `json:"operator" validate:"required"`
	Shift         string `json:"shift" validate:"required"`
}

type LoadUpdateBody struct {
	Id        string `json:"id" validate:"required"`
	Material  string `json:"material" validate:"required"`
	CreatedAt string `json:"createdAt" validate:"required"`
	Operator  string `json:"operator" validate:"required"`
	Shift     string `json:"shift" validate:"required"`
}

type LoadFilteredBody struct {
	Material   string `json:"material" validate:"required"`
	FirstDate  string `json:"firstDate" validate:"required"`
	SecondDate string `json:"secondDate" validate:"required"`
}

func (s *FiberServer) HandleGetLatest(c *fiber.Ctx) error {
	loads, err := s.db.GetLatest(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Could not get latestest records",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(loads)
}

func (s *FiberServer) HandleGetSummary(c *fiber.Ctx) error {
	now := pgtype.Timestamp{Time: time.Now(), Valid: true}

	loads, err := s.db.GetSummary(c.Context(), now)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Could not get summary data",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(loads)
}

func (s *FiberServer) HandleGetById(c *fiber.Ctx) error {
	loads, err := s.db.GetById(context.Background(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Load does not exists",
		})
	}

	return c.Status(fiber.StatusOK).JSON(loads)
}

func (s *FiberServer) HandleGetFiltered(c *fiber.Ctx) error {
	var body LoadFilteredBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	records, err := s.db.GetFiltered(c.Context(), database.GetFilteredParams{
		Material: body.Material,
		Column2:  body.FirstDate,
		Column3:  body.SecondDate,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(records)
}

func (s *FiberServer) HandleCreateLoad(c *fiber.Ctx) error {
	var body LoadBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errs := err.(validator.ValidationErrors)

		var messages []string

		for _, e := range errs {
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", e.Field())

			}

			messages = append(messages, msg)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": messages,
		})
	}

	parsedDateTime, err := utils.ParseDateTime(body.CreatedAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid datetime format. Expected: 'yyyy-MM-dd HH:mm:ss'",
			"details": err.Error(),
		})
	}

	err = s.db.CreateLoad(c.Context(), database.CreateLoadParams{
		ID:            ulid.Make().String(),
		Material:      body.Material,
		AverageWeight: body.AverageWeight,
		Unit:          body.Unit,
		CreatedAt:     pgtype.Timestamp{Time: parsedDateTime, Valid: true},
		Operator:      body.Operator,
		Shift:         body.Shift,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create load",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (s *FiberServer) HandleUpdateLoad(c *fiber.Ctx) error {
	_, err := s.db.GetById(context.Background(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Load does not exists",
		})
	}

	var body LoadUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errs := err.(validator.ValidationErrors)

		var messages []string

		for _, e := range errs {
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", e.Field())

			}

			messages = append(messages, msg)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": messages,
		})
	}

	parsedDateTime, err := utils.ParseDateTime(body.CreatedAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid datetime format. Expected: ISO 8601 'yyyy-MM-ddTHH:mm:ss'",
			"details": err.Error(),
		})
	}

	err = s.db.UpdateLoad(c.Context(), database.UpdateLoadParams{
		ID:        body.Id,
		Material:  body.Material,
		CreatedAt: pgtype.Timestamp{Time: parsedDateTime, Valid: true},
		Operator:  body.Operator,
		Shift:     body.Shift,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update load",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *FiberServer) HandleDeleteLoad(c *fiber.Ctx) error {
	_, err := s.db.GetById(context.Background(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Load does not exists",
		})
	}

	err = s.db.DeleteLoad(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Could not delete load",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
