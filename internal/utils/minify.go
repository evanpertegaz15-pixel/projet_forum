package utils

import (
    "os"
    "github.com/tdewolff/minify/v2"
    "github.com/tdewolff/minify/v2/css"
)

func MinifyCSSFile(inputPath, outputPath string) error {
    m := minify.New()
    m.AddFunc("text/css", css.Minify)
    input, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }
    output, err := m.Bytes("text/css", input)
    if err != nil {
        return err
    }
    return os.WriteFile(outputPath, output, 0644)
}