import React from 'react'
import './style.css'

function Audiobook({ id, title, author, image_url, genres, history }) {
    const handleClick = () => {
        history.push(`/audiobook/${id}`)
    }

    const formatAuthor = (author) => {
        if (typeof author === 'object' && author.name) {
            return author.name
        }
        return author || 'Unknown Author'
    }

    const formatGenres = (genres) => {
        if (!genres || !Array.isArray(genres)) return ''
        return genres.map(genre => 
            typeof genre === 'object' ? genre.name : genre
        ).join(', ')
    }

    return (
        <div className="audiobook-item" onClick={handleClick}>
            <div className="audiobook-cover">
                <img 
                    src={image_url || '/assets/default-book-cover.jpg'} 
                    alt={title}
                    onError={(e) => {
                        e.target.src = '/assets/default-book-cover.jpg'
                    }}
                />
            </div>
            <div className="audiobook-info">
                <h3 className="audiobook-title">{title}</h3>
                <p className="audiobook-author">{formatAuthor(author)}</p>
                {genres && genres.length > 0 && (
                    <p className="audiobook-genres">{formatGenres(genres)}</p>
                )}
            </div>
        </div>
    )
}

export default Audiobook