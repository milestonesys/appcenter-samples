import logging
from datetime import datetime


def service_get_server_status() -> dict:
    """	
    Get the status of the server
    res: StatusResponse: The status of the server
    """	
    logging.info("Received request for server status")
    message = "Server is running!"
    status = "OK"
    timestamp = datetime.now().isoformat()

    return {
        "msg": message,
        "status": status,
        "timestamp": timestamp
    }