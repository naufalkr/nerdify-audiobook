import axios from 'axios'

const BASE_URL = 'http://localhost:3163/api/v1'

class AudiobooksRepository {
    /**
     * Get all audiobooks
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @param {string} params.search - Search query
     */
    static async getAllAudiobooks(params = {}) {
        try {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)
            if (params.search) queryParams.append('q', params.search)

            const response = await axios.get(`${BASE_URL}/audiobooks?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching audiobooks:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get audiobook by ID
     * @param {number} id - Audiobook ID
     */
    static async getAudiobookById(id) {
        try {
            const response = await axios.get(`${BASE_URL}/audiobooks/${id}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching audiobook:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Create new audiobook (SUPERADMIN only)
     * @param {Object} audiobookData - Audiobook data
     */
    static async createAudiobook(audiobookData) {
        try {
            const response = await axios.post(`${BASE_URL}/audiobooks`, audiobookData)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error creating audiobook:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Update audiobook (SUPERADMIN only)
     * @param {number} id - Audiobook ID
     * @param {Object} audiobookData - Updated audiobook data
     */
    static async updateAudiobook(id, audiobookData) {
        try {
            const response = await axios.put(`${BASE_URL}/audiobooks/${id}`, audiobookData)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error updating audiobook:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Delete audiobook (SUPERADMIN only)
     * @param {number} id - Audiobook ID
     */
    static async deleteAudiobook(id) {
        try {
            const response = await axios.delete(`${BASE_URL}/audiobooks/${id}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error deleting audiobook:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Search audiobooks
     * @param {string} query - Search query
     * @param {Object} params - Additional parameters
     */
    static async searchAudiobooks(query, params = {}) {
        try {
            const queryParams = new URLSearchParams()
            queryParams.append('q', query)
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${BASE_URL}/audiobooks/search?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error searching audiobooks:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get all authors
     * @param {Object} params - Query parameters
     */
    static async getAllAuthors(params = {}) {
        try {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)
            if (params.search) queryParams.append('q', params.search)

            const response = await axios.get(`${BASE_URL}/authors?${queryParams}`)
            
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
        } catch (error) {
            console.error('Error fetching authors:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get all genres
     * @param {Object} params - Query parameters
     */
    static async getAllGenres(params = {}) {
        try {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)
            if (params.search) queryParams.append('q', params.search)

            const response = await axios.get(`${BASE_URL}/genres?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching genres:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }
}

export default AudiobooksRepository