import React, { createContext, useState, useRef, useEffect } from 'react'

const GlobalContext = createContext()

function GlobalProvider({ children }) {
    const [user, setUser] = useState(() => {
        const savedUser = localStorage.getItem('user')
        return savedUser ? JSON.parse(savedUser) : null
    })

    const [currentAudio, setCurrentAudio] = useState(() => {
        const savedAudio = localStorage.getItem('currentAudio')
        return savedAudio ? JSON.parse(savedAudio) : null
    })

    const [isPlaying, setIsPlaying] = useState(false)
    const [currentTime, setCurrentTime] = useState(0)
    const [duration, setDuration] = useState(0)
    const [volume, setVolume] = useState(() => {
        const savedVolume = localStorage.getItem('audioVolume')
        return savedVolume ? parseFloat(savedVolume) : 1
    })
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState(null)

    const audioRef = useRef(null)

    // Save current audio to localStorage when it changes
    useEffect(() => {
        if (currentAudio) {
            localStorage.setItem('currentAudio', JSON.stringify(currentAudio))
        } else {
            localStorage.removeItem('currentAudio')
        }
    }, [currentAudio])

    // Save volume to localStorage when it changes
    useEffect(() => {
        localStorage.setItem('audioVolume', volume.toString())
    }, [volume])

    const login = (userData) => {
        setUser(userData)
        localStorage.setItem('user', JSON.stringify(userData))
    }

    const logout = () => {
        setUser(null)
        setCurrentAudio(null)
        localStorage.removeItem('user')
        localStorage.removeItem('currentAudio')
        
        // Stop and reset audio
        if (audioRef.current) {
            audioRef.current.pause()
            audioRef.current.currentTime = 0
            audioRef.current.src = ''
        }
        setIsPlaying(false)
        setCurrentTime(0)
        setDuration(0)
        setError(null)
    }

    const playAudio = async (audioData) => {
        try {
            setIsLoading(true)
            setError(null)

            console.log('=== Playing Audio ===')
            console.log('Audio Data:', audioData)
            console.log('Chapter URL:', audioData.chapter?.url)

            // Validate audio data
            if (!audioData || !audioData.chapter || !audioData.chapter.url) {
                throw new Error('Invalid audio data: Missing URL')
            }

            let audioUrl = audioData.chapter.url

            // Transform archive.org URLs to be more reliable
            if (audioUrl.includes('www.archive.org')) {
                // Try HTTPS first
                audioUrl = audioUrl.replace('http://www.archive.org', 'https://archive.org')
                console.log('Transformed URL to HTTPS:', audioUrl)
            }

            // Alternative URL patterns for archive.org
            if (audioUrl.includes('archive.org')) {
                const originalUrl = audioUrl
                // Try direct download URL format
                if (audioUrl.includes('/download/')) {
                    audioUrl = originalUrl.replace('archive.org/download/', 'ia800604.us.archive.org/17/')
                    console.log('Trying alternative archive.org URL:', audioUrl)
                }
            }

            // Test URL accessibility with multiple methods
            console.log('Testing URL accessibility:', audioUrl)
            let urlAccessible = false
            
            try {
                // Method 1: Try HEAD request
                const headResponse = await fetch(audioUrl, { 
                    method: 'HEAD',
                    mode: 'no-cors'
                })
                console.log('HEAD request response:', headResponse)
                urlAccessible = true
            } catch (headError) {
                console.warn('HEAD request failed:', headError.message)
                
                try {
                    // Method 2: Try GET with range
                    const rangeResponse = await fetch(audioUrl, {
                        headers: { 'Range': 'bytes=0-1023' },
                        mode: 'no-cors'
                    })
                    console.log('Range request response:', rangeResponse)
                    urlAccessible = true
                } catch (rangeError) {
                    console.warn('Range request failed:', rangeError.message)
                    console.log('Proceeding anyway - some servers block preflight requests')
                }
            }

            // Set current audio data
            const newAudioData = {
                trackId: audioData.trackId,
                bookTitle: audioData.bookTitle,
                chapter: {
                    title: audioData.chapter.title,
                    url: audioUrl,
                    duration: audioData.chapter.duration
                }
            }

            setCurrentAudio(newAudioData)

            // If audio element exists, load the new source
            if (audioRef.current) {
                // Stop current playback
                audioRef.current.pause()
                audioRef.current.currentTime = 0

                // Clear previous source and all event listeners
                audioRef.current.src = ''
                audioRef.current.load()
                
                // Wait a bit for cleanup
                await new Promise(resolve => setTimeout(resolve, 100))

                // Set new source with better error handling
                console.log('Setting audio source:', audioUrl)
                audioRef.current.src = audioUrl
                
                // Important: Remove crossOrigin for external URLs
                audioRef.current.removeAttribute('crossorigin')
                
                // Set preload mode
                audioRef.current.preload = 'auto' // Changed from 'metadata' to 'auto'
                
                // Add comprehensive error handling
                const handleLoadSuccess = () => {
                    console.log('Audio loaded successfully, ready state:', audioRef.current.readyState)
                    audioRef.current.removeEventListener('canplay', handleLoadSuccess)
                    audioRef.current.removeEventListener('error', handleLoadError)
                    audioRef.current.removeEventListener('loadedmetadata', handleLoadSuccess)
                    
                    playWhenReady()
                }

                const handleLoadError = (e) => {
                    console.error('Audio load error details:', {
                        error: e,
                        audioError: audioRef.current.error,
                        errorCode: audioRef.current.error?.code,
                        errorMessage: audioRef.current.error?.message,
                        readyState: audioRef.current.readyState,
                        networkState: audioRef.current.networkState,
                        src: audioRef.current.src,
                        currentSrc: audioRef.current.currentSrc
                    })
                    
                    audioRef.current.removeEventListener('canplay', handleLoadSuccess)
                    audioRef.current.removeEventListener('error', handleLoadError)
                    audioRef.current.removeEventListener('loadedmetadata', handleLoadSuccess)
                    
                    let errorMessage = 'Failed to load audio'
                    
                    if (audioRef.current.error) {
                        switch (audioRef.current.error.code) {
                            case 1: // MEDIA_ERR_ABORTED
                                errorMessage = 'Audio loading was aborted. Try again.'
                                break
                            case 2: // MEDIA_ERR_NETWORK
                                errorMessage = 'Network error. Please check your connection and try again.'
                                break
                            case 3: // MEDIA_ERR_DECODE
                                errorMessage = 'Audio format error. The file may be corrupted or in an unsupported format.'
                                break
                            case 4: // MEDIA_ERR_SRC_NOT_SUPPORTED
                                errorMessage = 'Audio format not supported by this browser. Try a different browser.'
                                break
                            default:
                                errorMessage = audioRef.current.error.message || 'Unknown audio error'
                        }
                    }
                    
                    // Try fallback URL if original fails
                    if (audioUrl.includes('ia800604.us.archive.org') && !audioUrl.includes('retry')) {
                        console.log('Trying original archive.org URL as fallback...')
                        const fallbackUrl = audioData.chapter.url.replace('http://www.archive.org', 'https://archive.org') + '?retry=1'
                        audioRef.current.src = fallbackUrl
                        audioRef.current.load()
                        return
                    }
                    
                    setError(errorMessage)
                    setIsLoading(false)
                }

                // Add event listeners before loading
                audioRef.current.addEventListener('canplay', handleLoadSuccess)
                audioRef.current.addEventListener('loadedmetadata', handleLoadSuccess)
                audioRef.current.addEventListener('error', handleLoadError)
                
                // Load the audio
                console.log('Loading audio...')
                audioRef.current.load()

                // Wait for audio to be ready and then play
                const playWhenReady = () => {
                    audioRef.current.play()
                        .then(() => {
                            console.log('Audio started playing successfully')
                            setIsPlaying(true)
                            setIsLoading(false)
                        })
                        .catch((err) => {
                            console.error('Error playing audio:', err)
                            
                            // Handle specific play errors
                            let playErrorMessage = 'Failed to play audio'
                            if (err.name === 'NotAllowedError') {
                                playErrorMessage = 'Playback blocked by browser. Please click play again.'
                            } else if (err.name === 'NotSupportedError') {
                                playErrorMessage = 'Audio format not supported by this browser'
                            } else if (err.name === 'AbortError') {
                                playErrorMessage = 'Playback was aborted'
                            } else {
                                playErrorMessage = `Playback error: ${err.message}`
                            }
                            
                            setError(playErrorMessage)
                            setIsLoading(false)
                        })
                }
            }
        } catch (err) {
            console.error('Error in playAudio:', err)
            setError(`Failed to play audio: ${err.message}`)
            setIsLoading(false)
        }
    }

    const pauseAudio = () => {
        if (audioRef.current && !audioRef.current.paused) {
            audioRef.current.pause()
        }
    }

    const resumeAudio = () => {
        if (audioRef.current) {
            audioRef.current.play()
                .then(() => {
                    setIsPlaying(true)
                })
                .catch((err) => {
                    console.error('Error resuming audio:', err)
                    setError(`Failed to resume audio: ${err.message}`)
                })
        }
    }

    const seekTo = (time) => {
        if (audioRef.current && !isNaN(time) && time >= 0) {
            audioRef.current.currentTime = time
            setCurrentTime(time)
        }
    }

    const changeVolume = (newVolume) => {
        if (audioRef.current && newVolume >= 0 && newVolume <= 1) {
            audioRef.current.volume = newVolume
            setVolume(newVolume)
        }
    }

    const stopAudio = () => {
        if (audioRef.current) {
            audioRef.current.pause()
            audioRef.current.currentTime = 0
        }
        setIsPlaying(false)
        setCurrentTime(0)
        setCurrentAudio(null)
        setError(null)
    }

    // Format time helper
    const formatTime = (seconds) => {
        if (!seconds || isNaN(seconds)) return '0:00'
        
        const hours = Math.floor(seconds / 3600)
        const minutes = Math.floor((seconds % 3600) / 60)
        const secs = Math.floor(seconds % 60)
        
        if (hours > 0) {
            return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
        }
        return `${minutes}:${secs.toString().padStart(2, '0')}`
    }

    const value = {
        // User state
        user,
        setUser, 
        login,
        logout,
        
        // Audio state
        currentAudio,
        setCurrentAudio,
        isPlaying,
        currentTime,
        duration,
        volume,
        isLoading,
        error,
        audioRef,
        
        // Audio controls
        playAudio,
        pauseAudio,
        resumeAudio,
        seekTo,
        changeVolume,
        stopAudio,
        
        // Setters for audio events
        setIsPlaying,
        setCurrentTime,
        setDuration,
        setError,
        setIsLoading,
        
        // Helpers
        formatTime
    }

    return (
        <GlobalContext.Provider value={value}>
            {children}
        </GlobalContext.Provider>
    )
}

export { GlobalContext, GlobalProvider }