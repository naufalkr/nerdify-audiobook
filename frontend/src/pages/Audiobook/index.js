import React, { useEffect, useState, useContext } from 'react'
import WrapperPage from '../../components/WrapperPage'
import { GlobalContext } from '../../contexts'
import { CatalogRepository } from '../../repositories'

import './style.css'

function Audiobook({ match, history }) {
    document.title = "Audiobook Details | The Book Hub"

    const audiobookId = match.params.id
    const { setCurrentAudio } = useContext(GlobalContext)

    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [audiobook, setAudiobook] = useState(null)

    const updateCurrentAudio = (track) => {
        setCurrentAudio({
            bookTitle: audiobook.title,
            chapter: {
                title: track.title,
                url: track.url,
                duration: track.duration
            }
        })
    }

    const formatDuration = (duration) => {
        if (!duration) return 'N/A'
        return duration
    }

    const formatGenres = (genres) => {
        if (!genres || genres.length === 0) return 'No genres'
        return genres.map(genre => genre.name).join(', ')
    }

    useEffect(() => {
        const loadAudiobook = async () => {
            try {
                setLoading(true)
                setError(null)

                console.log(`Loading audiobook with ID: ${audiobookId}`)

                const response = await CatalogRepository.getAudiobookById(audiobookId)

                if (response.success) {
                    setAudiobook(response.data)
                    console.log('Audiobook loaded:', response.data)
                } else {
                    setError(response.error || 'Failed to load audiobook')
                }
            } catch (err) {
                console.error('Error loading audiobook:', err)
                setError('Network error occurred while loading audiobook')
            } finally {
                setLoading(false)
            }
        }

        if (audiobookId) {
            loadAudiobook()
        } else {
            setError('No audiobook ID provided')
            setLoading(false)
        }
    }, [audiobookId])

    const handleBackClick = () => {
        history.goBack()
    }

    if (loading) {
        return (
            <WrapperPage>
                <div className="audiobook-page rest-page">
                    <div className="loading">
                        <div></div>
                        <div></div>
                        <div></div>
                    </div>
                </div>
            </WrapperPage>
        )
    }

    if (error) {
        return (
            <WrapperPage>
                <div className="audiobook-page rest-page">
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
                            onClick={handleBackClick}
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
                            Go Back
                        </button>
                    </div>
                </div>
            </WrapperPage>
        )
    }

    if (!audiobook) {
        return (
            <WrapperPage>
                <div className="audiobook-page rest-page">
                    <div style={{ 
                        textAlign: 'center', 
                        padding: '3rem',
                        color: '#94a3b8'
                    }}>
                        <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>📚</div>
                        <h3 style={{ margin: '0 0 0.5rem 0' }}>Audiobook not found</h3>
                        <p style={{ margin: 0, opacity: 0.8 }}>The requested audiobook could not be found.</p>
                        <button 
                            onClick={handleBackClick}
                            style={{
                                marginTop: '1rem',
                                padding: '0.5rem 1rem',
                                background: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
                                color: 'white',
                                border: 'none',
                                borderRadius: '6px',
                                cursor: 'pointer'
                            }}
                        >
                            Go Back
                        </button>
                    </div>
                </div>
            </WrapperPage>
        )
    }

    return (
        <WrapperPage>
            <div className="audiobook-page rest-page">
                {/* Back Button */}
                <button 
                    onClick={handleBackClick}
                    className="back-button"
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        background: 'rgba(255, 255, 255, 0.1)',
                        border: '1px solid rgba(255, 255, 255, 0.2)',
                        color: '#cbd5e1',
                        padding: '0.5rem 1rem',
                        borderRadius: '8px',
                        cursor: 'pointer',
                        fontSize: '0.875rem',
                        marginBottom: '2rem',
                        transition: 'all 0.2s ease'
                    }}
                    onMouseEnter={(e) => {
                        e.target.style.background = 'rgba(255, 255, 255, 0.15)'
                        e.target.style.borderColor = 'rgba(255, 255, 255, 0.3)'
                    }}
                    onMouseLeave={(e) => {
                        e.target.style.background = 'rgba(255, 255, 255, 0.1)'
                        e.target.style.borderColor = 'rgba(255, 255, 255, 0.2)'
                    }}
                >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19 12H5l7-7-1.41-1.41L2.17 12l8.42 8.41L12 19l-7-7z"/>
                    </svg>
                    Back
                </button>

                {/* Audiobook Header */}
                <div className="audiobook-header">
                    <div className="bg-img" style={{ 
                        backgroundImage: `url("${audiobook.image_url}")`,
                        backgroundSize: 'cover',
                        backgroundPosition: 'center',
                        filter: 'blur(20px)',
                        opacity: 0.3
                    }} />
                    
                    <div className="audiobook-info">
                        <div className="cover-art">
                            <img
                                src={audiobook.image_url || '/assets/default-book-cover.jpg'}
                                alt={audiobook.title}
                                onError={(e) => {
                                    e.target.src = '/assets/default-book-cover.jpg'
                                }}
                            />
                        </div>
                        
                        <div className="details">
                            <h1 className="title">{audiobook.title}</h1>
                            <h3 className="author">By {audiobook.author?.name || 'Unknown Author'}</h3>
                            <p className="reader">Narrated by {audiobook.reader?.name || 'Unknown Reader'}</p>
                            
                            <div className="metadata">
                                <div className="meta-item">
                                    <span className="label">Duration:</span>
                                    <span className="value">{formatDuration(audiobook.total_duration)}</span>
                                </div>
                                <div className="meta-item">
                                    <span className="label">Language:</span>
                                    <span className="value">{audiobook.language || 'Unknown'}</span>
                                </div>
                                <div className="meta-item">
                                    <span className="label">Published:</span>
                                    <span className="value">{audiobook.year_of_publishing || 'Unknown'}</span>
                                </div>
                                <div className="meta-item">
                                    <span className="label">Genres:</span>
                                    <span className="value">{formatGenres(audiobook.genres)}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Description */}
                {audiobook.description && (
                    <div className="description-section">
                        <h3>Description</h3>
                        <p>{audiobook.description}</p>
                    </div>
                )}

                {/* Tracks/Chapters */}
                <div className="tracks-section">
                    <h3>
                        Chapters 
                        {audiobook.tracks && audiobook.tracks.length > 0 && (
                            <span className="chapter-count">({audiobook.tracks.length} chapters)</span>
                        )}
                    </h3>
                    
                    {audiobook.tracks && audiobook.tracks.length > 0 ? (
                        <div className="chapters">
                            {audiobook.tracks.map((track, index) => (
                                <div
                                    key={track.id}
                                    className="chapter"
                                    onClick={() => updateCurrentAudio(track)}
                                >
                                    <div className="chapter-number">
                                        {String(index + 1).padStart(2, '0')}
                                    </div>
                                    <div className="chapter-info">
                                        <div className="chapter-title">{track.title}</div>
                                        <div className="chapter-duration">{formatDuration(track.duration)}</div>
                                    </div>
                                    <div className="play-button">
                                        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                            <path d="M8 5v14l11-7z"/>
                                        </svg>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="no-chapters">
                            <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>🎧</div>
                            <p>No chapters available for this audiobook.</p>
                        </div>
                    )}
                </div>
            </div>
        </WrapperPage>
    )
}

export default Audiobook