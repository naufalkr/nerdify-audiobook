import axios from 'axios'
import BaseRepository from './BaseRepository'
import SingletonLoggerUtil from '../utils/singletonLogger'

class AudiobooksRepository extends BaseRepository {
    constructor() {
        super('AudiobooksRepository')
        
        // Log singleton instance creation
        const instanceId = `AudiobooksRepository_${Date.now()}`
        const estimatedMemorySize = 768 // Estimated memory footprint in bytes
        SingletonLoggerUtil.logInstanceCreation('AudiobooksRepository', instanceId, estimatedMemorySize)
        
        this.contentBaseURL = 'http://localhost:3163/api/v1' // Content Management Service
        this.cache = new Map()
        this.cacheTimeout = 300000 // 5 minutes
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
                'AudiobooksRepository',
                'set',
                key,
                false,
                JSON.stringify(data).length
            )
        } catch (error) {
            console.error('AudiobooksRepository: Cache set error:', error)
        }
    }

    getFromCache(key) {
        try {
            const cached = this.cache.get(key)
            const isValid = cached && (Date.now() - cached.timestamp < this.cacheTimeout)
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'AudiobooksRepository',
                'get',
                key,
                !!isValid,
                isValid ? JSON.stringify(cached.data).length : 0
            )
            
            return isValid ? cached.data : null
        } catch (error) {
            console.error('AudiobooksRepository: Cache get error:', error)
            return null
        }
    }

    /**
     * Get all audiobooks
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @param {string} params.search - Search query
     */
    async getAllAudiobooks(params = {}) {
        const startTime = SingletonLoggerUtil.logMethodCall('AudiobooksRepository', 'getAllAudiobooks', params)
        
        try {
            // Check cache first
            const cacheKey = `audiobooks_${JSON.stringify(params)}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAllAudiobooks', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getAllAudiobooks', async () => {
                const queryParams = new URLSearchParams()
                
                if (params.page) queryParams.append('page', params.page)
                if (params.limit) queryParams.append('limit', params.limit)
                if (params.search) queryParams.append('q', params.search)

                const response = await axios.get(`${this.contentBaseURL}/audiobooks?${queryParams}`)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, params)
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAllAudiobooks', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAllAudiobooks', startTime, 'error', error.message)
            throw error
        }
    }

    /**
     * Get audiobook by ID
     * @param {number} id - Audiobook ID
     */
    async getAudiobookById(id) {
        const startTime = SingletonLoggerUtil.logMethodCall('AudiobooksRepository', 'getAudiobookById', { id })
        
        try {
            // Check cache first
            const cacheKey = `audiobook_${id}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAudiobookById', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getAudiobookById', async () => {
                const response = await axios.get(`${this.contentBaseURL}/audiobooks/${id}`)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, { id })
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAudiobookById', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAudiobookById', startTime, 'error', error.message)
            throw error
        }
    }

    /**
     * Create new audiobook
     * @param {Object} audiobookData - Audiobook data
     */
    async createAudiobook(audiobookData) {
        return this.loggedCall('createAudiobook', async () => {
            const response = await axios.post(
                `${this.contentBaseURL}/audiobooks`,
                audiobookData,
                { headers: this.getHeaders() }
            )
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, { audiobookData })
    }

    /**
     * Update audiobook
     * @param {number} id - Audiobook ID
     * @param {Object} audiobookData - Updated audiobook data
     */
    async updateAudiobook(id, audiobookData) {
        return this.loggedCall('updateAudiobook', async () => {
            const response = await axios.put(
                `${this.contentBaseURL}/audiobooks/${id}`,
                audiobookData,
                { headers: this.getHeaders() }
            )
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, { id, audiobookData })
    }

    /**
     * Delete audiobook (SUPERADMIN only)
     * @param {number} id - Audiobook ID
     */
    async deleteAudiobook(id) {
        return this.loggedCall('deleteAudiobook', async () => {
            const response = await axios.delete(`${this.contentBaseURL}/audiobooks/${id}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        })
    }

    /**
     * Search audiobooks
     * @param {string} query - Search query
     * @param {Object} params - Additional parameters
     */
    async searchAudiobooks(query, params = {}) {
        return this.loggedCall('searchAudiobooks', async () => {
            const queryParams = new URLSearchParams()
            queryParams.append('q', query)
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${this.contentBaseURL}/audiobooks/search?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        })
    }

    /**
     * Get all authors
     * @param {Object} params - Query parameters
     */
    async getAllAuthors(params = {}) {
        return this.loggedCall('getAllAuthors', async () => {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)
            if (params.search) queryParams.append('q', params.search)

            const response = await axios.get(`${this.contentBaseURL}/authors?${queryParams}`)
            
            // DEBUG: Log response untuk melihat struktur data
            console.log('Authors API Raw Response:', response.data)
            console.log('Authors Response Structure:')
            console.log('- response.data:', response.data)
            console.log('- response.data.data:', response.data?.data)
            console.log('- response.data.data.items:', response.data?.data?.items)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        })
    }

    /**
     * Get all genres
     * @param {Object} params - Query parameters
     */
    async getAllGenres(params = {}) {
        const startTime = SingletonLoggerUtil.logMethodCall('AudiobooksRepository', 'getAllGenres', params)
        
        try {
            // Check cache first
            const cacheKey = `genres_${JSON.stringify(params)}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAllGenres', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getAllGenres', async () => {
                const queryParams = new URLSearchParams()
                
                if (params.page) queryParams.append('page', params.page)
                if (params.limit) queryParams.append('limit', params.limit)
                if (params.search) queryParams.append('q', params.search)

                const response = await axios.get(`${this.contentBaseURL}/genres?${queryParams}`)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            })
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAllGenres', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('AudiobooksRepository', 'getAllGenres', startTime, 'error', error.message)
            throw error
        }
    }
}

export default new AudiobooksRepository()