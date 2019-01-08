package common

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

//向文件追加行
func AppendLineToFile(line bytes.Buffer, filename string) {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintln(w, line.String())
	w.Flush()
}
