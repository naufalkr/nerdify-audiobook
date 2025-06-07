import React, { useContext, useEffect, useRef } from 'react'
import { GlobalContext } from '../../contexts'
// import { CatalogRepository } from '../../repositories'

import './style.css'

const checkAudioSupport = () => {
    const audio = new Audio()
    const formats = {
        mp3: audio.canPlayType('audio/mpeg'),
        ogg: audio.canPlayType('audio/ogg'),
        wav: audio.canPlayType('audio/wav'),
        m4a: audio.canPlayType('audio/mp4'),
        webm: audio.canPlayType('audio/webm')
    }
    
    console.log('Browser audio format support:', formats)
    return formats
}

function Player() {
    const {
        currentAudio,
        isPlaying,
        currentTime,
        duration,
        volume,
        isLoading,
        error,
        audioRef,
        changeVolume,
        // stopAudio,
        setIsPlaying,
        setCurrentTime,
        setDuration,
        setError,
        setIsLoading,
        formatTime
    } = useContext(GlobalContext)

    const progressRef = useRef(null)
    const volumeRef = useRef(null)

    // Initialize audio element and event listeners
    useEffect(() => {
        if (!audioRef.current) {
            audioRef.current = new Audio()
            // Don't set crossOrigin for external URLs
            audioRef.current.preload = "metadata"
        }

        const audio = audioRef.current

        // Audio event listeners
        const handleTimeUpdate = () => {
            if (audio.currentTime) {
                setCurrentTime(audio.currentTime)
            }
        }

        const handleLoadedMetadata = () => {
            console.log('Audio metadata loaded, duration:', audio.duration)
            setDuration(audio.duration || 0)
            setIsLoading(false)
        }

        const handleCanPlay = () => {
            console.log('Audio can play, ready state:', audio.readyState)
            setIsLoading(false)
        }

        const handleEnded = () => {
            console.log('Audio ended')
            setIsPlaying(false)
            setCurrentTime(0)
        }

        const handleError = (e) => {
            console.error('Player audio error:', {
                event: e,
                error: audio.error,
                src: audio.src,
                readyState: audio.readyState,
                networkState: audio.networkState
            })
            
            let errorMsg = 'Failed to load audio'
            
            if (audio.error) {
                switch (audio.error.code) {
                    case MediaError.MEDIA_ERR_ABORTED:
                        errorMsg = 'Audio loading was aborted'
                        break
                    case MediaError.MEDIA_ERR_NETWORK:
                        errorMsg = 'Network error. Please check your connection and try again.'
                        break
                    case MediaError.MEDIA_ERR_DECODE:
                        errorMsg = 'Audio format error. This file may be corrupted or unsupported.'
                        break
                    case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
                        errorMsg = 'Audio format not supported by this browser'
                        break
                    default:
                        errorMsg = audio.error.message || 'Unknown audio error'
                }
            }
            
            setError(errorMsg)
            setIsLoading(false)
            setIsPlaying(false)
        }

        const handleLoadStart = () => {
            console.log('Audio load start')
            setIsLoading(true)
            setError(null)
        }

        const handlePlay = () => {
            console.log('Audio started playing')
            setIsPlaying(true)
        }

        const handlePause = () => {
            console.log('Audio paused')
            setIsPlaying(false)
        }

        const handleWaiting = () => {
            console.log('Audio waiting for data')
            setIsLoading(true)
        }

        const handleCanPlayThrough = () => {
            console.log('Audio can play through')
            setIsLoading(false)
        }

        // Add event listeners
        audio.addEventListener('timeupdate', handleTimeUpdate)
        audio.addEventListener('loadedmetadata', handleLoadedMetadata)
        audio.addEventListener('canplay', handleCanPlay)
        audio.addEventListener('canplaythrough', handleCanPlayThrough)
        audio.addEventListener('ended', handleEnded)
        audio.addEventListener('error', handleError)
        audio.addEventListener('loadstart', handleLoadStart)
        audio.addEventListener('play', handlePlay)
        audio.addEventListener('pause', handlePause)
        audio.addEventListener('waiting', handleWaiting)

        // Set initial volume
        audio.volume = volume

        // Cleanup
        return () => {
            audio.removeEventListener('timeupdate', handleTimeUpdate)
            audio.removeEventListener('loadedmetadata', handleLoadedMetadata)
            audio.removeEventListener('canplay', handleCanPlay)
            audio.removeEventListener('canplaythrough', handleCanPlayThrough)
            audio.removeEventListener('ended', handleEnded)
            audio.removeEventListener('error', handleError)
            audio.removeEventListener('loadstart', handleLoadStart)
            audio.removeEventListener('play', handlePlay)
            audio.removeEventListener('pause', handlePause)
            audio.removeEventListener('waiting', handleWaiting)
        }
    }, [volume, setCurrentTime, setDuration, setIsPlaying, setError, setIsLoading, audioRef])

    // Check audio support on component mount
    useEffect(() => {
        checkAudioSupport()
    }, [])

    // Handle play/pause toggle
    const handlePlayPause = async () => {
        if (!audioRef.current) return

        try {
            if (isPlaying) {
                audioRef.current.pause()
            } else {
                // Ensure we have a valid source
                if (!audioRef.current.src || audioRef.current.src === window.location.href) {
                    if (currentAudio && currentAudio.chapter && currentAudio.chapter.url) {
                        console.log('Setting source from current audio:', currentAudio.chapter.url)
                        audioRef.current.src = currentAudio.chapter.url
                        
                        // Remove crossOrigin for external URLs
                        audioRef.current.removeAttribute('crossorigin')
                        audioRef.current.preload = 'auto'
                        
                        // Load and wait for ready state
                        audioRef.current.load()
                        
                        // Wait for loadedmetadata before playing
                        const waitForLoad = new Promise((resolve, reject) => {
                            const timeout = setTimeout(() => {
                                reject(new Error('Audio load timeout'))
                            }, 10000) // 10 second timeout
                            
                            const onLoad = () => {
                                clearTimeout(timeout)
                                audioRef.current.removeEventListener('loadedmetadata', onLoad)
                                audioRef.current.removeEventListener('error', onError)
                                resolve()
                            }
                            
                            const onError = (e) => {
                                clearTimeout(timeout)
                                audioRef.current.removeEventListener('loadedmetadata', onLoad)
                                audioRef.current.removeEventListener('error', onError)
                                reject(new Error(`Audio load failed: ${audioRef.current.error?.message || 'Unknown error'}`))
                            }
                            
                            audioRef.current.addEventListener('loadedmetadata', onLoad)
                            audioRef.current.addEventListener('error', onError)
                        })
                        
                        await waitForLoad
                    } else {
                        setError('No audio source available')
                        return
                    }
                }
                
                const playPromise = audioRef.current.play()
                if (playPromise !== undefined) {
                    await playPromise
                }
            }
        } catch (err) {
            console.error('Error in play/pause:', err)
            setError(`Playback error: ${err.message}`)
        }
    }

    // Handle previous track
    const handlePrevious = () => {
        // For now, just restart current track if time > 3 seconds
        // Later you can implement actual previous track logic
        if (audioRef.current) {
            if (currentTime > 3) {
                audioRef.current.currentTime = 0
                setCurrentTime(0)
            } else {
                // TODO: Implement previous track logic
                console.log('Previous track - not implemented yet')
                // You can add logic to play previous track from audiobook
            }
        }
    }

    // Handle next track
    const handleNext = () => {
        // TODO: Implement next track logic
        console.log('Next track - not implemented yet')
        // You can add logic to play next track from audiobook
    }

    // Handle progress bar click
    const handleProgressClick = (e) => {
        if (!duration || !progressRef.current || !audioRef.current) return

        const progressBar = progressRef.current
        const rect = progressBar.getBoundingClientRect()
        const clickX = e.clientX - rect.left
        const percentage = clickX / rect.width
        const newTime = percentage * duration

        audioRef.current.currentTime = newTime
        setCurrentTime(newTime)
    }

    // Handle volume change
    const handleVolumeChange = (e) => {
        const newVolume = parseFloat(e.target.value)
        if (audioRef.current) {
            audioRef.current.volume = newVolume
        }
        changeVolume(newVolume)
    }

    // Calculate progress percentage
    const progressPercentage = duration > 0 ? (currentTime / duration) * 100 : 0

    // Don't render if no current audio
    if (!currentAudio) {
        return null
    }

    // Add this test function
    const testAudioUrl = async (url) => {
        console.log('Testing audio URL:', url)
        
        const testAudio = new Audio()
        testAudio.crossOrigin = null
        
        return new Promise((resolve, reject) => {
            const timeout = setTimeout(() => {
                reject(new Error('Test timeout'))
            }, 5000)
            
            testAudio.addEventListener('loadedmetadata', () => {
                clearTimeout(timeout)
                console.log('Test SUCCESS - Audio loadable:', {
                    duration: testAudio.duration,
                    readyState: testAudio.readyState
                })
                resolve(true)
            })
            
            testAudio.addEventListener('error', (e) => {
                clearTimeout(timeout)
                console.log('Test FAILED - Audio error:', {
                    error: e,
                    audioError: testAudio.error
                })
                reject(new Error(testAudio.error?.message || 'Test failed'))
            })
            
            testAudio.src = url
            testAudio.load()
        })
    }

    return (
        <div className="player">
            <div className="player-content">
                {/* Track Info */}
                <div className="player-info">
                    <div className="track-title">{currentAudio.chapter?.title || 'Unknown Track'}</div>
                    <div className="book-title">{currentAudio.bookTitle || 'Unknown Book'}</div>
                </div>

                {/* Controls */}
                <div className="player-controls">
                    <div className="control-buttons">
                        {/* Previous Button */}
                        <button 
                            className="control-btn prev-btn"
                            onClick={handlePrevious}
                            title="Previous / Restart"
                        >
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M6 6h2v12H6zm3.5 6l8.5 6V6z"/>
                            </svg>
                        </button>

                        {/* Play/Pause Button */}
                        <button 
                            className="control-btn play-pause-btn"
                            onClick={handlePlayPause}
                            disabled={isLoading}
                            title={isPlaying ? 'Pause' : 'Play'}
                        >
                            {isLoading ? (
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor" className="loading-icon">
                                    <path d="M12,4V2A10,10 0 0,0 2,12H4A8,8 0 0,1 12,4Z" />
                                </svg>
                            ) : isPlaying ? (
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z"/>
                                </svg>
                            ) : (
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M8 5v14l11-7z"/>
                                </svg>
                            )}
                        </button>

                        {/* Next Button */}
                        <button 
                            className="control-btn next-btn"
                            onClick={handleNext}
                            title="Next"
                        >
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M16 18h2V6h-2v12zM6 6v12l8.5-6z"/>
                            </svg>
                        </button>
                    </div>

                    {/* Progress Bar */}
                    <div className="progress-container">
                        <span className="time-display">{formatTime(currentTime)}</span>
                        <div 
                            className="progress-bar" 
                            ref={progressRef}
                            onClick={handleProgressClick}
                        >
                            <div 
                                className="progress-fill"
                                style={{ width: `${progressPercentage}%` }}
                            />
                            <div 
                                className="progress-handle"
                                style={{ left: `${progressPercentage}%` }}
                            />
                        </div>
                        <span className="time-display">{formatTime(duration)}</span>
                    </div>
                </div>

                {/* Volume Control */}
                <div className="volume-control">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>
                    </svg>
                    <input
                        type="range"
                        ref={volumeRef}
                        min="0"
                        max="1"
                        step="0.1"
                        value={volume}
                        onChange={handleVolumeChange}
                        className="volume-slider"
                    />
                </div>

                {/* Error Display */}
                {error && (
                    <div className="player-error">
                        <span>{error}</span>
                        <button onClick={() => setError(null)} className="error-close">×</button>
                    </div>
                )}
            </div>

            {/* Debug info - keep existing debug components */}
            {process.env.NODE_ENV === 'development' && (
                <div style={{
                    position: 'absolute',
                    top: '-60px',
                    right: '10px',
                    background: 'rgba(0,0,0,0.8)',
                    color: 'white',
                    padding: '5px',
                    fontSize: '10px',
                    borderRadius: '3px'
                }}>
                    URL: {currentAudio?.chapter?.url ? 'Available' : 'None'} | 
                    Playing: {isPlaying ? 'Yes' : 'No'} | 
                    Duration: {duration ? formatTime(duration) : 'Unknown'}
                </div>
            )}

            {/* Test URL button - keep for debugging */}
            {process.env.NODE_ENV === 'development' && currentAudio && (
                <button 
                    onClick={() => testAudioUrl(currentAudio.chapter.url)}
                    style={{
                        padding: '4px 8px',
                        fontSize: '10px',
                        background: '#333',
                        color: 'white',
                        border: 'none',
                        borderRadius: '3px',
                        cursor: 'pointer'
                    }}
                >
                    Test URL
                </button>
            )}
        </div>
    )
}

export default Player