package grafana

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"strconv"
)

const (
	LoginURL = "http://%s:%d/login"
	PanelSelectorFmt = "#panel-%s > div > div:nth-child(1) > div > div.panel-content > div > plugin-component > panel-plugin-graph > grafana-panel > ng-transclude > div > div.graph-panel__chart > canvas.flot-overlay"
)

func GetLoginTasks(url string, username, password string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#login-view > div > form > div:nth-child(2) > div:nth-child(2) > div > div > input`, chromedp.BySearch),
		chromedp.SetValue(`#login-view > div > form > div:nth-child(1) > div:nth-child(2) > div > div > input`, username, chromedp.BySearch),
		chromedp.SetValue(`#login-view > div > form > div:nth-child(2) > div:nth-child(2) > div > div > input`, password, chromedp.BySearch),
		chromedp.Click("#login-view > div > form > button"),
		chromedp.WaitVisible("#panel-1 > div > div:nth-child(1) > div > div.panel-content.panel-content--no-padding > div > h1", chromedp.BySearch),
		//chromedp.CaptureScreenshot(&buf),
		//chromedp.ActionFunc(func(context.Context) error {
		//	return ioutil.WriteFile("myfb.png", buf, 0644)
		//}
	}
}

func SaveSnapshot(panelID int, url, saveAsFileURL string) chromedp.Tasks {
	var buf []byte
	panelSelector := fmt.Sprintf(PanelSelectorFmt, strconv.Itoa(panelID))
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(panelSelector, chromedp.BySearch),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context) error {
			return ioutil.WriteFile(saveAsFileURL, buf, 0644)
		}),
	}
}