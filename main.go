package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type wikipediaResponse struct {
	Query struct {
		Pages map[string]struct {
			Pageid  int    `json:"pageid"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Add /description?name=<name>&lang=<lang> in your browser, such that name and language are query parameters")
	}).Methods(http.MethodGet)

	r.HandleFunc("/description", getDescription).Methods(http.MethodGet)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getDescription(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	lang := query.Get("lang")

	// Check for name input
	if name == "" {
		http.Error(w, "Please provide a name.", http.StatusBadRequest)
		return
	}

	// Get short description from Wikipedia API
	desc, err := getShortDescription(name, lang)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching short description: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the short description is empty
	if desc == "" {
		http.Error(w, "Short description not found.", http.StatusInternalServerError)
		return
	}

	// Return the result
	json.NewEncoder(w).Encode(map[string]string{
		"name": name,
		"desc": desc,
	})
}

func getShortDescription(name string, lang string) (string, error) {
	// Get clean escaped apiurl
	apiURL := getWikipediaAPIURL(name, lang)

	// Make request to Wikipedia API
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("error making http request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response
	var wikiResponse wikipediaResponse
	err = json.Unmarshal(body, &wikiResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling json response: %v", err)
	}

	// Check if the short description field is present
	pages := wikiResponse.Query.Pages
	if len(pages) == 0 {
		return "", errors.New("no page found")
	}

	for _, page := range pages {
		if len(page.Extract) > 0 {
			// Extract only the first sentence of the short description so it's just strings
			reg, err := regexp.Compile("<[^>]*>")
			if err != nil {
				return "", fmt.Errorf("error compiling regex: %v", err)
			}
			desc := reg.ReplaceAllString(page.Extract, "")
			desc = strings.Split(desc, ". ")[0]
			return desc, nil
		}
	}

	return "", errors.New("short description not found")
}

func getWikipediaAPIURL(name string, lang string) string {
	// Convert name capital first
	name = strings.Title(strings.ToLower(name))

	// Construct the Wikipedia API query URL with the "extracts" value for the "prop" parameter
	query := url.QueryEscape(name)
	return fmt.Sprintf("https://%s.wikipedia.org/w/api.php?action=query&titles=%s&prop=extracts&exsentences=1&exintro=1&format=json", lang, query)
}
