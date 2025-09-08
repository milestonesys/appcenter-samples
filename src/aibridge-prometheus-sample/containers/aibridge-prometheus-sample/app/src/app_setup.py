from fastapi import FastAPI
from src.routers import router
from prometheus_client import make_asgi_app


def create_app() -> FastAPI:
    """
    Create and configure the FastAPI application.
    """
    app = FastAPI(debug=True)
    app.include_router(router.router, prefix="/api")

    # Add Prometheus endpoint to scrape metrics
    metrics_app = make_asgi_app()
    app.mount("/metrics", metrics_app)
    
    return app