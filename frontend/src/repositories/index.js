// Repository Pattern Export
export { default as AdminRepository } from './AdminRepository'
export { default as UserRepository } from './UserRepository'
export { default as SystemRepository } from './SystemRepository'

// Helper function to create repository instances with custom config
export const createRepositoryWithConfig = (RepositoryClass, config = {}) => {
    const instance = new RepositoryClass()
    
    // Override default config if provided
    if (config.baseURL) instance.baseURL = config.baseURL
    if (config.apiKey) instance.apiKey = config.apiKey
    
    return instance
}

// Repository pattern validation
export const validateRepositoryPattern = (repository) => {
    const requiredMethods = ['getHeaders']
    const missingMethods = requiredMethods.filter(method => 
        typeof repository[method] !== 'function'
    )
    
    if (missingMethods.length > 0) {
        throw new Error(`Repository missing required methods: ${missingMethods.join(', ')}`)
    }
    
    return true
}