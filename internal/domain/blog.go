package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Core Models ---

type Blog struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	AuthorID       primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorUsername string             `bson:"author_username" json:"author_username"`
	Tags           []string           `bson:"tags" json:"tags"`
	ViewCount      int64              `bson:"view_count" json:"view_count"`
	LikeCount      int64              `bson:"like_count" json:"like_count"`
	CommentCount   int64              `bson:"comment_count" json:"comment_count"`
	Comments       []EmbeddedComment  `bson:"comments" json:"comments"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// EmbeddedComment represents a comment stored directly within a Blog document.
type EmbeddedComment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AuthorID       primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorUsername string             `bson:"author_username" json:"author_username"`
	Content        string             `bson:"content" json:"content"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// Reaction represents a like or dislike in the 'reactions' collection.
type Reaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BlogID       primitive.ObjectID `bson:"blog_id" json:"blog_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	ReactionType string             `bson:"reaction_type" json:"reaction_type"` // e.g., "like" or "dislike"
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// Tag represents a single tag in the 'tags' collection.
type Tag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Constants for reaction types to avoid magic strings in the code.
const (
	ReactionLike    = "like"
	ReactionDislike = "dislike"
)

// --- Repository Interfaces (Data Layer Contracts) ---

type BlogRepository interface {
	Create(blog *Blog) error
	GetByID(id primitive.ObjectID) (*Blog, error)
	Update(blog *Blog) error
	Delete(id primitive.ObjectID) error
	List(params ListBlogParams) ([]*Blog, *PaginationMeta, error)
	AddComment(blogID primitive.ObjectID, comment EmbeddedComment) error
	UpdateComment(blogID, commentID primitive.ObjectID, newContent string) error
	DeleteComment(blogID, commentID primitive.ObjectID) error
	UpdateCounts(blogID primitive.ObjectID, likeDelta, commentDelta int) error
	IncrementViewCount(blogID primitive.ObjectID) error
}

type ReactionRepository interface {
	Create(reaction *Reaction) error
	GetByBlogAndUser(blogID, userID primitive.ObjectID) (*Reaction, error)
	Update(reaction *Reaction) error
	Delete(id primitive.ObjectID) error
}

type TagRepository interface {
	Create(tag *Tag) error
	GetByName(name string) (*Tag, error)
	List() ([]*Tag, error)
}

// --- UseCase Interfaces (Business Logic Contracts) ---

type BlogUseCase interface {
	CreateBlog(authorID primitive.ObjectID, req CreateBlogRequest) (*Blog, error)
	GetBlogByID(id primitive.ObjectID, viewerID *primitive.ObjectID) (*Blog, error)
	UpdateBlog(userID, blogID primitive.ObjectID, req UpdateBlogRequest) (*Blog, error)
	DeleteBlog(userID, blogID primitive.ObjectID) error
	ListBlogs(params ListBlogParams) ([]*Blog, *PaginationMeta, error)
	AddComment(authorID, blogID primitive.ObjectID, req CreateCommentRequest) (*EmbeddedComment, error)
	UpdateComment(userID, blogID, commentID primitive.ObjectID, req UpdateCommentRequest) error
	DeleteComment(userID, blogID, commentID primitive.ObjectID) error
	ReactToBlog(userID, blogID primitive.ObjectID, reactionType string) error
}

// --- Data Transfer Objects (DTOs) & Helpers ---

type CreateBlogRequest struct {
	Title   string   `json:"title" validate:"required,min=5,max=255"`
	Content string   `json:"content" validate:"required,min=20"`
	Tags    []string `json:"tags" validate:"omitempty,dive,alphanum,min=2,max=20"`
}

type UpdateBlogRequest struct {
	Title   *string   `json:"title" validate:"omitempty,min=5,max=255"`
	Content *string   `json:"content" validate:"omitempty,min=20"`
	Tags    *[]string `json:"tags" validate:"omitempty,dive,alphanum,min=2,max=20"`
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type ReactToBlogRequest struct {
	ReactionType string `json:"reaction_type" validate:"required,oneof=like dislike"`
}

type ListBlogParams struct {
	Page       int
	Limit      int
	SortBy     string   // e.g., "newest", "popularity"
	Tags       []string // Filter by tags
	Author     string   // Filter by author username
	SearchTerm string   // For text search on title/content
}

type PaginationMeta struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
}
