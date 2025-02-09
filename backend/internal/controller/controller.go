package controller

import (
	"backend/internal/db/queries"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type pingSchema struct {
	ContainerIp string `json:"container_ip" binding:"required"`
	PingTimeMs  *int32 `json:"ping_time_ms" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type ContainerController struct {
	queries *queries.Queries
	ctx     context.Context
}

func NewContainerController(queries *queries.Queries, ctx context.Context) *ContainerController {
	return &ContainerController{queries, ctx}
}

func (c *ContainerController) ListContainerStatuses(ctx *gin.Context) {
	statuses, err := c.queries.ListContainerStatuses(c.ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, statuses)
}

func (c *ContainerController) PingHandler(ctx *gin.Context) {
	var payload *pingSchema

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Printf(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := c.queries.GetContainerStatusByIp(c.ctx, payload.ContainerIp)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status.ContainerIp = payload.ContainerIp
	status.PingTimeMs = *payload.PingTimeMs
	status.Status = queries.ContainerStatus(payload.Status)
	if payload.Status == string(queries.ContainerStatusUP) {
		status.LastSuccessfulPing.Scan(time.Now().UTC())
	} else if payload.Status == string(queries.ContainerStatusDOWN) {
		status.PingTimeMs = -1
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	err = c.queries.UpsertContainerStatus(c.ctx, queries.UpsertContainerStatusParams(status))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
