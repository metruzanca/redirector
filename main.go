package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func buildRedirects(baseURL string, redirectMappings []string) map[string]string {
	redirects := make(map[string]string)

	for i := 0; i < len(redirectMappings); i += 2 {
		if i+1 < len(redirectMappings) {
			path := redirectMappings[i]
			redirectPath := redirectMappings[i+1]
			redirects[path] = baseURL + redirectPath
		}
	}

	return redirects
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func setupServer(baseURL string, redirectMappings []string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	redirects := buildRedirects(baseURL, redirectMappings)
	for path, redirectTemplate := range redirects {
		path := path
		redirectTemplate := redirectTemplate
		e.GET(path, func(ctx echo.Context) error {
			// Build the redirect URL by substituting parameters
			redirectURL := redirectTemplate

			// Get all parameter names from the path pattern
			pathSegments := strings.Split(path, "/")
			for _, segment := range pathSegments {
				if strings.HasPrefix(segment, ":") {
					paramName := segment[1:] // Remove the ":"
					paramValue := ctx.Param(paramName)
					// Replace :paramName with actual value in redirect URL
					redirectURL = strings.ReplaceAll(redirectURL, ":"+paramName, paramValue)
				}
			}

			// Preserve query parameters
			if queryString := ctx.Request().URL.RawQuery; queryString != "" {
				separator := "?"
				if strings.Contains(redirectURL, "?") {
					separator = "&"
				}
				redirectURL += separator + queryString
			}

			fmt.Printf("Redirecting %s -> %s\n", ctx.Request().URL.Path, redirectURL)
			return ctx.Redirect(301, redirectURL)
		})
	}

	return e
}

func main() {
	godotenv.Load()

	// Parse command line flags
	var baseURL = flag.String("base", "", "Base URL for redirects (e.g., https://example.com)")
	var port = flag.String("port", getEnv("PORT", "8080"), "Port to run the server on")
	flag.Parse()

	// Get remaining arguments as redirect mappings
	args := flag.Args()

	// Validate arguments
	if *baseURL == "" {
		fmt.Fprintf(os.Stderr, "Error: base URL is required. Use -base flag or provide as first argument.\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -base <baseURL> [path1 redirect1 path2 redirect2 ...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -base https://example.com \"/path/:user\" \"/redirecthere/:user\"\n", os.Args[0])
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least one redirect mapping is required.\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -base <baseURL> [path1 redirect1 path2 redirect2 ...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -base https://example.com \"/path/:user\" \"/redirecthere/:user\"\n", os.Args[0])
		os.Exit(1)
	}

	if len(args)%2 != 0 {
		fmt.Fprintf(os.Stderr, "Error: redirect mappings must be in pairs (path, redirect).\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -base <baseURL> [path1 redirect1 path2 redirect2 ...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -base https://example.com \"/path/:user\" \"/redirecthere/:user\"\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("Starting redirect service with base URL: %s\n", *baseURL)
	fmt.Printf("Redirect mappings:\n")
	for i := 0; i < len(args); i += 2 {
		fmt.Printf("  %s -> %s%s\n", args[i], *baseURL, args[i+1])
	}

	e := setupServer(*baseURL, args)
	e.Start(":" + *port)
}
