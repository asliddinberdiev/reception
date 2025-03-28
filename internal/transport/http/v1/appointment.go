package v1

import (
	"time"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) initAppointmentRoutes(router fiber.Router) {
	appointment := router.Group("appointments", h.middlewareUser)
	{
		appointment.Post("", h.appointmentCreate)
		appointment.Put("/:id", h.appointmentUpdateStatus)
	}
}

func (h *Handler) appointmentCreate(c *fiber.Ctx) error {
	userID, err := h.MwGetUserID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "not found user")
	}

	var req models.AppointmentCreateInput
	req.UserID = userID

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.valid.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	date, err := helper.StringToTime(req.AppointmentDate, "2006-01-02")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid date format")
	}

	startTime, err := helper.StringToTime(req.AppointmentTime, "15:04:05")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid time format")
	}

	endTime, err := helper.StringToTime(req.AppointmentTime, "15:04:05")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid time format")
	}

	hasTimeData := models.AppointmentRangeTime{
		DoctorID:  req.DoctorID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime.Add(30 * time.Minute),
	}

	hasEmptyTime, err := h.sve.Appointment.GetByRangeTime(c.Context(), hasTimeData)
	if err != nil {
		return err
	}

	if !hasEmptyTime {
		return fiber.NewError(fiber.StatusNotFound, "already taken")
	}

	res, err := h.sve.Appointment.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *Handler) appointmentUpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	}

	if err := h.sve.Appointment.UpdateStatus(
		c.Context(),
		models.CommonGetByID{ID: id},
	); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString("successfully update")

}
