import asyncio
from typing import List

from PIL import Image

from app.model import Detection, InvalidImageError, YoloDetector, unique_preserve_order


class DetectionService:
    def __init__(self, detector: YoloDetector, conf_threshold: float = 0.1) -> None:
        self._detector = detector
        self._conf_threshold = conf_threshold

    async def detect_ingredients_from_bytes(self, image_bytes: bytes) -> List[str]:
        img: Image.Image = await asyncio.to_thread(self._detector.decode_image, image_bytes)
        detections: List[Detection] = await asyncio.to_thread(
            self._detector.detect, img, self._conf_threshold
        )
        return unique_preserve_order([d.name for d in detections])


__all__ = ["DetectionService", "InvalidImageError"]

