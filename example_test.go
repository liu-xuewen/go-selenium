package selenium_test

import (
	"fmt"

	"sourcegraph.com/sourcegraph/go-selenium"

	"image"
	"os"


	"image/png"
)

func ExampleFindElement() {
	var webDriver selenium.WebDriver
	var err error
	//caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox"})

	var caps = make(selenium.Capabilities)
	if webDriver, err = selenium.NewRemote(caps, "http://rahulghose2:EQhNn9ipXqVWWqpN9mmM@hub-cloud.browserstack.com/wd/hub"); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()

	err = webDriver.Get("https:/www.baidu.com")//https://github.com/sourcegraph/go-selenium
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}

	//保存截图
	reader,err := webDriver.Screenshot();
	if err != nil {
		fmt.Println("Save Image Error!")
	}

	// 转换成png格式的图像，需要导入：_“image/png”
	m, _, _ := image.Decode(reader)
	// 输出到磁盘里
	wt, err := os.Create("test.png")
	if err != nil {
		fmt.Println("Save Image Error!")
	}
	defer wt.Close()

	png.Encode(wt, m)






	if title, err := webDriver.Title(); err == nil {
		fmt.Printf("Page title: %s\n", title)
	} else {
		fmt.Printf("Failed to get page title: %s", err)
		return
	}

	var elem selenium.WebElement
	elem, err = webDriver.FindElement(selenium.ByCSSSelector, ".author")
	if err != nil {
		fmt.Printf("Failed to find element: %s\n", err)
		return
	}

	if text, err := elem.Text(); err == nil {
		fmt.Printf("Author: %s\n", text)
	} else {
		fmt.Printf("Failed to get text of element: %s\n", err)
		return
	}

	// output:
	// Page title: GitHub - sourcegraph/go-selenium: Selenium WebDriver client for Go
	// Author: sourcegraph
}
