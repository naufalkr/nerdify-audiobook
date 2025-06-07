import React, { useEffect, useState } from 'react'

import WrapperPage from '../../components/WrapperPage'
import Audiobook from '../../components/Audiobook'
import { CatalogRepository } from '../../repositories'

function GenreAudiobooks({ history, match }) {
    document.title = "Genre Audiobooks | The Book Hub"

    const genreName = match.params.genre
    const [loadedResults, setLoadedResults] = useState(false)
    const [searchResults, setSearchResults] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [genreId, setGenreId] = useState(null)
    const [currentPage, setCurrentPage] = useState(1)
    const [totalPages, setTotalPages] = useState(1)
    const [totalItems, setTotalItems] = useState(0)

    // First, get all genres to find the genre ID by name
    useEffect(() => {
        const findGenreId = async () => {
            try {
                setLoading(true)
                setError(null)
                
                // Get all genres to find the matching genre ID
                const genresResponse = await CatalogRepository.getAllGenres({ limit: 100 })
                
                if (genresResponse.success) {
                    const genres = genresResponse.data.items || []
                    const matchedGenre = genres.find(genre => 
                        genre.name.toLowerCase() === decodeURIComponent(genreName).toLowerCase()
                    )
                    
                    if (matchedGenre) {
                        setGenreId(matchedGenre.id)
                        console.log(`Found genre: ${matchedGenre.name} with ID: ${matchedGenre.id}`)
                    } else {
                        setError(`Genre "${decodeURIComponent(genreName)}" not found`)
                        setLoading(false)
                    }
                } else {
                    setError('Failed to load genres')
                    setLoading(false)
                }
            } catch (err) {
                console.error('Error finding genre:', err)
                setError('Network error occurred while finding genre')
                setLoading(false)
            }
        }
        
        findGenreId()
    }, [genreName])

    // Then fetch audiobooks by genre ID
    useEffect(() => {
        if (!genreId) return

        const fetchAudiobooksByGenre = async () => {
            try {
                setLoading(true)
                setError(null)
                
                console.log(`Fetching audiobooks for genre ID: ${genreId}, page: ${currentPage}`)
                
                // Use the getAudiobooksByGenre method
                const response = await CatalogRepository.getAudiobooksByGenre(genreId, {
                    page: currentPage,
                    limit: 20
                })
                
                if (response.success) {
                    const audiobooks = response.data.items || []
                    setSearchResults(audiobooks)
                    
                    // Handle pagination
                    if (response.data.pagination) {
                        setTotalItems(response.data.pagination.total)
                        setTotalPages(response.data.pagination.total_pages)
                    } else {
                        setTotalItems(audiobooks.length)
                        setTotalPages(1)
                    }
                    
                    console.log(`Loaded ${audiobooks.length} audiobooks for genre`)
                } else {
                    setError(response.error || 'Failed to load audiobooks')
                    setSearchResults([])
                }
            } catch (err) {
                console.error('Error loading audiobooks by genre:', err)
                setError('Network error occurred while loading audiobooks')
                setSearchResults([])
            } finally {
                setLoading(false)
                setLoadedResults(true)
            }
        }

        fetchAudiobooksByGenre()
    }, [genreId, currentPage])

    const handlePageChange = (newPage) => {
        if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
            setCurrentPage(newPage)
            window.scrollTo({ top: 0, behavior: 'smooth' })
        }
    }

    const generatePaginationNumbers = () => {
        if (totalPages <= 7) {
            return Array.from({ length: totalPages }, (_, i) => i + 1)
        }

        const delta = 2
        const range = []
        const rangeWithDots = []

        for (let i = Math.max(2, currentPage - delta); 
             i <= Math.min(totalPages - 1, currentPage + delta); 
             i++) {
            range.push(i)
        }

        if (currentPage - delta > 2) {
            rangeWithDots.push(1, '...')
        } else {
            rangeWithDots.push(1)
        }

        rangeWithDots.push(...range)

        if (currentPage + delta < totalPages - 1) {
            rangeWithDots.push('...', totalPages)
        } else {
            if (totalPages > 1) {
                rangeWithDots.push(totalPages)
            }
        }

        return rangeWithDots
    }

    return (
        <WrapperPage>
            <div className="search-page rest-page">
                <div style={{ marginBottom: '2rem' }}>
                    <h2 style={{ 
                        margin: 0, 
                        fontSize: '2rem', 
                        fontWeight: '700',
                        background: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
                        WebkitBackgroundClip: 'text',
                        WebkitTextFillColor: 'transparent',
                        backgroundClip: 'text'
                    }}>
                        {decodeURIComponent(genreName)} Audiobooks
                    </h2>
                    {!loading && !error && (
                        <p style={{ 
                            margin: '0.5rem 0 0 0', 
                            color: '#94a3b8', 
                            fontSize: '1rem'
                        }}>
                            {totalItems} audiobook{totalItems !== 1 ? 's' : ''} found
                        </p>
                    )}
                </div>

                {loading && (
                    <div className="loading">
                        <div></div>
                        <div></div>
                        <div></div>
                    </div>
                )}

                {error && (
                    <div style={{ 
                        textAlign: 'center', 
                        padding: '3rem',
                        background: 'rgba(239, 68, 68, 0.1)',
                        border: '1px solid rgba(239, 68, 68, 0.2)',
                        borderRadius: '12px',
                        color: '#fca5a5'
                    }}>
                        <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>😔</div>
                        <h3 style={{ margin: '0 0 0.5rem 0', color: '#ef4444' }}>Error</h3>
                        <p style={{ margin: 0, opacity: 0.8 }}>{error}</p>
                        <button 
                            onClick={() => window.location.reload()}
                            style={{
                                marginTop: '1rem',
                                padding: '0.5rem 1rem',
                                background: '#ef4444',
                                color: 'white',
                                border: 'none',
                                borderRadius: '6px',
                                cursor: 'pointer'
                            }}
                        >
                            Try Again
                        </button>
                    </div>
                )}

                {loadedResults && !loading && !error && searchResults.length > 0 && (
                    <>
                        <div className="group">
                            <div className="items audiobooks">
                                {searchResults.map(audiobook => (
                                    <Audiobook 
                                        key={audiobook.id}
                                        id={audiobook.id}
                                        title={audiobook.title}
                                        author={audiobook.author}
                                        image_url={audiobook.image_url}
                                        genres={audiobook.genres}
                                        history={history} 
                                    />
                                ))}
                            </div>
                        </div>

                        {/* Pagination */}
                        {totalPages > 1 && (
                            <div style={{
                                display: 'flex',
                                justifyContent: 'space-between',
                                alignItems: 'center',
                                marginTop: '3rem',
                                padding: '1.5rem 0',
                                borderTop: '1px solid rgba(255, 255, 255, 0.1)'
                            }}>
                                <div style={{
                                    fontSize: '0.875rem',
                                    color: '#94a3b8'
                                }}>
                                    Showing {((currentPage - 1) * 20) + 1} to {Math.min(currentPage * 20, totalItems)} of {totalItems} results
                                </div>
                                
                                <div style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: '0.5rem'
                                }}>
                                    <button 
                                        onClick={() => handlePageChange(currentPage - 1)}
                                        disabled={currentPage === 1}
                                        style={{
                                            padding: '0.5rem 1rem',
                                            border: '1px solid rgba(255, 255, 255, 0.1)',
                                            background: 'rgba(255, 255, 255, 0.05)',
                                            color: currentPage === 1 ? '#64748b' : '#cbd5e1',
                                            borderRadius: '8px',
                                            cursor: currentPage === 1 ? 'not-allowed' : 'pointer',
                                            fontSize: '0.875rem',
                                            transition: 'all 0.2s ease'
                                        }}
                                    >
                                        ← Previous
                                    </button>
                                    
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                        {totalPages <= 7 ? (
                                            Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                                                <button
                                                    key={page}
                                                    onClick={() => handlePageChange(page)}
                                                    style={{
                                                        padding: '0.5rem 0.75rem',
                                                        border: '1px solid rgba(255, 255, 255, 0.1)',
                                                        background: currentPage === page 
                                                            ? 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)' 
                                                            : 'rgba(255, 255, 255, 0.05)',
                                                        color: currentPage === page ? 'white' : '#cbd5e1',
                                                        borderRadius: '8px',
                                                        cursor: 'pointer',
                                                        fontSize: '0.875rem',
                                                        fontWeight: currentPage === page ? '600' : '400',
                                                        transition: 'all 0.2s ease'
                                                    }}
                                                >
                                                    {page}
                                                </button>
                                            ))
                                        ) : (
                                            generatePaginationNumbers().map((page, index) => (
                                                page === '...' ? (
                                                    <span key={`dots-${index}`} style={{ padding: '0 0.5rem', color: '#6b7280' }}>
                                                        ...
                                                    </span>
                                                ) : (
                                                    <button
                                                        key={page}
                                                        onClick={() => handlePageChange(page)}
                                                        style={{
                                                            padding: '0.5rem 0.75rem',
                                                            border: '1px solid rgba(255, 255, 255, 0.1)',
                                                            background: currentPage === page 
                                                                ? 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)' 
                                                                : 'rgba(255, 255, 255, 0.05)',
                                                            color: currentPage === page ? 'white' : '#cbd5e1',
                                                            borderRadius: '8px',
                                                            cursor: 'pointer',
                                                            fontSize: '0.875rem',
                                                            fontWeight: currentPage === page ? '600' : '400',
                                                            transition: 'all 0.2s ease'
                                                        }}
                                                    >
                                                        {page}
                                                    </button>
                                                )
                                            ))
                                        )}
                                    </div>
                                    
                                    <button 
                                        onClick={() => handlePageChange(currentPage + 1)}
                                        disabled={currentPage === totalPages}
                                        style={{
                                            padding: '0.5rem 1rem',
                                            border: '1px solid rgba(255, 255, 255, 0.1)',
                                            background: 'rgba(255, 255, 255, 0.05)',
                                            color: currentPage === totalPages ? '#64748b' : '#cbd5e1',
                                            borderRadius: '8px',
                                            cursor: currentPage === totalPages ? 'not-allowed' : 'pointer',
                                            fontSize: '0.875rem',
                                            transition: 'all 0.2s ease'
                                        }}
                                    >
                                        Next →
                                    </button>
                                </div>
                            </div>
                        )}
                    </>
                )}

                {loadedResults && !loading && !error && searchResults.length === 0 && (
                    <div style={{ 
                        textAlign: 'center', 
                        padding: '4rem 2rem',
                        background: 'rgba(255, 255, 255, 0.02)',
                        borderRadius: '12px',
                        border: '1px solid rgba(255, 255, 255, 0.05)'
                    }}>
                        <div style={{ fontSize: '4rem', marginBottom: '1rem' }}>📚</div>
                        <h3 style={{ 
                            margin: '0 0 0.5rem 0', 
                            color: '#e2e8f0',
                            fontSize: '1.5rem',
                            fontWeight: '600'
                        }}>
                            No audiobooks found
                        </h3>
                        <p style={{ 
                            margin: '0 0 1.5rem 0', 
                            color: '#94a3b8',
                            fontSize: '1rem'
                        }}>
                            We couldn't find any audiobooks in the "{decodeURIComponent(genreName)}" genre.
                        </p>
                        <button 
                            onClick={() => history.push('/')}
                            style={{
                                background: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
                                border: 'none',
                                color: 'white',
                                padding: '0.75rem 1.5rem',
                                borderRadius: '8px',
                                cursor: 'pointer',
                                fontSize: '0.9rem',
                                fontWeight: '600',
                                transition: 'all 0.3s ease'
                            }}
                            onMouseEnter={(e) => {
                                e.target.style.transform = 'translateY(-2px)'
                                e.target.style.boxShadow = '0 8px 25px rgba(139, 92, 246, 0.3)'
                            }}
                            onMouseLeave={(e) => {
                                e.target.style.transform = 'translateY(0)'
                                e.target.style.boxShadow = 'none'
                            }}
                        >
                            Browse All Categories
                        </button>
                    </div>
                )}
            </div>
        </WrapperPage>
    )
}

export default GenreAudiobooks