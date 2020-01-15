package funcs

import (
	"os"
	"bufio"
	"strings"
	"os/exec"
)

func GetCrypt(user string) (string, error) {
	shadow, err := os.Open("/etc/shadow")

	if err != nil {
		return "", err
	}

	defer shadow.Close()

	scanner := bufio.NewScanner(shadow)

	split := []string{}

	for scanner.Scan() {
		split = strings.Split(scanner.Text(), ":")

		if split[0] == user {
			return split[1], scanner.Err()
		}
	}

	return "", nil
}

func MatchCrypt(pass string, crypt string) (bool, error) {
	split := strings.Split(crypt, "$")
	cmd := exec.Command("openssl", "passwd", "-" + split[1], "-salt", split[2], pass)

	out, err := cmd.Output()

	if err != nil {
		return false, err
	}

	if string(out)[0:len(out)-1] == crypt {
		return true, nil
	} else {
		return false, nil
	}
}
