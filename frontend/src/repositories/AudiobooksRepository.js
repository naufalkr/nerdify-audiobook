import axios from 'axios'
import BaseRepository from './BaseRepository'

class AudiobooksRepository extends BaseRepository {
    constructor() {
        super('AudiobooksRepository')
        this.contentBaseURL = 'http://localhost:3163/api/v1' // Content Management Service
    }

    /**
     * Get all audiobooks
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @param {string} params.search - Search query
     */
    async getAllAudiobooks(params = {}) {
        return this.loggedCall('getAllAudiobooks', async () => {
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
    }

    /**
     * Get audiobook by ID
     * @param {number} id - Audiobook ID
     */
    async getAudiobookById(id) {
        return this.loggedCall('getAudiobookById', async () => {
            const response = await axios.get(`${this.contentBaseURL}/audiobooks/${id}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        }, { id })
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
        return this.loggedCall('getAllGenres', async () => {
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
    }
}

export default new AudiobooksRepository()