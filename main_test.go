package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func init_file(root string) {
	fileNames := []string{"hoge.cpp", "hoge.hpp", "hoge.cu", "hoge.bmp"}
	for i, v := range fileNames {
		path := filepath.Join(root, v)
		f, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		if i == 1 {
			var bom = []byte{0xEF, 0xBB, 0xBF}
			f.Write(bom)
		}
		f.WriteString("korekore")
		defer f.Close()
	}
}

func Test_search(t *testing.T) {

	root := "./test"
	if _, err := os.Stat(root); err == nil {
		os.RemoveAll(root)
	}
	err := os.Mkdir(root, 0777)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan interface{})
	file := make(chan string)
	init_file(root)

	type args struct {
		root string
		done chan<- interface{}
		path chan<- string
		nest int
	}
	tests := []struct {
		name  string
		args  args
		count int
	}{
		// TODO: Add test cases.
		{"test", args{root, done, file, 0}, 4},
	}
	for _, tt := range tests {
		count := 0
		t.Run(tt.name, func(t *testing.T) {
			go search(tt.args.root, tt.args.done, tt.args.path, tt.args.nest)

		revRoop:
			for {
				select {
				case <-file:
					count = count + 1
				case <-done:
					break revRoop
				}
			}
		})
		if count != tt.count {
			t.Error(fmt.Printf("count err:%d,%d", count, tt.count))
		}
	}

	os.RemoveAll(root)
}

func Test_extentionCheck(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"cpp test", args{"hoge.cpp"}, true},
		{"hpp test", args{"hoge.hpp"}, true},
		{"cu test", args{"hoge.cu"}, true},
		{"bmp test", args{"hoge.bmp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extentionCheck(tt.args.file); got != tt.want {
				t.Errorf("extentionCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkBom(t *testing.T) {
	root := "checkBom"
	os.Mkdir(root, 0777)
	init_file(root)
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"cpp test", args{filepath.Join(root, "hoge.cpp")}, false, false},
		{"hpp test", args{filepath.Join(root, "hoge.hpp")}, true, false},
		{"cu test", args{filepath.Join(root, "hoge.cu")}, false, false},
		{"bmp test", args{filepath.Join(root, "hoge.bmp")}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkBom(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkBom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkBom() = %v, want %v", got, tt.want)
			}
		})
	}
	os.RemoveAll(root)
}
