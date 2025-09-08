# external libs
import os
from fastapi import APIRouter, HTTPException
from fastapi.responses import JSONResponse
import logging
from datetime import datetime

# internal
from src.services.api.service import service_get_server_status
from src.services.aib.service import s_get_about_query, s_get_cameras_query

logger = logging.getLogger(__name__)
router = APIRouter()


@router.get("/get-server-status")
async def get_server_status() -> JSONResponse:
    """
    Get server status
    :return: JSONResponse with server status
    """
    return service_get_server_status()


@router.get("/configuration-info")
async def get_info() -> JSONResponse:
    """
    Get ENV variables for prometheus and aibridge
    :return: JSONResponse with ENV variables
    """
    try:
        env_variables = {
            key: value for key, value in os.environ.items()
            if key in ["AIB_URL", "PROMETHEUS_URL"]
        }
        
        res = {
            "date": datetime.now().isoformat(),
            "status": "OK",
            "env_variables": env_variables
        }
        return res
    
    except Exception as e:
        logger.error(f"Error retrieving environment variables: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve environment variables")
    

@router.get("/vms-info")
async def get_vms_info() -> JSONResponse:
    """
    Get VMS information from AIB
    :return: JSONResponse with VMS information
    """
    try:
        vms_info = s_get_about_query()
        if vms_info is None:
            raise HTTPException(status_code=500, detail="Failed to retrieve VMS information")
        res = {
            "date": datetime.now().isoformat(),
            "status": "OK",
            "vms_info": vms_info
        }
        return res
    except Exception as e:
        logger.error(f"Error retrieving cameras information: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve cameras information")
    

@router.get("/cameras-info")
async def get_cameras_info() -> JSONResponse:
    """
    Get cameras information from AIB
    :return: JSONResponse with cameras information
    """
    try:
        cameras_info = s_get_cameras_query()
        if cameras_info is None:
            raise HTTPException(status_code=500, detail="Failed to retrieve cameras information")
        res = {
            "date": datetime.now().isoformat(),
            "status": "OK",
            "cameras_info": cameras_info
        }
        return res
    except Exception as e:
        logger.error(f"Error retrieving cameras information: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve cameras information")