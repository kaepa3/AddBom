package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	root := "./"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	done := make(chan interface{})
	path := make(chan string)
	go search(root, done, path, 0)

Loop:
	for {
		select {
		case v := <-path:
			if extentionCheck(v) {
				log.Println(v)
				err := enc(v)
				if err != nil {
					log.Println(err)
				}
			}
		case <-done:
			break Loop
		}
	}
}

func Read(file string) (string, error) {
	// ShiftJISファイルを開く
	sorceFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer sorceFile.Close()
	bytes, err := ioutil.ReadAll(sorceFile)
	if err != nil {
		return "", err
	}
	// ShiftJISのデコーダーを噛ませたReaderを作成する
	return string(bytes), nil
}

var bom = []byte{0xEF, 0xBB, 0xBF}

func Write(file, text string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	// bom
	_, err = f.Write(bom)
	if err != nil {
		return err
	}
	_, err = f.WriteString(text)
	return err
}

func enc(file string) error {
	text, err := Read(file)
	if err != nil {
		return err
	}
	err = Write(file, text)
	if err != nil {
		return err
	}
	return nil

}

func extentionCheck(file string) bool {
	list := []string{".cpp", ".hpp", ".cu"}
	for _, v := range list {
		ex := filepath.Ext(file)
		if ex == v {
			flg, err := checkBom(file)
			if flg && err == nil {
				return true
			}
		}
	}
	return false
}

func checkBom(file string) (bool, error) {
	sorceFile, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer sorceFile.Close()
	buf := make([]byte, 3)
	n, err := sorceFile.Read(buf)
	if n != 3 || err != nil {
		return false, err
	}
	for i := range bom {
		if buf[i] != bom[i] {
			return false, nil
		}
	}

	return true, nil

}

func search(root string, done chan<- interface{}, path chan<- string, nest int) {
	files, _ := ioutil.ReadDir(root)
	for _, v := range files {
		childPath := filepath.Join(root, v.Name())
		if v.IsDir() {
			search(childPath, done, path, nest+1)
		} else {
			path <- childPath
		}
	}
	if nest == 0 {
		close(done)
	}
}
