# For database construction
1. Create tables with `create_tables.sql`
   - Defines the schema for the database (tables, relationships, etc.).
2. Run `get_ids.py`
   - Fetches a list of beer IDs from the API and saves them to a file.
3. Run `fetch_all_beers.py`
   - Fetches details for all beers using the beer IDs and saves them to `beers.json`.
4. Run `populate_brewers_and_beers.py`
   - Populates the database with data from `beers.json`.
