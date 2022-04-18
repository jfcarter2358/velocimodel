import requests
import sys
from pprint import pprint

obj = sys.argv[1]

r = requests.get(f'http://localhost:9004/api/{obj}')
pprint(r.json())