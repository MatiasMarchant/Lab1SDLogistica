package chat

import (
	//"fmt"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

//Registro es
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

//Cola es
type Cola struct {
	Idpaquete   string
	Seguimiento string
	Tipo        string
	Valor       string
	Intentos    string
	Estado      string
	Origen      string
	Destino     string
}

//PaqueteEnMarcha es
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

//Server es
type Server struct {
	ListaRegistro    []Registro
	Seguimiento      string
	ColaRetail       []Cola
	ColaPrioritario  []Cola
	ColaNormal       []Cola
	PaquetesEnMarcha []PaqueteEnMarcha
}

//ActualizarRegistroPaqueteCamionRetail es
func (s *Server) ActualizarRegistroPaqueteCamionRetail(ctx context.Context, message *PaqueteEnviado) (*PaqueteEnviado, error) {
	// Actualizar PaquetesEnMarcha

	// Buscar indice
	indice := 0
	contador := 0
	for _, elem := range s.PaquetesEnMarcha {
		if message.Idpaquete == elem.Idpaquete {
			indice = contador
		} else {
			contador++
		}
	}

	// Agregar paquete modificado al final
	paqueteModificado := PaqueteEnMarcha{
		Idpaquete:     message.Idpaquete,
		Estado:        message.Estado,
		Idcamion:      message.Idcamion,
		Idseguimiento: message.Seguimiento,
		Intentos:      message.Intentos,
		Origen:        message.Origen,
		Destino:       message.Destino,
		Timestamp:     time.Now(),
	}

	s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, paqueteModificado)

	// Reemplazar indice por el ultimo, y acortar lista
	s.PaquetesEnMarcha[indice] = s.PaquetesEnMarcha[len(s.PaquetesEnMarcha)-1]
	s.PaquetesEnMarcha = s.PaquetesEnMarcha[:len(s.PaquetesEnMarcha)-1]

	// Notificar a financiero
	return message, nil
}

