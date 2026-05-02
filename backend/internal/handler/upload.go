package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/service"
)

// Предоставляет HTTP-эндпоинт POST /upload.
type UploadHandler struct {
	svc service.UploadService
}

// Создает новый экземпляр UploadHandler.
func NewUploadHandler(svc service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		JSONError(c, http.StatusBadRequest, "file is required")
		return
	}
	if fh.Size > 5*1024*1024 {
		JSONError(c, http.StatusBadRequest, "file too large")
		return
	}

	f, err := fh.Open()
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid file")
		return
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid file")
		return
	}

	out, err := h.svc.UploadAndMatch(c.Request.Context(), b)
	if err != nil {
		RespondError(c, err)
		return
	}

	JSONOK(c, out)
}

