package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("image/cache.tar")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	// gzf, err := gzip.NewReader(f)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	tarReader := tar.NewReader(f)
	targetName := "manifest.json"
	var m Manifest
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := header.Name
		if name == targetName {
			m = handleManifest(tarReader, header)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			fmt.Println("File: ", name)
			// show the contents
			// io.Copy(os.Stdout, tarReader)
		default:
			fmt.Printf("%s : %c %s %s\n",
				"hmmm?",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
	fmt.Printf("%+v\n", m)
}

type Manifest struct {
	Config   string
	RepoTags []string
	Layers   []string
}

func handleManifest(r *tar.Reader, header *tar.Header) Manifest {
	size := header.Size
	manifestBytes := make([]byte, size)
	_, err := r.Read(manifestBytes)
	if err != nil {
		panic(err)
	}
	var m [1]Manifest
	err = json.Unmarshal(manifestBytes, &m)
	if err != nil {
		panic(err)
	}
	return m[0]
}
