package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DBUsername = "dias"
	DBPassword = ""
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "news_db"
)

var db *sql.DB
var templatesDir string

type Article struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
}

type VideoLink struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
	Img string `json:"img"`
}

type ViewData struct {
	Articles   []Article
	VideoLinks []VideoLink
}

func main() {
	initDB()

	mux := http.NewServeMux()

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	assetsDir := filepath.Join(baseDir, "shahala/assets")
	templatesDir = filepath.Join(baseDir, "shahala/assets")

	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/about", http.HandlerFunc(aboutHandler))
	mux.Handle("/student", http.HandlerFunc(studentHandler))
	mux.Handle("/map", http.HandlerFunc(mapHandler))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(assetsDir, "css")))))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(filepath.Join(assetsDir, "fonts")))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(assetsDir, "images")))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(assetsDir, "js")))))

	log.Println("Server is running on :5000...")
	err = http.ListenAndServe("localhost:5000", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func initDB() {
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName)
	var err error
	db, err = sql.Open("mysql", dbSource)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			category VARCHAR(255) NOT NULL,
			image_url VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS video_links (
			id INT AUTO_INCREMENT PRIMARY KEY,
			url VARCHAR(255) NOT NULL,
			img VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticlesFromDB()
	if err != nil {
		log.Println("Error getting articles from DB:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoLinks, err := getVideoLinksFromDB()
	if err != nil {
		log.Println("Error getting video links from DB:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := ViewData{
		Articles:   articles,
		VideoLinks: videoLinks,
	}

	renderHTML(w, "index.html", templatesDir, data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticlesFromDB()
	if err != nil {
		log.Println("Error getting articles from DB:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoLinks, err := getVideoLinksFromDB()
	if err != nil {
		log.Println("Error getting video links from DB:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := ViewData{
		Articles:   articles,
		VideoLinks: videoLinks,
	}

	renderHTML(w, "article.html", templatesDir, data)
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticlesFromDB()
	if err != nil {
		log.Println("Error getting articles from DB:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := ViewData{
		Articles: articles,
	}
	renderHTML(w, "first.html", templatesDir, data)
}
func mapHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticlesFromDB()
	if err != nil {
		log.Println("Error getting articles from DB:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := ViewData{
		Articles: articles,
	}
	renderHTML(w, "map.html", templatesDir, data)
}
func renderHTML(w http.ResponseWriter, filename string, templatesDir string, data ViewData) {
	absPath := filepath.Join(templatesDir, filename)
	tmpl, err := template.ParseFiles(absPath)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func getArticlesFromDB() ([]Article, error) {
	rows, err := db.Query("SELECT id, title, category, image_url FROM articles")
	if err != nil {
		log.Println("Error querying articles:", err)
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Category, &article.ImageURL)
		if err != nil {
			log.Println("Error scanning article:", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return articles, nil
}

func getVideoLinksFromDB() ([]VideoLink, error) {
	rows, err := db.Query("SELECT id, url, img FROM video_links")
	if err != nil {
		log.Println("Error querying video links:", err)
		return nil, err
	}
	defer rows.Close()

	var videoLinks []VideoLink
	for rows.Next() {
		var videoLink VideoLink
		err := rows.Scan(&videoLink.ID, &videoLink.URL, &videoLink.Img)
		if err != nil {
			log.Println("Error scanning video link:", err)
			return nil, err
		}
		videoLinks = append(videoLinks, videoLink)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return videoLinks, nil
}
