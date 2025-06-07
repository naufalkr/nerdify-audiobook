import React, { useState, useEffect } from 'react'
import WrapperPage from '../../components/WrapperPage'
import Audiobook from '../../components/Audiobook'
import { CatalogRepository } from '../../repositories'
import './style.css'

function ListingPage({ history }){
    document.title = "Curated Audiobooks | The Book Hub"

    const [audiobooks, setAudiobooks] = useState([])
    const [genres, setGenres] = useState([])
    const [loaded, setLoaded] = useState(false)
    const [loadingAudiobooks, setLoadingAudiobooks] = useState(true)
    const [loadingGenres, setLoadingGenres] = useState(true)

    // Load genres on component mount
    useEffect(() => {
        const loadGenres = async() => {
            try {
                setLoadingGenres(true)
                const response = await CatalogRepository.getAllGenres({ 
                    page: 1, 
                    limit: 10
                })
                
                if (response.success) {
                    console.log('Genres loaded:', response.data)
                    setGenres(response.data.items || [])
                } else {
                    console.error('Failed to load genres:', response.error)
                    setGenres([])
                }
            } catch (error) {
                console.error('Error loading genres:', error)
                setGenres([])
            } finally {
                setLoadingGenres(false)
            }
        }
        
        loadGenres()
    }, [])

    

    // Load audiobooks on component mount
    useEffect(() => {
        const loadAudiobooks = async() => {
            try {
                setLoadingAudiobooks(true)
                const response = await CatalogRepository.getAllAudiobooks({ 
                    page: 1, 
                    limit: 20 
                })
                
                if (response.success) {
                    console.log('Audiobooks loaded:', response.data)
                    setAudiobooks(response.data.items || [])
                } else {
                    console.error('Failed to load audiobooks:', response.error)
                    setAudiobooks([])
                }
            } catch (error) {
                console.error('Error loading audiobooks:', error)
                setAudiobooks([])
            } finally {
                setLoadingAudiobooks(false)
                setLoaded(true)
            }
        }
        
        loadAudiobooks()
    }, [])

    const handleGenreClick = (genreName) => {
        history.push(`/genre/${encodeURIComponent(genreName)}`)
    }

    const handleExploreMore = () => {
        console.log('Navigate to full catalog')
    }

    const handleMoreGenres = () => {
        console.log('Navigate to all genres')
    }

    return (
        <WrapperPage>
            <div className="listing-page">
                {/* Genres Section */}
                <div className="group">
                    <div style={{ marginBottom: '2rem' }}>
                        <h2 className="heading discover" style={{ 
                            margin: 0, 
                            fontSize: '2rem', 
                            fontWeight: '700',
                            background: 'linear-gradient(135deg, #ffffff 0%, #e0e7ff 100%)',
                            WebkitBackgroundClip: 'text',
                            WebkitTextFillColor: 'transparent',
                            backgroundClip: 'text'
                        }}>
                            Browse Genres
                        </h2>
                        <p style={{ 
                            margin: '0.5rem 0 0 0', 
                            color: '#94a3b8', 
                            fontSize: '1rem',
                            fontWeight: '400'
                        }}>
                            Discover audiobooks by category
                        </p>
                    </div>
                    
                    <div className="items genres">
                        {loadingGenres ? (
                            <div style={{ gridColumn: '1/-1', textAlign: 'center', padding: '3rem' }}>
                                <div style={{ color: '#94a3b8', fontSize: '1rem' }}>Loading genres...</div>
                            </div>
                        ) : genres.length > 0 ? (
                            genres.map(genre => 
                                <div 
                                    key={genre.id}
                                    className="genre-item" 
                                    onClick={() => handleGenreClick(genre.name)}
                                >
                                    <h3>{genre.name}</h3>
                                </div>
                            )
                        ) : (
                            <div style={{ gridColumn: '1/-1', textAlign: 'center', padding: '3rem' }}>
                                <div style={{ color: '#94a3b8', fontSize: '1rem' }}>No genres available</div>
                            </div>
                        )}
                    </div>

                     {/* Explore More Button - Bottom */}
                    {audiobooks.length >= 20 && (
                        <div style={{ textAlign: 'center', marginTop: '2.5rem' }}>
                            <button 
                                onClick={handleExploreMore}
                                style={{
                                    background: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(168, 85, 247, 0.15) 100%)',
                                    border: '1px solid rgba(139, 92, 246, 0.3)',
                                    color: '#c4b5fd',
                                    padding: '0.6rem 1.5rem',
                                    borderRadius: '8px',
                                    cursor: 'pointer',
                                    fontSize: '0.8rem',
                                    fontWeight: '600',
                                    transition: 'all 0.3s ease',
                                    backdropFilter: 'blur(10px)'
                                }}
                                onMouseEnter={(e) => {
                                    e.target.style.background = 'linear-gradient(135deg, rgba(139, 92, 246, 0.3) 0%, rgba(168, 85, 247, 0.2) 100%)'
                                    e.target.style.borderColor = 'rgba(139, 92, 246, 0.5)'
                                    e.target.style.color = '#e0e7ff'
                                    e.target.style.transform = 'translateY(-2px)'
                                    e.target.style.boxShadow = '0 8px 25px rgba(139, 92, 246, 0.3)'
                                }}
                                onMouseLeave={(e) => {
                                    e.target.style.background = 'linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(168, 85, 247, 0.15) 100%)'
                                    e.target.style.borderColor = 'rgba(139, 92, 246, 0.3)'
                                    e.target.style.color = '#c4b5fd'
                                    e.target.style.transform = 'translateY(0)'
                                    e.target.style.boxShadow = 'none'
                                }}
                            >
                                More Audiobook Categories →
                            </button>
                        </div>
                    )}
                </div>
                
                <div style={{ margin: '4rem 0' }}></div>
                
                {/* Audiobooks Section */}
                <div className="group">
                    <div style={{ marginBottom: '2rem' }}>
                        <h2 className="heading discover" style={{ 
                            margin: 0, 
                            fontSize: '2rem', 
                            fontWeight: '700',
                            background: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
                            WebkitBackgroundClip: 'text',
                            WebkitTextFillColor: 'transparent',
                            backgroundClip: 'text'
                        }}>
                            Featured Audiobooks
                        </h2>
                        <p style={{ 
                            margin: '0.5rem 0 0 0', 
                            color: '#94a3b8', 
                            fontSize: '1rem',
                            fontWeight: '400'
                        }}>
                            Handpicked selections for your listening pleasure
                        </p>
                    </div>
                    
                    <div className="items audiobooks">
                        {loadingAudiobooks ? (
                            <div style={{ gridColumn: '1/-1' }} className="loading">
                                <div></div>
                                <div></div>
                                <div></div>
                            </div>
                        ) : audiobooks.length > 0 ? (
                            audiobooks.map(audiobook => 
                                <Audiobook 
                                    key={audiobook.id}
                                    id={audiobook.id}
                                    title={audiobook.title}
                                    author={audiobook.author}
                                    image_url={audiobook.image_url}
                                    genres={audiobook.genres}
                                    history={history} 
                                />
                            )
                        ) : (
                            <div style={{ gridColumn: '1/-1', textAlign: 'center', padding: '3rem' }}>
                                <div style={{ color: '#94a3b8', fontSize: '1rem' }}>No audiobooks available</div>
                            </div>
                        )}
                    </div>

                    {/* Explore More Button - Bottom */}
                    {audiobooks.length >= 20 && (
                        <div style={{ textAlign: 'center', marginTop: '2.5rem' }}>
                            <button 
                                onClick={handleExploreMore}
                                style={{
                                    background: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(168, 85, 247, 0.15) 100%)',
                                    border: '1px solid rgba(139, 92, 246, 0.3)',
                                    color: '#c4b5fd',
                                    padding: '0.6rem 1.5rem',
                                    borderRadius: '8px',
                                    cursor: 'pointer',
                                    fontSize: '0.8rem',
                                    fontWeight: '600',
                                    transition: 'all 0.3s ease',
                                    backdropFilter: 'blur(10px)'
                                }}
                                onMouseEnter={(e) => {
                                    e.target.style.background = 'linear-gradient(135deg, rgba(139, 92, 246, 0.3) 0%, rgba(168, 85, 247, 0.2) 100%)'
                                    e.target.style.borderColor = 'rgba(139, 92, 246, 0.5)'
                                    e.target.style.color = '#e0e7ff'
                                    e.target.style.transform = 'translateY(-2px)'
                                    e.target.style.boxShadow = '0 8px 25px rgba(139, 92, 246, 0.3)'
                                }}
                                onMouseLeave={(e) => {
                                    e.target.style.background = 'linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(168, 85, 247, 0.15) 100%)'
                                    e.target.style.borderColor = 'rgba(139, 92, 246, 0.3)'
                                    e.target.style.color = '#c4b5fd'
                                    e.target.style.transform = 'translateY(0)'
                                    e.target.style.boxShadow = 'none'
                                }}
                            >
                                Explore All Audiobooks →
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </WrapperPage>
    )
}

export default ListingPage