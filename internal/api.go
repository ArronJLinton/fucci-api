package main

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	// returns data as bytes so we can write it in a binary format directly to the http response
// 	data, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Printf("Failed to marshal JSON response: %v", payload)
// 		w.WriteHeader(code)
// 	}

// 	w.Header().Add("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(data)
// }

// func handlerReadiness(w http.ResponseWriter, r *http.Request) {
// 	respondWithJSON(w, http.StatusAccepted, struct{}{})
// }