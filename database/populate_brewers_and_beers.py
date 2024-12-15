import json
import psycopg2
from tqdm import tqdm
from dotenv import load_dotenv
import os

load_dotenv("../.env")

DB_CONFIG = {
    "dbname": os.getenv("DB_NAME"),
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASSWORD"),
    "host": os.getenv("DB_HOST"),
    "port": os.getenv("DB_PORT"),
}

with open("beers.json", "r") as file:
    beers = json.load(file)

def insert_data():
    connection = psycopg2.connect(**DB_CONFIG)
    cursor = connection.cursor()
    inserted_brewers = set()
    print("Starting data insertion...")
    for beer in tqdm(beers, desc="Processing beers"):
        brewer = beer.get("brewer", {})
        brewer_id = brewer.get("id")
        if brewer_id and brewer_id not in inserted_brewers:
            cursor.execute("""
                INSERT INTO brewers (
                    id, name, description, short_description, url, 
                    bp_verified, brewer_verified, facebook_url, twitter_url, 
                    instagram_url, last_modified
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (id) DO NOTHING
            """, (
                brewer_id,
                brewer.get("name"),
                brewer.get("description"),
                brewer.get("short_description"),
                brewer.get("url"),
                brewer.get("bp_verified"),
                brewer.get("brewer_verified"),
                brewer.get("facebook_url"),
                brewer.get("twitter_url"),
                brewer.get("instagram_url"),
                brewer.get("last_modified")
            ))
            inserted_brewers.add(brewer_id)

        cursor.execute("""
            INSERT INTO beers (
                id, name, style, description, abv, ibu, 
                bp_verified, brewer_verified, last_modified, brewer_id
            )
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (id) DO NOTHING
        """, (
            beer["id"],
            beer["name"],
            beer.get("style"),
            beer.get("description"),
            beer.get("abv"),
            beer.get("ibu"),
            beer.get("bp_verified"),
            beer.get("brewer_verified"),
            beer.get("last_modified"),
            brewer_id
        ))

    connection.commit()
    cursor.close()
    connection.close()
    print("All data inserted successfully!")

if __name__ == "__main__":
    insert_data()
