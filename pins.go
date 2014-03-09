package gogadgets

//The beaglebone black GPIO pins
var (
	Pins = map[string]map[string]map[string]string{
		"gpio": map[string]map[string]string{
			"8": map[string]string{
				"7":  "66",
				"8":  "67",
				"9":  "69",
				"10": "68",
				"11": "45",
				"12": "44",
				"14": "26",
				"15": "47",
				"16": "46",
				"26": "61",
			},
			"9": map[string]string{
				"12": "60",
				"14": "50",
				"15": "48",
				"16": "51",
			},
		},
	}
)
