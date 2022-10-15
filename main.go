package main

import (
	"github.com/bakyazi/gopass/sheetsapi"
)

func main() {
	db, err := sheetsapi.NewPasswordDB("creds.json", "1u8W_mNDPl8xRsl2YPrqk46I2oOZ_NDnZUpFThRGMo5A")
	if err != nil {
		panic(err)
	}

	err = db.DeletePassword("twitter.com", "berkay15")
	if err != nil {
		panic(err)
	}

}
