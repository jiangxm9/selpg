package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

type selpgArgs struct {
	start_page  int    //开始页码
	end_page    int    //结束页码
	page_length int    //每页行数
	dest        string //输出管道
	filename    string //文件名
	page_type   bool   //文件类型

}

func Parser(args *selpgArgs) {
	pflag.IntVarP(&args.start_page, "start", "s", -1, " start page number")
	pflag.IntVarP(&args.end_page, "end", "e", -1, "end page number")
	pflag.IntVarP(&args.page_length, "length", "l", 72, "lines per page")
	pflag.StringVarP(&args.dest, "dest", "d", "", "select the output file")
	pflag.BoolVarP(&args.page_type, "type", "f", false, "divede pages by /f")
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: selpg -s startPage -e endPage [-l linePerPage | -f ][-d dest] filename\n\n")
		pflag.PrintDefaults()
	}
	pflag.Parse()

}

func check(args *selpgArgs) {
	if args.start_page < 1 || args.end_page < 1 || args.start_page > args.end_page {
		pflag.Usage()
	}

	if args.start_page < 1 || args.end_page < 1 {
		fmt.Fprintf(os.Stderr, "\nThe start page or end page is necessary and must be greater than 0!\n")
		os.Exit(1)
	}

	if args.start_page > args.end_page {
		fmt.Fprintf(os.Stderr, "Start page should be less than end page!\n")
		os.Exit(2)
	}

	if args.page_length < 1 {
		fmt.Fprintf(os.Stderr, "The line number must be greater than 0!\n")
		os.Exit(3)
	}

	if pflag.NArg() > 0 {
		args.filename = pflag.Arg(0)
		_, err := os.Stat(args.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Please check your filepath!\n")
			os.Exit(4)
		}
	}
}

func handle(args *selpgArgs) {
	filein := os.Stdin
	fileout := os.Stdout
	lineCount := 0
	pageCount := 1
	//读取文件并判断读取是否出错
	if args.filename != "" {
		err := errors.New("")
		filein, err = os.Open(args.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file failed!\n")
			os.Exit(5)
		}
		defer filein.Close()
	}
	//开始根据是否以换页符分页进行分页
	readLine := bufio.NewReader(filein)
	if args.page_type == false {
		for {
			line, err := readLine.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Read file error!\n")
				os.Exit(6)
			}
			lineCount++
			if lineCount > args.page_length {
				pageCount++
				lineCount = 1
			}
			if pageCount >= args.start_page && pageCount <= args.end_page {
				fmt.Fprintf(fileout, "%s", line)
			}
		}
	} else {
		for {
			page, err := readLine.ReadString('\f')
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Read file error!\n")
				os.Exit(6)
			}
			pageCount++
			if pageCount >= args.start_page && pageCount <= args.end_page {
				fmt.Fprintf(fileout, "%s", page)
			}
		}
	}
	cmd := exec.Command("cat", "-n")
	_, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create pipe error\n")
		os.Exit(7)
	}
	if args.dest != "" {
		cmd.Stdout = fileout
		cmd.Run()
	}
	filein.Close()
	fileout.Close()
}

func main() {
	args := selpgArgs{0, 0, 72, "", "", false}
	Parser(&args)
	check(&args)
	handle(&args)
}
