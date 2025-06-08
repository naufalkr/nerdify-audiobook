# Database Seeders

This directory contains comprehensive database seeders for the audiobook content management system. The seeders populate the database with realistic test data for all entities.

## Seeder Files

### Individual Seeders
- `author_seeder.go` - Seeds the authors table with 100+ famous authors
- `genre_seeder.go` - Seeds the genres table with 100+ diverse genres
- `reader_seeder.go` - Seeds the readers table with 100+ voice actors/readers
- `user_seeder.go` - Seeds the users table with various user types
- `audiobook_seeder.go` - Seeds the audiobooks table with comprehensive book data
- `track_seeder.go` - Seeds the tracks table with chapter/track information
- `analytics_seeder.go` - Seeds the analytics table with user interaction data
- `audiobook_genre_seeder.go` - Creates many-to-many relationships between audiobooks and genres

### Main Seeder
- `main_seeder.go` - Orchestrates all seeders and provides utility functions

## Features

### Comprehensive Data
- **100+ Authors**: From classic literature authors to contemporary writers
- **100+ Genres**: Covering all major literary genres and categories
- **100+ Readers**: Voice actors and narrators for audiobooks
- **80+ Users**: Various user types (admin, moderator, regular users)
- **50+ Audiobooks**: Complete audiobook records with metadata
- **500+ Tracks**: Individual chapters/tracks for each audiobook
- **10,000+ Analytics**: Realistic user interaction and engagement data

### Realistic Relationships
- Proper foreign key relationships between all entities
- Many-to-many relationships between audiobooks and genres
- Intelligent genre assignment based on audiobook content
- Realistic user behavior patterns in analytics

### Smart Data Generation
- Pattern-based genre assignment (e.g., Shakespeare plays get "Tragedy" and "Drama")
- Realistic track naming (Chapter X, Act X Scene Y, etc.)
- Varied track durations and audiobook lengths
- Time-based analytics with realistic user engagement patterns

## Usage

### Using the CLI Tool

```bash
# Run all migrations and seeders
go run cmd/seeder/main.go

# Run only migrations
go run cmd/seeder/main.go -migrate

# Run only seeders
go run cmd/seeder/main.go -seed

# Run specific seeder
go run cmd/seeder/main.go -seed-specific authors
go run cmd/seeder/main.go -seed-specific genres
go run cmd/seeder/main.go -seed-specific audiobooks

# Clear all seeded data
go run cmd/seeder/main.go -clear

# Show seeding statistics
go run cmd/seeder/main.go -stats

# Show help
go run cmd/seeder/main.go -help
```

### Using in Code

```go
import "catalog-service/data_layer/migration"

// Run migration and seeding
err := migration.AutoMigrateAndSeed(db)

// Run only seeding
err := migration.SeedDatabase(db)

// Run specific seeder
err := migration.SeedSpecific(db, "authors")

// Clear all data
err := migration.ClearSeededData(db)

// Get statistics
stats := migration.GetSeedingStatistics(db)
```

## Seeding Order

The seeders run in a specific order to maintain referential integrity:

1. **Authors** - Independent table
2. **Genres** - Independent table
3. **Readers** - Independent table
4. **Users** - Independent table
5. **Audiobooks** - Depends on Authors and Readers
6. **Audiobook-Genre Relationships** - Depends on Audiobooks and Genres
7. **Tracks** - Depends on Audiobooks
8. **Analytics** - Depends on Users and Audiobooks

## Data Volume

After running all seeders, you'll have:
- ~100 authors
- ~100 genres
- ~100 readers
- ~80 users
- ~50 audiobooks
- ~500 tracks
- ~10,000 analytics events
- ~150 audiobook-genre relationships

## Batch Processing

All seeders use batch processing for optimal performance:
- Authors: 20 per batch
- Genres: 25 per batch
- Readers: 20 per batch
- Users: 15 per batch
- Audiobooks: 10 per batch
- Tracks: 50 per batch
- Analytics: 100 per batch

## Duplicate Prevention

The seeders include duplicate prevention:
- Checks existing record count before seeding
- Skips seeding if data already exists
- Safe to run multiple times

## Data Sources

The seed data is inspired by:
- Public domain literature from Project Gutenberg
- LibriVox audiobook catalog
- Classic and contemporary authors
- Diverse literary genres
- Realistic user interaction patterns

## Customization

You can easily customize the seeders by:
- Modifying the data arrays in each seeder file
- Adjusting batch sizes for performance
- Adding new data patterns
- Customizing relationship logic

## Testing

The seeders are designed for:
- Development environment setup
- Testing with realistic data volumes
- Performance testing with substantial datasets
- Feature testing with proper relationships

## Performance

Seeding performance (approximate):
- Authors: ~1 second
- Genres: ~1 second
- Readers: ~1 second
- Users: ~2 seconds
- Audiobooks: ~3 seconds
- Tracks: ~10 seconds
- Analytics: ~30 seconds
- Total: ~50 seconds

## Troubleshooting

### Common Issues

1. **Foreign Key Errors**: Ensure seeders run in the correct order
2. **Duplicate Errors**: The seeders should handle this, but you can clear data first
3. **Performance Issues**: Adjust batch sizes if needed
4. **Memory Issues**: For very large datasets, consider streaming inserts

### Solutions

```bash
# Clear and reseed
go run cmd/seeder/main.go -clear
go run cmd/seeder/main.go -seed

# Check current data
go run cmd/seeder/main.go -stats

# Seed incrementally
go run cmd/seeder/main.go -seed-specific authors
go run cmd/seeder/main.go -seed-specific genres
# ... continue with each seeder
```
