package main

import (
	"encoding/csv"
	"github.com/MatiasMarchant/Prueba1/tree/master/chat"
	"log"
	"net"
	"os"
	"time"

	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"strconv"
)

//----------------------------------------------------------------
// Struct Entrega, contiene la información que se envía a Finanzas
// mediante RabbitMQ
type Entrega struct {
	Id_paquete string
	Tipo       string
	Valor      int
	Intentos   int
	Estado     string
}

//----------------------------------------------------------------
// Struct que contiene la información de chat.Server con los
// paquetes en marcha.
type PaqueteEnMarcha struct {
	Idpaquete     string
	Estado        string
	Idcamion      string
	Idseguimiento string
	Intentos      string
	Origen        string
	Destino       string
	Timestamp     time.Time
}

//----------------------------------------------------------------
// Struct que contiene el registro de información de chat.Server
type Registro struct {
	Timestamp   time.Time
	Idpaquete   string
	Tipo        string
	Nombre      string
	Valor       string
	Origen      string
	Destino     string
	Seguimiento string
}

//----------------------------------------------------------------
// Retorna true si se encuentra el string id en el arreglo arr
func InArr(id string, arr []string) bool {
	for _, i := range arr {
		if i == id {
			return true
		}
	}
	return false
}

//----------------------------------------------------------------
// Recibe un idpaquete y la listaRegistro, retorna el tipo y el valor
// del paquete con el id idpaquete.
func tipoYvalor(idpaquete string,
	listaRegistro []chat.Registro,
) (string, int) {

	var tipo string
	var valor int
	for _, registro := range listaRegistro {
		if registro.Idpaquete == idpaquete {
			tipo = registro.Tipo
			v, _ := strconv.Atoi(registro.Valor)
			valor = v
			return tipo, valor
		}
	}

	return tipo, valor
}

//----------------------------------------------------------------
// Recibe los paquetesProcesados, paqueteEnMarcha y la listaRegistro. Se filtra paqueteEnMarcha
// mediante si el estado es Recibido o No recibido y si no se ha procesado previamente, es decir,
// es primera vez que se lee. Luego se procesa el paquete y se obtiene tipo y valor de listaRegistro.
// Un paquete es procesado al crear un struct Entrega con esos datos.
func procesarEntregas(paquetesProcesados []string,
	paqueteEnMarcha []chat.PaqueteEnMarcha,
	listaRegistro []chat.Registro) ([]string, []Entrega) {

	var entregasProcesadas []Entrega

	for _, Paquete := range paqueteEnMarcha {

		IntIntentos, _ := strconv.Atoi(Paquete.Intentos)
		print(IntIntentos)
		if (Paquete.Estado == "Recibido" || Paquete.Estado == "No Recibido") && !InArr(Paquete.Idpaquete, paquetesProcesados) {

			paquetesProcesados = append(paquetesProcesados, Paquete.Idpaquete)

			var tipo string
			var valor int
			tipo, valor = tipoYvalor(Paquete.Idpaquete, listaRegistro)

			ent := &Entrega{Id_paquete: Paquete.Idpaquete,
				Tipo:     tipo,
				Valor:    valor,
				Intentos: IntIntentos,
				Estado:   Paquete.Estado}

			entregasProcesadas = append(entregasProcesadas, *ent)
		}
	}
	return paquetesProcesados, entregasProcesadas
}

//----------------------------------------------------------------
// Recibe las entregasProcesadas para ser enviadas a Finanzas mediante RabbitMQ
func enviarRabbit(entregasProcesadas []Entrega) {
	// Iniciar conexion con Finanzas
	conn, err := amqp.Dial("amqp://mqadmin:mqadminpassword@10.6.40.180:5672/")
	if err != nil {
		fmt.Println("Falla inicializando conección")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	fmt.Println(q)
	if err != nil {
		fmt.Println(err)
	}

	// Se itera para enviar cada entrega
	for _, entrega := range entregasProcesadas {
		// Se pasa a Json la entrega, para ser enviada
		js, err := json.Marshal(entrega)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        js,
			},
		)

		if err != nil {
			fmt.Println(err)
		}
	}
}

//----------------------------------------------------------------
func main() {
	// Conexion gRPC
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Falle al escuchar puerto 9000: %v", err)
	}

	var listaRegistro []chat.Registro
	s := chat.Server{}
	s.ListaRegistro = listaRegistro
	s.Seguimiento = "0"

	// Servidor gRPC
	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)

	var paquetesProcesados []string
	var entregasProcesadas []Entrega

	fmt.Println("Corriendo el sistema logístico...\n")

	go func() {
		for {
			// Cada 2 segundos se consulta son el servidor (s), y se obtiene
			// PaquetesEnMarcha y ListaRegistro, para ser procesados.
			time.Sleep(2 * time.Second)

			paquetesProcesados, entregasProcesadas = procesarEntregas(paquetesProcesados,
				s.PaquetesEnMarcha,
				s.ListaRegistro)

			// Se envía a Finanzas las entregas procesadas
			enviarRabbit(entregasProcesadas)
		}
	}()

	go func() {
		csvfile, err := os.Create("registrologistica.csv")
		if err != nil {
			log.Fatalf("Error creando registrologistica.csv")
		}
		csvwriter := csv.NewWriter(csvfile)
		defer csvwriter.Flush()

		for {
			csvfile.Truncate(0)
			csvfile.Seek(0, 0)
			primeralinea := []string{
				"timestamp",
				"id-paquete",
				"tipo",
				"nombre",
				"valor",
				"origen",
				"destino",
				"seguimiento",
			}
			csvwriter.Write(primeralinea)
			csvwriter.Flush()
			time.Sleep(3 * time.Second)
			for _, elem := range s.ListaRegistro {
				elem := []string{
					elem.Timestamp.String(),
					elem.Idpaquete,
					elem.Tipo,
					elem.Nombre,
					elem.Valor,
					elem.Origen,
					elem.Destino,
					elem.Seguimiento,
				}
				csvwriter.Write(elem)
				csvwriter.Flush()
			}
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falle siendo un servidor gRPC en el puerto 9000: %v", err)
	}

}
