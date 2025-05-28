package auto

import (
	"bufio"
	"fmt"
	"os"
)

func enterToContinue() {
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
