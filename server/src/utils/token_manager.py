""" Token Manager for QingPing API """

import base64
import httpx

from datetime import datetime, timedelta, timezone

from fastapi import HTTPException
from decouple import config

QP_APP_KEY = config("QP_APP_KEY")
QP_APP_SECRET = config("QP_APP_SECRET")
QP_AUTH_API_BASE = config("QP_AUTH_API_BASE")

class TokenManager:
    def __init__(self):
        self.token = None
        self.expiry = datetime.now(timezone.utc)

    async def fetch_token(self):
        url = f"{QP_AUTH_API_BASE}/token"
        headers = {
            "Content-Type": "application/x-www-form-urlencoded",
            "Authorization": f"Basic {base64.b64encode(f'{QP_APP_KEY}:{QP_APP_SECRET}'.encode()).decode()}",
        }
        data = {"grant_type": "client_credentials", "scope": "device_full_access"}
        async with httpx.AsyncClient() as client:
            response = await client.post(url, headers=headers, data=data)
            if response.status_code != 200:
                raise HTTPException(
                    status_code=response.status_code, detail=response.text
                )
            token_data = response.json()
            self.token = token_data["access_token"]
            self.expiry = datetime.now(timezone.utc) + timedelta(
                seconds=token_data["expires_in"]
            )

    async def get_token(self):
        if self.token is None or datetime.now(timezone.utc) >= self.expiry:
            await self.fetch_token()
        return self.token


token_manager = TokenManager()
