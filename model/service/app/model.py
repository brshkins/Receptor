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

        # Низкий conf: иначе «помидор» часто отрезается, остаётся только яблоко (красное круглое).
        raw_conf = 0.02
        with self._lock:
            results = self._model.predict(arr, verbose=False, conf=raw_conf)

        if not results:
            return []

        r0 = results[0]
        boxes = getattr(r0, "boxes", None)
        if boxes is None or boxes.cls is None or boxes.conf is None:
            return []

        class_ids: Sequence[float] = boxes.cls.tolist()
        confs: Sequence[float] = boxes.conf.tolist()

        names = getattr(r0, "names", None) or getattr(self._model, "names", None) or {}

        def bucket_class(raw: str) -> str:
            n = (raw or "").strip().lower()
            if n in ("apple -red-", "apple -green-", "apples"):
                return "apple"
            if n in ("tomatoes", "tomato-", "tomate"):
                return "tomato"
            return n

        by_class: dict[str, float] = {}
        for cls_id, conf in zip(class_ids, confs):
            if conf < raw_conf:
                continue
            raw = names.get(int(cls_id))
            if not raw:
                continue
            key = bucket_class(str(raw))
            c = float(conf)
            if c > by_class.get(key, 0.0):
                by_class[key] = c

        apple_c = by_class.pop("apple", 0.0)
        tomato_c = by_class.pop("tomato", 0.0)
        # В датасете лапша/макароны → instant_noodle; оранжевые penne часто ошибочно бьют в apple.
        noodle_c = by_class.pop("instant_noodle", 0.0)

        out: List[Detection] = []

        strong = (
            apple_c >= conf_threshold
            or tomato_c >= conf_threshold
            or noodle_c >= conf_threshold
        )
        if strong:
            if apple_c >= conf_threshold and tomato_c >= conf_threshold:
                if tomato_c + 0.08 >= apple_c:
                    out.append(Detection("tomato", tomato_c))
                else:
                    out.append(Detection("apple", apple_c))
            elif apple_c >= conf_threshold and noodle_c >= conf_threshold:
                if noodle_c + 0.06 >= apple_c:
                    out.append(Detection("instant_noodle", noodle_c))
                else:
                    out.append(Detection("apple", apple_c))
            elif tomato_c >= conf_threshold and noodle_c >= conf_threshold:
                if noodle_c + 0.05 >= tomato_c:
                    out.append(Detection("instant_noodle", noodle_c))
                else:
                    out.append(Detection("tomato", tomato_c))
            elif apple_c >= conf_threshold:
                alts: List[tuple[str, float]] = []
                if tomato_c > 0.03 and (apple_c - tomato_c) < 0.65:
                    alts.append(("tomato", tomato_c))
                # Паста чаще сильно проигрывает яблоку по conf — допускаем больший разрыв.
                if noodle_c > 0.025 and (apple_c - noodle_c) < 0.78:
                    alts.append(("instant_noodle", noodle_c))
                if alts:
                    name, c = max(alts, key=lambda x: x[1])
                    out.append(Detection(name, max(c, 0.05)))
                else:
                    out.append(Detection("apple", apple_c))
            elif tomato_c >= conf_threshold:
                out.append(Detection("tomato", tomato_c))
            elif noodle_c >= conf_threshold:
                out.append(Detection("instant_noodle", noodle_c))
        elif tomato_c > raw_conf:
            if apple_c <= raw_conf or tomato_c + 0.05 >= apple_c:
                out.append(Detection("tomato", max(tomato_c, 0.05)))
        elif noodle_c > raw_conf:
            if apple_c <= raw_conf or noodle_c + 0.05 >= apple_c:
                out.append(Detection("instant_noodle", max(noodle_c, 0.05)))
        elif apple_c > raw_conf:
            out.append(Detection("apple", max(apple_c, 0.05)))

        for cls, c in sorted(by_class.items(), key=lambda x: -x[1]):
            if c >= conf_threshold:
                out.append(Detection(cls, c))

        return out


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

