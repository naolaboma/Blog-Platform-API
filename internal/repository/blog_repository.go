package repository

import (
	"Blog-API/internal/domain"
	"Blog-API/internal/infrastructure/database"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepo struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewBlogRepo(db *database.MongoDB) *BlogRepo {
	collection := db.GetCollection("blogs")

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "author_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "author_username", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "view_count", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		// Log error do not fail but indexes might already exist
		log.Printf("Warning: Failed to create blog indexes: %v", err)
	}

	return &BlogRepo{db, collection}
}

func (Br *BlogRepo) Create(blog *domain.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blog.ID = primitive.NewObjectID()
	_, err := Br.collection.InsertOne(ctx, blog)
	if err != nil {

		// Handle different types of errors
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("blog with this ID already exists: %w", err)
		}

		if mongo.IsTimeout(err) {
			return fmt.Errorf("database operation timed out: %w", err)
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("database context deadline exceeded: %w", err)
		}

		return fmt.Errorf("failed to create blog: %w", err)
	}

	return nil

}

func (Br *BlogRepo) GetByID(id primitive.ObjectID) (*domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var blog domain.Blog

	filter := bson.D{{Key: "_id", Value: id}}

	err := Br.collection.FindOne(ctx, filter).Decode(&blog)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Document not found")
		}
		return nil, err
	}

	return &blog, nil
}

func (Br *BlogRepo) GetAll(page, limit int, sort string) ([]*domain.Blog, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blogs []*domain.Blog
	filter := bson.D{{}}
	options := options.Find()
	options.SetLimit(int64(limit))
	options.SetSkip(int64(page-1) * int64(limit))

	curr, err := Br.collection.Find(ctx, filter, options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, int64(page), errors.New("No document found")
		}
		return nil, 0, fmt.Errorf("failed to find blogs: %w", err)
	}

	if err := curr.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	total, err := Br.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

func (Br *BlogRepo) Update(blog *domain.Blog) error {
	update := bson.D{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check each field and add to update if not empty/zero
	if blog.Title != "" {
		update = append(update, bson.E{Key: "title", Value: blog.Title})
	}

	if blog.Content != "" {
		update = append(update, bson.E{Key: "content", Value: blog.Content})
	}

	// if blog.AuthorUsername != "" {
	//     update = append(update, bson.E{Key: "author_username", Value: blog.AuthorUsername})
	// }

	if len(blog.Tags) > 0 {
		update = append(update, bson.E{Key: "tags", Value: blog.Tags})
	}

	if blog.ViewCount != 0 {
		update = append(update, bson.E{Key: "view_count", Value: blog.ViewCount})
	}

	if blog.LikeCount != 0 {
		update = append(update, bson.E{Key: "like_count", Value: blog.LikeCount})
	}

	if blog.CommentCount != 0 {
		update = append(update, bson.E{Key: "comment_count", Value: blog.CommentCount})
	}

	if len(blog.Likes) > 0 {
		update = append(update, bson.E{Key: "likes", Value: blog.Likes})
	}

	if len(blog.Dislikes) > 0 {
		update = append(update, bson.E{Key: "dislikes", Value: blog.Dislikes})
	}

	// If no fields to update, return early
	if len(update) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Build the update document
	updateDoc := bson.D{
		{Key: "$set", Value: update},
	}

	// Execute the update
	_, err := Br.collection.UpdateByID(ctx, blog.ID, updateDoc)
	if err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}

	return nil
}

func (Br *BlogRepo) Delete(id primitive.ObjectID) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	_, err := Br.collection.DeleteOne(ctx, bson.D{{"_id", id}})

	if err != nil {
		return err
	}
	return nil
}

func (Br *BlogRepo) SearchByTitle(title string, page, limit int) ([]*domain.Blog, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blogs []*domain.Blog
	// Use text search with proper index
	filter := bson.D{{"$text", bson.D{{"$search", title}}}}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}

	options := options.Find().SetSort(sort)
	options.SetLimit(int64(limit))
	options.SetSkip(int64(page-1) * int64(limit))

	curr, err := Br.collection.Find(ctx, filter, options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, int64(page), errors.New("No document found")
		}
		return nil, 0, fmt.Errorf("failed to search blogs by title: %w", err)
	}

	if err := curr.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	total, err := Br.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

func (Br *BlogRepo) SearchByAuthor(author string, page, limit int) ([]*domain.Blog, int64, error) {

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	var blogs []*domain.Blog
	filter := bson.D{{"author", author}}
	sort := bson.D{{"author", 1}}

	options := options.Find().SetSort(sort)
	options.SetLimit(int64(limit))
	options.SetSkip(int64(page-1) * int64(limit))

	curr, err := Br.collection.Find(ctx, filter, options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, int64(page), errors.New("No  document found")
		}
		log.Fatal(err)
	}

	if err := curr.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	total, err := Br.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil

}

// FilterByTags(tags []string, page, limit int) ([]*Blog, int64, error)
// GetPopular(limit int) ([]*Blog, error)
// IncrementViewCount(id primitive.ObjectID) error
// AddComment(blogID primitive.ObjectID, comment *Comment) error
// DeleteComment(blogID, commentID primitive.ObjectID) error
// UpdateComment(blogID, commentID primitive.ObjectID, content string) error
// AddLike(blogID primitive.ObjectID, userID string) error
// RemoveLike(blogID primitive.ObjectID, userID string) error
// AddDislike(blogID primitive.ObjectID, userID string) error
// RemoveDislike(blogID primitive.ObjectID, userID string) error
// GetTagIDByName(name string) (primitive.ObjectID, error)
