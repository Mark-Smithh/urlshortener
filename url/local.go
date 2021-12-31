package url

import (
	"bufio"
	"fmt"
	"os"
)

//Local run local instance
func Local() {
	for {
		txt := takeInput("Enter Text:")
		i := Input{URL: txt}
		val512, _ := hash512(c, i)
		val := val512.([64]byte)
		fmt.Printf("Hash512:\t%x\n\n", val[:4])
	}
}

func takeInput(msg string) string {
	fmt.Println(msg)
	reader := bufio.NewReader(os.Stdin)
	textVal, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("got err %s", err.Error())
	}
	return textVal[:len(textVal)-1] //remove carriage return from string
}
