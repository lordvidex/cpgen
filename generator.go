package cpgen

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var (
	tmpl *template.Template
)

var (
	//go:embed templates
	fs         embed.FS
	folderName = "solution"
)

func init() {
	var err error
	tmpl, err = template.ParseFS(fs, "**/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
}

// Config contains the parameters used to generate the files to be used
type Config struct {
	FileIO *IO  // input.txt and output.txt files instead of stdin and stdout respectively
	Pq     bool // include priority queue struct and methods
	Uf     bool // include unionfind struct and methods
	Sv     bool // include sieve of eratosthenes
	Cf     bool // codeforces style template with `t` testcases
}

type IO struct {
	Input, Output string
}

// Generate takes the files and their configs to generate the project.
// files should not contain file extension.
// folderName is "solution" by default when empty
func Generate(files []string, c Config, folder string) <-chan float64 {
	ch := make(chan float64)
	if folder == "" {
		folder = folderName
	}
	go func() {
		defer close(ch)
		for i := 0; i < len(files); i++ {
			err := write(files[i], folder, c)
			ch <- float64(i) / float64(len(files))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	return ch
}

func write(file, folderName string, config Config) error {
	fileName := fmt.Sprintf("%s/%s/%s.go", folderName, file, file)
	err := os.MkdirAll(fmt.Sprintf("%s/%s", folderName, file), 0755)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("error closing file", f.Name())
		}
	}(f)
	if config.FileIO != nil {
		input, err := filepath.Abs(filepath.Join(folderName, file, config.FileIO.Input))
		if err != nil {
			return err
		}
		output, err := filepath.Abs(filepath.Join(folderName, file, config.FileIO.Output))
		if err != nil {
			return err
		}
		f1, err := os.Create(input)
		if err != nil {
			return err
		}
		f1.Close()
		f1, err = os.Create(output)
		if err != nil {
			return err
		}
		f1.Close()

		config.FileIO.Input, input = input, config.FileIO.Input
		config.FileIO.Output, output = output, config.FileIO.Output
		defer func() {
			// after template has used config, revert back to filename instead of abs path
			config.FileIO.Input = input
			config.FileIO.Output = output
		}()
	}
	writeTmpl(f, config)
	log.Println("written to file", fileName)
	return nil
}

func writeTmpl(f *os.File, config Config) {
	err := tmpl.ExecuteTemplate(f, "main.go.tmpl", config)
	if err != nil {
		log.Fatal(err)
	}
}
