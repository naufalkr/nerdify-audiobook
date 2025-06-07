import axios from 'axios'

const BASE_URL = 'http://localhost:3163/api/v1'

class PlaybackRepository {
    /**
     * Get track details by ID for playback
     * @param {number} trackId - Track ID
     */
    static async getTrackById(trackId) {
        try {
            console.log(`🎵 Fetching track data for ID: ${trackId}`)
            
            const response = await axios.get(`${BASE_URL}/tracks/${trackId}`)
            
            console.log('✅ Track data received:', response.data)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('❌ Error fetching track:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get all tracks for an audiobook (for playlist functionality)
     * @param {number} audiobookId - Audiobook ID
     */
    static async getTracksByAudiobookId(audiobookId) {
        try {
            console.log(`🎵 Fetching tracks for audiobook ID: ${audiobookId}`)
            
            const response = await axios.get(`${BASE_URL}/tracks/audiobook/${audiobookId}`)
            
            console.log('✅ Tracks data received:', response.data)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('❌ Error fetching tracks:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
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
        } catch (error) {
            console.error('Error clearing playback progress:', error)
        }
    }
}

export default PlaybackRepository