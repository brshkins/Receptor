from fastapi import APIRouter, File, HTTPException, Request, UploadFile

from app.service import DetectionService, InvalidImageError


router = APIRouter()


@router.post("/detect")
async def detect(request: Request, file: UploadFile = File(...)):
    if file is None:
        raise HTTPException(status_code=400, detail="Missing file")

    content_type = (file.content_type or "").lower()
    if content_type and not content_type.startswith("image/"):
        raise HTTPException(status_code=400, detail="File must be an image")

    try:
        data = await file.read()
    except Exception:
        raise HTTPException(status_code=400, detail="Failed to read file")

    svc: DetectionService | None = getattr(request.app.state, "detection_service", None)
    if svc is None:
        raise HTTPException(status_code=503, detail="Model is not ready")

    try:
        ingredients = await svc.detect_ingredients_from_bytes(data)
    except InvalidImageError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception:
        raise HTTPException(status_code=500, detail="Detection failed")

    return {"ingredients": ingredients}

