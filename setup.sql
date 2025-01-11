CREATE TABLE IF NOT EXISTS musica(
    nome TEXT PRIMARY KEY,
    artista TEXT,
    caminho TEXT
);


CREATE TABLE IF NOT EXISTS playlist (
    nome TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS possui_musica (
    playlist TEXT,
    musica TEXT,
    PRIMARY KEY(playlist, musica),
    FOREIGN KEY (musica) REFERENCES musica(nome),
    FOREIGN KEY (playlist) REFERENCES playlist(nome)
);

