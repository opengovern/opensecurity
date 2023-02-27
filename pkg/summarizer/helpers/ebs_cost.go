package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	io2Thresholds = []int{32000, 64000}
)

type EBSCostDescription struct {
	Region        string `json:"region"`
	Gp2Size       int    `json:"gp2Size"`
	Gp3Size       int    `json:"gp3Size"`
	Gp3Throughput int    `json:"gp3Throughput"`
	Gp3IOPS       int    `json:"gp3IOPS"`
	Io1Size       int    `json:"io1Size"`
	Io1IOPS       int    `json:"io1IOPS"`
	Io2Size       int    `json:"io2Size"`
	Io2IOPS       int    `json:"io2IOPS"`
	Sc1Size       int    `json:"sc1Size"`
	St1Size       int    `json:"st1Size"`
	StandardSize  int    `json:"standardSize"`
	StandardIOPS  int    `json:"standardIOPS"`
}

func (e EBSCostDescription) GetCost() float64 {
	costManifest, ok := GetEbsCosts()[strings.ToLower(e.Region)]
	if !ok {
		return 0
	}
	total := float64(0)

	// GP2
	total += costManifest.Gp2.PricePerGBMonth.GetFloat64() * float64(e.Gp2Size)

	// GP3
	total += costManifest.Gp3.PricePerGBMonth.GetFloat64() * float64(e.Gp3Size)
	total += costManifest.Gp3.PricePerIOPSMonth.GetFloat64() * float64(e.Gp3IOPS)
	total += costManifest.Gp3.PricePerGiBpsMonth.GetFloat64() * float64(e.Gp3Throughput)

	//Io1
	total += costManifest.Io1.PricePerGBMonth.GetFloat64() * float64(e.Io1Size)
	total += costManifest.Io1.PricePerIOPSMonth.GetFloat64() * float64(e.Io1IOPS)

	//Io2
	total += costManifest.Io2.PricePerGBMonth.GetFloat64() * float64(e.Io2Size)
	switch {
	case e.Io2IOPS <= io2Thresholds[0]:
		total += costManifest.Io2.PricePerTier1IOPSMonth.GetFloat64() * float64(e.Io2IOPS)
	case io2Thresholds[0] < e.Io2IOPS && e.Io2IOPS <= io2Thresholds[1]:
		total += costManifest.Io2.PricePerTier2IOPSMonth.GetFloat64() * float64(e.Io2IOPS)
	case io2Thresholds[1] < e.Io2IOPS:
		total += costManifest.Io2.PricePerTier3IOPSMonth.GetFloat64() * float64(e.Io2IOPS)
	}

	//Sc1
	total += costManifest.Sc1.PricePerGBMonth.GetFloat64() * float64(e.Sc1Size)

	//St1
	total += costManifest.St1.PricePerGBMonth.GetFloat64() * float64(e.St1Size)

	//Standard
	total += costManifest.Standard.PricePerGBMonth.GetFloat64() * float64(e.StandardSize)
	total += costManifest.Standard.PricePerIOs.GetFloat64() * float64(e.StandardIOPS)

	return total
}

type PricePerMonth struct {
	USD string `json:"USD"`
}

func (p PricePerMonth) GetFloat64() float64 {
	res, _ := strconv.ParseFloat(p.USD, 64)
	return res
}

type EbsCost struct {
	Gp2 struct {
		PricePerGBMonth PricePerMonth `json:"pricePerGBMonth"`
	} `json:"gp2,omitempty"`
	Gp3 struct {
		PricePerGBMonth    PricePerMonth `json:"pricePerGBMonth"`
		PricePerGiBpsMonth PricePerMonth `json:"pricePerGiBpsMonth"`
		PricePerIOPSMonth  PricePerMonth `json:"pricePerIOPSMonth"`
	} `json:"gp3,omitempty"`
	Io1 struct {
		PricePerGBMonth   PricePerMonth `json:"pricePerGBMonth"`
		PricePerIOPSMonth PricePerMonth `json:"pricePerIOPSMonth"`
	} `json:"io1,omitempty"`
	Io2 struct {
		PricePerGBMonth        PricePerMonth `json:"pricePerGBMonth"`
		PricePerTier1IOPSMonth PricePerMonth `json:"pricePerTier1IOPSMonth"`
		PricePerTier2IOPSMonth PricePerMonth `json:"pricePerTier2IOPSMonth"`
		PricePerTier3IOPSMonth PricePerMonth `json:"pricePerTier3IOPSMonth"`
	} `json:"io2,omitempty"`
	Sc1 struct {
		PricePerGBMonth PricePerMonth `json:"pricePerGBMonth"`
	} `json:"sc1,omitempty"`
	St1 struct {
		PricePerGBMonth PricePerMonth `json:"pricePerGBMonth"`
	} `json:"st1,omitempty"`
	Standard struct {
		PricePerGBMonth PricePerMonth `json:"pricePerGBMonth"`
		PricePerIOs     PricePerMonth `json:"pricePerIOs"`
	} `json:"standard,omitempty"`
}

type JSONEbsCosts struct {
	EbsPrices EbsCost `json:"ebs_prices"`
	Location  string  `json:"location"`
	Partition string  `json:"partition"`
	RzCode    string  `json:"rzCode"`
	RzType    string  `json:"rzType"`
}

// RegionCode to EBS cost map
var ebsCosts = map[string]EbsCost{}

// Singleton pattern for getting ebsCosts
func GetEbsCosts() map[string]EbsCost {
	if len(ebsCosts) == 0 {
		err := initEbsCosts()
		if err != nil {
			fmt.Printf("Error initializing EBS costs: %v", err)
		}
	}
	return ebsCosts
}

func initEbsCosts() error {
	// read from file
	jsonFile, err := os.Open("/config/ebs-costs.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	costsArr := make([]JSONEbsCosts, 0)
	err = json.Unmarshal(jsonBytes, &costsArr)
	if err != nil {
		return err
	}
	for _, cost := range costsArr {
		ebsCosts[strings.ToLower(cost.RzCode)] = cost.EbsPrices
	}
	return nil
}
