package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func openFile(t *testing.T, file string, receipt *Receipt) {
	jsonFile, err := os.Open(file)
	if err != nil {
		t.Error("error opening file " + file)
		t.Fail()
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Error("error reading file" + file)
		t.Fail()
	}

	err = json.Unmarshal(bytes, &receipt)
	if err != nil {
		t.Error("error decoding JSON target.json")
		t.Fail()
	}
}

func TestCalculateTargetPoints(t *testing.T) {
	var receipt Receipt
	openFile(t, "testjsons/target.json", &receipt)

	expected := 28
	actual := calculatePoints(receipt)
	if expected != actual {
		t.Errorf("expected %d, got %d", expected, actual)
		t.Fail()
	}
}

func TestCalculateCornerMartPoints(t *testing.T) {
	var receipt Receipt
	openFile(t, "testjsons/cornermart.json", &receipt)

	expected := 109
	actual := calculatePoints(receipt)

	if expected != actual {
		t.Errorf("expected %d, got %d", expected, actual)
		t.Fail()
	}
}

func TestCalculateWalmartPoints(t *testing.T) {
	var receipt Receipt
	openFile(t, "testjsons/walmart.json", &receipt)

	expected := 53
	actual := calculatePoints(receipt)

	if expected != actual {
		t.Errorf("expected %d, got %d", expected, actual)
		t.Fail()
	}
}
