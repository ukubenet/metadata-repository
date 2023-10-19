package main

import (
	"net/http"

	parcel "github.com/tshaddix/parcel"
	"github.com/tshaddix/parcel/encoding"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		next.ServeHTTP(w, r)
	})
}

// Encoders/Decoders will be called in the order
// they are registered. The setup below:
// Request ->
// 1. Query Strings
// 2. Json
// 3. Xml
// Response ->
// 1. Json
// 2. Xml
// Notice that the Query codec only provides a decoder,
// so it will not be added to response chain
func CreateFactory() *parcel.Factory {
	factory := parcel.NewFactory()

	factory.Use(encoding.Query())
	factory.Use(encoding.JSON())
	factory.Use(encoding.XML())

	return factory
}

func getParcel(rw http.ResponseWriter, r *http.Request) *parcel.Parcel {
	output := CreateFactory()
	return output.Parcel(rw, r)
}