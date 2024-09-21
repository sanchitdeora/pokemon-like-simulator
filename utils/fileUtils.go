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

	// jsonData, _ := os.ReadFile(filename)
	// slog.Info(string(jsonData))

	// Decode the JSON data
	err = json.NewDecoder(file).Decode(&jsonObj)
	if err != nil {
		slog.Error("error decoding json from file", "filename", filename, "error", err)
		return jsonObj, err
	}
	return jsonObj, nil
}

func WriteJsonToFile[T any](filename string, data T) error {
	// Marshal the data into a JSON byte slice
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("error while marshalling json", "filename", filename, "error", err)
		return err
	}

	// Open the file for writing, overwriting any existing content
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("error while opening file", "filename", filename, "error", err)
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		slog.Error("failed to write JSON to file", "filename", filename, "error", err)
		return err
	}

	return nil
}

func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
