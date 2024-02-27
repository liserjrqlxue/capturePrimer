package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"

	"PrimerDesigner/util"

	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
// NewApp 创建一个新的 App 应用程序
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

func (a *App) SelectFile(title string) (filepath string, err error) {
	filepath, err = runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
		ShowHiddenFiles: true,
	})
	fmt.Printf("[%s]:[%v]\n", filepath, err)
	log.Printf("[%s]:[%v]\n", filepath, err)
	if err != nil {
		fmt.Println("cancel?")
		return "", err
	}
	return filepath, err
}

type PrimerPair struct {
	Primer5F *util.Primer `json:"primer5F"`
	Primer3R *util.Primer `json:"primer3R"`
}

type PrimerPairs []*PrimerPair

type Result struct {
	Index          int         `json:"index"`  // 序号
	Status         string      `json:"status"` // 状态
	Name           string      `json:"name"`
	Seq            string      `json:"seq"`
	CapturePrimers PrimerPairs `json:"capturePrimers"`
}

type AllResult struct {
	Status  string    `json:"status"`
	Results []*Result `json:"results"`
}

func (a *App) RunCapture(path string, totalLength int) (allResult AllResult, err error) {
	var data [][]string
	allResult = AllResult{
		Status:  "fail",
		Results: []*Result{},
	}
	data, err = ExcelToSlice(path, "Sheet1")
	if err != nil {
		slog.Error("load xlsx failed", "file", path, "err", err)
		return
	}

	// 创建等待组和结果通道
	var wg sync.WaitGroup
	resultCh := make(chan *Result, len(data))

	// loop Run
	for i, item := range Slice2MapArray(data) {
		wg.Add(1)
		slog.Info("RunCapturePrimer", "i", i, "item", item)

		go func(i int, item map[string]string) {
			defer wg.Done()

			geneName := item["基因名称"]
			rawSeq := item["DNA序列"]

			singelResult := &Result{
				Index:  i + 1,
				Name:   geneName,
				Status: "fail",
			}
			defer func() {
				resultCh <- singelResult
			}()

			Seq := util.NewSeq(geneName, rawSeq, geneName, false)
			slog.Info("Seq", "Name", Seq.Name, "Length", Seq.Length, "extendLength", totalLength)
			Seq.CalAll()
			err := Seq.FindCapturePrimers(totalLength, 200)
			primers := Seq.CapturePrimers
			if err != nil || len(primers) == 0 {
				slog.Error("FindCapturePrimer", "err", err)
				return
			}

			var pairs []*PrimerPair
			for i, primer := range primers {
				left := util.NewPrimer(primer.Name+"-5F-"+strconv.Itoa(i+1), Seq.Seq, primer.End-totalLength, primer.End)
				right := util.NewPrimer(primer.Name+"-3R-"+strconv.Itoa(i+1), Seq.Seq, primer.Start, primer.Start+totalLength)
				pair := &PrimerPair{
					Primer5F: left,
					Primer3R: right,
				}
				pairs = append(pairs, pair)
			}

			singelResult.Status = "success"
			singelResult.CapturePrimers = pairs
		}(i, item)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var result []*Result
	allResult.Status = "success"
	// 从结果通道中读取结果
	for r := range resultCh {
		result = append(result, r)
		if r.Status == "fail" {
			allResult.Status = "fail"
		}
	}

	// 按照索引排序结果
	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})

	allResult.Results = result

	return
}

func ExcelToSlice(filename, sheetName string) ([][]string, error) {
	// Open the Excel file
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer simpleUtil.DeferClose(file)

	// Get all the rows from the specified sheet
	return file.GetRows(sheetName)
}

func Slice2MapArray(s [][]string) (data []map[string]string) {
	var key = s[0]
	for i := 1; i < len(s); i++ {
		var item = make(map[string]string)
		for j := 0; j < len(s[i]); j++ {
			item[key[j]] = s[i][j]
		}
		data = append(data, item)
	}
	return
}

func (a *App) SaveRows(rows []*Result) (filepath string, err error) {
	filepath, err = runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存结果",
		DefaultFilename: "capture.primer.xlsx",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
		ShowHiddenFiles: true,
	})
	slog.Info("SaveFileDialog", "filepath", filepath, "err", err)
	if err != nil || filepath == "" {
		slog.Info("cancel?", "err", err)
		return
	}
	if !strings.HasSuffix(filepath, ".xlsx") {
		filepath += ".xlsx"
	}

	newXlsx := excelize.NewFile()
	newXlsx.SetSheetRow("Sheet1", "A1", &[]interface{}{
		// "index",
		// "status",
		"name",
		"seq",
	})
	for i, result := range rows {
		if result.Status == "success" {
			newXlsx.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+2), &[]interface{}{
				// result.Index,
				// result.Status,
				result.Name,
				result.Seq,
			})
		}
	}
	err = newXlsx.SaveAs(filepath)
	if err != nil {
		slog.Error("SaveAs Failed", "filepath", filepath, "err", err)
		a.DialogWarning("保存失败:" + err.Error())
		return
	}
	slog.Info("SaveFile", "filepath", filepath)

	// 弹窗提醒保存完成
	msg, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "提示",
		Message: "保存完成",
	})
	if err != nil {
		slog.Error("MessageDialog", "msg", msg, "error", err)
		return filepath, nil
	}

	return
}

func (a *App) DialogWarning(msg string) {
	msg, err := runtime.MessageDialog(
		a.ctx,
		runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: msg,
		},
	)
	if err != nil {
		slog.Error("DialogWarning", "msg", msg, "error", err)
	} else {
		slog.Info("DialogWarning", "msg", msg)
	}
}
