package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	godror "github.com/godror/godror"
	UUID "github.com/google/uuid"
)

var username, password, dataBaseIp, dataBaseName, dataBasePort = "", "", "", "", ""
var dbIdentifier, dbClientInfo, dbModuloGo, dbOperation = "", "", "", ""
var errors []interface{}

func GetEnvDefault(key string) string {
	val, ex := os.LookupEnv(key)
	if !ex || val == "" {
		errors = append(errors, fmt.Sprintf("\nVariable %s no tiene valor establecido.", key))
		return ""
	}
	return val
}

func GetEnvHistory() {

	username = GetEnvDefault("ENV_BAP_GO_DB_USER")
	password = GetEnvDefault("ENV_BAP_GO_DB_PASS")
	dataBaseIp = GetEnvDefault("ENV_BAP_GO_DB_IP")
	dataBasePort = GetEnvDefault("ENV_BAP_GO_DB_PORT")
	dataBaseName = GetEnvDefault("ENV_BAP_GO_DB_NAME") //en el resultado vemos el nombre de esta conexion

	// recupera variables para creacion del contexto para realizar el trace del api que levanta sesion con la bd, asi asignamos el nombre a la conexion
	dbIdentifier = GetEnvDefault("ENV_BAP_GO_DB_IDENTIFIER") //Este es el nombre que vemos en la tabla para identificar
	dbClientInfo = GetEnvDefault("ENV_BAP_GO_DB_CLIENTINFO")
	dbOperation = GetEnvDefault("ENV_BAP_GO_DB_OPERATION")
	dbModuloGo = GetEnvDefault("ENV_BAP_GO_DB_MODULO")

	if len(errors) > 0 {
		log.Fatalln(errors...)
	}
}

func main() {
	GetEnvHistory()           //Validar que las variables de entorno se carguen correctamente
	mux := http.NewServeMux() //estoy usando un multiplexor para mostar en host
	mux.HandleFunc("/", Inicio)
	log.Println("server corriendo...")
	http.ListenAndServe(":4000", mux)
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	//stringConection es nuestra variable con los datos del server para conectar
	stringConection := username + "/" + password + "@" + dataBaseIp + ":" + dataBasePort + "/" + dataBaseName

	//Connection Params
	var newConnParams godror.ConnectionParams
	newConnParams.Username = username
	newConnParams.Password = godror.NewPassword(password)
	newConnParams.ConnectString = dataBaseIp + ":" + dataBasePort + "/" + dataBaseName
	//newConnParams.SetSessionParamOnInit("MACHINE", "'' || MY-LAPTOP-JEJE || ''")
	/*s := [][2]string{
		{"MACHINE", "MY-LAPTOP-eee"},
	}*/
	//newConnParams.AlterSession = s

	log.Println(stringConection, "este")

	// inicializamos el pool y la conexion a oracle
	db, err := sql.Open("godror", stringConection)
	if err != nil {
		log.Fatalln(err)
	}

	//db := sql.OpenDB(godror.NewConnector(newConnParams))

	// configuracion del pool de conexiones
	// numero maximo de conexiones concurrentes abiertas
	db.SetMaxOpenConns(5)
	// numero maximo de conexiones retenidas y reusadas
	db.SetMaxIdleConns(5)
	// tiempo de vida de las conexiones creadas
	//db.SetConnMaxLifetime(1 * time.Minute)

	// se crea contexto para realizar el ping a la base de datos
	ctxPing := godror.ContextWithTraceTag(context.Background(), godror.TraceTag{
		ClientIdentifier: dbIdentifier,
		ClientInfo:       dbClientInfo,
		DbOp:             dbOperation,
		Module:           dbModuloGo,
		Action:           dbOperation,
	})

	// solicitud de PING hacia la base de datos
	if err := db.PingContext(ctxPing); err != nil {
		log.Println(err.Error())
		log.Panic("La DB esta : abajo")
	}
	log.Println("La DB esta : arriba   <----> ", ctxPing) //contexto de conexion

	ctxQuery := godror.ContextWithTraceTag(context.Background(), godror.TraceTag{
		ClientIdentifier: dbIdentifier,
		ClientInfo:       dbClientInfo,
		DbOp:             dbOperation,
		Module:           dbModuloGo,
		Action:           "select",
	})

	consultaDB := "SELECT * FROM canales_de_venta" //ejemplo
	//consultaDB := "GRANT SELECT * from v$session" //ejemplo de muestra de tabla //no lo pongan me salio una lista demaciado larga y puede causar problemas
	//consultaDB := "SELECT sess.username, sess.client_identifier, sess.module, sess.action, area.sql_text FROM v$session sess, v$sqlarea areaarea WHERE sess.sql_address = area.address"

	//for i := 0; i < 5; i++ {

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
		fmt.Println(idUUID, nombre)

		fmt.Fprintf(w, nombre) //en localhost
		fmt.Fprint(w, "\n")
		//fmt.Fprintf(w, nombre)

	}

	//}

}
