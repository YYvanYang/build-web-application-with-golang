package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 从环境变量获取 GitHub token
var token = os.Getenv("GITHUB_TOKEN")

// 定义一个访问者结构体
type Visitor struct{}

func (v *Visitor) md2html(arg map[string]string) error {
	// 检查 token 是否为空
	if token == "" {
		fmt.Println("错误：未设置 GITHUB_TOKEN 环境变量")
		os.Exit(1)
	}

	from := arg["from"]
	to := arg["to"]

	fmt.Printf("处理目录：%s -> %s\n", from, to)

	s := `<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
`
	err := filepath.Walk(from+"/", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if (f.Mode() & os.ModeSymlink) > 0 {
			return nil
		}
		if !strings.HasSuffix(f.Name(), ".md") {
			return nil
		}

		fmt.Printf("处理文件：%s\n", path)

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("打开文件失败：%v\n", err)
			return err
		}

		input_byte, _ := io.ReadAll(file)
		input := string(input_byte)
		input = regexp.MustCompile(`\[(.*?)\]\(<?(.*?)\.md>?\)`).ReplaceAllString(input, "[$1](<$2.html>)")

		if f.Name() == "README.md" {
			input = regexp.MustCompile(`https:\/\/github\.com\/astaxie\/build-web-application-with-golang\/blob\/master\/`).ReplaceAllString(input, "")
		}

		// 以#开头的行，在#后增加空格
		// 以#开头的行, 删除多余的空格
		input = FixHeader(input)

		// 删除页面链接
		input = RemoveFooterLink(input)

		// remove image suffix
		input = RemoveImageLinkSuffix(input)

		var out *os.File
		filename := strings.Replace(f.Name(), ".md", ".html", -1)
		fmt.Printf("生成文件：%s\n", to+"/"+filename)
		if out, err = os.Create(to + "/" + filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", f.Name(), err)
			os.Exit(-1)
		}
		defer out.Close()
		client := &http.Client{}

		req, err := http.NewRequest("POST", "https://api.github.com/markdown/raw", strings.NewReader(input))
		if err != nil {
			fmt.Printf("创建请求失败：%v\n", err)
			return err
		}

		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("charset", "utf-8")
		req.Header.Set("Authorization", "token "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("发送请求失败：%v\n", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("GitHub API 返回错误：%s\n", string(body))
			return fmt.Errorf("GitHub API 返回状态码：%d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("读取响应失败：%v\n", err)
			return err
		}

		w := bufio.NewWriter(out)
		n4, err := w.WriteString(s + string(body))
		if err != nil {
			fmt.Printf("写入文件失败：%v\n", err)
			return err
		}
		fmt.Printf("写入 %d 字节\n", n4)
		w.Flush()

		return nil
	})
	return err
}

func FixHeader(input string) string {
	re_header := regexp.MustCompile(`(?m)^#.+$`)
	re_sub := regexp.MustCompile(`^(#+)\s*(.+)$`)
	fixer := func(header string) string {
		s := re_sub.FindStringSubmatch(header)
		return s[1] + " " + s[2]
	}
	return re_header.ReplaceAllStringFunc(input, fixer)
}

func RemoveFooterLink(input string) string {
	re_footer := regexp.MustCompile(`(?m)^#{2,} links.*?\n(.+\n)*`)
	return re_footer.ReplaceAllString(input, "")
}

func RemoveImageLinkSuffix(input string) string {
	re_footer := regexp.MustCompile(`png\?raw=true`)
	return re_footer.ReplaceAllString(input, "png")
}

func main() {
	workdir := os.Getenv("WORKDIR")
	if workdir == "" {
		workdir = ".."
	}

	// 使用本地的 output 目录
	outputDir := filepath.Join(workdir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败：%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("工作目录：%s\n", workdir)
	fmt.Printf("输出目录：%s\n", outputDir)

	arg := map[string]string{
		"from": workdir,
		"to":   outputDir,
	}

	v := &Visitor{}
	err := v.md2html(arg)
	if err != nil {
		fmt.Printf("处理失败：%v\n", err)
		os.Exit(1)
	}
}
