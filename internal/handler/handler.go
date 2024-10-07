package handler

import (
	"net/http"

	"github.com/babakgh/ticket_allocation_coding_test/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
    service service.TicketService
}

func NewHandler(s service.TicketService) *Handler {
    return &Handler{service: s}
}

func (h *Handler) Register(e *echo.Echo) {
    e.POST("/ticket_options", h.CreateTicketOption)
    e.GET("/ticket_options/:id", h.GetTicketOption)
    e.POST("/ticket_options/:id/purchases", h.PurchaseTickets)
}

func (h *Handler) CreateTicketOption(c echo.Context) error {
    var req service.TicketOptionRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    res, err := h.service.CreateTicketOption(c.Request().Context(), req)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetTicketOption(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket option ID")
    }

    res, err := h.service.GetTicketOption(c.Request().Context(), id)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Ticket option not found")
    }

    return c.JSON(http.StatusOK, res)
}

func (h *Handler) PurchaseTickets(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket option ID")
    }

    var req service.PurchaseRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    err = h.service.PurchaseTickets(c.Request().Context(), id, req)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.NoContent(http.StatusCreated)
}
