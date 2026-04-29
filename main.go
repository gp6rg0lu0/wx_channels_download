package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// Commit is set at build time via ldflags
	Commit = "none"
	// Date is set at build time via ldflags
	Date = "unknown"
)

func main() {
	app := &cli.App{
		Name:    "wx_channels_download",
		Usage:   "Download videos from WeChat Channels (微信视频号)",
		Version: fmt.Sprintf("%s (commit: %s, built at: %s)", Version, Commit, Date),
		Authors: []*cli.Author{
			{
				Name: "wx_channels_download contributors",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "8888", // changed from 8080 to avoid conflicts with other local dev servers
				Usage:   "Port for the proxy server to listen on",
				EnvVars: []string{"WX_PROXY_PORT"},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "./wx_videos", // personal preference: more descriptive folder name than "downloads"
				Usage:   "Directory to save downloaded videos",
				EnvVars: []string{"WX_OUTPUT_DIR"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose logging",
				EnvVars: []string{"WX_VERBOSE"},
			},
		},
		Action: runProxy,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// runProxy starts the MITM proxy server that intercepts WeChat Channels video requests.
func runProxy(c *cli.Context) error {
	port := c.String("port")
	outputDir := c.String("output")
	verbose := c.Bool("verbose")

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
	}

	if verbose {
		log.Printf("Starting wx_channels_download proxy on port %s", port)
		log.Printf("Videos will be saved to: %s", outputDir)
	}

	fmt.Printf("wx_channels_download v%s\n", Version)
	fmt.Printf("Proxy listening on :%s\n", port)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println("Configure your WeChat to use this proxy, then browse WeChat Channels to capture video URLs.")
	// Reminder: after capturing, check the downloads folder — files are named by video ID
	fmt.Println("Tip: press Ctrl+C to stop the proxy when done.")

	server := NewProxyServer(port, outputDir, verbose)
	return server.Start()
}
