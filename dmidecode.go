package shell

import (
	"fmt"

	"github.com/dselans/dmidecode"
)

func Dmidecode() (*dmidecode.DMI, error) {
	dmi := dmidecode.New()

	if err := dmi.Run(); err != nil {
		return dmi, fmt.Errorf("Unable to get dmidecode information. Error: %v\n", err)
	}
	return dmi, nil
}
