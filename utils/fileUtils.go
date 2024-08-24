package utils

import (
	"encoding/json"
	"log/slog"
	"os"
)

func ReadJsonFromFile[T any](filename string) (T, error) {
	var jsonObj T
	
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		slog.Error("error while opening file", "filename", filename, "error", err)
		return jsonObj, err
	}
	defer file.Close()

	jsonData, _ := os.ReadFile("pokemon.json")
	slog.Info(string(jsonData))


	// Decode the JSON data
	err = json.NewDecoder(file).Decode(&jsonObj)
	if err != nil {
		slog.Error("error decoding json from file", "filename", filename, "error", err)
		return jsonObj, err
	}
	return jsonObj, nil
}