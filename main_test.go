package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// prepareTemplates charge le template HTML avant les tests.
func prepareTemplates(t *testing.T) {
	t.Helper()

	parsedTemplates, err := template.ParseFiles(
		"templates/index.html",
	)
	if err != nil {
		t.Fatalf(
			"impossible de charger le template : %v",
			err,
		)
	}

	previousTemplates := templates
	templates = parsedTemplates

	t.Cleanup(func() {
		templates = previousTemplates
	})
}

// newFormRequest crée une requête POST contenant un formulaire.
func newFormRequest(
	path string,
	values url.Values,
) *http.Request {
	request := httptest.NewRequest(
		http.MethodPost,
		path,
		strings.NewReader(values.Encode()),
	)

	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	return request
}

// TestHomeHandler vérifie le fonctionnement de la page principale.
func TestHomeHandler(t *testing.T) {
	prepareTemplates(t)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "page principale",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "route inexistante",
			method:         http.MethodGet,
			path:           "/inconnue",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "mauvaise méthode",
			method:         http.MethodPost,
			path:           "/",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(
				test.method,
				test.path,
				nil,
			)

			response := httptest.NewRecorder()

			homeHandler(response, request)

			if response.Code != test.expectedStatus {
				t.Errorf(
					"statut obtenu = %d, statut attendu = %d",
					response.Code,
					test.expectedStatus,
				)
			}
		})
	}
}

// TestASCIIArtHandler vérifie les requêtes envoyées au formulaire.
func TestASCIIArtHandler(t *testing.T) {
	prepareTemplates(t)

	tests := []struct {
		name           string
		method         string
		values         url.Values
		expectedStatus int
	}{
		{
			name:           "mauvaise méthode",
			method:         http.MethodGet,
			values:         nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "texte valide",
			method: http.MethodPost,
			values: url.Values{
				"text":   {"Hello"},
				"banner": {"standard"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "texte vide",
			method: http.MethodPost,
			values: url.Values{
				"text":   {""},
				"banner": {"standard"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "texte avec espaces seulement",
			method: http.MethodPost,
			values: url.Values{
				"text":   {"   "},
				"banner": {"standard"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "bannière invalide",
			method: http.MethodPost,
			values: url.Values{
				"text":   {"Hello"},
				"banner": {"inconnue"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "caractère non pris en charge",
			method: http.MethodPost,
			values: url.Values{
				"text":   {"é"},
				"banner": {"standard"},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var request *http.Request

			if test.method == http.MethodPost {
				request = newFormRequest(
					"/ascii-art",
					test.values,
				)
			} else {
				request = httptest.NewRequest(
					test.method,
					"/ascii-art",
					nil,
				)
			}

			response := httptest.NewRecorder()

			asciiArtHandler(response, request)

			if response.Code != test.expectedStatus {
				t.Errorf(
					"statut obtenu = %d, statut attendu = %d",
					response.Code,
					test.expectedStatus,
				)
			}
		})
	}
}

// TestIsValidBanner vérifie les noms de bannière autorisés.
func TestIsValidBanner(t *testing.T) {
	tests := []struct {
		banner   string
		expected bool
	}{
		{
			banner:   "standard",
			expected: true,
		},
		{
			banner:   "shadow",
			expected: true,
		},
		{
			banner:   "thinkertoy",
			expected: true,
		},
		{
			banner:   "inconnue",
			expected: false,
		},
		{
			banner:   "",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.banner, func(t *testing.T) {
			result := isValidBanner(test.banner)

			if result != test.expected {
				t.Errorf(
					"isValidBanner(%q) = %v, résultat attendu = %v",
					test.banner,
					result,
					test.expected,
				)
			}
		})
	}
}

// TestGenerateASCII vérifie que la génération produit huit lignes.
func TestGenerateASCII(t *testing.T) {
	result, err := generateASCII("A", "standard")
	if err != nil {
		t.Fatalf(
			"generateASCII a retourné une erreur : %v",
			err,
		)
	}

	if result == "" {
		t.Fatal("le résultat ASCII ne doit pas être vide")
	}

	lineCount := strings.Count(result, "\n")

	if lineCount != 8 {
		t.Errorf(
			"nombre de lignes obtenu = %d, nombre attendu = 8",
			lineCount,
		)
	}
}

// TestGenerateASCIIWithNewLine vérifie les retours à la ligne.
func TestGenerateASCIIWithNewLine(t *testing.T) {
	result, err := generateASCII(
		"A\\nB",
		"standard",
	)
	if err != nil {
		t.Fatalf(
			"generateASCII a retourné une erreur : %v",
			err,
		)
	}

	lineCount := strings.Count(result, "\n")

	if lineCount != 16 {
		t.Errorf(
			"nombre de lignes obtenu = %d, nombre attendu = 16",
			lineCount,
		)
	}
}

// TestGenerateASCIIInvalidBanner vérifie une bannière inexistante.
func TestGenerateASCIIInvalidBanner(t *testing.T) {
	_, err := generateASCII(
		"Hello",
		"inconnue",
	)

	if err == nil {
		t.Fatal(
			"une erreur était attendue pour une bannière inexistante",
		)
	}
}

// TestGenerateASCIIUnsupportedCharacter vérifie les caractères refusés.
func TestGenerateASCIIUnsupportedCharacter(t *testing.T) {
	_, err := generateASCII(
		"é",
		"standard",
	)

	if err == nil {
		t.Fatal(
			"une erreur était attendue pour un caractère non pris en charge",
		)
	}
}

// TestReadBanner vérifie la lecture d'un fichier de bannière.
func TestReadBanner(t *testing.T) {
	lines, err := readBanner("standard")
	if err != nil {
		t.Fatalf(
			"impossible de lire la bannière standard : %v",
			err,
		)
	}

	if len(lines) == 0 {
		t.Fatal(
			"le fichier de bannière ne doit pas être vide",
		)
	}
}

// TestReadBannerNotFound vérifie une bannière inexistante.
func TestReadBannerNotFound(t *testing.T) {
	_, err := readBanner("inconnue")

	if err == nil {
		t.Fatal(
			"une erreur était attendue pour un fichier inexistant",
		)
	}
}
