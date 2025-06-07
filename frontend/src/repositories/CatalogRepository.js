import axios from 'axios'

const BASE_URL = 'http://localhost:3163/api/v1'

class CatalogRepository {
    /**
     * Get all audiobooks for catalog
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     */
    static async getAllAudiobooks(params = {}) {
        try {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${BASE_URL}/audiobooks?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching audiobooks for catalog:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get all genres for catalog
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number  
     * @param {number} params.limit - Items per page
     */
    static async getAllGenres(params = {}) {
        try {
            const queryParams = new URLSearchParams()
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${BASE_URL}/genres?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching genres for catalog:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get audiobooks by genre ID
     * @param {number} genreId - Genre ID
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     */
    static async getAudiobooksByGenre(genreId, params = {}) {
        try {
            const queryParams = new URLSearchParams()
            queryParams.append('genre_id', genreId)
            
            if (params.page) queryParams.append('page', params.page)
            if (params.limit) queryParams.append('limit', params.limit)

            const response = await axios.get(`${BASE_URL}/audiobooks?${queryParams}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching audiobooks by genre:', error)
            return {
                success: false,
                error: error.response?.data?.error || error.message,
                status: error.response?.status || 500
            }
        }
    }

    /**
     * Get genre by ID
     * @param {number} genreId - Genre ID
     */
    static async getGenreById(genreId) {
        try {
            const response = await axios.get(`${BASE_URL}/genres/${genreId}`)
            
            return {
                success: true,
                data: response.data,
                status: response.status
            }
        } catch (error) {
            console.error('Error fetching genre by ID:', error)
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
            console.error('Error fetching audiobook by ID:', error)
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
}

export default CatalogRepository