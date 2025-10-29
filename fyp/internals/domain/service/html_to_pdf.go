package service

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func (s *Service) ConvertHTMLToPDF(htmlContent string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var pdfBuf []byte

	er := chromedp.Run(chromeCtx,
		chromedp.Navigate("data:text/html;charset=utf-8,"+htmlContent),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var er error
			pdfBuf, _, er = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0.5).
				WithMarginBottom(0.5).
				WithMarginLeft(0.5).
				WithMarginRight(0.5).
				Do(ctx)
			return er
		}),
	)
	if er != nil {
		return nil, er
	}

	return pdfBuf, nil
}
