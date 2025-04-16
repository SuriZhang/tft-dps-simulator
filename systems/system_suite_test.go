package systems_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/suriz/tft-dps-simulator/data"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestSystems(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Systems Suite")
}

var _ = ginkgo.BeforeSuite(func() {
    // Load item data once for the entire manager suite
    // Adjust the path to your actual item data file
	dataDir := "../assets"
	fileName := "en_us_14.1b.json"
	filePath := filepath.Join(dataDir, fileName)
	tftData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
	if err != nil {
		fmt.Printf("Error loading set data: %v\n", err)
		os.Exit(1)
	}
	data.InitializeChampions(tftData)

	data.InitializeTraits(tftData)
    data.InitializeSetActiveAugments(tftData, filePath)

	data.InitializeSetActiveItems(tftData, filePath)

    gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to load item data")
})
