package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	temp, tempErr := template.ParseGlob("./Templates/*.html")
	if tempErr != nil {
		fmt.Printf("oups erreur avec le chargement du Template: %s", tempErr.Error())
		os.Exit(02)
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

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		data := Promo{
			Nom:     "B1 Informatique",
			Filiere: "Informatique",
			Niveau:  "Bachelor 1",
			Nombre:  10,
			Etudiants: []InfoEtudiants{
				{Nom: "Massinissa", Prenom: "AHFIR", Age: 20, Sexe: "M"},
				{Nom: "Antony", Prenom: "FONTAINE", Age: 19, Sexe: "M"},
				{Nom: "Jérémie", Prenom: "JULLEMIER", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
				{Nom: "Moussa", Prenom: "KONATE", Age: 20, Sexe: "M"},
			},
		}
		err := temp.ExecuteTemplate(w, "promo", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe("localhost:8080", nil)
}
