package gogadgets

//The beaglebone black GPIO pins that are available by default.
//You can use the device tree overlay to get more.
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
		"pwm": map[string]map[string]string{
			"8": map[string]string{
				"13": "bone_pwm_p8_13",
				"19": "bone_pwm_p8_19",
			},
			"9": map[string]string{
				"14": "bone_pwm_p9_14",
				"16": "bone_pwm_p9_16",
				"21": "bone_pwm_p9_21",
				"22": "bone_pwm_p9_22",
			},
		},
	}
	PiPins = map[string]string{
		"11": "17",
		"13": "27",
		"15": "22",
		"16": "23",
		"18": "24",
		"22": "25",
	}
)
