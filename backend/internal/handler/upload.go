package handler

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
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

const (
	uploadMaxFileBytes = 5 * 1024 * 1024
	uploadMaxFiles     = 12
)

func (h *UploadHandler) Upload(c *gin.Context) {
	var headers []*multipart.FileHeader

	// Явный разбор: иначе на части запросов c.MultipartForm() без maxMemory даёт ошибку и поля теряются.
	if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
		if form := c.Request.MultipartForm; form != nil {
			// Несколько файлов с одним именем `file` — обычный способ для браузеров.
			if fs := form.File["file"]; len(fs) > 0 {
				headers = fs
			}
			if len(headers) == 0 {
				if fs := form.File["files"]; len(fs) > 0 {
					headers = fs
				}
			}
		}
	}
	if len(headers) == 0 {
		fh, err := c.FormFile("file")
		if err != nil || fh == nil {
			JSONError(c, http.StatusBadRequest, "file is required")
			return
		}
		headers = []*multipart.FileHeader{fh}
	}
	if len(headers) > uploadMaxFiles {
		JSONError(c, http.StatusBadRequest, "too many files")
		return
	}

	blobs := make([][]byte, 0, len(headers))
	for _, fh := range headers {
		if fh.Size > uploadMaxFileBytes {
			JSONError(c, http.StatusBadRequest, "file too large")
			return
		}
		f, err := fh.Open()
		if err != nil {
			JSONError(c, http.StatusBadRequest, "invalid file")
			return
		}
		b, err := io.ReadAll(f)
		_ = f.Close()
		if err != nil {
			JSONError(c, http.StatusBadRequest, "invalid file")
			return
		}
		if _, _, err = image.Decode(bytes.NewReader(b)); err != nil {
			JSONError(c, http.StatusBadRequest, "File must be an image")
			return
		}
		blobs = append(blobs, b)
	}

	out, err := h.svc.UploadAndMatch(c.Request.Context(), blobs)
	if err != nil {
		RespondError(c, err)
		return
	}

	JSONOK(c, out)
}
