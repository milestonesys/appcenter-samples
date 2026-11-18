"""
Backend-for-frontend (BFF) for the RESTful Config API microfrontend.

The microfrontend runs inside the Smart Client. The shell adds the bearer access token of the
current Smart Client session to every request it makes, and this service forwards that token to
the XProtect API Gateway. No username or password is ever handled here.

The service also serves the built microfrontend (remote-entry.js, manifest.json, ...).
"""
import json
import os
from pathlib import Path
from typing import Optional

import requests
import uvicorn
from fastapi import Depends, FastAPI, Header, HTTPException, Response
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel

from api_gateway import Gateway

RESOURCE = "userDefinedEvents"

# Environment variables available for a container
legacy_use_tls = os.getenv("LEGACY_USE_TLS", "true")  # Whether the XProtect server uses TLS
legacy_management_server = os.getenv("LEGACY_MANAGEMENT_SERVER", "")  # Hostname of the Management Server

scheme = "https" if legacy_use_tls.lower() == "true" else "http"
server_url = f"{scheme}://{legacy_management_server}"

# Location of the built microfrontend. Defaults to ../client/dist relative to this file's app dir.
client_dir = Path(os.getenv("CLIENT_DIR", Path(__file__).resolve().parent.parent / "client"))

api_gateway = Gateway(server_url)
app = FastAPI(title="RESTful Config API sample")


class EventBody(BaseModel):
    name: str


def get_token(authorization: Optional[str] = Header(default=None)) -> str:
    """Extract the Smart Client session's bearer token from the request.

    The token is not validated here; the API Gateway validates it when we forward it.
    """
    scheme_name, _, token = (authorization or "").partition(" ")
    if scheme_name.lower() != "bearer" or not token.strip():
        raise HTTPException(status_code=401, detail="Missing bearer access token")
    return token.strip()


def relay(response: requests.Response) -> Response:
    """Pass the API Gateway's status code and JSON body on to the microfrontend."""
    return Response(
        content=response.content,
        status_code=response.status_code,
        media_type=response.headers.get("content-type", "application/json"),
    )


def call_gateway(call):
    """Run a gateway call on a fresh session, mapping connection problems to 502."""
    with requests.Session() as session:
        try:
            return relay(call(session))
        except requests.exceptions.RequestException as err:
            raise HTTPException(status_code=502, detail=f"Could not reach the XProtect server: {err}")


@app.get("/api/server-info")
def server_info():
    return {"serverUrl": server_url}


@app.post("/api/user-defined-events")
def create_event(body: EventBody, token: str = Depends(get_token)):
    payload = json.dumps({"name": body.name})
    return call_gateway(lambda s: api_gateway.create_item(s, RESOURCE, payload, token))


@app.get("/api/user-defined-events/{event_id}")
def get_event(event_id: str, token: str = Depends(get_token)):
    return call_gateway(lambda s: api_gateway.get_single(s, RESOURCE, event_id, token))


@app.put("/api/user-defined-events/{event_id}")
def update_event(event_id: str, body: EventBody, token: str = Depends(get_token)):
    payload = json.dumps({"name": body.name})
    return call_gateway(lambda s: api_gateway.update_item(s, RESOURCE, payload, event_id, token))


@app.delete("/api/user-defined-events/{event_id}")
def delete_event(event_id: str, token: str = Depends(get_token)):
    return call_gateway(lambda s: api_gateway.delete_item(s, RESOURCE, event_id, token))


# Registered last so it does not shadow the API routes above.
if client_dir.is_dir():
    app.mount("/", StaticFiles(directory=client_dir), name="client")


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=int(os.getenv("PORT", "8080")))
