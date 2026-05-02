import os
import threading
from dataclasses import dataclass
from io import BytesIO
from typing import Iterable, List, Sequence

import numpy as np
from PIL import Image, UnidentifiedImageError
from ultralytics import YOLO


class ModelLoadError(RuntimeError):
    pass


class InvalidImageError(ValueError):
    pass


@dataclass(frozen=True)
class Detection:
    name: str
    confidence: float


class YoloDetector:
    def __init__(self, model_path: str) -> None:
        if not model_path:
            raise ModelLoadError("MODEL_PATH is empty")
        if not os.path.exists(model_path):
            raise ModelLoadError(f"MODEL_PATH does not exist: {model_path}")

        try:
            self._model = YOLO(model_path)
        except Exception as e:  # ultralytics can raise various exceptions
            raise ModelLoadError(f"Failed to load YOLO model from {model_path}") from e

        self._lock = threading.Lock()

    @staticmethod
    def decode_image(image_bytes: bytes) -> Image.Image:
        if not image_bytes:
            raise InvalidImageError("Empty file")
        try:
            img = Image.open(BytesIO(image_bytes))
            img = img.convert("RGB")
            return img
        except UnidentifiedImageError as e:
            raise InvalidImageError("Unsupported or corrupted image") from e
        except Exception as e:
            raise InvalidImageError("Failed to decode image") from e

    def detect(self, img: Image.Image, conf_threshold: float) -> List[Detection]:
        if conf_threshold < 0.0 or conf_threshold > 1.0:
            raise ValueError("conf_threshold must be in [0, 1]")

        # Ultralytics works with numpy arrays; keep everything in memory.
        arr = np.asarray(img)

        with self._lock:
            results = self._model.predict(arr, verbose=False)

        if not results:
            return []

        r0 = results[0]
        boxes = getattr(r0, "boxes", None)
        if boxes is None or boxes.cls is None or boxes.conf is None:
            return []

        class_ids: Sequence[float] = boxes.cls.tolist()
        confs: Sequence[float] = boxes.conf.tolist()

        names = getattr(r0, "names", None) or getattr(self._model, "names", None) or {}

        detections: List[Detection] = []
        for cls_id, conf in zip(class_ids, confs):
            if conf < conf_threshold:
                continue
            name = names.get(int(cls_id))
            if not name:
                continue
            detections.append(Detection(name=name, confidence=float(conf)))

        return detections


def get_model_path_from_env() -> str:
    return os.getenv("MODEL_PATH", "").strip()


def unique_preserve_order(items: Iterable[str]) -> List[str]:
    seen = set()
    out: List[str] = []
    for x in items:
        if x in seen:
            continue
        seen.add(x)
        out.append(x)
    return out