//EntregarPaqueteCamionRetail es
func (s *Server) EntregarPaqueteCamionRetail(ctx context.Context, message *IdCamion) (*ColaPaquete, error) {
	messageColaPaqueteError := ColaPaquete{
		Idpaquete:   "NoPaquetes",
		Seguimiento: "",
		Tipo:        "",
		Valor:       "",
		Intentos:    "",
		Estado:      "",
		Origen:      "",
		Destino:     "",
	}
	if len(s.ColaRetail) < 1 && len(s.ColaPrioritario) < 1 { // NoPaquetes ninguna cola
		return &messageColaPaqueteError, nil
	}

	indicemascaro := 0
	elemmascaro := 0
	iterador := 0
	/*
		fmt.Println("---")
		fmt.Println(len(s.ColaPrioritario))
		fmt.Println("---")
	*/

	if len(s.ColaRetail) >= 1 {
		//fmt.Println("Entre a ColaRetail")
		// Busqueda paquete mas caro
		for _, elem := range s.ColaRetail {
			value, _ := strconv.Atoi(elem.Valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete Retail mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaRetail[indicemascaro].Idpaquete,
			Seguimiento: s.ColaRetail[indicemascaro].Seguimiento,
			Tipo:        s.ColaRetail[indicemascaro].Tipo,
			Valor:       s.ColaRetail[indicemascaro].Valor,
			Intentos:    s.ColaRetail[indicemascaro].Intentos,
			Estado:      s.ColaRetail[indicemascaro].Estado,
			Origen:      s.ColaRetail[indicemascaro].Origen,
			Destino:     s.ColaRetail[indicemascaro].Destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			Idpaquete:     s.ColaRetail[indicemascaro].Idpaquete,
			Estado:        s.ColaRetail[indicemascaro].Estado,
			Idcamion:      message.Idcamion,
			Idseguimiento: s.ColaRetail[indicemascaro].Seguimiento,
			Intentos:      s.ColaRetail[indicemascaro].Intentos,
			Origen:        s.ColaRetail[indicemascaro].Origen,
			Destino:       s.ColaRetail[indicemascaro].Destino,
			Timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaRetail[indicemascaro] = s.ColaRetail[len(s.ColaRetail)-1]
		s.ColaRetail = s.ColaRetail[:len(s.ColaRetail)-1]
		//log.Println(&messageColaPaquete)

		//fmt.Println(s.ColaRetail)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}

	indicemascaro = 0
	elemmascaro = 0
	iterador = 0
	if len(s.ColaPrioritario) >= 1 {
		// Busqueda paquete mas caro
		for _, elem := range s.ColaPrioritario {
			value, _ := strconv.Atoi(elem.Valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete Prioritario mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaPrioritario[indicemascaro].Idpaquete,
			Seguimiento: s.ColaPrioritario[indicemascaro].Seguimiento,
			Tipo:        s.ColaPrioritario[indicemascaro].Tipo,
			Valor:       s.ColaPrioritario[indicemascaro].Valor,
			Intentos:    s.ColaPrioritario[indicemascaro].Intentos,
			Estado:      s.ColaPrioritario[indicemascaro].Estado,
			Origen:      s.ColaPrioritario[indicemascaro].Origen,
			Destino:     s.ColaPrioritario[indicemascaro].Destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			Idpaquete:     s.ColaPrioritario[indicemascaro].Idpaquete,
			Estado:        s.ColaPrioritario[indicemascaro].Estado,
			Idcamion:      message.Idcamion,
			Idseguimiento: s.ColaPrioritario[indicemascaro].Seguimiento,
			Intentos:      s.ColaPrioritario[indicemascaro].Intentos,
			Origen:        s.ColaPrioritario[indicemascaro].Origen,
			Destino:       s.ColaPrioritario[indicemascaro].Destino,
			Timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaPrioritario[indicemascaro] = s.ColaPrioritario[len(s.ColaPrioritario)-1]
		s.ColaPrioritario = s.ColaPrioritario[:len(s.ColaPrioritario)-1]
		//log.Println(&messageColaPaquete)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}
	return &messageColaPaqueteError, nil
}

//ActualizarRegistroPaqueteCamionNormal es
func (s *Server) ActualizarRegistroPaqueteCamionNormal(ctx context.Context, message *PaqueteEnviado) (*PaqueteEnviado, error) {
	// Actualizar PaquetesEnMarcha

	// Buscar indice
	indice := 0
	contador := 0
	for _, elem := range s.PaquetesEnMarcha {
		if message.Idpaquete == elem.Idpaquete {
			indice = contador
		} else {
			contador++
		}
	}

	// Agregar paquete modificado al final
	paqueteModificado := PaqueteEnMarcha{
		Idpaquete:     message.Idpaquete,
		Estado:        message.Estado,
		Idcamion:      message.Idcamion,
		Idseguimiento: message.Seguimiento,
		Intentos:      message.Intentos,
		Origen:        message.Origen,
		Destino:       message.Destino,
		Timestamp:     time.Now(),
	}

	s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, paqueteModificado)

	// Reemplazar indice por el ultimo, y acortar lista
	s.PaquetesEnMarcha[indice] = s.PaquetesEnMarcha[len(s.PaquetesEnMarcha)-1]
	s.PaquetesEnMarcha = s.PaquetesEnMarcha[:len(s.PaquetesEnMarcha)-1]

	// Notificar a financiero
	return message, nil
}

//EntregarPaqueteCamionNormal es
func (s *Server) EntregarPaqueteCamionNormal(ctx context.Context, message *IdCamion) (*ColaPaquete, error) {
	messageColaPaqueteError := ColaPaquete{
		Idpaquete:   "NoPaquetes",
		Seguimiento: "",
		Tipo:        "",
		Valor:       "",
		Intentos:    "",
		Estado:      "",
		Origen:      "",
		Destino:     "",
	}
	if len(s.ColaPrioritario) < 1 && len(s.ColaNormal) < 1 { // NoPaquetes ninguna cola
		return &messageColaPaqueteError, nil
	}

	indicemascaro := 0
	elemmascaro := 0
	iterador := 0
	/*
		fmt.Println("---")
		fmt.Println(len(s.ColaPrioritario))
		fmt.Println("---")
	*/

	if len(s.ColaPrioritario) >= 1 {
		//fmt.Println("Entre a ColaPrioritario")
		// Busqueda paquete mas caro
		for _, elem := range s.ColaPrioritario {
			value, _ := strconv.Atoi(elem.Valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete prioritario mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaPrioritario[indicemascaro].Idpaquete,
			Seguimiento: s.ColaPrioritario[indicemascaro].Seguimiento,
			Tipo:        s.ColaPrioritario[indicemascaro].Tipo,
			Valor:       s.ColaPrioritario[indicemascaro].Valor,
			Intentos:    s.ColaPrioritario[indicemascaro].Intentos,
			Estado:      s.ColaPrioritario[indicemascaro].Estado,
			Origen:      s.ColaPrioritario[indicemascaro].Origen,
			Destino:     s.ColaPrioritario[indicemascaro].Destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			Idpaquete:     s.ColaPrioritario[indicemascaro].Idpaquete,
			Estado:        s.ColaPrioritario[indicemascaro].Estado,
			Idcamion:      message.Idcamion,
			Idseguimiento: s.ColaPrioritario[indicemascaro].Seguimiento,
			Intentos:      s.ColaPrioritario[indicemascaro].Intentos,
			Origen:        s.ColaPrioritario[indicemascaro].Origen,
			Destino:       s.ColaPrioritario[indicemascaro].Destino,
			Timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaPrioritario[indicemascaro] = s.ColaPrioritario[len(s.ColaPrioritario)-1]
		s.ColaPrioritario = s.ColaPrioritario[:len(s.ColaPrioritario)-1]
		//log.Println(&messageColaPaquete)

		//fmt.Println(s.ColaPrioritario)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}

	indicemascaro = 0
	elemmascaro = 0
	iterador = 0
	if len(s.ColaNormal) >= 1 {
		// Busqueda paquete mas caro
		for _, elem := range s.ColaNormal {
			value, _ := strconv.Atoi(elem.Valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete normal mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaNormal[indicemascaro].Idpaquete,
			Seguimiento: s.ColaNormal[indicemascaro].Seguimiento,
			Tipo:        s.ColaNormal[indicemascaro].Tipo,
			Valor:       s.ColaNormal[indicemascaro].Valor,
			Intentos:    s.ColaNormal[indicemascaro].Intentos,
			Estado:      s.ColaNormal[indicemascaro].Estado,
			Origen:      s.ColaNormal[indicemascaro].Origen,
			Destino:     s.ColaNormal[indicemascaro].Destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			Idpaquete:     s.ColaNormal[indicemascaro].Idpaquete,
			Estado:        s.ColaNormal[indicemascaro].Estado,
			Idcamion:      message.Idcamion,
			Idseguimiento: s.ColaNormal[indicemascaro].Seguimiento,
			Intentos:      s.ColaNormal[indicemascaro].Intentos,
			Origen:        s.ColaNormal[indicemascaro].Origen,
			Destino:       s.ColaNormal[indicemascaro].Destino,
			Timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaNormal[indicemascaro] = s.ColaNormal[len(s.ColaNormal)-1]
		s.ColaNormal = s.ColaNormal[:len(s.ColaNormal)-1]
		//log.Println(&messageColaPaquete)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}
	return &messageColaPaqueteError, nil
}

//RecibirOrdenPymes es
func (s *Server) RecibirOrdenPymes(ctx context.Context, message *Ordenclientepymes) (*Ordenseguimiento, error) {
	var tipo string
	messageOrdenseguimiento := Ordenseguimiento{
		Nordenseguimiento: s.Seguimiento,
	}
	if message.Prioritario == true {
		tipo = "prioritario"
	} else {
		tipo = "normal"
	}
	nuevaEntrada := Registro{
		Timestamp:   time.Now(),
		Idpaquete:   message.Id,
		Tipo:        tipo,
		Nombre:      message.Producto,
		Valor:       strconv.Itoa(int(message.Valor)),
		Origen:      message.Tienda,
		Destino:     message.Destino,
		Seguimiento: s.Seguimiento,
	}

	nuevaCola := Cola{
		Idpaquete:   message.Id,
		Seguimiento: s.Seguimiento,
		Tipo:        tipo,
		Valor:       strconv.Itoa(int(message.Valor)),
		Intentos:    "0",
		Estado:      "En bodega",
		Origen:      message.Tienda,
		Destino:     message.Destino,
	}
	// Agregar paquete a una de las tres colas de Server
	if tipo == "prioritario" {
		// Add cola prioritario
		s.ColaPrioritario = append(s.ColaPrioritario, nuevaCola)
	} else {
		// Add cola normal
		s.ColaNormal = append(s.ColaNormal, nuevaCola)
	}

	// Actualizacion codigo de seguimiento (atributo Server)
	numerosiguiente, _ := strconv.Atoi(s.Seguimiento)
	numerosiguiente++
	s.Seguimiento = strconv.Itoa(numerosiguiente)
	s.ListaRegistro = append(s.ListaRegistro, nuevaEntrada)

	//log.Println(s.ColaNormal)
	//log.Println(s.ListaRegistro)
	return &messageOrdenseguimiento, nil
}

//RedecirOrdenPymes es
func (s *Server) RedecirOrdenPymes(ctx context.Context, message *Ordenclientepymes) (*Ordenclientepymes, error) {
	log.Println("Recibi la orden")
	log.Printf("Producto pymes: %s", message.Producto)
	return message, nil
}

//RecibirOrdenRetail es
func (s *Server) RecibirOrdenRetail(ctx context.Context, message *Ordenclienteretail) (*Ordenseguimiento, error) {
	messageOrdenseguimiento := Ordenseguimiento{
		Nordenseguimiento: "9999",
	}

	nuevaEntrada := Registro{
		Timestamp:   time.Now(),
		Idpaquete:   message.Id,
		Tipo:        "retail",
		Nombre:      message.Producto,
		Valor:       strconv.Itoa(int(message.Valor)),
		Origen:      message.Tienda,
		Destino:     message.Destino,
		Seguimiento: "9999",
	}

	nuevaCola := Cola{
		Idpaquete:   message.Id,
		Seguimiento: "9999",
		Tipo:        "retail",
		Valor:       strconv.Itoa(int(message.Valor)),
		Intentos:    "0",
		Estado:      "En bodega",
		Origen:      message.Tienda,
		Destino:     message.Destino,
	}

	// Agregar paquete a una de las tres colas de Server
	s.ColaRetail = append(s.ColaRetail, nuevaCola)
	//log.Println(s.ColaRetail)

	s.ListaRegistro = append(s.ListaRegistro, nuevaEntrada)
	//log.Println(s.ListaRegistro)
	return &messageOrdenseguimiento, nil
}

//RedecirOrdenRetail es
func (s *Server) RedecirOrdenRetail(ctx context.Context, message *Ordenclienteretail) (*Ordenclienteretail, error) {
	log.Println("Recibi la orden retail")
	log.Printf("Producto retail: %s", message.Producto)
	return message, nil
}

//CodigoSeguimiento es
func (s *Server) CodigoSeguimiento(ctx context.Context, message *Ordenseguimiento) (*Estado, error) {
	mensajeError := Estado{
		Estado: "En bodega",
	}
	for _, elem := range s.PaquetesEnMarcha {
		if message.Nordenseguimiento == elem.Idseguimiento {
			mensajeFinal := Estado{
				Estado: elem.Estado,
			}
			return &mensajeFinal, nil
		}
	}

	return &mensajeError, nil
}
