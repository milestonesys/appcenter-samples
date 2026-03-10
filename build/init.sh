#!/bin/bash
set -e

# Check for main tools
YELLOW='\033[1;33m'
for tool in docker kubectl helm curl unzip chmod; do
    if ! command -v $tool &>/dev/null; then
        echo -e "${YELLOW}Warning: $tool is not installed. Please install $tool before continuing." >&2
    fi
done

# Download app-builder.zip
if curl -L -o appbuilder.zip "https://doc.developer.milestonesys.com/appen/App-Builder/app-builder.zip"; then
    unzip -o appbuilder.zip -d .
    chmod +x app-builder.sh
    rm -f appbuilder.zip
    echo "Initialization complete."
else
    echo "Error: Failed to download app-builder.zip" >&2
    exit 1
fi