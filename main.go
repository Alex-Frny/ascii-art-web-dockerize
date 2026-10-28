package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	Result string
	Text   string
	Banner string
	Error  string
}

var templates *template.Template

func main() {
	var err error

	templates, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Template error:", err)
		return
	}

	// Permet au serveur d'accéder au fichier static/style.css.
	staticFiles := http.FileServer(http.Dir("static"))

	http.Handle(
		"/static/",
		http.StripPrefix("/static/", staticFiles),
	)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	fmt.Println("Server running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

// homeHandler affiche la page principale.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(
			w,
			"404 Not Found",
			http.StatusNotFound,
		)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"400 Bad Request",
			http.StatusBadRequest,
		)
		return
	}

	renderTemplate(
		w,
		http.StatusOK,
		PageData{
			Banner: "standard",
		},
	)
}

// asciiArtHandler récupère le formulaire et génère l'ASCII Art.
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"400 Bad Request",
			http.StatusBadRequest,
		)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(
			w,
			"400 Bad Request",
			http.StatusBadRequest,
		)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// Supprime uniquement les espaces pour vérifier si le texte est vide.
	if strings.TrimSpace(text) == "" {
		renderTemplate(
			w,
			http.StatusBadRequest,
			PageData{
				Text:   text,
				Banner: banner,
				Error:  "Le texte ne peut pas être vide.",
			},
		)
		return
	}

	if !isValidBanner(banner) {
		renderTemplate(
			w,
			http.StatusBadRequest,
			PageData{
				Text:   text,
				Banner: "standard",
				Error:  "La bannière sélectionnée est invalide.",
			},
		)
		return
	}

	result, err := generateASCII(text, banner)
	if err != nil {
		renderTemplate(
			w,
			http.StatusInternalServerError,
			PageData{
				Text:   text,
				Banner: banner,
				Error:  err.Error(),
			},
		)
		return
	}

	renderTemplate(
		w,
		http.StatusOK,
		PageData{
			Result: result,
			Text:   text,
			Banner: banner,
		},
	)
}

// isValidBanner vérifie que la bannière existe.
func isValidBanner(banner string) bool {
	switch banner {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}

// generateASCII transforme le texte en ASCII Art.
func generateASCII(text, banner string) (string, error) {
	bannerLines, err := readBanner(banner)
	if err != nil {
		return "", fmt.Errorf(
			"impossible de lire la bannière %q",
			banner,
		)
	}

	// Normalise les retours à la ligne provenant du formulaire.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Accepte aussi les caractères \n écrits directement.
	text = strings.ReplaceAll(text, "\\n", "\n")

	textLines := strings.Split(text, "\n")

	var result strings.Builder

	for _, textLine := range textLines {
		if textLine == "" {
			result.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			for _, character := range textLine {
				characterLine, err := getCharacterLine(
					character,
					row,
					bannerLines,
				)
				if err != nil {
					return "", err
				}

				result.WriteString(characterLine)
			}

			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// readBanner lit le fichier correspondant à la bannière.
func readBanner(banner string) ([]string, error) {
	fileName := filepath.Join(
		"banners",
		banner+".txt",
	)

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	content := string(data)

	// Normalise les retours à la ligne Windows et Linux.
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	lines := strings.Split(content, "\n")

	return lines, nil
}

// getCharacterLine récupère une ligne précise du dessin d'un caractère.
func getCharacterLine(
	character rune,
	row int,
	bannerLines []string,
) (string, error) {
	if character < 32 || character > 126 {
		return "", fmt.Errorf(
			"le caractère %q n’est pas pris en charge",
			character,
		)
	}

	characterIndex := int(character - 32)

	// Chaque caractère occupe neuf lignes :
	// une ligne vide et huit lignes de dessin.
	startLine := characterIndex*9 + 1
	lineIndex := startLine + row

	if lineIndex < 0 || lineIndex >= len(bannerLines) {
		return "", fmt.Errorf(
			"le fichier de bannière est invalide",
		)
	}

	return bannerLines[lineIndex], nil
}

// renderTemplate envoie les données au fichier HTML.
func renderTemplate(
	w http.ResponseWriter,
	statusCode int,
	data PageData,
) {
	var buffer bytes.Buffer

	err := templates.ExecuteTemplate(
		&buffer,
		"index.html",
		data,
	)
	if err != nil {
		http.Error(
			w,
			"500 Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(statusCode)

	_, err = buffer.WriteTo(w)
	if err != nil {
		fmt.Println("Response writing error:", err)
	}
}
