# xbee
A package for parsing automatic xbee packets.

This a single-purposed xbee package - when you have an xbee(s) that sleeps, wakes
up, and pushes messages to the coordinator the message take the form of:

    []byte{0x7E, 0x00, 0x16, 0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x03, 0x00, 0x10, 0x02, 0x2F, 0x01, 0xFE, 0x49}

Where

    0x7E                                      - delimiter
    0x00 0x16                                 - length (from after these two bytes until the checksum)
    0x92                                      - frame type
    0x00 0x13 0xA2 0x00 0x40 0x4C 0x0E 0xBE   - long address of the sender
    0x61 0x59                                 - short address of the sender
    0x01                                      - receive options
    0x01                                      - number of samples
    0x00 0x18                                 - digital channel mask
    0x03                                      - analog channel mask
    0x00 0x10                                 - digital samples (not present if digital channel mask == 0)
    0x02 0x2F                                 - first analog sample
    0x01 0xFE                                 - second analog sample
    0x49                                      - checksum

To parse this message:

    func main() {
    	data := []byte{0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x03, 0x00, 0x10, 0x02, 0x2F, 0x01, 0xFE, 0x49}
    	x, err := xbee.NewMessage(data)
    	if err != nil {
    		log.Fatal(err)
    	}

    	a, err := x.GetAnalog()
    	if err != nil {
    		log.Fatal(err)
    	}

    	d, err := x.GetDigital()
    	if err != nil {
    		log.Fatal(err)
    	}

    	for k, v := range a {
    		fmt.Println(k, v)
    	}

    	for k, v := range d {
    		fmt.Println(k, v)
    	}
    }

The output will be:
    
    adc0: 655.7184750733138
    adc1: 598.2404692082112
    dio3: false
    dio4: true


## Setting Up a Pair XBees

Using two Pro S2B Xbees, I was able to get a proof of concept setup working.
You need to set up a coordinator and a router.  For the coordinator load the
"Coordinator API" firmware and apply these settings:

    <?xml version="1.0" encoding="UTF-8"?>

    <data>
      <profile>
        <description_file>XBP24-ZB_21A7_S2B.xml</description_file>
        <settings>
          <setting command="ID">4443</setting>
          <setting command="SC">7FFF</setting>
          <setting command="SD">3</setting>
          <setting command="ZS">0</setting>
          <setting command="NJ">FF</setting>
          <setting command="DH">0</setting>
          <setting command="DL">FFFF</setting>
          <setting command="NI">0x20</setting>
          <setting command="NH">1E</setting>
          <setting command="BH">0</setting>
          <setting command="AR">FF</setting>
          <setting command="DD">30000</setting>
          <setting command="NT">3C</setting>
          <setting command="NO">0</setting>
          <setting command="CR">3</setting>
          <setting command="PL">4</setting>
          <setting command="PM">1</setting>
          <setting command="EE">0</setting>
          <setting command="EO">0</setting>
          <setting command="KY"></setting>
          <setting command="NK"></setting>
          <setting command="BD">3</setting>
          <setting command="NB">0</setting>
          <setting command="SB">0</setting>
          <setting command="D7">1</setting>
          <setting command="D6">0</setting>
          <setting command="AP">1</setting>
          <setting command="AO">0</setting>
          <setting command="SP">7D0</setting>
          <setting command="SN">1</setting>
          <setting command="D0">1</setting>
          <setting command="D1">0</setting>
          <setting command="D2">0</setting>
          <setting command="D3">0</setting>
          <setting command="D4">0</setting>
          <setting command="D5">1</setting>
          <setting command="P0">1</setting>
          <setting command="P1">0</setting>
          <setting command="P2">0</setting>
          <setting command="PR">1FFF</setting>
          <setting command="LT">0</setting>
          <setting command="RP">28</setting>
          <setting command="DO">1</setting>
          <setting command="IR">0</setting>
          <setting command="IC">0</setting>
          <setting command="V+">0</setting>
        </settings>
      </profile>
    </data>

For the router load the "Router API" firmware and apply these settings:

    <?xml version="1.0" encoding="UTF-8"?>

    <data>
      <profile>
        <description_file>XBP24-ZB_23A7_S2B.xml</description_file>
        <settings>
          <setting command="ID">4443</setting>
          <setting command="SC">7FFF</setting>
          <setting command="SD">3</setting>
          <setting command="ZS">0</setting>
          <setting command="NJ">FF</setting>
          <setting command="NW">0</setting>
          <setting command="JV">0</setting>
          <setting command="JN">0</setting>
          <setting command="DH">0</setting>
          <setting command="DL">0</setting>
          <setting command="NI">GRASS</setting>
          <setting command="NH">1E</setting>
          <setting command="BH">0</setting>
          <setting command="AR">FF</setting>
          <setting command="DD">30000</setting>
          <setting command="NT">3C</setting>
          <setting command="NO">0</setting>
          <setting command="CR">3</setting>
          <setting command="PL">4</setting>
          <setting command="PM">1</setting>
          <setting command="EE">0</setting>
          <setting command="EO">0</setting>
          <setting command="KY"></setting>
          <setting command="BD">3</setting>
          <setting command="NB">0</setting>
          <setting command="SB">0</setting>
          <setting command="D7">1</setting>
          <setting command="D6">0</setting>
          <setting command="AP">1</setting>
          <setting command="AO">1</setting>
          <setting command="SM">5</setting>
          <setting command="SN">1</setting>
          <setting command="SO">4</setting>
          <setting command="SP">3E8</setting>
          <setting command="ST">157C</setting>
          <setting command="PO">0</setting>
          <setting command="D0">1</setting>
          <setting command="D1">2</setting>
          <setting command="D2">2</setting>
          <setting command="D3">0</setting>
          <setting command="D4">0</setting>
          <setting command="D5">1</setting>
          <setting command="P0">1</setting>
          <setting command="P1">0</setting>
          <setting command="P2">0</setting>
          <setting command="PR">1FFF</setting>
          <setting command="LT">0</setting>
          <setting command="RP">28</setting>
          <setting command="DO">1</setting>
          <setting command="IR">200</setting>
          <setting command="IC">0</setting>
          <setting command="V+">0</setting>
        </settings>
      </profile>
    </data>

The above xml files can be loaded onto the XBees by using XCTU.  Once you add the radio
module to the user interface you can click on the "Profile" dropdown then click "Load
configuration profile".


## Command line tool

You can watch the values on the command line by:

    cd cmd/xbee
    go install
    xbee /dev/<path to tty serial device your xbee is connected to>
