package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// ReadFile lee el contenido de un archivo y devuelve su contenido como una cadena.
func ReadFile(filename string) string {
	content, err := os.ReadFile(filename)

	// Verificar si hay errores al abrir el archivo.
	if err != nil {
		fmt.Println("Failed to open file")
		return ""
	}

	var result api.TransformResult

	if strings.HasSuffix(filename, ".ts") {
		result = api.Transform(string(content), api.TransformOptions{
			Loader: api.LoaderTS,
			TsconfigRaw: `{
				"experimentalDecorators": true,
				"emitDecoratorMetadata": true,
				"allowJs": true,
			}`,
			Format: api.FormatCommonJS,
		})
		if len(result.Errors) != 0 {
			os.Exit(1)
		}
		return string(result.Code)
	}

	result = api.Transform(string(content), api.TransformOptions{
		Loader: api.LoaderJS,
		Format: api.FormatCommonJS,
	})

	return string(result.Code)
}
