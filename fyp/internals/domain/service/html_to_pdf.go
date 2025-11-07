package service

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func (s *Service) ConvertHTMLToPDF(htmlContent string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	tmpFile, err := os.CreateTemp("", "certificate-*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write([]byte(htmlContent))
	tmpFile.Close()

	var pdfBuf []byte

	err = chromedp.Run(chromeCtx,

		chromedp.Navigate("file://"+tmpFile.Name()),

		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0.5).
				WithMarginBottom(0.5).
				WithMarginLeft(0.5).
				WithMarginRight(0.5).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, err
	}
	return pdfBuf, nil
}
