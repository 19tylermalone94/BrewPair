CREATE TABLE brewers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    short_description TEXT,
    url TEXT,
    bp_verified BOOLEAN,
    brewer_verified BOOLEAN,
    facebook_url TEXT,
    twitter_url TEXT,
    instagram_url TEXT,
    last_modified BIGINT
);

CREATE TABLE beers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    style TEXT,
    description TEXT,
    abv FLOAT,
    ibu INT,
    bp_verified BOOLEAN,
    brewer_verified BOOLEAN,
    last_modified BIGINT,
    brewer_id UUID REFERENCES brewers(id) ON DELETE CASCADE
);
