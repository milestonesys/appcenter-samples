"""
Get a bearer access token from the identity provider
"""
import requests
from requests_ntlm import HttpNtlmAuth

def get_token_secret(
    session: requests.Session, clientId: str, clientSecret: str, serverUrl: str, verify: bool
) -> str:
    """
    Requests an OAuth 2.0 access token from the identity provider on a VMS server for a VMS user.
    The API Gateway forwards the request to the identity provider

    :param session: A requests.Session object which will be used for the duration of the
        integration to maintain logged-in state
    :param clientId: The client ID of an XProtect user with the XProtect Administrators role
    :param clientSecret: The secret of the client logging in
    :param serverUrl: The hostname of the machine hosting the identity provider, e.g. "vms.example.com"
    :param verify: Whether to verify the server's TLS certificate (passed to requests). Set to True to verify, False to skip verification.
    :returns: session.Response object. The value of the 'access_token' property is the bearer token.

        Note the "expires_in" property; if you're planning on making a larger integration, you will
        have to renew before it has elapsed.
    """
    url = f"{serverUrl}/API/IDP/connect/token"
    headers = {"Content-Type": "application/x-www-form-urlencoded"}
    payload = f"grant_type=client_credentials&client_id={clientId}&client_secret={clientSecret}"
    return session.request("POST", url, headers=headers, data=payload, verify=verify)
