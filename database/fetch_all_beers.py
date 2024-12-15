import requests
import json
from concurrent.futures import ThreadPoolExecutor, as_completed
from dotenv import load_dotenv
import os

load_dotenv(dotenv_path="../.env")

API_KEY = os.getenv("CB_API_KEY")
BASE_URL = "https://api.catalog.beer/beer/"
HEADERS = {
    "Accept": "application/json",
    "Authorization": f"Basic {API_KEY}:"
}

INPUT_FILE = "cb_beer_ids.txt"
OUTPUT_FILE = "beers.json"

def fetch_beer(beer_id):
    """Fetch a single beer object."""
    url = BASE_URL + beer_id
    try:
        response = requests.get(url, headers=HEADERS)
        if response.status_code == 200:
            return response.json()
        else:
            print(f"Failed to fetch beer ID {beer_id}: {response.status_code}")
            return None
    except Exception as e:
        print(f"Error fetching beer ID {beer_id}: {e}")
        return None

def fetch_all_beers_concurrently():
    """Fetch all beer objects using multiple threads."""
    beer_objects = []
    
    with open(INPUT_FILE, "r") as f:
        beer_ids = [line.strip() for line in f.readlines()]

    print(f"Fetching details for {len(beer_ids)} beers...")

    with ThreadPoolExecutor(max_workers=10) as executor:  # Adjust max_workers as needed
        future_to_id = {executor.submit(fetch_beer, beer_id): beer_id for beer_id in beer_ids}
        for i, future in enumerate(as_completed(future_to_id)):
            result = future.result()
            if result:
                beer_objects.append(result)

            if (i + 1) % 10 == 0 or (i + 1) == len(beer_ids):
                print(f"Fetched {i + 1}/{len(beer_ids)} beers")

    with open(OUTPUT_FILE, "w") as f:
        json.dump(beer_objects, f, indent=2)

    print(f"Successfully retrieved {len(beer_objects)} beer objects.")
    print(f"Beer objects saved to {OUTPUT_FILE}")

if __name__ == "__main__":
    fetch_all_beers_concurrently()
