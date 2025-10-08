# Use official Python 3.10 image
FROM python:3.10-slim

WORKDIR /restful-config-api

# Copy requirements file and install dependencies
COPY containers/restful-config-api/requirements.txt /restful-config-api/requirements.txt

RUN pip3 install -r requirements.txt --no-cache-dir

# Copy application code
COPY containers/restful-config-api/ /restful-config-api/

# Expose default Streamlit port
EXPOSE 8501

HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health

ENTRYPOINT ["streamlit", "run", "restful_communication.py", "--server.port=8501", "--server.address=0.0.0.0"]