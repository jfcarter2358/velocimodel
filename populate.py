import requests
from cartils import generators
import random

num_models = 4
num_assets = 3
num_snapshots = 2
num_releases = 1

tags = [
    "foo",
    "bar",
    "linear",
    "regression",
    "xgboost",
    "nlp",
    "credit",
    "consumer"
]
departments = [
    "sales",
    "marketing",
    "infrastructure",
    "platform",
    "applications"
]
languages = [
    "Python",
    "R",
    "Go",
    "Java",
    "C++"
]

for i in range(0, num_models):
    model_name = generators.name()
    model_tags = random.sample(tags, 2)
    model_department = random.sample(departments, 1)[0]
    model_language = random.sample(languages, 1)[0]
    data = {
        "name": model_name,
        "type": "raw",
        "tags": model_tags,
        "metadata": {"department": model_department},
        "assets": [],
        "snapshots": [],
        "releases": [],
        "language": model_language
    }

    r = requests.post('http://localhost:9004/api/model', json=data)
    model_id = r.json()['id']

    print(f'model: {i}')
    print(r.status_code)

    for j in range(0, num_assets):
        asset_name = generators.name()
        data = {
            "name": asset_name,
            "type": "file",
            "tags": model_tags,
            "metadata": {"department": model_department},
            "models": [],
            "url": f'{asset_name}'
        }

        r = requests.post('http://localhost:9004/api/asset', json=data)
        asset_id = r.json()['id']

        print(f'asset: {i} {j}')
        print(r.status_code)

        data = {
            'model': model_id,
            'asset': asset_id
        }
        r = requests.post('http://localhost:9004/api/model/asset', json=data)

        print(f'model asset: {i} {j}')
        print(r.status_code)

    r = requests.get('http://localhost:9004/api/model', params={"filter": f"id = \"{model_id}\""})
    model = r.json()[0]
    
    for j in range(0, num_snapshots):
        r = requests.post('http://localhost:9004/api/snapshot/model', json=model)
        snapshot_id = r.json()['id']

        print(f'snapshot model: {i} {j}')
        print(r.status_code)

        r = requests.get('http://localhost:9004/api/snapshot', params={"filter": f"id = \"{snapshot_id}\""})
        snapshot = r.json()[0]

        print(f'snapshot: {i} {j}')
        print(r.status_code)

        for k in range(0, num_releases):
            r = requests.post('http://localhost:9004/api/release/snapshot', json=snapshot)
            release_id = r.json()['id']

            print(f'release: {i} {j} {k}')
            print(r.status_code)
    