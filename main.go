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

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

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
				{Nom: "AHFIR", Prenom: "Massinissa", Age: 20, Sexe: "M"},
				{Nom: "FONTAINE", Prenom: "Antony", Age: 19, Sexe: "M"},
				{Nom: "JULLEMIER", Prenom: "Jérémie", Age: 20, Sexe: "M"},
				{Nom: "CHECKAL", Prenom: "Abdel", Age: 20, Sexe: "M"},
				{Nom: "Azilis", Prenom: "KONATE", Age: 20, Sexe: "F"},
				{Nom: "WEHBE", Prenom: "Edwin", Age: 20, Sexe: "M"},
				{Nom: "BAGNEAU", Prenom: "Emma", Age: 20, Sexe: "F"},
				{Nom: "BENKIRANE", Prenom: "Yassine", Age: 20, Sexe: "M"},
				{Nom: "AIT", Prenom: "Rania", Age: 20, Sexe: "F"},
				{Nom: "VELAZQUEZ", Prenom: "Léo", Age: 20, Sexe: "M"},
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
