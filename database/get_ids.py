import requests
import json
from dotenv import load_dotenv
import os

load_dotenv(dotenv_path="../.env")

API_KEY = os.getenv("CB_API_KEY")
BASE_URL = "https://api.catalog.beer/beer/"
HEADERS = {
    "Accept": "application/json",
    "Authorization": f"Basic {API_KEY}:"
}

OUTPUT_FILE = "cb_beer_ids.txt"

def fetch_beer_ids():
    beer_ids = []
    cursor = None

    while True:
        url = BASE_URL if not cursor else f"{BASE_URL}?cursor={cursor}"

        response = requests.get(url, headers=HEADERS)
        if response.status_code != 200:
            print(f"Error fetching data: {response.status_code}")
            print(response.text)
            break

        data = response.json()
        if "data" not in data:
            print("Unexpected response format")
            break

        for beer in data["data"]:
            beer_ids.append(beer["id"])

        if data.get("has_more"):
            cursor = data.get("next_cursor")
        else:
            break

    with open(OUTPUT_FILE, "w") as f:
        for beer_id in beer_ids:
            f.write(beer_id + "\n")

    print(f"Successfully retrieved {len(beer_ids)} beer IDs.")
    print(f"IDs saved to {OUTPUT_FILE}")

if __name__ == "__main__":
    fetch_beer_ids()
