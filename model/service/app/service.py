import asyncio
from typing import List

from PIL import Image

from app.model import Detection, InvalidImageError, YoloDetector, unique_preserve_order

# Имена классов YOLO → те же токены, что ждёт Go normalizeMatchIngredientNames / БД.
_ML_LABEL_FIX: dict[str, str] = {
    "capsicum": "pepper",
    "bell pepper": "pepper",
    "sweet potato": "potato",
    "floury potato": "potato",
    "solid potato": "potato",
    "potatoes -package-": "potato",
    "tomatoes": "tomato",
    "tomato-": "tomato",
    "meat -red-": "meat",
    "apple -green-": "apple",
    "apple -red-": "apple",
    "egg": "eggs",
    "oranges": "orange",
}


def _normalize_ml_label(name: str) -> str:
    n = (name or "").strip().lower()
    if not n:
        return n
    if n in _ML_LABEL_FIX:
        return _ML_LABEL_FIX[n]
    return n


class DetectionService:
    def __init__(self, detector: YoloDetector, conf_threshold: float = 0.1) -> None:
        self._detector = detector
        self._conf_threshold = conf_threshold

    async def detect_ingredients_from_bytes(self, image_bytes: bytes) -> List[str]:
        img: Image.Image = await asyncio.to_thread(self._detector.decode_image, image_bytes)
        detections: List[Detection] = await asyncio.to_thread(
            self._detector.detect, img, self._conf_threshold
        )
        return unique_preserve_order([_normalize_ml_label(d.name) for d in detections])


__all__ = ["DetectionService", "InvalidImageError"]

