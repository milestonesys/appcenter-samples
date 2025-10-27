"""
Get a bearer access token from the identity provider
"""
import requests
from requests_ntlm import HttpNtlmAuth

class IdentityProvider:
    def __init__(self, identity_provider_url: str):
        """
        Initialize the IdentityProvider with the server URL.
        
        :param identity_provider_url: The URL of the identity provider on the VMS server
        """
        self.identity_provider_url = identity_provider_url

    def get_token(
        self,
        session: requests.Session,
        username: str,
        password: str,
        is_basic_user: bool,
    ) -> requests.Response:
        """
        Requests an OAuth 2.0 access token from the identity provider on a VMS server for a VMS user.
        The API Gateway forwards the request to the identity provider

        :param session: A requests.Session object which will be used for the duration of the
            integration to maintain logged-in state
        :param username: The username of an XProtect user with the XProtect Administrators role
        :param password: The password of the user logging in
        :param is_basic_user: Defines whether the login should be done using basic authentication

        :returns: session.Response object. The value of the 'access_token' property is the bearer token.

        Note the "expires_in" property; if you're planning on making a larger integration, you will
        have to renew before it has elapsed.
        """

        if is_basic_user:
            return self._get_token_basic(session, username, password)
        else:
            return self._get_token_windows(session, username, password)


    def _get_token_basic(
        self, session: requests.Session, username: str, password: str
    ) -> requests.Response:
        url = f"{self.identity_provider_url}/connect/token"
        headers = {"Content-Type": "application/x-www-form-urlencoded"}
        payload = f"grant_type=password&username={username}&password={password}&client_id=GrantValidatorClient"
        return session.request("POST", url, headers=headers, data=payload, verify=False)


    def _get_token_windows(
        self, session: requests.Session, username: str, password: str
    ) -> requests.Response:
        url = f"{self.identity_provider_url}/connect/token"
        headers = {"Content-Type": "application/x-www-form-urlencoded"}
        payload = f"grant_type=windows_credentials&client_id=GrantValidatorClient"
        return session.request(
            "POST",
            url,
            headers=headers,
            data=payload,
            verify=False,
            auth=HttpNtlmAuth(username, password),
        )