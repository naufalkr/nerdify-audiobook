import axios from 'axios'
import BaseRepository from './BaseRepository'
import SingletonLoggerUtil from '../utils/singletonLogger'

class CatalogRepository extends BaseRepository {
    constructor() {
        super('CatalogRepository')
        
        // Log singleton instance creation
        const instanceId = `CatalogRepository_${Date.now()}`
        const estimatedMemorySize = 640 // Estimated memory footprint in bytes
        SingletonLoggerUtil.logInstanceCreation('CatalogRepository', instanceId, estimatedMemorySize)
        
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
                'CatalogRepository',
                'set',
                key,
                false,
                JSON.stringify(data).length
            )
        } catch (error) {
            console.error('CatalogRepository: Cache set error:', error)
        }
    }

    getFromCache(key) {
        try {
            const cached = this.cache.get(key)
            const isValid = cached && (Date.now() - cached.timestamp < this.cacheTimeout)
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'CatalogRepository',
                'get',
                key,
                !!isValid,
                isValid ? JSON.stringify(cached.data).length : 0
            )
            
            return isValid ? cached.data : null
        } catch (error) {
            console.error('CatalogRepository: Cache get error:', error)
            return null
        }
    }

    /**
     * Get all audiobooks for catalog
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     */
    async getAllAudiobooks(params = {}) {
        const startTime = SingletonLoggerUtil.logMethodCall('CatalogRepository', 'getAllAudiobooks', params)
        
        try {
            // Check cache first
            const cacheKey = `catalog_audiobooks_${JSON.stringify(params)}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAllAudiobooks', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getAllAudiobooks', async () => {
                const queryParams = new URLSearchParams()
                
                if (params.page) queryParams.append('page', params.page)
                if (params.limit) queryParams.append('limit', params.limit)

                const response = await axios.get(`${this.contentBaseURL}/audiobooks?${queryParams}`)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, params)
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAllAudiobooks', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAllAudiobooks', startTime, 'error', error.message)
            throw error
        }
    }

    /**
     * Get all genres for catalog
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number  
     * @param {number} params.limit - Items per page
     */
    async getAllGenres(params = {}) {
        const startTime = SingletonLoggerUtil.logMethodCall('CatalogRepository', 'getAllGenres', params)
        
        try {
            // Check cache first
            const cacheKey = `catalog_genres_${JSON.stringify(params)}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAllGenres', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getAllGenres', async () => {
                const queryParams = new URLSearchParams()
                
                if (params.page) queryParams.append('page', params.page)
                if (params.limit) queryParams.append('limit', params.limit)

                const response = await axios.get(`${this.contentBaseURL}/genres?${queryParams}`)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, params)
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAllGenres', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAllGenres', startTime, 'error', error.message)
            throw error
        }
    }

    /**
     * Get audiobooks by genre
     * @param {number} genreId - Genre ID
     * @param {Object} params - Query parameters
     */
    static async getAudiobooksByGenre(genreId, params = {}) {
        const instance = new CatalogRepository()
        return instance.loggedCall('getAudiobooksByGenre', async () => {
            const queryParams = new URLSearchParams()
            queryParams.append('genre_id', genreId)
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${instance.contentBaseURL}/audiobooks?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, { genreId, ...params })
    }

    /**
     * Get audiobook by ID
     * @param {number} audiobookId - Audiobook ID
     */
    async getAudiobookById(audiobookId) {
        const startTime = SingletonLoggerUtil.logMethodCall('CatalogRepository', 'getAudiobookById', { audiobookId })
        
        try {
            // Check cache first
            const cacheKey = `catalog_audiobook_${audiobookId}`
            const cached = this.getFromCache(cacheKey)
            if (cached) {
                SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAudiobookById', startTime, 'success', cached)
                return cached
            }
            
            const result = await this.loggedCall('getAudiobookById', async () => {
                const response = await axios.get(`${this.contentBaseURL}/audiobooks/${audiobookId}`)
                
                return {
                    success: true,
                    data: response.data,
                    status: response.status
                }
            }, { audiobookId })
            
            // Cache the result
            this.setCache(cacheKey, result)
            
            SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAudiobookById', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('CatalogRepository', 'getAudiobookById', startTime, 'error', error.message)
            throw error
        }
    }
}

export default new CatalogRepository()