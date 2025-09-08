import requests
import pandas as pd
from prometheus_client import Gauge, Counter, Histogram
from datetime import datetime, timedelta

# Prometheus custom metrics
AIB_QUERIES_TOTAL = Counter(
    'aib_queries_total',
    'Total number of queries to AIB',
    ['status']  # Label: success or failure
)

AIB_QUERY_DURATION_SECONDS = Histogram(
    'aib_query_duration_seconds',
    'Duration of AIB queries in seconds'
)
