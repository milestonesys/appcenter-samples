import os
import argparse
import uvicorn

from src.utils.config import settings
from src.app_setup import create_app

app = create_app()


def main():
    parser = argparse.ArgumentParser(description="Start the FastAPI server")
    parser.add_argument("--host", type=str, default=os.getenv("HOST", "0.0.0.0"), help="Host to run the server on")
    parser.add_argument("--port", type=int, default=int(os.getenv("PORT", 9090)), help="Port to run the server on")
    parser.add_argument("--aib-url", type=str, default=os.getenv("AIB_URL", settings.AIB_URL), help="URL of the AIB server")
    args = parser.parse_args()

    # Update settings with parsed arguments
    settings.AIB_URL = args.aib_url

    uvicorn.run(app, port=args.port)

if __name__ == "__main__":
    main()
