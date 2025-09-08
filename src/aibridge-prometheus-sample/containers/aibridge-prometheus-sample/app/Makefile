HOST=0.0.0.0
PORT=9091
VENV = venv

# virtualenv
init: create-env
	if [ ! -d $(VENV) ]; then \
		echo "Virtual env does not exist, creating..."; \
		python3 -m venv $(VENV); \
	fi
	echo "Activating virtual environment..." ; \
	. $(VENV)/bin/activate && \
	pip3 install -r requirements.txt; \

# Create .env file with environment variables
create-env:
	@echo "AIB_URL=http://aibridge-webservice.processing-server:4000/api/bridge/graphql" > .env
	@echo "PROMETHEUS_URL=http://prometheus-server.prometheus:80" >> .env

# Run the application locally
run:
	. $(VENV)/bin/activate && \
	uvicorn src.main:app --host $(HOST) --port $(PORT) --reload