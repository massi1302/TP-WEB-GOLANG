package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	PORT               = ":8080"
	TEMPLATE_DIR       = "./Templates/*.html"
	ASSETS_DIR         = "assets"
	DATE_FORMAT        = "2006-01-02"
	ERROR_REDIRECT     = "/erreur"
	ERR_NOM_INVALID    = "Le nom doit contenir uniquement des lettres (1-32 caractères)"
	ERR_PRENOM_INVALID = "Le prénom doit contenir uniquement des lettres (1-32 caractères)"
	ERR_DATE_INVALID   = "La date de naissance est invalide"
	ERR_SEXE_INVALID   = "Le sexe doit être 'M', 'F' ou 'A'"
	ERR_METHOD_INVALID = "Méthode non autorisée"
	ERR_MISSING_FIELDS = "Tous les champs sont obligatoires"
)

// Structures de données
type Form struct {
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
	Check         bool
	Errors        map[string]string
}

type ViewData struct {
	Message string
	Class   string
}

type InfoEtudiants struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string
}

type Promo struct {
	Nom       string
	Filiere   string
	Niveau    string
	Nombre    int
	Etudiants []InfoEtudiants
}

type PageAffiche struct {
	Check         bool
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
	IsEmpty       bool
	Errors        map[string]string
}

// Variables globales
var (
	stockageForm = Form{}
	templates    *template.Template
	viewCounter  int
)

// Validateurs
var (
	nomRegex  = regexp.MustCompile(`^[a-zA-ZÀ-ÿ-]+$`)
	dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	sexeRegex = regexp.MustCompile(`^[MFA]$`)
)

func init() {
	var err error
	templates, err = template.ParseGlob(TEMPLATE_DIR)
	if err != nil {
		log.Fatalf("Erreur lors du chargement des templates: %v", err)
	}
}

// Gestionnaires HTTP
func handlePromo(w http.ResponseWriter, r *http.Request) {
	data := Promo{
		Nom:     "B1 Informatique",
		Filiere: "Informatique",
		Niveau:  "Bachelor 1",
		Nombre:  10,
		Etudiants: []InfoEtudiants{
			{Nom: "AHFIR", Prenom: "Massinissa", Age: 20, Sexe: "M"},
			{Nom: "FONTAINE", Prenom: "Antony", Age: 19, Sexe: "M"},
			{Nom: "JULLEMIER", Prenom: "Jérémie", Age: 20, Sexe: "M"},
			{Nom: "CHECKAL", Prenom: "Abdel", Age: 20, Sexe: "M"},
			{Nom: "KONATE", Prenom: "Azilis", Age: 20, Sexe: "F"},
			{Nom: "WEHBE", Prenom: "Edwin", Age: 20, Sexe: "M"},
			{Nom: "BAGNEAU", Prenom: "Emma", Age: 20, Sexe: "F"},
			{Nom: "BENKIRANE", Prenom: "Yassine", Age: 20, Sexe: "M"},
			{Nom: "AIT", Prenom: "Rania", Age: 20, Sexe: "F"},
			{Nom: "VELAZQUEZ", Prenom: "Léo", Age: 20, Sexe: "M"},
		},
	}

	if err := templates.ExecuteTemplate(w, "promo", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCounter(w http.ResponseWriter, r *http.Request) {
	var data ViewData

	switch {
	case viewCounter == 0:
		data = ViewData{
			Message: fmt.Sprintf("Le compteur démarre : %d", viewCounter),
			Class:   "even",
		}
	case viewCounter%2 == 0:
		data = ViewData{
			Message: fmt.Sprintf("Le nombre de vues est pair : %d", viewCounter),
			Class:   "even",
		}
	default:
		data = ViewData{
			Message: fmt.Sprintf("Le nombre de vues est impair : %d", viewCounter),
			Class:   "odd",
		}
	}

	if err := templates.ExecuteTemplate(w, "change", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	viewCounter++
}

func handleUserForm(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "userform", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleUserTreatment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		redirectError(w, r, "405", ERR_METHOD_INVALID)
		return
	}

	form := Form{
		Nom:           r.FormValue("nom"),
		Prenom:        r.FormValue("prenom"),
		DateNaissance: r.FormValue("dateNaissance"),
		Sexe:          r.FormValue("sexe"),
	}

	// Validate the form
	if err := validateForm(form); err != nil {
		// If there's an error, redirect to the error page
		redirectError(w, r, "400", err.Error())
		return
	}

	// If all validations pass, set Check to true and store the form
	form.Check = true
	stockageForm = form
	http.Redirect(w, r, "/user/display", http.StatusSeeOther)
}

func handleUserDisplay(w http.ResponseWriter, r *http.Request) {
	data := PageAffiche{
		Check:         stockageForm.Check,
		Nom:           stockageForm.Nom,
		Prenom:        stockageForm.Prenom,
		DateNaissance: stockageForm.DateNaissance,
		Sexe:          stockageForm.Sexe,
		IsEmpty:       isFormEmpty(stockageForm),
	}

	if err := templates.ExecuteTemplate(w, "userdisplay", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func isFormEmpty(form Form) bool {
	return form.Nom == "" && form.Prenom == "" && form.DateNaissance == "" && form.Sexe == ""
}

// Fonctions utilitaires
func validateForm(form Form) error {
	if form.Nom != "" && (!nomRegex.MatchString(form.Nom) || len(form.Nom) > 32) {
		return fmt.Errorf(ERR_NOM_INVALID)
	}
	if form.Prenom != "" && (!nomRegex.MatchString(form.Prenom) || len(form.Prenom) > 32) {
		return fmt.Errorf(ERR_PRENOM_INVALID)
	}
	if form.DateNaissance != "" {
		if !dateRegex.MatchString(form.DateNaissance) {
			return fmt.Errorf(ERR_DATE_INVALID)
		}
		_, err := time.Parse(DATE_FORMAT, form.DateNaissance)
		if err != nil {
			return fmt.Errorf(ERR_DATE_INVALID)
		}
	}
	if form.Sexe != "" && !sexeRegex.MatchString(form.Sexe) {
		return fmt.Errorf(ERR_SEXE_INVALID)
	}
	return nil
}

func redirectError(w http.ResponseWriter, r *http.Request, code, message string) {
	http.Redirect(w, r, fmt.Sprintf("%s?code=%s&message=%s",
		ERROR_REDIRECT, code, message), http.StatusMovedPermanently)
}

func setupRoutes() {
	// Serveur de fichiers statiques
	fs := http.FileServer(http.Dir(ASSETS_DIR))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Routes
	http.HandleFunc("/promo", handlePromo)
	http.HandleFunc("/change", handleCounter)
	http.HandleFunc("/user/form", handleUserForm)
	http.HandleFunc("/user/treatement", handleUserTreatment)
	http.HandleFunc("/user/display", handleUserDisplay)
	http.HandleFunc("/erreur", handleError)
}

type ErrorData struct {
	Code    string
	Message string
}

func handleError(w http.ResponseWriter, r *http.Request) {
	data := ErrorData{
		Code:    r.FormValue("code"),
		Message: r.FormValue("message"),
	}

	// Si pas de code ou message, utiliser des valeurs par défaut
	if data.Code == "" {
		data.Code = "500"
	}
	if data.Message == "" {
		data.Message = "Une erreur inattendue s'est produite"
	}

	if err := templates.ExecuteTemplate(w, "error", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	setupRoutes()
	log.Printf("Serveur démarré sur http://localhost%s", PORT)
	log.Fatal(http.ListenAndServe("localhost"+PORT, nil))
}
