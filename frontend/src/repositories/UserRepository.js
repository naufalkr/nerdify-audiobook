import axios from 'axios'
import BaseRepository from './BaseRepository'
import SingletonLoggerUtil from '../utils/singletonLogger'

class UserRepository extends BaseRepository {
    constructor() {
        super('UserRepository')
        
        // Log singleton instance creation
        const instanceId = `UserRepository_${Date.now()}`
        const estimatedMemorySize = 512 // Estimated memory footprint in bytes
        SingletonLoggerUtil.logInstanceCreation('UserRepository', instanceId, estimatedMemorySize)
        
        // Remove apiKey since it's causing CORS issues
    }

    getHeaders() {
        const baseHeaders = super.getHeaders()
        return {
            ...baseHeaders
            // Remove X-API-Key header that's causing CORS issues
        }
    }

    // Get all users for admin dashboard (FIXED METHOD)
    async getAllUsersForAdmin(page = 1, limit = 20, search = '') {
        const startTime = SingletonLoggerUtil.logMethodCall('UserRepository', 'getAllUsersForAdmin', { page, limit, search })
        
        try {
            const result = await this.loggedCall('getAllUsersForAdmin', async () => {
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
            }, { page, limit, search })
            
            SingletonLoggerUtil.logMethodEnd('UserRepository', 'getAllUsersForAdmin', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('UserRepository', 'getAllUsersForAdmin', startTime, 'error', error.message)
            throw error
        }
    }

    async updateUser(userId, userData) {
        return this.loggedCall('updateUser', async () => {
            console.log(`UserRepository: Updating user ${userId}`, userData)
            
            const response = await axios.put(
                `${this.baseURL}/api/admin/users/${userId}`,
                userData,
                { headers: this.getHeaders() }
            )

            console.log(`UserRepository: User ${userId} updated successfully`)
            return response.data
        }, { userId, userData })
    }

    async getCurrentUser() {
        return this.loggedCall('getCurrentUser', async () => {
            console.log('UserRepository: Fetching current user profile')
            
            const response = await axios.get(
                `${this.baseURL}/api/users/profile`,
                { headers: this.getHeaders() }
            )

            const user = response.data.user || response.data.data || null
            console.log('UserRepository: Current user profile retrieved')
            return user
        })
    }

    async updateUserProfile(profileData) {
        const startTime = SingletonLoggerUtil.logMethodCall('UserRepository', 'updateUserProfile', { profileData })
        
        try {
            const result = await this.loggedCall('updateUserProfile', async () => {
                console.log('UserRepository: Updating user profile', profileData)
                
                const response = await axios.put(
                    `${this.baseURL}/api/users/profile`,
                    profileData,
                    { headers: this.getHeaders() }
                )

                console.log('UserRepository: Profile updated successfully')
                return response.data
            }, { profileData })
            
            SingletonLoggerUtil.logMethodEnd('UserRepository', 'updateUserProfile', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('UserRepository', 'updateUserProfile', startTime, 'error', error.message)
            throw error
        }
    }
}

export default new UserRepository()