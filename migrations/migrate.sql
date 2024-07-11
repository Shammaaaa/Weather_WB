CREATE TABLE IF NOT EXISTS cities
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    lat     FLOAT        NOT NULL,
    lon     FLOAT        NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cities_name ON cities (name);

CREATE TABLE IF NOT EXISTS weather_forecasts
(
    id       SERIAL PRIMARY KEY,
    city_id  INTEGER   NOT NULL REFERENCES cities (id),
    datetime timestamp NOT NULL,
    temp     FLOAT     NOT NULL,
    data     JSONB     NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_weather_forecasts_city_date ON weather_forecasts (city_id, datetime);



CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE
);


CREATE TABLE favorite_cities
(
    id        SERIAL PRIMARY KEY,
    user_id   INT          NOT NULL,
    city_id   INT          NOT NULL,
    city_name VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE (user_id, city_id)
);