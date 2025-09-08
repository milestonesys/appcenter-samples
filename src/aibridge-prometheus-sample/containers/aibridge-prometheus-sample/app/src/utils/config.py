import os
from dotenv import load_dotenv
from pydantic_settings import BaseSettings

load_dotenv()

class Settings(BaseSettings):
    # url for the k8s services
    AIB_URL: str = os.getenv("AIB_URL", "http://aibridge-webservice.processing-server:4000/api/bridge/graphql")
    PROMETHEUS_URL: str = os.getenv("PROMETHEUS_URL", "http://prometheus-server.prometheus:80")

settings = Settings()
