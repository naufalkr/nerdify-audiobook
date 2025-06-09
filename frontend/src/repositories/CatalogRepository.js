import axios from 'axios'
import BaseRepository from './BaseRepository'

class CatalogRepository extends BaseRepository {
    constructor() {
        super('CatalogRepository')
        this.contentBaseURL = 'http://localhost:3163/api/v1' // Content Management Service
    }

    /**
     * Get all audiobooks for catalog
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     */
    static async getAllAudiobooks(params = {}) {
        const instance = new CatalogRepository()
        return instance.loggedCall('getAllAudiobooks', async () => {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${instance.contentBaseURL}/audiobooks?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, params)
    }

    /**
     * Get all genres for catalog
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number  
     * @param {number} params.limit - Items per page
     */
    static async getAllGenres(params = {}) {
        const instance = new CatalogRepository()
        return instance.loggedCall('getAllGenres', async () => {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${instance.contentBaseURL}/genres?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, params)
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
    static async getAudiobookById(audiobookId) {
        const instance = new CatalogRepository()
        return instance.loggedCall('getAudiobookById', async () => {
            const response = await axios.get(`${instance.contentBaseURL}/audiobooks/${audiobookId}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, { audiobookId })
    }
}

export default CatalogRepository