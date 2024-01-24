package fetch

import (
	"os"
)

func File(path string) ([]byte, error) {
	return os.ReadFile(path)
}
