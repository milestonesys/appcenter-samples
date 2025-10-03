# Use official Python 3.10 image
FROM python:3.10-slim

# Set working directory
WORKDIR /restful-config-api

# Copy project files
COPY /restful-config-api .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose default Streamlit port
EXPOSE 8501

HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health

ENTRYPOINT ["streamlit", "run", "restful_communication.py", "--server.port=8501", "--server.address=0.0.0.0"]