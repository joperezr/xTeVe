package src

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// TestAddPPVFilters_AutomatedFilterCleanup tests that automated filters are removed before adding new ones
func TestAddPPVFilters_AutomatedFilterCleanup(t *testing.T) {
	// Create a settings struct with existing filters
	settings := SettingsStruct{
		Filter: make(map[int64]interface{}),
	}

	// Add some existing filters including automated ones
	automatedFilter := FilterStruct{
		Type:   "group-title",
		Name:   "Automated",
		Filter: "Old PPV Category",
		Active: true,
	}

	manualFilter := FilterStruct{
		Type:   "group-title",
		Name:   "Manual",
		Filter: "Manual Category",
		Active: true,
	}

	settings.Filter[1] = jsonToMap(mapToJSON(automatedFilter))
	settings.Filter[2] = jsonToMap(mapToJSON(manualFilter))

	// Call addPPVFilters (it will fail to fetch from server, but will clean up filters)
	// We expect it to remove automated filters
	newSettings, _ := addPPVFilters(settings)

	// Check that manual filter is preserved
	foundManual := false
	foundAutomated := false

	for _, f := range newSettings.Filter {
		var filter FilterStruct
		err := json.Unmarshal([]byte(mapToJSON(f)), &filter)
		if err != nil {
			t.Fatalf("Failed to unmarshal filter: %v", err)
		}

		if filter.Name == "Manual" {
			foundManual = true
		}
		if filter.Name == "Automated" && filter.Filter == "Old PPV Category" {
			foundAutomated = true
		}
	}

	if !foundManual {
		t.Error("Manual filter should be preserved")
	}

	if foundAutomated {
		t.Error("Old automated filter should be removed")
	}
}

// TestDateFormat tests that today's date is formatted correctly as mm/dd
func TestDateFormat(t *testing.T) {
	todayDate := time.Now().Format("01/02")

	// Check format: should be exactly 5 characters with a slash in the middle
	if len(todayDate) != 5 {
		t.Errorf("Date format should be 5 characters, got %d: %s", len(todayDate), todayDate)
	}

	if todayDate[2] != '/' {
		t.Errorf("Date format should have '/' at position 2, got: %s", todayDate)
	}

	// Verify format matches mm/dd pattern
	now := time.Now()
	expectedDate := fmt.Sprintf("%02d/%02d", now.Month(), now.Day())
	if todayDate != expectedDate {
		t.Errorf("Date format mismatch. Got %s, expected %s", todayDate, expectedDate)
	}

	t.Logf("Today's date in mm/dd format: %s", todayDate)
}

// TestFilterStruct_NameField tests that FilterStruct has a Name field
func TestFilterStruct_NameField(t *testing.T) {
	filter := FilterStruct{
		Type:          "group-title",
		Name:          "Automated",
		Filter:        "Test Category",
		Active:        true,
		CaseSensitive: false,
	}

	if filter.Name != "Automated" {
		t.Errorf("FilterStruct Name field should be 'Automated', got: %s", filter.Name)
	}

	// Test JSON marshaling includes Name field
	jsonBytes, err := json.Marshal(filter)
	if err != nil {
		t.Fatalf("Failed to marshal FilterStruct: %v", err)
	}

	jsonStr := string(jsonBytes)
	if jsonStr == "" {
		t.Error("JSON string should not be empty")
	}

	// Unmarshal and verify Name field is preserved
	var unmarshaled FilterStruct
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal FilterStruct: %v", err)
	}

	if unmarshaled.Name != "Automated" {
		t.Errorf("Unmarshaled Name should be 'Automated', got: %s", unmarshaled.Name)
	}
}
