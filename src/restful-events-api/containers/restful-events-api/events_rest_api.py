import json
import requests
import identity_provider
import urllib3
import os
import streamlit as st


def main():
    clientId = os.getenv("CCF_CLIENT_ID")  # The client ID read from the environment
    clientSecret = os.getenv("CCF_CLIENT_SECRET")  # The client secret read from the environment
    legacyUseTls = os.getenv("LEGACY_USE_TLS")  # Whether the XProtect server uses TLS
    legacyManagementServer = os.getenv("LEGACY_MANAGEMENT_SERVER")  # The hostname of the management server, assuming that the API Gateway has been installed on the same host
    verify = legacyUseTls.lower() == "true"  

    if verify:
        serverUrl = "https://" + legacyManagementServer
    else:
        serverUrl = "http://" + legacyManagementServer

    if not verify:
        urllib3.disable_warnings(
            urllib3.exceptions.InsecureRequestWarning
        )  # Remove this line if verifying the certificate (which is recommended)

    # First we need a session to ensure that we stay logged in
    session = requests.Session()

    # Now authenticate using the identity provider and get access token
    response = identity_provider.get_token_secret(
        session, clientId, clientSecret, serverUrl, verify
    )
    if response.status_code == 200:
        token_response = response.json()
        st.info(f"IDP access token response:\n{token_response}\n\n")
        access_token = token_response[
            "access_token"
        ]  # The token that we'll use for RESTful API calls
    else:
        error = response.json()["error"]
        st.info(error)
        return

    headers = {
        "Authorization": f"Bearer {access_token}",
        "Content-type": "application/json",
    }

    # Get an existing user defined event type
    response = requests.get(
        f"{serverUrl}/api/rest/v1/userDefinedEvents", headers=headers, verify=verify
    )
    if response.status_code != 200:
        error = response.json()["error"]
        st.info(error)
        return

    user_defined_events = response.json()["array"]
    
    #Default ids present in the management client that should be ignored
    default_ids = ['85867627-b287-4439-9e55-a63701e1715b', '77b1e70d-ba8d-4bb8-9ee8-43b09746d82a', '7605f8b0-7f5f-4432-b223-0bb2dc3f1f5c']
    user_defined_events = [ids for ids in user_defined_events if ids['id'] not in default_ids]
 
    st.info(f"Retrieved {len(user_defined_events)} user defined event types")

    if len(user_defined_events) == 0:
        st.info(f"You need to have at least one user defined event to run this sample")
        return

    event_type = user_defined_events[0]
    st.info(f"Triggering an event for event type {event_type['id']}")

    # Trigger an event
    response = requests.post(
        f"{serverUrl}/api/rest/v1/events",
        headers=headers,
        data=json.dumps({"type": event_type["id"]}),
        verify=verify,
    )
    if response.status_code != 202:
        error = response.json()
        st.info(error)
        return

    event = response.json()["data"]
    st.info(f"Triggered an event: {event}")

    # Retrieve the first 10 events with additional event data
    response = requests.get(
        f"{serverUrl}/api/rest/v1/events?page=0&size=10&include=data", headers=headers, verify=verify
    )
    if response.status_code != 200:
        error = response.json()
        st.info("Unable to retrieve events.", error)
        return

    events = response.json()["array"]
    st.info(f"Retrieved first page of events: {events}")

    # Retrieve an event by id
    response = requests.get(
        f"{serverUrl}/api/rest/v1/events/{events[-1]['id']}",
        headers=headers,
        verify=verify,
    )
    if response.status_code != 200:
        error = response.json()
        st.info(
            "Unable to retrieve event.",
            f"Make sure that the event type retention of {event_type['name']!r} is greater than zero.",
            error,
        )
        return

    event = response.json()["data"]
    st.info(f"Retrieved an event: {event}")


if __name__ == "__main__":
    main()
