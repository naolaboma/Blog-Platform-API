// MongoDB Setup Script for Blog API
// Usage: mongosh < mongodb_setup.js

use blog_db;

print("Setting up Blog API Database...");

// Create users collection with indexes
db.createCollection("users");
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "role": 1 });

print("Users collection created with indexes");

// Create blogs collection with indexes
db.createCollection("blogs");
db.blogs.createIndex({ "author_id": 1 });
db.blogs.createIndex({ "author_username": 1 });
db.blogs.createIndex({ "tags": 1 });
db.blogs.createIndex({ "created_at": -1 });
db.blogs.createIndex({ "view_count": -1 });
db.blogs.createIndex({ "title": "text", "content": "text" });

print("Blogs collection created with indexes");

// Create sessions collection with indexes
db.createCollection("sessions");
db.sessions.createIndex({ "user_id": 1 }, { unique: true });
db.sessions.createIndex({ "username": 1 });
db.sessions.createIndex({ "refresh_token": 1 });
db.sessions.createIndex({ "verification_token": 1 });
db.sessions.createIndex({ "password_reset_token": 1 });
db.sessions.createIndex({ "expires_at": 1 });
db.sessions.createIndex({ "verification_expires_at": 1 });
db.sessions.createIndex({ "reset_expires_at": 1 });

print("Sessions collection created with indexes");

print("Inserting sample data...");

// Sample user
db.users.insertOne({
    username: "admin",
    email: "admin@example.com",
    password: "$2a$10$GCVKOmay2jG0Gi/zDJ2phOCLRxuba4aSkLwZzdjjBEn9eSKfX2fpy",
    role: "admin",
    email_verified: false,
    profile_picture: {
        filename: "",
        file_path: "",
        public_id: "",
        uploaded_at: new Date()
    },
    bio: "System Administrator",
    created_at: new Date(),
    updated_at: new Date()
});

// Sample blog
db.blogs.insertOne({
    title: "Welcome to Blog API",
    content: "This is a sample blog post to test the API.",
    author_id: db.users.findOne({username: "admin"})._id,
    author_username: "admin",
    tags: ["welcome", "sample"],
    view_count: 0,
    like_count: 0,
    comment_count: 0,
    likes: [],
    dislikes: [],
    comments: [],
    created_at: new Date(),
    updated_at: new Date()
});

print("Sample data inserted");

// Show collections
print("\nDatabase Collections:");
db.getCollectionNames().forEach(function(collection) {
    print("  - " + collection);
});

print("\nMongoDB setup completed successfully!");