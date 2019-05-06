package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//Recipe it's a type that represent a recipe structure
type Recipe struct {
	ID   string `json:"idrecipe"`
	Name string `json:"namerp"`
	Proc string `json:"process"`
}

//-----------------------------FUNCTIONS THAT MODIFIES THE DATABASE-----------------------------

//Function that insert a row in the table recipes from the database recipes_db
func insertDB(id string, name string, process string, database *sql.DB) {

	sql := `INSERT INTO recipes VALUES (` + id + `,'` + name + `', '` + process + `');`

	if _, err := database.Exec(sql); err != nil {
		log.Fatal(err)
	}
}

//Function that delete a row in the table recipes from the database recipes_db
func deleteDB(id string, database *sql.DB) {

	sql := `DELETE FROM recipes WHERE idrecipe = ` + id + `;`

	if _, err := database.Exec(sql); err != nil {
		log.Fatal(err)
	}
}

//Function that modifies a row in the table recipes from the database recipes_db
func updateDB(id string, name string, process string, database *sql.DB) {

	sql := `UPDATE recipes SET namerp = '` + name + `', process = '` + process + `' WHERE idrecipe = ` + id + `;`

	if _, err := database.Exec(sql); err != nil {
		log.Fatal(err)
	}
}

//Function that select all the rows of the table recipes in the database recipes_db
func selectAllDB(database *sql.DB) ([]Recipe, error) {

	sql := `SELECT * FROM recipes;`

	rows, err := database.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	recipes := []Recipe{}

	for rows.Next() {
		var r Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.Proc); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil

}

//------------------------------------END PONIT FUNCTIONS------------------------------------

//function that returns all the recipes in recipes_db to thw browser
func getRecipes(res http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/recipes_db?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	re, er := selectAllDB(db)

	if er != nil || len(re) == 0 {
		http.Error(res, "Error en la consulta o no hay ningun objeto", 4)
	}

	json.NewEncoder(res).Encode(re)

}

//THis function filters the database to find a recipe given the id or the name
func getOneRecipe(res http.ResponseWriter, req *http.Request) {

	param := mux.Vars(req)

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/recipes_db?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	re, er := selectAllDB(db)

	if er != nil || len(re) == 0 {
		http.Error(res, "Error en la consulta o no hay ningun objeto", 4)
	} else {
		for _, item := range re {
			if item.ID == param["idrecipe"] || item.Name == param["idrecipe"] {
				json.NewEncoder(res).Encode(item)
				return
			}
		}
		json.NewEncoder(res).Encode(&Recipe{})
	}

}

//This function creates in the database a recipe if this don't exists
func createRecipe(res http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/recipes_db?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	param := mux.Vars(req)
	var rec Recipe

	re, er := selectAllDB(db)

	if er != nil || len(re) == 0 {
		http.Error(res, "Error en la consulta o no hay ningun objeto", 4)
	} else {
		for _, item := range re {
			if item.ID == param["idrecipe"] {
				http.Error(res, "Ya existe el objeto con identificador: "+param["idrecipe"], 5)
				return
			}
		}
		_ = json.NewDecoder(req.Body).Decode(&rec)
		rec.ID = param["idrecipe"]
		insertDB(rec.ID, rec.Name, rec.Proc, db)
		json.NewEncoder(res).Encode(rec)
	}

}

//This function updates an existing recipe in the database
func updateRecipe(res http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/recipes_db?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	param := mux.Vars(req)
	var rec Recipe

	re, er := selectAllDB(db)

	if er != nil || len(re) == 0 {
		http.Error(res, "Error en la consulta o no hay ningun objeto", 4)
	} else {
		for _, item := range re {
			if item.ID == param["idrecipe"] {
				_ = json.NewDecoder(req.Body).Decode(&rec)
				updateDB(param["idrecipe"], rec.Name, rec.Proc, db)
				json.NewEncoder(res).Encode(re)
				return
			}
		}
		http.Error(res, "No existe el objeto con identificador: "+param["idrecipe"], 5)
	}
}

//This function deletes an existing recipe in the database
func deleteRecipe(res http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/recipes_db?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	param := mux.Vars(req)

	re, er := selectAllDB(db)

	if er != nil || len(re) == 0 {
		http.Error(res, "Error en la consulta", 4)
	} else {
		for _, item := range re {
			if item.ID == param["idrecipe"] {
				deleteDB(param["idrecipe"], db)
				json.NewEncoder(res).Encode(re)
				return
			}
		}
		http.Error(res, "No existe el objeto con identificador: "+param["idrecipe"], 5)
	}

}

//---------------------------------------MAIN FUNCTION---------------------------------------

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/recipes", getRecipes).Methods("GET")
	router.HandleFunc("/recipe/{idrecipe}", getOneRecipe).Methods("GET", "OPTIONS")
	router.HandleFunc("/recipe/{idrecipe}", createRecipe).Methods("POST")
	router.HandleFunc("/recipe/{idrecipe}", updateRecipe).Methods("PUT")
	router.HandleFunc("/recipe/{idrecipe}", deleteRecipe).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))

}
