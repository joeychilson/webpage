package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/joeychilson/webpage"
)

var (
	timeout      time.Duration
	userAgent    string
	output       string
	landscape    bool
	background   bool
	scale        float64
	paperWidth   float64
	paperHeight  float64
	marginTop    float64
	marginBottom float64
	marginLeft   float64
	marginRight  float64
	pageRanges   string
	format       string
	quality      int64
)

var rootCmd = &cobra.Command{
	Use:   "webpage",
	Short: "A CLI tool for capturing webpages as PDFs or screenshots",
	Long:  `webshot is a command line interface for capturing webpages as PDFs or screenshots.`,
}

var pdfCmd = &cobra.Command{
	Use:   "pdf [url]",
	Short: "Capture a webpage as PDF",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		opts := []webpage.BrowserOption{
			webpage.WithTimeout(timeout),
		}
		if userAgent != "" {
			opts = append(opts, webpage.WithUserAgent(userAgent))
		}

		page := webpage.New(url, opts...)

		pdfOpts := []webpage.PDFOption{
			webpage.WithLandscape(landscape),
			webpage.WithBackground(background),
			webpage.WithScale(scale),
			webpage.WithPaperWidth(paperWidth),
			webpage.WithPaperHeight(paperHeight),
			webpage.WithMarginTop(marginTop),
			webpage.WithMarginBottom(marginBottom),
			webpage.WithMarginLeft(marginLeft),
			webpage.WithMarginRight(marginRight),
		}
		if pageRanges != "" {
			pdfOpts = append(pdfOpts, webpage.WithPageRanges(pageRanges))
		}

		pdf, err := page.PDF(context.Background(), pdfOpts...)
		if err != nil {
			return fmt.Errorf("failed to generate PDF: %w", err)
		}

		if err := os.WriteFile(output, pdf, 0644); err != nil {
			return fmt.Errorf("failed to write PDF file: %w", err)
		}

		fmt.Printf("PDF saved to: %s\n", output)
		return nil
	},
}

var screenshotCmd = &cobra.Command{
	Use:   "screenshot [url]",
	Short: "Capture a webpage screenshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		opts := []webpage.BrowserOption{
			webpage.WithTimeout(timeout),
		}
		if userAgent != "" {
			opts = append(opts, webpage.WithUserAgent(userAgent))
		}

		page := webpage.New(url, opts...)

		screenshotOpts := []webpage.ScreenshotOption{
			webpage.WithFormat(format),
			webpage.WithQuality(quality),
		}

		screenshot, err := page.Screenshot(context.Background(), screenshotOpts...)
		if err != nil {
			return fmt.Errorf("failed to capture screenshot: %w", err)
		}

		if err := os.WriteFile(output, screenshot, 0644); err != nil {
			return fmt.Errorf("failed to write screenshot file: %w", err)
		}

		fmt.Printf("Screenshot saved to: %s\n", output)
		return nil
	},
}

func init() {
	// browser flags
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 30*time.Second, "timeout for the operation")
	rootCmd.PersistentFlags().StringVarP(&userAgent, "user-agent", "u", "", "user agent string")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "output", "output file path")

	// pdf flags
	pdfCmd.Flags().BoolVarP(&landscape, "landscape", "l", false, "enable landscape mode")
	pdfCmd.Flags().BoolVarP(&background, "background", "b", false, "include background graphics")
	pdfCmd.Flags().Float64VarP(&scale, "scale", "s", 1.0, "scale factor for the PDF")
	pdfCmd.Flags().Float64Var(&paperWidth, "width", 8.5, "paper width in inches")
	pdfCmd.Flags().Float64Var(&paperHeight, "height", 11.0, "paper height in inches")
	pdfCmd.Flags().Float64Var(&marginTop, "margin-top", 0.4, "top margin in inches")
	pdfCmd.Flags().Float64Var(&marginBottom, "margin-bottom", 0.4, "bottom margin in inches")
	pdfCmd.Flags().Float64Var(&marginLeft, "margin-left", 0.4, "left margin in inches")
	pdfCmd.Flags().Float64Var(&marginRight, "margin-right", 0.4, "right margin in inches")
	pdfCmd.Flags().StringVar(&pageRanges, "pages", "", "page ranges to print (e.g., '1-5, 8, 11-13')")

	// screenshot flags
	screenshotCmd.Flags().StringVarP(&format, "format", "f", "png", "screenshot format (png or jpeg)")
	screenshotCmd.Flags().Int64VarP(&quality, "quality", "q", 100, "image quality (0-100, only for jpeg)")

	rootCmd.AddCommand(pdfCmd)
	rootCmd.AddCommand(screenshotCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
