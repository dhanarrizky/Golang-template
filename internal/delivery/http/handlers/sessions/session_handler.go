package sessions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"
)

type SessionHandler struct {
	usecase auth.SessionUsecase
}

func NewSessionHandler(usecase auth.SessionUsecase) *SessionHandler {
	return &SessionHandler{usecase: usecase}
}

// GET /sessions
func (h *SessionHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	currentFamily := c.GetString("family_id") // dari middleware JWT

	sessions, err := h.usecase.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to fetch sessions"})
		return
	}

	resp := make([]dto.SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		resp = append(resp, dto.SessionResponse{
			ID:        s.ID,
			Device:    s.Device,
			IP:        s.IP,
			LastUsed:  s.LastUsed,
			CreatedAt:s.CreatedAt,
			Current:   s.FamilyID == currentFamily,
		})
	}

	c.JSON(http.StatusOK, dto.ListSessionResponse{Sessions: resp})
}

// DELETE /sessions/:id
func (h *SessionHandler) Revoke(c *gin.Context) {
	userID := c.GetString("user_id")
	currentFamily := c.GetString("family_id")
	sessionID := c.Param("id")

	err := h.usecase.Revoke(c.Request.Context(), userID, sessionID, currentFamily)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RevokeSessionResponse{
		Message: "Session revoked successfully",
	})
}
