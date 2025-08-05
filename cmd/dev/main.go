package main

import (
	"context"
	"fmt"
	"go-clean-arch/internal/adapter"
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"log/slog"
	"os"
	"time"
)

// thinkDifferentContent contains the "Think Different" content for our blog article
const thinkDifferentContent = `Here's to the crazy ones. The misfits. The rebels. The troublemakers. The round pegs in the square holes. The ones who see things differently.

They're not fond of rules. And they have no respect for the status quo.

You can praise them, disagree with them, quote them, disbelieve them, glorify or vilify them. About the only thing you can't do is ignore them. Because they change things.

They invent. They imagine. They heal. They explore. They create. They inspire. They push the human race forward.

Maybe they have to be crazy. How else can you stare at an empty canvas and see a work of art? Or sit in silence and hear a song that's never been written? Or gaze at a red planet and see a laboratory on wheels?

We make tools for these kinds of people. While some may see them as the crazy ones, we see genius. Because the people who are crazy enough to think they can change the world, are the ones who do.`

func main() {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Initialize repository
	blogRepo := adapter.NewMemoryBlogRepository()

	// Create application
	app, err := usecase.NewApplication(usecase.NewApplicationParams{
		Logger:   logger,
		BlogRepo: blogRepo,
	})
	if err != nil {
		logger.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Println("🚀 Clean Blog - Development Environment")
	fmt.Println("=======================================")

	// Test Case 1: Create a blog article
	fmt.Println("\n📝 Creating a new blog article...")
	article, err := app.CreateArticle(ctx, "Think Different", thinkDifferentContent, "steve-jobs")
	if err != nil {
		logger.Error("Failed to create article", "error", err)
		return
	}
	fmt.Printf("✅ Created article: %s (ID: %s)\n", article.Title, article.ID)
	fmt.Printf("   Status: %s\n", getArticleStatus(article))

	// Test Case 2: List drafts
	fmt.Println("\n📋 Listing draft articles...")
	draftFilter := usecase.BlogFilter{DraftsOnly: boolPtr(true)}
	drafts, err := app.ListArticles(ctx, draftFilter)
	if err != nil {
		logger.Error("Failed to list drafts", "error", err)
		return
	}
	fmt.Printf("✅ Found %d draft(s)\n", len(drafts))
	for _, draft := range drafts {
		fmt.Printf("   - %s (%s)\n", draft.Title, draft.ID)
	}

	// Test Case 3: Publish the article
	fmt.Println("\n📢 Publishing the article...")
	publishedArticle, err := app.PublishArticle(ctx, article.ID)
	if err != nil {
		logger.Error("Failed to publish article", "error", err)
		return
	}
	fmt.Printf("✅ Published article: %s\n", publishedArticle.Title)
	fmt.Printf("   Status: %s\n", getArticleStatus(publishedArticle))
	fmt.Printf("   Published at: %s\n", publishedArticle.PublishedAt.Format(time.RFC3339))

	// Test Case 4: List published articles
	fmt.Println("\n📖 Listing published articles...")
	publishedFilter := usecase.BlogFilter{PublishedOnly: boolPtr(true)}
	published, err := app.ListArticles(ctx, publishedFilter)
	if err != nil {
		logger.Error("Failed to list published articles", "error", err)
		return
	}
	fmt.Printf("✅ Found %d published article(s)\n", len(published))
	for _, pub := range published {
		fmt.Printf("   - %s (%s) - Published: %s\n",
			pub.Title, pub.ID, pub.PublishedAt.Format("2006-01-02 15:04:05"))
	}

	// Test Case 5: Get specific article
	fmt.Println("\n🔍 Retrieving specific article...")
	retrievedArticle, err := app.GetArticle(ctx, article.ID)
	if err != nil {
		logger.Error("Failed to get article", "error", err)
		return
	}
	fmt.Printf("✅ Retrieved article: %s\n", retrievedArticle.Title)
	fmt.Printf("   Content preview: %.100s...\n", retrievedArticle.Content)

	// Test Case 6: Update article
	fmt.Println("\n✏️  Updating article content...")
	updatedContent := thinkDifferentContent + "\n\n--- Updated with additional inspiration ---"
	updatedArticle, err := app.UpdateArticle(ctx, article.ID, "Think Different - Updated", updatedContent)
	if err != nil {
		logger.Error("Failed to update article", "error", err)
		return
	}
	fmt.Printf("✅ Updated article: %s\n", updatedArticle.Title)
	fmt.Printf("   Updated at: %s\n", updatedArticle.UpdatedAt.Format(time.RFC3339))

	// Test Case 7: Schedule a future article
	fmt.Println("\n⏰ Creating and scheduling a future article...")
	futurePost, err := app.CreateArticle(ctx, "Future Innovation", "This article is scheduled for the future!", "future-author")
	if err != nil {
		logger.Error("Failed to create future article", "error", err)
		return
	}

	futureTime := time.Now().Add(24 * time.Hour) // Schedule for tomorrow
	scheduledPost, err := app.ScheduleArticle(ctx, futurePost.ID, futureTime)
	if err != nil {
		logger.Error("Failed to schedule article", "error", err)
		return
	}
	fmt.Printf("✅ Scheduled article: %s\n", scheduledPost.Title)
	fmt.Printf("   Status: %s\n", getArticleStatus(scheduledPost))
	fmt.Printf("   Scheduled for: %s\n", scheduledPost.PublishedAt.Format(time.RFC3339))

	// Test Case 8: List all articles
	fmt.Println("\n📚 Listing all articles...")
	allPosts, err := app.ListArticles(ctx, usecase.BlogFilter{})
	if err != nil {
		logger.Error("Failed to list all articles", "error", err)
		return
	}
	fmt.Printf("✅ Found %d total article(s)\n", len(allPosts))
	for _, article := range allPosts {
		fmt.Printf("   - %s (%s) - Status: %s\n",
			article.Title, article.ID, getArticleStatus(article))
	}

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("Clean Architecture blog implementation is working correctly.")
}

// Helper function to get article status as string
func getArticleStatus(article *entity.Article) string {
	if article.IsPublished() {
		return "Published"
	}
	if article.IsScheduled() {
		return "Scheduled"
	}
	return "Draft"
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}
