package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	godror "github.com/godror/godror"
	UUID "github.com/google/uuid"
)

var username, password, dataBaseIp, dataBaseName, dataBasePort = "brandonj", "Daniel22**", "https://apex.oracle.com/pls/apex/", "WKSP_BRANDONJ", ""
var errors []interface{}

func GetEnvDefault(key string) string {
	val, ex := os.LookupEnv(key)
	if !ex || val == "" {
		errors = append(errors, fmt.Sprintf("\nVariable %s no tiene valor establecido.", key))
		return ""
	}
	return val
}

func main() {
	//stringConection es nuestra variable con los datos del server para conectar
	stringConection := username + "/" + password + "@" + dataBaseIp + ":" + dataBasePort + "/" + dataBaseName

	log.Println(stringConection, "este")

	// inicializamos el pool y la conexion a oracle
	db, err := sql.Open("godror", stringConection)
	if err != nil {
		log.Fatalln(err)
	}

	ctxQuery := godror.ContextWithTraceTag(context.Background(), godror.TraceTag{
		Action: "select",
	})

	consultaDB := "SELECT * FROM canales_de_venta" //ejemplo

	//log.Println("Consulta # ", i+1)
	rows, err := db.QueryContext(ctxQuery, consultaDB)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//log.Println("... Parsing query results.")

	var id []byte
	var nombre string
	for rows.Next() {
		rows.Scan(&id, &nombre)

		idUUID, _ := UUID.FromBytes(id) //en consola
		log.Println(idUUID, nombre)

		log.Println(nombre) //en localhost
		log.Println("\n")

	}
}

