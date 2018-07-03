package shell

import "strings"

func Sn() (string, error) {
	//dmidecode -s system-serial-number | grep -v \# | head -1 | sed s/[^a-zA-Z0-9]//g | tr a-z A-Z
	str, err := Out(`dmidecode -s system-serial-number | grep -v \# | head -1 | sed s/[^a-zA-Z0-9]//g`)
	if err != nil {
		return "", err
	}
	str = strings.ToUpper(str)
	return strings.Trim(str, "\n"), err
}
