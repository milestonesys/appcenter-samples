import json
import os
import requests
from identity_provider import IdentityProvider
from api_gateway import Gateway
import streamlit as st

def main():
    st.title("RESTful Config API - Python")

    # Environment variables available for a container
    identity_provider_url = os.getenv("SYSTEM_IDENTITY_PROVIDER") # The URL of the Identity Provider of the System
    legacy_use_tls = os.getenv("LEGACY_USE_TLS") # Whether the XProtect server uses TLS
    legacy_management_server = os.getenv("LEGACY_MANAGEMENT_SERVER") # The hostname of the Management Server in XProtect

    if legacy_use_tls.lower() == "true":
        server_url = "https://" + legacy_management_server
    else:
        server_url = "http://" + legacy_management_server

    st.info(f"XProtect server URL: {server_url}")

    # Input data
    username = st.text_input("Username")
    password = st.text_input("Password", type="password")
    is_basic_user = st.selectbox("Use basic authentication?", ("True", "False"))
    event_name = st.text_input("Name of user defined event to create")

    # Run demo flow
    if st.button("Run Demo"):
        # Clear any previous delete response when starting a new demo
        if 'delete_response' in st.session_state:
            del st.session_state.delete_response
        run_demo(username, password, is_basic_user == "True", server_url, identity_provider_url, event_name)
    
    update_user_defined_event_section()
    delete_user_defined_event_section()
    display_delete_response_section()

def run_demo(username: str, password: str, is_basic_user: bool, server_url: str, identity_provider_url: str, event_name: str) -> None:
    
    # Initialize session in session_state to ensure that we stay logged in during demo flow
    st.session_state.session = requests.Session()
    
    # Now authenticate using the identity provider and get access token
    identity_provider = IdentityProvider(identity_provider_url)
    response = identity_provider.get_token(
        st.session_state.session, username, password, is_basic_user
    )

    with st.expander("View IDP access token response"):
        if response.status_code == 200:
            token_response = response.json()
            st.json(token_response)
            access_token = token_response[
                "access_token"
            ]  # The token that we'll use for RESTful API calls
        else:
            error = response.json()["error"]
            st.error("Failed to get access token")
            st.write(error)
            return

    # Create an API Gateway
    api_gateway = Gateway(server_url)

    # Demo of creating, updating, and deleting a user-defined event through the API Gateway
    create_user_defined_event(api_gateway, access_token, event_name)

def create_user_defined_event(api_gateway: Gateway, token: str, event_name: str) -> None:
    """Create a user-defined event"""
    
    # Create a user defined event
    payload = json.dumps({"name": event_name})
    response = api_gateway.create_item(st.session_state.session, "userDefinedEvents", payload, token)
    with st.expander("View create item result"):
        if response.status_code == 201:
            create_result = response.json()["result"]
            st.json(create_result)
        else:
            error = response.json()["error"]
            st.error("Failed to create event")
            st.json(error)
            return

    # Get the user defined event that we just created
    event_id = create_result["id"]

    # Store essential data in session state for button callbacks
    st.session_state.event_id = event_id
    st.session_state.api_gateway = api_gateway
    st.session_state.token = token

    response = api_gateway.get_single(st.session_state.session, "userDefinedEvents", event_id, token)
    with st.expander("View get item result"):
        if response.status_code == 200:
            event_data = response.json()["data"]
            st.json(event_data)
        else:
            error = response.json()["error"]
            st.error("Failed to get event")
            st.json(error)
            return

def update_user_defined_event_section() -> None:
    """Update user-defined event"""

    if 'event_id' not in st.session_state:
        return

    new_event_name = st.text_input("Enter new event name:")
    if st.button("Update Event"):
        payload = json.dumps({"name": new_event_name})
        response = st.session_state.api_gateway.update_item(
            st.session_state.session, "userDefinedEvents", payload, st.session_state.event_id, st.session_state.token
        )

        with st.expander("Update item data"):
            if response.status_code == 200:
                update_data = response.json()["data"]
                st.success("Event updated successfully!")
                st.json(update_data)
            else:
                error = response.json()["error"]
                st.error("Failed to update event")
                st.write(error)

def delete_user_defined_event_section() -> None:
    """Delete user-defined event"""

    if 'event_id' not in st.session_state:
        return
    
    if st.button("Delete Event"):
        response = st.session_state.api_gateway.delete_item(
            st.session_state.session, "userDefinedEvents", st.session_state.event_id, st.session_state.token
        )
        
        if response.status_code == 200:
            delete_item_state = response.json()
            
            # Store delete response in session state to persist it
            st.session_state.delete_response = delete_item_state
            
            # Clear the event-related session state
            if 'event_id' in st.session_state:
                del st.session_state.event_id
            if 'api_gateway' in st.session_state:
                del st.session_state.api_gateway  
            if 'token' in st.session_state:
                del st.session_state.token
            # Note: Keep the st.session_state.session for potential future operations
            
            # Refresh the page 
            st.rerun()
        else:
            error = response.json()["error"]
            st.error("Failed to delete event")
            st.write(error)

def display_delete_response_section() -> None:
    """Display the delete response if available."""
    if 'delete_response' in st.session_state:
        st.success("Event deleted successfully!")
        with st.expander("View delete response"):
            st.json(st.session_state.delete_response)

if __name__ == "__main__":
    main()