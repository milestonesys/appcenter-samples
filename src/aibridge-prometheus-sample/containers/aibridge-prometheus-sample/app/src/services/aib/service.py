import time
import requests

from src.services.aib.graphql_queries import *
from src.utils.config import settings
from src.services.prometheus.service import AIB_QUERIES_TOTAL, AIB_QUERY_DURATION_SECONDS


def perform_query(query: str, variables: dict = dict(), headers: dict = None) -> dict:
    """
    Perform a query to the AIB API and update Prometheus metrics.
    :param query: query to perform
    :param variables: variables to pass to the query
    :param headers: headers to pass to the query
    :return: response from the query
    """
    url = settings.AIB_URL
    start_time = time.time()  # Start timer for query duration

    try:
        response = requests.post(url=url, json={'query': query, 'variables': variables}, headers=headers)
        response.raise_for_status()

        response_data = response.json()

        # Validate the response
        if 'errors' in response_data:
            AIB_QUERIES_TOTAL.labels(status='failure').inc()  
            raise ValueError(f"GraphQL query failed with errors: {response_data['errors']}")

        if 'data' not in response_data:
            AIB_QUERIES_TOTAL.labels(status='failure').inc() 
            raise ValueError("GraphQL query did not return any data")

        # Update Prometheus metrics for successful query
        AIB_QUERIES_TOTAL.labels(status='success').inc() 
        AIB_QUERY_DURATION_SECONDS.observe(time.time() - start_time) 

        return response_data['data']

    except requests.exceptions.RequestException as e:
        AIB_QUERIES_TOTAL.labels(status='failure').inc() 
        AIB_QUERY_DURATION_SECONDS.observe(time.time() - start_time)  
        raise SystemExit(f"Request failed: {e}")

    except ValueError as e:
        AIB_QUERIES_TOTAL.labels(status='failure').inc()  
        AIB_QUERY_DURATION_SECONDS.observe(time.time() - start_time)  
        raise SystemExit(f"Validation failed: {e}")


def s_get_cameras_query() -> dict:
    """
    Get all cameras from the AIB
    :return: response from the query
    """
    query = get_all_cameras()
    return perform_query(query)


def s_get_about_query() -> dict:
    """
    Get VMS information from the AIB
    :return: response from the query
    """
    query = get_about()
    return perform_query(query)
