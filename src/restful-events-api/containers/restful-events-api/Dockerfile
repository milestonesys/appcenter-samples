# Use official Python 3.10 image
FROM python:3.10-slim

# Set working directory
WORKDIR /restful-events-api

# Copy project files
COPY src/. .
COPY requirements.txt requirements.txt
# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose default Streamlit port
EXPOSE 8501

HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health

ENTRYPOINT ["streamlit", "run", "events_rest_api.py", "--server.port=8501", "--server.address=0.0.0.0"]