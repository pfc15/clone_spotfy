package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)



type Manager struct {

	db *sql.DB
}

func new_manager() *Manager{
	db, err := sql.Open("sqlite3", "mydb.db")
	if err!=nil {
		log.Println("erro na leitura do banco de dados, ",err)
		return nil
	}

	m := &Manager{
		db:db,
	}
	m.setupDB()
	return m

}


func (m *Manager) setupDB() {
	content, err := os.ReadFile("setup.sql")
	if err!=nil{
		log.Println("error reading file: ", err)
		return
	}
	sql := strings.Split(string(content), ";")
	for _, promt:= range sql {
		m.db.Exec(promt+";")
	}
}

func (m *Manager) listMusic(w http.ResponseWriter, r *http.Request)  {
	var resp_rows *sql.Rows;
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing from", http.StatusBadRequest)
			return
		}

		artista := r.FormValue("artista")
		nome_musica := r.FormValue("nome")

		resp_rows, _ = m.db.Query(
			fmt.Sprintf("SELECT nome, artista, caminho FROM musica where nome LIKE '%s' AND artista LIKE '%s';",
				 nome_musica, artista))

	}else {
		resp_rows, _ = m.db.Query("SELECT nome, artista, caminho FROM musica;")
	}
	
	var lista_musicas []send_musica;

	for resp_rows.Next() {
		var value send_musica
		if err := resp_rows.Scan(&value.Nome, &value.Artista, &value.Caminho); err!=nil {
			log.Println(err)
			return 
		}
		lista_musicas = append(lista_musicas, value)
	}
	
	log.Println(lista_musicas)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(lista_musicas); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		log.Println(err)
		return
	}

}


func (m *Manager) uploadMusic(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1 << 30)

	if err := r.ParseMultipartForm(1 << 30); err != nil {
		http.Error(w, "File size limit exceded hu3", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to retrive the file", http.StatusBadRequest)
		log.Println("FormFile", err)
		return
	}

	defer file.Close()

	dirpath := "./musicas/" + r.FormValue("artista")

	_ , err = os.Stat(dirpath)
	if os.IsNotExist(err) {
		os.Mkdir(dirpath, os.ModePerm)
	}

	destino, err := os.Create( dirpath +"/" + header.Filename)
	if err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		log.Println("file Create error: ", err)
		return
	}

	if _, err := io.Copy(destino, file); err != nil {
		http.Error(w, "failed to saved the file", http.StatusInternalServerError)
		log.Println("File copy error: ", err)
		return
	}

	nova_musica := musica{
		nome: r.FormValue("nome"),
		artista: r.FormValue("artista"),
		path: fmt.Sprintf("./musicas/%s/%s", r.FormValue("artista"),header.Filename),
	}
	

	_, err = m.db.Exec(
		fmt.Sprintf("INSERT INTO musica(nome, artista, caminho) VALUES ('%s', '%s', '%s');", 
		nova_musica.nome, nova_musica.artista, nova_musica.path,
		))
	if err!= nil{
		log.Println(err)
		return
	}else {
		log.Println("musica add com sucesso,", nova_musica)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File %s uploaded successfully!", header.Filename)
}

func (m *Manager) streamMusic(w http.ResponseWriter, r *http.Request) {
	if err:= r.ParseForm();err!= nil {
		log.Println(err)
		return
	}
	log.Println(r.FormValue("musica"))
	var path string
	path_row := m.db.QueryRow(fmt.Sprintf("SELECT caminho FROM musica WHERE '%s' = nome", r.FormValue("musica")))
	path_row.Scan(&path)

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "could not open the music file", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		if err!= nil {
			http.Error(w, "Error reading the file", http.StatusInternalServerError)
			return
		}

		w.Write(buffer[:n])
	}
}

func (m *Manager) addPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var playlist playlist

	err = json.Unmarshal(body, &playlist)
	if err!=nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received Playlist: %+v\n", playlist)
	log.Printf("Parsed Playlist: %+v\n", playlist)
	
	
}