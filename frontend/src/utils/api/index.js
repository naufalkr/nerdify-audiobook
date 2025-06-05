import axios from 'axios'
import { mockBooks, mockChapters, mockUser } from './mockData'

const apiPreifx = "/api"
// Update URL untuk sesuai dengan BE-LecSens
const authApiUrl = "http://localhost:3160/api/auth"
const userApiUrl = "http://localhost:3160/api/users"

// Force menggunakan real API untuk auth
const isDevelopment = process.env.NODE_ENV === 'development'
const useMockData = isDevelopment && !process.env.REACT_APP_USE_REAL_API
const useRealAuthAPI = true // Force real API untuk auth endpoints

// Mock delay to simulate network
const mockDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

// Set up axios interceptor to include token
axios.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// Add response interceptor for debugging
axios.interceptors.response.use(
    (response) => {
        console.log('✅ Response received from:', response.config.url)
        console.log('📋 Response status:', response.status)
        console.log('📊 Response data:', response.data)
        return response
    },
    (error) => {
        console.error('❌ Request failed to:', error.config?.url)
        console.error('❌ Error status:', error.response?.status)
        console.error('❌ Error data:', error.response?.data)
        return Promise.reject(error)
    }
)

// Mock API functions
const mockAPI = {
    listAllBooks: async () => {
        await mockDelay()
        return {
            status: 200,
            data: {
                books: mockBooks
            }
        }
    },
    
    readBook: async (id) => {
        await mockDelay()
        const book = mockBooks.find(b => b.book_id === id)
        if (!book) {
            return { status: 404, data: null }
        }
        
        return {
            status: 200,
            data: {
                bookDetails: {
                    title: book.title,
                    author: book.author,
                    image_url: book.image_url
                },
                chapters: mockChapters[id] || []
            }
        }
    },
    
    search: async (params) => {
        await mockDelay()
        let results = mockBooks
        
        if (params.keyword) {
            results = mockBooks.filter(book => 
                book.title.String.toLowerCase().includes(params.keyword.toLowerCase()) ||
                book.author.String.toLowerCase().includes(params.keyword.toLowerCase())
            )
        }
        
        if (params.genre) {
            results = mockBooks.filter(book => book.genre === params.genre)
        }
        
        return {
            status: 200,
            data: results
        }
    },
    
    getCurrentUser: async () => {
        await mockDelay()
        const user = localStorage.getItem('mockUser')
        if (user) {
            return {
                status: 200,
                data: JSON.parse(user)
            }
        }
        throw new Error('No user logged in')
    },
    
    updateSeekTime: async (userId, encodedChapterURL, seekPosition) => {
        await mockDelay()
        const key = `seek_${userId}_${encodedChapterURL}`
        localStorage.setItem(key, seekPosition.toString())
        return { status: 200, data: { success: true } }
    },
    
    getSeek: async (userId, encodedChapterURL) => {
        await mockDelay()
        const key = `seek_${userId}_${encodedChapterURL}`
        const seekTime = localStorage.getItem(key) || 0
        return { status: 200, data: parseInt(seekTime) }
    },
    
    registerUser: async (userData) => {
        await mockDelay()
        // Simulate registration
        const users = JSON.parse(localStorage.getItem('mockUsers') || '[]')
        const existingUser = users.find(u => u.email === userData.email)
        
        if (existingUser) {
            throw {
                response: {
                    status: 400,
                    data: { message: 'User already exists' }
                }
            }
        }
        
        const newUser = {
            id: Date.now().toString(),
            ...userData,
            password: undefined // Don't store password
        }
        
        users.push(newUser)
        localStorage.setItem('mockUsers', JSON.stringify(users))
        
        return {
            status: 201,
            data: {
                success: true,
                message: 'User created successfully'
            }
        }
    },
    
    loginUser: async (credentials) => {
        await mockDelay()
        const users = JSON.parse(localStorage.getItem('mockUsers') || '[]')
        const user = users.find(u => u.email === credentials.email)
        
        if (!user) {
            throw {
                response: {
                    status: 401,
                    data: { message: 'Invalid email or password' }
                }
            }
        }
        
        // In real scenario, we'd verify password
        const mockToken = `mock-token-${Date.now()}`
        const userData = {
            id: user.id,
            email: user.email,
            user_name: user.user_name,
            full_name: user.full_name
        }
        
        return {
            status: 200,
            data: {
                token: mockToken,
                user: userData
            }
        }
    },
    
    getCurrentUserNew: async () => {
        await mockDelay()
        const token = localStorage.getItem('token')
        const user = localStorage.getItem('user')
        
        if (!token || !user) {
            throw new Error('No token')
        }
        
        return {
            status: 200,
            data: JSON.parse(user)
        }
    }
}

