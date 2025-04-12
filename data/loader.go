package data

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadSetDataFromFile extracts and returns data containing only the object
// with the specified mutator from the JSON file
func LoadSetDataFromFile(filePath string, targetMutator string) (*TFTSetData, error) {
	// Read the file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Parse the JSON data
	var fullData TFTSetData
	if err := json.Unmarshal(file, &fullData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Filter for the target mutator
	var filteredSetData []Set
	var availableMutators []string

	for _, set := range fullData.SetData {
		availableMutators = append(availableMutators, set.Mutator)

		if set.Mutator == targetMutator {
			filteredSetData = append(filteredSetData, set)
		}
	}

	// Check if we found any matching data
	if len(filteredSetData) == 0 {
		// Create a string of available mutators
		var mutatorsStr string
		if len(availableMutators) > 0 {
			mutatorsStr = fmt.Sprintf(" Available mutators: %v", availableMutators)
		}

		return nil, fmt.Errorf("no data with mutator '%s' found.%s",
			targetMutator, mutatorsStr)
	}

	// Return the filtered data
	return &TFTSetData{SetData: filteredSetData}, nil
}
