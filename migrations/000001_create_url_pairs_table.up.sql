CREATE TABLE url_pairs (
    short VARCHAR(255) NOT NULL PRIMARY KEY ,
    long TEXT NOT NULL
);

CREATE INDEX idx_url_pairs_short ON url_pairs (short);
CREATE INDEX idx_url_pairs_long ON url_pairs (long);
