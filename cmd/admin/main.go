package main

import (
	"context"
	"fmt"
	"go-clean-arch/internal/adapter"
	"go-clean-arch/internal/config"
	"go-clean-arch/internal/usecase"
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Initialize repository and application (same as API server)
	blogRepo, err := adapter.NewPostgresBlogRepository(cfg.DatabaseDSN)
	if err != nil {
		logger.Error("Failed to create blog repository", "error", err)
		os.Exit(1)
	}

	app, err := usecase.NewApplication(usecase.NewApplicationParams{
		Logger:   logger,
		BlogRepo: blogRepo,
	})
	if err != nil {
		logger.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	command := os.Args[1]

	switch command {
	case "publish":
		handlePublish(ctx, app, os.Args[2:])
	case "delete":
		handleDelete(ctx, app, os.Args[2:])
	case "list":
		handleList(ctx, app, os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handlePublish(ctx context.Context, app *usecase.Application, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Article ID is required")
		fmt.Println("Usage: go run cmd/admin/main.go publish <article-id>")
		os.Exit(1)
	}

	articleID := args[0]
	article, err := app.PublishArticle(ctx, articleID)
	if err != nil {
		fmt.Printf("Error: Failed to publish article %s: %v\n", articleID, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully published article: %s\n", article.Title)
	fmt.Printf("   ID: %s\n", article.ID)
	fmt.Printf("   Published at: %s\n", article.PublishedAt.Format("2006-01-02 15:04:05"))
}

func handleDelete(ctx context.Context, app *usecase.Application, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Article ID is required")
		fmt.Println("Usage: go run cmd/admin/main.go delete <article-id>")
		os.Exit(1)
	}

	articleID := args[0]
	err := app.DeleteArticle(ctx, articleID)
	if err != nil {
		fmt.Printf("Error: Failed to delete article %s: %v\n", articleID, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully deleted article: %s\n", articleID)
}

func handleList(ctx context.Context, app *usecase.Application, args []string) {
	var filter usecase.BlogFilter

	// Parse optional flags
	for _, arg := range args {
		switch arg {
		case "--drafts":
			draftsOnly := true
			filter.DraftsOnly = &draftsOnly
		case "--published":
			publishedOnly := true
			filter.PublishedOnly = &publishedOnly
		}
	}

	articles, err := app.ListArticles(ctx, filter)
	if err != nil {
		fmt.Printf("Error: Failed to list articles: %v\n", err)
		os.Exit(1)
	}

	if len(articles) == 0 {
		fmt.Println("No articles found")
		return
	}

	fmt.Printf("Found %d article(s):\n\n", len(articles))
	for _, article := range articles {
		status := "Draft"
		if article.IsPublished() {
			status = "Published"
		} else if article.IsScheduled() {
			status = "Scheduled"
		}

		fmt.Printf("📄 %s\n", article.Title)
		fmt.Printf("   ID: %s\n", article.ID)
		fmt.Printf("   Author: %s\n", article.AuthorID)
		fmt.Printf("   Status: %s\n", status)
		if article.PublishedAt != nil {
			fmt.Printf("   Published: %s\n", article.PublishedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("   Created: %s\n", article.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
}

func printUsage() {
	fmt.Println("Clean Blog Admin CLI")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/admin/main.go <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  publish <article-id>     Publish a draft article")
	fmt.Println("  delete <article-id>      Delete any article")
	fmt.Println("  list [--drafts|--published]  List articles (all by default)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/admin/main.go publish abc-123-def")
	fmt.Println("  go run cmd/admin/main.go delete abc-123-def")
	fmt.Println("  go run cmd/admin/main.go list")
	fmt.Println("  go run cmd/admin/main.go list --drafts")
	fmt.Println("  go run cmd/admin/main.go list --published")
}
