//go:build gui

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"xCelFlow/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type fileInfo struct {
	fullFilename string
	filename     string
	mkTime       time.Time
}

var (
	searchEntryWidget    *widget.Entry
	fileListWidget       *widget.List
	allButtonWidget      *widget.Button
	clientButtonWidget   *widget.Button
	serverButtonWidget   *widget.Button
	checkButtonWidget    *widget.Check
	multiLineEntryWidget *MultiLineEntryEx
)

var (
	topBorder    *canvas.Rectangle
	leftBorder   *canvas.Rectangle
	rightBorder  *canvas.Rectangle
	bottomBorder *canvas.Rectangle
)

var (
	logBuffers     bytes.Buffer
	filenameList   []fileInfo
	selectFilename string
)

type MultiLineEntryEx struct {
	*widget.Entry
}

func NewMultiLineEntry() *MultiLineEntryEx {
	e := &MultiLineEntryEx{&widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapBreak}}
	e.ExtendBaseWidget(e)
	return e
}

func (e *MultiLineEntryEx) Append(text string) {
	if len(multiLineEntryWidget.Text) > 1024*50 { // 50KB 就需要清空一下
		logBuffers.Reset()
		logBuffers.WriteString("清空日志缓存...\n")
	} else {
		logBuffers.WriteString(text + "\n")
	}
	onChanged := e.OnChanged
	e.OnChanged = nil
	multiLineEntryWidget.SetText(logBuffers.String())
	e.OnChanged = onChanged
	multiLineEntryWidget.CursorRow = len(multiLineEntryWidget.Text) - 1
}

func (e *MultiLineEntryEx) CleanText() {
	logBuffers.Reset()
	onChanged := e.OnChanged
	e.OnChanged = nil
	multiLineEntryWidget.SetText(logBuffers.String())
	e.OnChanged = onChanged
	multiLineEntryWidget.CursorRow = len(multiLineEntryWidget.Text) - 1
}

func startUI() {
	fyneApp := app.New()
	window := fyneApp.NewWindow("配置表导出工具")
	window.Resize(fyne.NewSize(800, 700))
	window.SetFixedSize(true)
	// 输入栏
	initSearchEntry()

	// 文件搜索结果列表
	initList()

	// 初始化按钮
	initBottom()

	// 初始化文本框
	initText()

	left := createLeft(window)
	right := createRight(window)

	content := container.NewHBox(left, right)
	window.SetContent(content)

	window.ShowAndRun()
}

func createLeft(window fyne.Window) *fyne.Container {
	// 创建带有颜色的矩形作为边框
	topBorder := canvas.NewRectangle(color.Black)                    // 蓝色边框
	topBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10)) // 设置边框高度

	leftBorder := canvas.NewRectangle(color.Black)
	leftBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20)) // 减去上下边框的高度

	rightBorder := canvas.NewRectangle(color.Black)
	rightBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20))

	bottomBorder := canvas.NewRectangle(color.Black)
	bottomBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10))

	searchEntryContainer := container.NewGridWrap(fyne.NewSize(300, 40), searchEntryWidget)
	fileListContainer := container.NewGridWrap(fyne.NewSize(300, 635), fileListWidget)
	left := container.NewBorder(topBorder, bottomBorder, leftBorder, rightBorder, container.NewVBox(searchEntryContainer, fileListContainer))
	return left
}

func initSearchEntry() {
	searchEntryWidget = widget.NewEntry()
	searchEntryWidget.SetPlaceHolder("请输入配置文件名，支持模糊搜索")
	// 绑定搜索功能到输入栏
	searchEntryWidget.OnChanged = func(s string) {
		searchFiles(searchEntryWidget.Text)
	}

}

