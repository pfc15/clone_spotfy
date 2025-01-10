package main

type musica struct {
	nome string;
	path string;
	artista string;
};

type send_musica struct {
	Nome string `json:"nome"`;
	Artista string `json:"artista"`;
	Caminho string `json:"path"`;
}
