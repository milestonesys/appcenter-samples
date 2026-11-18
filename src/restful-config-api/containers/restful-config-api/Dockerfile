# Use official Python 3.10 image
FROM python:3.10.12-slim

WORKDIR /restful-config-api

# Install dependencies
COPY app/requirements.txt /restful-config-api/app/requirements.txt
RUN pip3 install -r app/requirements.txt --no-cache-dir

# Copy application code (backend-for-frontend)
COPY app/src /restful-config-api/app/src

# Copy the built microfrontend (npm run build:federation, done by `make build`)
COPY client/dist /restful-config-api/app/client

EXPOSE 8080

HEALTHCHECK CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8080/api/server-info')"

WORKDIR /restful-config-api/app/src
ENTRYPOINT ["python", "main.py"]