// Export functions that switch between mock and real API
export const listAllBooks = () => {
    if (useMockData) {
        console.log('🔧 Using mock data for listAllBooks')
        return mockAPI.listAllBooks()
    }
    return axios.get(`${apiPreifx}/books`)
}

export const readBook = (id) => {
    if (useMockData) {
        console.log('🔧 Using mock data for readBook:', id)
        return mockAPI.readBook(id)
    }
    return axios.get(`${apiPreifx}/books/${id}`)
}

export const search = (params) => {
    if (useMockData) {
        console.log('🔧 Using mock data for search:', params)
        return mockAPI.search(params)
    }
    return axios.get(`${apiPreifx}/search`, { params })
}

export const getCurrentUser = () => {
    if (useMockData) {
        console.log('🔧 Using mock data for getCurrentUser')
        return mockAPI.getCurrentUser()
    }
    return axios.get(`/currentuser`)
}

export const updateSeekTime = (userId, encodedChapterURL, seekPosition) => {
    if (useMockData) {
        console.log('🔧 Using mock data for updateSeekTime')
        return mockAPI.updateSeekTime(userId, encodedChapterURL, seekPosition)
    }
    return axios.post(`${apiPreifx}/user/${userId}/bookchapter/${encodedChapterURL}/seek`, seekPosition)
}

export const getSeek = (userId, encodedChapterURL) => {
    if (useMockData) {
        console.log('🔧 Using mock data for getSeek')
        return mockAPI.getSeek(userId, encodedChapterURL)
    }
    return axios.get(`${apiPreifx}/user/${userId}/bookchapter/${encodedChapterURL}/seek`)
}

// AUTH ENDPOINTS - FORCE REAL API
export const registerUser = (userData) => {
    console.log('🚀 Using REAL API for registerUser:', authApiUrl + '/register')
    console.log('📤 Sending data:', userData)
    return axios.post(`${authApiUrl}/register`, userData)
}

export const loginUser = (credentials) => {
    console.log('🚀 Using REAL API for loginUser:', authApiUrl + '/login')
    console.log('📤 Sending credentials:', { email: credentials.email, password: '***' })
    return axios.post(`${authApiUrl}/login`, credentials)
}

// USER PROFILE ENDPOINTS
export const getUserProfile = () => {
    console.log('🚀 Using REAL API for getUserProfile:', userApiUrl + '/profile')
    const token = localStorage.getItem('token')
    if (!token) return Promise.reject(new Error('No token'))
    return axios.get(`${userApiUrl}/profile`)
}

export const logoutUser = () => {
    console.log('🚀 Using REAL API for logoutUser:', userApiUrl + '/logout')
    const token = localStorage.getItem('token')
    if (!token) return Promise.reject(new Error('No token'))
    return axios.post(`${userApiUrl}/logout`)
}

// Legacy function - keep for compatibility
export const getCurrentUserNew = () => {
    if (useMockData && !useRealAuthAPI) {
        console.log('🔧 Using mock data for getCurrentUserNew')
        return mockAPI.getCurrentUserNew()
    }
    return getUserProfile() // Use the correct profile endpoint
}