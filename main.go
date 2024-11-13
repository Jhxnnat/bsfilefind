package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type File struct {
	name  string
	path  string
	isDir bool
}

func clock(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func new_file_arr() File {
	return File{
		name:  "root",
		path:  "./",
		isDir: true,
	}
}

func read_path(path string, files *[]File) error {
	defer clock(time.Now(), "read_path")
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if d.IsDir() {
			fmt.Println(">> Dir : ", d.Name())
		} else {
			fmt.Println(">> File : ", d.Name())
		}
		*files = append(*files, File{d.Name(), path, d.IsDir()})

		return nil
	})
	if err != nil {
		fmt.Printf("read_path: %s", err)
		return err
	}

	return nil
}

func swap(x *File, y *File) {
	temp := *x
	*x = *y
	*y = temp
}

func partition(files []File, low int, hight int) int {
	var pivot string = files[hight].name
	var i int = low - 1

	for j := low; j < hight; j++ {
		if files[j].name < pivot {
			i++
			swap(&files[i], &files[j])
		}
	}

	swap(&files[i+1], &files[hight])
	return i + 1
}

func f_quick_sort(files *[]File, low int, hight int) {
	if low < hight {
		var pi int = partition(*files, low, hight)

		f_quick_sort(files, low, pi-1)
		f_quick_sort(files, pi+1, hight)
	}
}

func bs_left(files []File, query string) int {
	defer clock(time.Now(), "bs_left")
	var low int = 0
	var hight int = len(files)

	for low < hight {
		mid := (low + hight) / 2
		if files[mid].name < query {
			low = mid + 1
		} else {
			hight = mid
		}
	}

	return low
}

func lowcap(str string) string {
	return strings.ToLower(str)
}

func match_str(target string, query string) int {
	score := 0

	if strings.HasPrefix(target, query) {
		score++
	}
	if strings.HasSuffix(target, query) {
		score++
	}
	if strings.Contains(target, query) {
		score++
	}
	matched, _ := regexp.MatchString(query, target)
	if matched {
		score++
	}

	return score
}

func s_query(files []File, query string) []int { //TODO: catch sniky queries like '.', ';', etc, etc
	defer clock(time.Now(), "s_query")
	var index int = bs_left(files, lowcap(query))
	var matches []int
	var scores []int

	if index < len(files) && lowcap(files[index].name) == lowcap(query) {
		matches = append(matches, index)
		scores = append(scores, 4)
		index++
	}

	for index < len(files) { //TODO: order matched results by score
		match_score := match_str(lowcap(files[index].name), lowcap(query))
		if match_score > 0 {
			matches = append(matches, index)
			scores = append(scores, match_score)
		}
		index++
	}

	return matches
}

func main() {
	///home/jhxnnat/Uni/semestre-6/algoritmos
	var root string = "/home/jhxnnat/Uni/semestre-6/algoritmos"
	var files []File
	err := read_path(root, &files)
	f_len := len(files)
	start := time.Now()
	f_quick_sort(&files, 0, f_len-1)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("main read path err: %s", err)
	}

	for _, file := range files {
		fmt.Printf("    >> %s\n", file.name)
	}
	log.Printf("quick_sort took %s", elapsed)

	var input string
	for {
		fmt.Print("\nquery: ")
		fmt.Scan(&input)

		if input == "-e" {
			return
		}

		matches := s_query(files, input)
    if len(matches) <= 0 {
      fmt.Println("No matches found...")
    } else {
      fmt.Println("\nMatches:")
      for _, m := range matches {
        fmt.Printf("	>_ %s\n", files[m].name)
      }
    }

	}
}
