package webpage

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	// DefaultTimeout is the default timeout for the browser
	DefaultTimeout = 30 * time.Second
)

// BrowserOptions is a struct that contains the options for the browser
type BrowserOptions struct {
	timeout   time.Duration
	userAgent string
}

// BrowserOption defines a function to modify BrowserOptions
type BrowserOption func(*BrowserOptions)

// WithTimeout sets the timeout for the browser
func WithTimeout(d time.Duration) BrowserOption {
	return func(o *BrowserOptions) { o.timeout = d }
}

// WithUserAgent sets the user agent for the browser
func WithUserAgent(ua string) BrowserOption {
	return func(o *BrowserOptions) { o.userAgent = ua }
}

// Webpage is a struct that contains the URL and Options for the webpage
type Webpage struct {
	url     string
	browser *BrowserOptions
}

// New creates a new Webpage instance with the given URL and Options
func New(url string, opts ...BrowserOption) *Webpage {
	browser := &BrowserOptions{
		timeout: DefaultTimeout,
	}
	for _, opt := range opts {
		opt(browser)
	}
	return &Webpage{url: url, browser: browser}
}

// PDFOptions is a struct that contains the options for the PDF generation
type PDFOptions struct {
	landscape    bool
	background   bool
	scale        float64
	pagerWidth   float64
	pagerHeight  float64
	marginTop    float64
	marginBottom float64
	marginLeft   float64
	marginRight  float64
	pageRanges   string
}

// PDFOption defines a function to modify PDFOptions
type PDFOption func(*PDFOptions)

// WithLandscape sets the orientation of the PDF to landscape
func WithLandscape(enable bool) PDFOption {
	return func(o *PDFOptions) { o.landscape = enable }
}

// WithBackground enables or disables the background for the PDF
func WithBackground(enable bool) PDFOption {
	return func(o *PDFOptions) { o.background = enable }
}

// WithScale sets the scale of the PDF
func WithScale(scale float64) PDFOption {
	return func(o *PDFOptions) { o.scale = scale }
}

// WithPaperWidth sets the width of the paper for the PDF
func WithPaperWidth(width float64) PDFOption {
	return func(o *PDFOptions) { o.pagerWidth = width }
}

// WithPaperHeight sets the height of the paper for the PDF
func WithPaperHeight(height float64) PDFOption {
	return func(o *PDFOptions) { o.pagerHeight = height }
}

// WithMarginTop sets the top margin for the PDF
func WithMarginTop(top float64) PDFOption {
	return func(o *PDFOptions) { o.marginTop = top }
}

// WithMarginBottom sets the bottom margin for the PDF
func WithMarginBottom(bottom float64) PDFOption {
	return func(o *PDFOptions) { o.marginBottom = bottom }
}

// WithMarginLeft sets the left margin for the PDF
func WithMarginLeft(left float64) PDFOption {
	return func(o *PDFOptions) { o.marginLeft = left }
}

// WithMarginRight sets the right margin for the PDF
func WithMarginRight(right float64) PDFOption {
	return func(o *PDFOptions) { o.marginRight = right }
}

// WithPageRanges sets the page ranges for the PDF
func WithPageRanges(ranges string) PDFOption {
	return func(o *PDFOptions) { o.pageRanges = ranges }
}

// PDF returns the PDF of the webpage as a byte slice
func (w *Webpage) PDF(ctx context.Context, opts ...PDFOption) ([]byte, error) {
	pdfOpts := &PDFOptions{}
	for _, opt := range opts {
		opt(pdfOpts)
	}

	execOpts := []chromedp.ExecAllocatorOption{
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Headless,
	}

	if w.browser.userAgent != "" {
		execOpts = append(execOpts, chromedp.UserAgent(w.browser.userAgent))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, execOpts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(taskCtx, w.browser.timeout)
	defer cancel()

	var pdfBuf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate(w.url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error

			pdfBuf, _, err = page.
				PrintToPDF().
				WithLandscape(pdfOpts.landscape).
				WithPrintBackground(pdfOpts.background).
				WithScale(pdfOpts.scale).
				WithPaperWidth(pdfOpts.pagerWidth).
				WithPaperHeight(pdfOpts.pagerHeight).
				WithMarginTop(pdfOpts.marginTop).
				WithMarginBottom(pdfOpts.marginBottom).
				WithMarginLeft(pdfOpts.marginLeft).
				WithMarginRight(pdfOpts.marginRight).
				WithPageRanges(pdfOpts.pageRanges).
				Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate PDF: %w", err)
			}
			return nil
		}),
	}

	if err := chromedp.Run(timeoutCtx, tasks); err != nil {
		return nil, fmt.Errorf("failed to execute tasks: %w", err)
	}
	return pdfBuf, nil
}

// ScreenshotOptions is a struct that contains the options for the screenshot
type ScreenshotOptions struct {
	format  string
	quality int64
}

// ScreenshotOption defines a function to modify ScreenshotOptions
type ScreenshotOption func(*ScreenshotOptions)

// WithFormat sets the format of the screenshot
func WithFormat(format string) ScreenshotOption {
	return func(o *ScreenshotOptions) { o.format = format }
}

// WithQuality sets the quality of the screenshot
func WithQuality(quality int64) ScreenshotOption {
	return func(o *ScreenshotOptions) { o.quality = quality }
}

// Screenshot returns the screenshot of the webpage as a byte slice
func (w *Webpage) Screenshot(ctx context.Context, opts ...ScreenshotOption) ([]byte, error) {
	screenshotOpts := &ScreenshotOptions{
		format:  "png",
		quality: 100,
	}
	for _, opt := range opts {
		opt(screenshotOpts)
	}

	execOpts := []chromedp.ExecAllocatorOption{
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Headless,
	}

	if w.browser.userAgent != "" {
		execOpts = append(execOpts, chromedp.UserAgent(w.browser.userAgent))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, execOpts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(taskCtx, w.browser.timeout)
	defer cancel()

	var screenshotBuf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate(w.url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			screenshotBuf, err = page.CaptureScreenshot().
				WithCaptureBeyondViewport(true).
				WithFromSurface(true).
				WithFormat(page.CaptureScreenshotFormat(screenshotOpts.format)).
				WithQuality(screenshotOpts.quality).
				Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to capture screenshot: %w", err)
			}
			return nil
		}),
	}

	if err := chromedp.Run(timeoutCtx, tasks); err != nil {
		return nil, fmt.Errorf("failed to execute tasks: %w", err)
	}
	return screenshotBuf, nil
}
