CREATE TABLE features (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  description TEXT,
  active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE variants (
  id SERIAL PRIMARY KEY,
  feature_id INT REFERENCES features(id),
  name TEXT NOT NULL,
  weight INT NOT NULL
);

CREATE TABLE events (
  id BIGSERIAL PRIMARY KEY,
  feature_id INT REFERENCES features(id),
  user_id TEXT,
  variant TEXT,
  event_type TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);
