package simulation_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"tft-dps-simulator/internal/core/data"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// TestSimulation is the entry point for the Ginkgo test runner in this package.
func TestSimulation(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail) // Connect Ginkgo's Fail function to Go's testing t.Fail
    ginkgo.RunSpecs(t, "Simulation Suite") // Run all specs in the package
}

// Optional: Add BeforeSuite/AfterSuite if needed for package-level setup/teardown
var _ = ginkgo.BeforeSuite(func() {
	// Load item data once for the entire manager suite
	// Adjust the path to your actual item data file
	dataDir := "../../../assets"
	fileName := "en_us_14.1b.json"
	filePath := filepath.Join(dataDir, fileName)
	tftData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
	if err != nil {
		log.Printf("Error loading set data: %v\n", err)
		os.Exit(1)
	}
	data.InitializeChampions(tftData)

	data.InitializeTraits(tftData)
	data.InitializeSetActiveAugments(tftData, filePath)

	data.InitializeSetActiveItems(tftData, filePath)

	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to load item data")
})

// var _ = AfterSuite(func() {
// 	// e.g., Clean up global resources
// })