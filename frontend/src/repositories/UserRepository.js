import axios from 'axios'

class UserRepository {
    constructor() {
        this.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:3160'
        // Remove apiKey since it's causing CORS issues
    }

    getHeaders() {
        const token = localStorage.getItem('token')
        console.log('UserRepository: Token found:', token ? 'YES' : 'NO')
        
        return {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
            // Remove X-API-Key header that's causing CORS issues
        }
    }

    // Get all users for admin dashboard (FIXED METHOD)
    async getAllUsersForAdmin(page = 1, limit = 20, search = '') {
        try {
            console.log('=== UserRepository: getAllUsersForAdmin START ===')
            console.log(`Page: ${page}, Limit: ${limit}, Search: "${search}"`)
            console.log(`BaseURL: ${this.baseURL}`)
            
            // Check token first
            const token = localStorage.getItem('token')
            if (!token) {
                throw new Error('No authentication token found')
            }
            
            // Build query params to match API response format
            const params = new URLSearchParams({
                page: page.toString(),
                page_size: limit.toString(),
                ...(search && { search })
            })

            const fullURL = `${this.baseURL}/api/admin/users?${params}`
            console.log(`UserRepository: Making request to: ${fullURL}`)
            console.log('UserRepository: Headers:', this.getHeaders())

            // Use the correct API endpoint from your backend
            const response = await axios.get(fullURL, { 
                headers: this.getHeaders(),
                timeout: 10000 // 10 second timeout
            })

            console.log('UserRepository: API Response Status:', response.status)
            console.log('UserRepository: API Response Data:', response.data)

            // Handle the response format from your API based on the console log
            const users = response.data.data || []
            const meta = response.data.meta || {}
            const total = meta.total || users.length || 0
            
            const result = {
                users: users,
                total: total,
                page: meta.page || page,
                totalPages: Math.ceil(total / (meta.page_size || limit))
            }

            console.log(`UserRepository: Processed result:`, result)
            console.log(`UserRepository: Found ${result.users.length} users of ${result.total} total`)
            console.log('=== UserRepository: getAllUsersForAdmin END ===')
            
            return result
        } catch (error) {
            console.error('=== UserRepository: getAllUsersForAdmin ERROR ===')
            console.error('Error details:', error)
            console.error('Error message:', error.message)
            console.error('Error response:', error.response?.data)
            console.error('Error status:', error.response?.status)
            
            // Return empty result instead of throwing to prevent component crash
            return {
                users: [],
                total: 0,
                page: page,
                totalPages: 0
            }
        }
    }

    // Update user role
    async updateUserRole(userId, newRole) {
        try {
            console.log(`UserRepository: Updating user ${userId} role to ${newRole}`)
            
            const response = await axios.put(
                `${this.baseURL}/api/admin/users/${userId}/role`,
                { role: newRole },
                { headers: this.getHeaders() }
            )

            console.log(`UserRepository: User ${userId} role updated to ${newRole}`)
            return response.data
        } catch (error) {
            console.error(`UserRepository: Error updating user role:`, error)
            throw new Error('Failed to update user role')
        }
    }

    // Toggle user status
    async toggleUserStatus(userId, isActive) {
        try {
            console.log(`UserRepository: Toggling user ${userId} status to ${isActive}`)
            
            const response = await axios.put(
                `${this.baseURL}/api/admin/users/${userId}/status`,
                { is_active: isActive },
                { headers: this.getHeaders() }
            )

            console.log(`UserRepository: User ${userId} status updated`)
            return response.data
        } catch (error) {
            console.error(`UserRepository: Error toggling user status:`, error)
            throw new Error('Failed to update user status')
        }
    }

    // Delete user
    async deleteUser(userId) {
        try {
            console.log(`UserRepository: Deleting user ${userId}`)
            
            const response = await axios.delete(
                `${this.baseURL}/api/admin/users/${userId}`,
                { headers: this.getHeaders() }
            )

            console.log(`UserRepository: User ${userId} deleted`)
            return response.data
        } catch (error) {
            console.error(`UserRepository: Error deleting user:`, error)
            throw new Error('Failed to delete user')
        }
    }

    // Get current user profile
    async getCurrentUser() {
        try {
            console.log('UserRepository: Fetching current user profile')
            
            const response = await axios.get(
                `${this.baseURL}/api/users/profile`,
                { headers: this.getHeaders() }
            )

            const user = response.data.user || response.data.data || null
            console.log('UserRepository: Current user profile retrieved')
            return user
        } catch (error) {
            console.error('UserRepository: Error fetching current user:', error)
            // Fallback to localStorage user data
            const userData = localStorage.getItem('user')
            return userData ? JSON.parse(userData) : null
        }
    }
}

export default new UserRepository()