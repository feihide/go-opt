package main

import "github.com/go-martini/martini"
import "github.com/martini-contrib/render"

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "test"
)

const (
	test_pwd    string = "testtest"
	product_pwd string = "feihide"
)

//MARTINI_ENV=production

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Env struct {
	Name   string
	Title  string
	Number int
}

func main() {
	m := martini.Classic()
	m.Use(martini.Static("../assets"))
	m.Use(render.Renderer(render.Options{
		Directory: "../templates", // Specify what path to load the templates from.
		Layout:    "layout",       // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		//		Extensions: []string{".tmpl", ".html"},  // Specify extensions to load for templates.
		//	Delims:     render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,    // Output human readable JSON
		IndentXML:  true,    // Output human readable XML
		//	HTMLContentType: "application/xhtml+xml",     // Output XHTML content type instead of default "text/html"
	}))
	m.Use(func(c martini.Context, log *log.Logger, res http.ResponseWriter, req *http.Request) {
		log.Println("before a request")

		fmt.Println(martini.Env)
		c.Next()

		log.Println("after a request")
	})
	m.Get("/", func() string {
		return "Hello world"
	})
	m.Get("/opt", func(r render.Render) {
		envs := []Env{{"dev", "开发环境", 1}, {"test", "测试环境", 1}, {"product", "生产环境", 2}}
		logData, _ := ioutil.ReadFile("log.txt")
		r.HTML(200, "opt", map[string]interface{}{"envs": envs, "log": string(logData)})
	})

	m.Get("/opt/run", func(req *http.Request, r render.Render) {
		queryForm, _ := url.ParseQuery(req.URL.RawQuery)
		fmt.Println(queryForm)
		fmt.Println(req.FormValue("name"))
		r.JSON(200, map[string]interface{}{"result": "ok"})
	})

	m.Post("/opt/run", func(w http.ResponseWriter, req *http.Request, r render.Render) {
		//fmt.Fprintln(w, req.PostFormValue("name"))
		name := req.PostFormValue("name")
		pwd := req.PostFormValue("pwd")
		tmp := strings.Split(name, "-")
		isAllow := 0
		if tmp[0] == "dev" {
			isAllow = 1
		}
		if tmp[0] == "test" {
			if pwd == test_pwd {
				isAllow = 1
			}
		}
		if tmp[0] == "product" {
			if pwd == product_pwd {
				isAllow = 1
			}
		}

		var data string
		if isAllow == 0 {
			data = "无权更新"
		} else {
			//command := "ls"
			//params := []string{"-l"}
			//执行cmd命令: ls -l
			command := "cd /work/kl/bin&&./auto.sh " + req.PostFormValue("name") + " update"
			//command := "sleep 3&& echo 'fk'"
			ret, err := execCmd(command)
			if err != nil {
				fmt.Println(err)
				data = "更新失败"
			} else {
				fmt.Println(ret)
				data = "更新成功"
			}
			comm := "sed  -i \"1i\\ `date`  opt:" + name + " result:" + data + " \r\"  log.txt"
			//写入日志
			execCmd(comm)
		}
		r.JSON(200, map[string]interface{}{"result": data})
	})

	m.Get("/test/html", func(r render.Render) {
		r.HTML(200, "index", "fjfjf")
	})

	m.Get("/test/json", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"hello": "world"})
	})

	m.Get("/test/text", func(r render.Render) {
		r.Text(200, "hello, world")
	})
	m.Get("/readfile", func() string {
		dat, err := ioutil.ReadFile("/work/test.pdf")
		check(err)
		return string(dat)
	})

	//路由分组
	m.Group("/books", func(r martini.Router) {
		r.Get("/list", getBooks)
		r.Post("/add", getBooks)
		r.Delete("/delete", getBooks)
	})

	m.Get("/test", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
	})
	m.Get("/hello/:name", func(params martini.Params) string {
		return "Hello " + params["name"]
	})
	m.NotFound(func(res http.ResponseWriter) string {
		res.WriteHeader(404)
		return "noFound"
	})

	m.RunOnAddr(":8400")
}

func execCmd(command string) (string, error) {

	cmd := exec.Command("sh", "-c", command)
	fmt.Println(cmd.Args)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s)", err, command)
		return "", err
	}
	//fmt.Fprintf(os.Stdout, "Result: %s", buf)
	return string(buf), nil
}

func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return true
}

func getBooks() string {
	return "books"
}
