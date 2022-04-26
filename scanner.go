package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Scanner struct {
	Count int
	File *os.File
	Method string
	Exact bool
}

func (s *Scanner) SearchForFunction(dir, ignore string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			//ignore error
			return nil
		}
		ignorePaths := strings.Split(strings.ToLower(ignore),";")
		path = strings.ToLower(path)
		ignorable := false
		for _, loc := range ignorePaths {
			if strings.Contains(path, loc) {
				ignorable = true
			}
		}
		if ignorable {
			return nil
		}
		if info.Mode().IsRegular() {
			if (!strings.Contains(info.Name(),"_test")) && strings.HasSuffix(info.Name(),".go") {
				//fmt.Printf("file name:%s\n",path)
				err, lineno := s.CheckFile(path)
				if err!=nil {
					s.Count += 1
					fmt.Printf("%v\n",err)
					if lineno != -1 {
						s.CaptureMethod(path, lineno)
					}
				}
			}

		}
		return nil
	})
}

func (s *Scanner) CheckFile(path string) (error, int) {

	contents, err := ioutil.ReadFile(path)
	if err!=nil {
		fmt.Printf("Error occurred while reading file:%v",err)
		return err, -1
	}
	source := string(contents)
	lineno := 0
	for _, text := range strings.Split(source,"\n") {
		lineno += 1
		//normalized := strings.ToLower(text)
		if strings.Contains(text,"func") && strings.Contains(text," "+s.Method+"(") {
			if s.Exact {
				f1 := strings.Index(text,"func") + 4
				f2 := strings.Index(text, " init(")
				ignorable := false
				for {
					if text[f1] != ' ' && f1 < f2 {
						ignorable = true
						break
					}
					f1++
					if f1 >= f2 {
						break
					}
				}
				if ignorable {
					return nil, -1
				}
			}

			return fmt.Errorf("The %v method found in : %v [%v] at line => %v\n",s.Method,path,text,lineno), lineno
		}
	}
	return nil, -1
}

func (s *Scanner) CaptureMethod(path string, lineno int) error {

	contents, _  := ioutil.ReadFile(path)
	source := string(contents)
	curline := 0
	methodFound := false
	openBlockCount := 0
	closingBlockCount := 0
	startSingleQuote := false
	startDoubleQuote := false
	startBackQuote := false
	for _, text := range strings.Split(source,"\n") {
		curline +=  1
		if curline == lineno {
			if strings.Contains(text,"func") && strings.Contains(text," "+s.Method+"(") {
				fmt.Fprintf(s.File,"===> %v <===\n\n",path)
				methodFound = true
			}
		}
		if methodFound {
			fmt.Fprintf(s.File,"%v\n",text)
			//openBlockCount += len(strings.Split(text, "{"))
			//closingBlockCount +=  len(strings.Split(text,"}"))
			for _, ch := range text {
				if startSingleQuote && ch == '\'' {
					startSingleQuote = false
				} else if startDoubleQuote && ch == '"' {
					startDoubleQuote = false
				} else if startBackQuote && ch == '`' {
					startBackQuote = false
				} else {
					if ch == '\'' && (!startDoubleQuote && !startBackQuote) {
						startSingleQuote = true
					} else if ch == '"' && (!startSingleQuote && !startBackQuote) {
						startDoubleQuote = true
					} else if ch == '`' && (!startSingleQuote && !startDoubleQuote) {
						startBackQuote = true
					} else if startSingleQuote || startDoubleQuote || startBackQuote {
						continue
					} else {
						if ch == '{' {
							openBlockCount ++
						} else if ch == '}' {
							closingBlockCount ++
						}
					}
				}

			}
			if openBlockCount == closingBlockCount {
				fmt.Fprintf(s.File,"\n===>=======================================================================<====\n\n")
				break
			}
		}
	}
	if !methodFound {
		return fmt.Errorf("%v method not found",s.Method)
	}
	return nil
}
