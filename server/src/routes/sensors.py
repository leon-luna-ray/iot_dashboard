import httpx
import datetime

from decouple import config
from fastapi import APIRouter, HTTPException, Depends, status

from ..utils.token_manager import token_manager

QP_API_BASE = config("QP_API_BASE")

router = APIRouter()

# /api/v1/sensors/
@router.get("/", status_code=status.HTTP_200_OK)
async def get_sensor_data(token: str = Depends(token_manager.get_token)):
    url = f"{QP_API_BASE}/devices?timestamp={int(datetime.datetime.now().timestamp())}"

    headers = {"Authorization": f"Bearer {token}"}

    async with httpx.AsyncClient() as client:
        response = await client.get(url, headers=headers)
        print(response)
        if response.status_code != 200:
            raise HTTPException(status_code=response.status_code, detail=response.text)
        return response.json()