func searchFiles(text string) {
	filenameList = make([]fileInfo, 0)
	// 调用filepath.Walk函数遍历目录
	err := filepath.Walk(config.Config.GetSource(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是文件
		if !info.IsDir() {
			if strings.Contains(info.Name(), text) && !strings.HasPrefix(filepath.Base(path), "~$") {
				filenameList = append(filenameList, fileInfo{fullFilename: path, filename: info.Name(), mkTime: info.ModTime()})
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	sort.Slice(filenameList, func(i, j int) bool {
		return filenameList[i].mkTime.After(filenameList[j].mkTime)
	})

	fileListWidget.Length = func() int { return len(filenameList) }
	fileListWidget.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		item.(*widget.Label).SetText(filenameList[id].filename)
	}
	fileListWidget.Refresh()
}

func initList() {
	fileListWidget = widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {},
	)
	searchFiles("")
	fileListWidget.OnSelected = func(id widget.ListItemID) {
		selectFilename = filenameList[id].fullFilename
	}
}

func createRight(window fyne.Window) *fyne.Container {
	// 创建带有颜色的矩形作为边框
	topBorder = canvas.NewRectangle(color.Black)                     // 红色边框
	topBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10)) // 设置边框高度

	leftBorder = canvas.NewRectangle(color.Black)
	leftBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20)) // 减去上下边框的高度

	rightBorder = canvas.NewRectangle(color.Black)
	rightBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20))

	bottomBorder = canvas.NewRectangle(color.Black)
	bottomBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10))

	buttonContainer := container.NewGridWrap(fyne.NewSize(100, 40), allButtonWidget, clientButtonWidget, serverButtonWidget, checkButtonWidget)
	textContainer := container.NewGridWrap(fyne.NewSize(470, 635), multiLineEntryWidget)

	right := container.NewBorder(topBorder, bottomBorder, leftBorder, rightBorder, container.NewVBox(buttonContainer, textContainer))
	return right
}

func initBottom() {
	allButtonWidget = widget.NewButton("全部导出", func() {
		multiLineEntryWidget.CleanText()
		runExporter([]string{"typescript", "erlang"})
	})

	clientButtonWidget = widget.NewButton("客户端导出", func() {
		multiLineEntryWidget.CleanText()
		runExporter([]string{"typescript"})
	})

	serverButtonWidget = widget.NewButton("服务端导出", func() {
		multiLineEntryWidget.CleanText()
		runExporter([]string{"erlang"})
	})

	checkButtonWidget = widget.NewCheck("是否检查", func(b bool) { config.Config.SetVerify(b) })
	checkButtonWidget.SetChecked(config.Config.GetVerify())
}

func runExporter(schemaNameList []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			printStack()
			borderToRed()
		}
	}()

	if selectFilename == "" {
		return
	}

	borderToBlack()
	for _, schemaName := range schemaNameList {
		// TODO 复用table数据
		if err := run(selectFilename, schemaName); err != nil {
			fmt.Println(err)
			borderToRed()
			break
		} else {
			borderToGreen()
		}
	}
}

func initText() {
	multiLineEntryWidget = NewMultiLineEntry()
	multiLineEntryWidget.OnChanged = func(s string) {
		before, found := strings.CutSuffix(multiLineEntryWidget.Text, s)
		if found {
			multiLineEntryWidget.SetText(before)
		}
	}

	r, w, _ := os.Pipe()

	//defer func() { _ = r.Close() }()
	//defer func() { _ = w.Close() }()

	os.Stdout = w
	os.Stderr = w
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			multiLineEntryWidget.Append(scanner.Text())
		}
	}()
}

func borderToBlack() {
	topBorder.FillColor = color.Black
	topBorder.Refresh()
	leftBorder.FillColor = color.Black
	leftBorder.Refresh()
	rightBorder.FillColor = color.Black
	rightBorder.Refresh()
	bottomBorder.FillColor = color.Black
	bottomBorder.Refresh()
}

func borderToGreen() {
	topBorder.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	topBorder.Refresh()
	leftBorder.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	leftBorder.Refresh()
	rightBorder.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	rightBorder.Refresh()
	bottomBorder.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	bottomBorder.Refresh()
}

func borderToRed() {
	topBorder.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	topBorder.Refresh()
	leftBorder.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	leftBorder.Refresh()
	rightBorder.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	rightBorder.Refresh()
	bottomBorder.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	bottomBorder.Refresh()
}

func main() {
	config.NewTomlConfig("config.toml")
	startUI()
}
