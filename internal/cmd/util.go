package cmd

import (
	"encoding/json"
	"fmt"
	"io"
)

func jsonPrint(w io.Writer, v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Fprintln(w, string(data))
}
