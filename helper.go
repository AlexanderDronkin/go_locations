package main

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/joho/godotenv"
)

// go to global environment
func env(name string, def string) string {
	godotenv.Load()
	value, exist := os.LookupEnv(name)
	if exist {
		return value
	}
	return def
}

// go obj to json
func toJson(obj jsonWriter) []byte {
	result, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	return result
}

// go collect garbage
func fgc(dname string, ttl time.Duration) {
	dir, err := os.ReadDir(DIR + dname)
	if err != nil {
		panic(dname + " not found")
	}
	for _, file := range dir {
		if file.IsDir() {
			fgc(dname+"/"+file.Name(), ttl)
		} else {
			now := time.Now()
			info, _ := file.Info()
			if now.Sub(info.ModTime()) > ttl {
				os.Remove(DIR + dname + "/" + file.Name())
			}
		}
	}
}

// go load file
func fget(fname string) []byte {
	content, err := os.ReadFile(DIR + fname)
	if err != nil {
		panic(fname + " not found")
	}
	return []byte(content)
}

// go write file
func fput(fname string, value []byte) {
	file, err := os.Create(DIR + fname)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(value)
}

// go make files smaller
func fput_br(fname string, value []byte) {
	if len(value) < 256*1024 {
		return
	}

	file, err := os.Create(DIR + fname + ".br")
	if err != nil {
		return
	}
	defer file.Close()

	b := bytes.Buffer{}
	bw := brotli.NewWriter(&b)
	bw.Write(value)
	bw.Close()
	file.Write(b.Bytes())
	file.Close()
}
