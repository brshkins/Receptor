import logging

from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse

from app.model import ModelLoadError, YoloDetector, get_model_path_from_env
from app.routes import router
from app.service import DetectionService


logger = logging.getLogger("ml_service")


def create_app() -> FastAPI:
    app = FastAPI()

    @app.on_event("startup")
    def _startup() -> None:
        model_path = get_model_path_from_env()
        detector = YoloDetector(model_path=model_path)
        # Ниже порог — чаще ловим объект на одном фото (помидор/апельсин и т.д.).
        app.state.detection_service = DetectionService(detector=detector, conf_threshold=0.05)
        logger.info("YOLO model loaded")

    @app.exception_handler(ModelLoadError)
    async def _model_load_error_handler(_: Request, exc: ModelLoadError):
        return JSONResponse(status_code=500, content={"detail": str(exc)})

    @app.exception_handler(HTTPException)
    async def _http_exception_handler(_: Request, exc: HTTPException):
        return JSONResponse(status_code=exc.status_code, content={"detail": exc.detail})

    app.include_router(router)
    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("app.main:app", host="0.0.0.0", port=8000, log_level="info")

