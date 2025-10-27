# Use official Python 3.10 image
FROM python:3.10.12-slim

WORKDIR /restful-config-api

# Copy requirements file and install dependencies

# Copy application code
COPY app/requirements.txt /restful-config-api/app/requirements.txt
COPY app/src /restful-config-api/app/src

RUN pip3 install -r app/requirements.txt --no-cache-dir

# Expose default Streamlit port
EXPOSE 8501

HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health

WORKDIR /restful-config-api/app/src
ENTRYPOINT ["streamlit", "run", "main.py", "--server.port=8501", "--server.address=0.0.0.0"]