import axios from 'axios'
import BaseRepository from './BaseRepository'
import SingletonLoggerUtil from '../utils/singletonLogger'

class PlaybackRepository extends BaseRepository {
    constructor() {
        super('PlaybackRepository')
        
        // Log singleton instance creation
        const instanceId = `PlaybackRepository_${Date.now()}`
        const estimatedMemorySize = 512 // Estimated memory footprint in bytes
        SingletonLoggerUtil.logInstanceCreation('PlaybackRepository', instanceId, estimatedMemorySize)
        
        this.contentBaseURL = 'http://localhost:3163/api/v1' // Content Management Service
        this.cache = new Map()
        this.cacheTimeout = 600000 // 10 minutes for track data
    }

    // Cache implementation with logging
    setCache(key, data) {
        try {
            const cacheData = {
                data,
                timestamp: Date.now()
            }
            this.cache.set(key, cacheData)
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'PlaybackRepository',
                'set',
                key,
                false,
                JSON.stringify(data).length
            )
        } catch (error) {
            console.error('PlaybackRepository: Cache set error:', error)
        }
    }

    getFromCache(key) {
        try {
            const cached = this.cache.get(key)
            const isValid = cached && (Date.now() - cached.timestamp < this.cacheTimeout)
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'PlaybackRepository',
                'get',
                key,
                !!isValid,
                isValid ? JSON.stringify(cached.data).length : 0
            )
            
            return isValid ? cached.data : null
        } catch (error) {
            console.error('PlaybackRepository: Cache get error:', error)
            return null
        }
    }

    /**
     * Get track details by ID for playback
     * @param {number} trackId - Track ID
     */
    async getTrackById(trackId) {
        const startTime = SingletonLoggerUtil.logMethodCall('PlaybackRepository', 'getTrackById', { trackId })
        
        try {
            // Check cache first
            const cacheKey = `track_${trackId}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('PlaybackRepository', 'getTrackById', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getTrackById', async () => {
                console.log(`🎵 Fetching track data for ID: ${trackId}`)
                
                const response = await axios.get(`${this.contentBaseURL}/tracks/${trackId}`)
                
                console.log('✅ Track data received:', response.data)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, { trackId })
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('PlaybackRepository', 'getTrackById', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('PlaybackRepository', 'getTrackById', startTime, 'error', error.message)
            throw error
        }
    }

    /**
     * Get all tracks for an audiobook (for playlist functionality)
     * @param {number} audiobookId - Audiobook ID
     */
    async getTracksByAudiobookId(audiobookId) {
        const startTime = SingletonLoggerUtil.logMethodCall('PlaybackRepository', 'getTracksByAudiobookId', { audiobookId })
        
        try {
            // Check cache first
            const cacheKey = `tracks_audiobook_${audiobookId}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('PlaybackRepository', 'getTracksByAudiobookId', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getTracksByAudiobookId', async () => {
                console.log(`🎵 Fetching tracks for audiobook ID: ${audiobookId}`)
                
                const response = await axios.get(`${this.contentBaseURL}/tracks/audiobook/${audiobookId}`)
                
                console.log('✅ Tracks data received:', response.data)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, { audiobookId })
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('PlaybackRepository', 'getTracksByAudiobookId', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('PlaybackRepository', 'getTracksByAudiobookId', startTime, 'error', error.message)
            throw error
        }
    }

    /**
     * Record playback progress
     * @param {number} trackId - Track ID
     * @param {number} currentTime - Current playback time in seconds
     * @param {number} duration - Total track duration in seconds
     */
    async recordProgress(trackId, currentTime, duration) {
        return this.loggedCall('recordProgress', async () => {
            const progressData = {
                track_id: trackId,
                current_time: currentTime,
                duration: duration,
                progress_percentage: Math.round((currentTime / duration) * 100)
            }

            const response = await axios.post(
                `${this.contentBaseURL}/playback/progress`,
                progressData,
                { headers: this.getHeaders() }
            )
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, { trackId, currentTime, duration })
    }

    /**
     * Parse duration string to seconds
     * @param {string} duration - Duration string like "00:21:10" or "21:10"
     */
    static parseDurationToSeconds(duration) {
        if (!duration || typeof duration !== 'string') return 0
        
        try {
            const parts = duration.split(':').map(part => parseInt(part, 10))
            
            if (parts.length === 2) {
                // Format: "mm:ss"
                return parts[0] * 60 + parts[1]
            } else if (parts.length === 3) {
                // Format: "hh:mm:ss"
                return parts[0] * 3600 + parts[1] * 60 + parts[2]
            }
            
            return 0
        } catch (error) {
            console.warn('⚠️ Failed to parse duration:', duration)
            return 0
        }
    }

    /**
     * Format seconds to duration string
     * @param {number} seconds - Seconds to format
     */
    static formatSecondsToTime(seconds) {
        if (!seconds || isNaN(seconds)) return '0:00'
        
        const hours = Math.floor(seconds / 3600)
        const minutes = Math.floor((seconds % 3600) / 60)
        const secs = Math.floor(seconds % 60)
        
        if (hours > 0) {
            return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
        }
        return `${minutes}:${secs.toString().padStart(2, '0')}`
    }

    /**
     * Create playback data object from track and audiobook info
     * @param {Object} track - Track data from API
     * @param {string} audiobookTitle - Audiobook title
     * @param {Array} playlist - All tracks for playlist
     */
    static createPlaybackData(track, audiobookTitle, playlist = []) {
        if (!track) {
            throw new Error('Track data is required')
        }

        return {
            trackId: track.id,
            bookTitle: audiobookTitle || 'Unknown Audiobook',
            chapter: {
                id: track.id,
                title: track.title || 'Unknown Track',
                url: track.url,
                duration: track.duration,
                durationInSeconds: this.parseDurationToSeconds(track.duration)
            },
            playlist: playlist,
            currentTrackIndex: playlist.findIndex(t => t.id === track.id)
        }
    }

    /**
     * Get next track in playlist
     * @param {Array} playlist - Array of tracks
     * @param {number} currentIndex - Current track index
     */
    static getNextTrack(playlist, currentIndex) {
        if (!playlist || playlist.length === 0) return null
        
        const nextIndex = currentIndex + 1
        if (nextIndex >= playlist.length) {
            return null // End of playlist
        }
        
        return { track: playlist[nextIndex], index: nextIndex }
    }

    /**
     * Get previous track in playlist
     * @param {Array} playlist - Array of tracks
     * @param {number} currentIndex - Current track index
     */
    static getPreviousTrack(playlist, currentIndex) {
        if (!playlist || playlist.length === 0) return null
        
        const prevIndex = currentIndex - 1
        if (prevIndex < 0) {
            return null // Beginning of playlist
        }
        
        return { track: playlist[prevIndex], index: prevIndex }
    }

    /**
     * Validate audio URL
     * @param {string} url - Audio URL to validate
     */
    static validateAudioUrl(url) {
        if (!url) return false
        
        try {
            new URL(url)
            return true
        } catch {
            return false
        }
    }

    /**
     * Log playback events (can be extended for analytics)
     * @param {string} eventType - Type of event (play, pause, complete, etc.)
     * @param {Object} trackData - Track data
     * @param {number} currentTime - Current playback time
     */
    static logPlaybackEvent(eventType, trackData, currentTime = 0) {
        console.log('🎵 Playback Event:', {
            event: eventType,
            trackId: trackData.trackId,
            title: trackData.title,
            audiobookTitle: trackData.audiobookTitle,
            currentTime: currentTime,
            timestamp: new Date().toISOString()
        })
        
        // Future implementation: send to analytics API
        // POST http://localhost:3163/api/v1/analytics
    }

    /**
     * Store playback progress in localStorage
     * @param {number} trackId - Track ID
     * @param {number} currentTime - Current playback time
     * @param {number} duration - Total duration
     */
    static savePlaybackProgress(trackId, currentTime, duration) {
        try {
            const progressData = {
                trackId,
                currentTime,
                duration,
                progress: duration > 0 ? (currentTime / duration) * 100 : 0,
                lastPlayed: new Date().toISOString()
            }
            
            localStorage.setItem(`playback_progress_${trackId}`, JSON.stringify(progressData))
        } catch (error) {
            console.error('Error saving playback progress:', error)
        }
    }

    /**
     * Get saved playback progress from localStorage
     * @param {number} trackId - Track ID
     */
    static getPlaybackProgress(trackId) {
        try {
            const saved = localStorage.getItem(`playback_progress_${trackId}`)
            if (saved) {
                return JSON.parse(saved)
            }
        } catch (error) {
            console.error('Error getting playback progress:', error)
        }
        return null
    }

    /**
     * Clear playback progress for a track
     * @param {number} trackId - Track ID
     */
    static clearPlaybackProgress(trackId) {
        try {
            localStorage.removeItem(`playback_progress_${trackId}`)
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'PlaybackRepository',
                'clear',
                `playback_progress_${trackId}`,
                false,
                0
            )
        } catch (error) {
            console.error('Error clearing playback progress:', error)
        }
    }
}

export default new PlaybackRepository()