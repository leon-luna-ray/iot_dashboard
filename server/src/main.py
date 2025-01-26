import base64
import httpx

from decouple import config, Csv
from fastapi import FastAPI, APIRouter, HTTPException

from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles


from .routes.sensors import router as sensor_router

ENV = config("ENV")
CORS_ALLOWED_ORIGINS = config("CORS_ALLOWED_ORIGINS", cast=Csv())

docs_url = None if ENV != "dev" else "/docs"
redoc_url = None if ENV != "dev" else "/redoc"

app = FastAPI(docs_url=docs_url, redoc_url=redoc_url)

app.add_middleware(
    CORSMiddleware,
    allow_origins=CORS_ALLOWED_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(sensor_router, prefix="/api/v1/sensors", tags=["sensors"])

app.mount("/", StaticFiles(directory="src/static", html=True), name="static")
