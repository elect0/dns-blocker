package dns

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadBlocklist(path string) (map[string]struct{}, error) {
	blocklist := make(map[string]struct{})

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open blocklist file %s, %w", path, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		blocklist[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading blocklist file %s, %w", path, err)
	}

	return blocklist, nil
}
