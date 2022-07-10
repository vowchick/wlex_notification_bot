package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
)

//JSONSaving interface that implements saving to file fast and easy
type JSONSaving interface {
	GetFile() string
	SetFile(string)
}

func writeFile(t JSONSaving) error {

	if t.GetFile() == "" {
		t.SetFile(fmt.Sprintf("%d", rand.Intn(1000)))
	}

	jsonBytes, err := json.MarshalIndent(t, "", "    ")

	if err != nil {
		log.Println("cant marshal ", err)
		return err
	}

	bs, err := ioutil.ReadFile(t.GetFile())
	if err != nil {
		log.Println("cant readfile ", err)

		return err
	}

	if !bytes.Equal(bs, jsonBytes) {
		err = ioutil.WriteFile(t.GetFile(), jsonBytes, 0644)
		if err != nil {
			log.Println("cant write ", err)
			return err
		}
	}

	return nil
}

//GetFile implements JsonSave
func (t *Settings) GetFile() string {
	return t.fileName
}

//SetFile implements JsonSave
func (t *Settings) SetFile(filename string) {
	t.fileName = filename
}

func readFile(t JSONSaving, fileName string) error {

	t.SetFile(fileName)
	bs, err := ioutil.ReadFile(t.GetFile())
	if err != nil {
		log.Println("cant readfile  ", err)

		return err
	}

	if len(bs) == 0 {
		return nil
	}
	err = json.Unmarshal(bs, t)
	if err != nil {
		log.Println("cant unmarshal  ", err)
		return err
	}
	return nil
}

func syncFile(t JSONSaving) error {

	err := writeFile(t)
	if err != nil {
		log.Println("cant write ", err)
		return err
	}

	return nil
}
