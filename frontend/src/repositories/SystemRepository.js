import axios from 'axios'
import BaseRepository from './BaseRepository'
import SingletonLoggerUtil from '../utils/singletonLogger'

class SystemRepository extends BaseRepository {
    constructor() {
        super('SystemRepository')
        
        // Log singleton instance creation
        const instanceId = `SystemRepository_${Date.now()}`
        const estimatedMemorySize = 1024 // Estimated memory footprint in bytes
        SingletonLoggerUtil.logInstanceCreation('SystemRepository', instanceId, estimatedMemorySize)
        
        this.apiKey = process.env.REACT_APP_API_KEY || 'system-monitoring-api-key'
    }

    getHeaders() {
        const baseHeaders = super.getHeaders()
        return {
            ...baseHeaders,
            'X-API-Key': this.apiKey
        }
    }

    // Cache implementation with logging
    setCache(key, data) {
        try {
            const serialized = JSON.stringify(data)
            localStorage.setItem(`system_cache_${key}`, serialized)
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'SystemRepository',
                'set',
                key,
                false,
                serialized.length
            )
        } catch (error) {
            console.error('SystemRepository: Cache set error:', error)
        }
    }

    getFromCache(key) {
        try {
            const cached = localStorage.getItem(`system_cache_${key}`)
            const result = cached ? JSON.parse(cached) : null
            
            // Log cache operation
            SingletonLoggerUtil.logCacheOperation(
                'SystemRepository',
                'get',
                key,
                !!result,
                result ? JSON.stringify(result).length : 0
            )
            
            return result
        } catch (error) {
            console.error('SystemRepository: Cache get error:', error)
            return null
        }
    }

    // Get comprehensive system health
    async getSystemHealth() {
        const startTime = SingletonLoggerUtil.logMethodCall('SystemRepository', 'getSystemHealth')
        
        try {
            // Check cache first
            const cacheKey = 'system_health'
            const cached = this.getFromCache(cacheKey)
            if (cached && (Date.now() - cached.timestamp < 60000)) { // 1 minute cache
                SingletonLoggerUtil.logMethodEnd('SystemRepository', 'getSystemHealth', startTime, 'success', cached.data)
                return cached.data
            }
            
            const result = await this.loggedCall('getSystemHealth', async () => {
                const response = await axios.get(
                    `${this.baseURL}/api/system/health`,
                    { 
                        headers: this.getHeaders(),
                        timeout: 10000 // 10 second timeout
                    }
                )

                return {
                    status: this.normalizeStatus(response.data.status),
                    uptime: response.data.uptime || 0,
                    version: response.data.version || '1.0.0',
                    environment: response.data.environment || 'production',
                    services: this.normalizeServices(response.data.services || []),
                    lastChecked: new Date().toISOString()
                }
            })
            
            // Cache the result
            this.setCache(cacheKey, { data: result, timestamp: Date.now() })
            
            SingletonLoggerUtil.logMethodEnd('SystemRepository', 'getSystemHealth', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('SystemRepository', 'getSystemHealth', startTime, 'error', error.message)
            throw error
        }
    }

    async getPerformanceMetrics() {
        const startTime = SingletonLoggerUtil.logMethodCall('SystemRepository', 'getPerformanceMetrics')
        
        try {
            // Check cache first
            const cacheKey = 'performance_metrics'
            const cached = this.getFromCache(cacheKey)
            if (cached && (Date.now() - cached.timestamp < 30000)) { // 30 seconds cache
                SingletonLoggerUtil.logMethodEnd('SystemRepository', 'getPerformanceMetrics', startTime, 'success', cached.data)
                return cached.data
            }
            
            const result = await this.loggedCall('getPerformanceMetrics', async () => {
                const response = await axios.get(
                    `${this.baseURL}/api/system/metrics`,
                    { headers: this.getHeaders() }
                )

                return {
                    cpu: {
                        usage: response.data.cpu?.usage || 0,
                        cores: response.data.cpu?.cores || 1
                    },
                    memory: {
                        used: response.data.memory?.used || 0,
                        total: response.data.memory?.total || 0,
                        percentage: response.data.memory?.percentage || 0
                    },
                    network: {
                        inbound: response.data.network?.inbound || 0,
                        outbound: response.data.network?.outbound || 0
                    },
                    requests: {
                        total: response.data.requests?.total || 0,
                        perSecond: response.data.requests?.perSecond || 0,
                        averageResponseTime: response.data.requests?.averageResponseTime || 0
                    },
                    errors: {
                        total: response.data.errors?.total || 0,
                        rate: response.data.errors?.rate || 0
                    }
                }
            })
            
            // Cache the result
            this.setCache(cacheKey, { data: result, timestamp: Date.now() })
            
            SingletonLoggerUtil.logMethodEnd('SystemRepository', 'getPerformanceMetrics', startTime, 'success', result)
            return result
        } catch (error) {
            SingletonLoggerUtil.logMethodEnd('SystemRepository', 'getPerformanceMetrics', startTime, 'error', error.message)
            throw error
        }
    }

    // Get database health and statistics
    async getDatabaseHealth() {
        return this.loggedCall('getDatabaseHealth', async () => {
            const response = await axios.get(
                `${this.baseURL}/api/system/database`,
                { headers: this.getHeaders() }
            )

            return {
                status: this.normalizeStatus(response.data.status),
                connections: {
                    active: response.data.connections?.active || 0,
                    idle: response.data.connections?.idle || 0,
                    max: response.data.connections?.max || 100
                },
                performance: {
                    queryTime: response.data.performance?.queryTime || 0,
                    slowQueries: response.data.performance?.slowQueries || 0
                },
                size: {
                    total: response.data.size?.total || 0,
                    used: response.data.size?.used || 0,
                    percentage: response.data.size?.percentage || 0
                }
            }
        })
    }

    // Get active sessions information
    async getActiveSessions() {
        return this.loggedCall('getActiveSessions', async () => {
            const response = await axios.get(
                `${this.baseURL}/api/system/sessions`,
                { headers: this.getHeaders() }
            )

            return {
                total: response.data.total || 0,
                authenticated: response.data.authenticated || 0,
                anonymous: response.data.anonymous || 0,
                byRole: response.data.byRole || {},
                recentActivity: response.data.recentActivity || [],
                averageDuration: response.data.averageDuration || 0
            }
        })
    }

    // Get system logs
    async getSystemLogs(level = 'all', limit = 100) {
        return this.loggedCall('getSystemLogs', async () => {
            const params = new URLSearchParams({
                level,
                limit: limit.toString()
            })

            const response = await axios.get(
                `${this.baseURL}/api/system/logs?${params}`,
                { headers: this.getHeaders() }
            )

            return response.data.logs || []
        })
    }

    // Get security events
    async getSecurityEvents(limit = 50) {
        try {
            const response = await axios.get(
                `${this.baseURL}/api/system/security/events?limit=${limit}`,
                { headers: this.getHeaders() }
            )

            return response.data.events || []
        } catch (error) {
            console.warn('SystemRepository: Security events unavailable:', error.message)
            return []
        }
    }

    // Helper methods
    normalizeStatus(status) {
        if (!status) return 'unknown'
        
        const statusMap = {
            'ok': 'operational',
            'healthy': 'operational',
            'up': 'operational',
            'running': 'operational',
            'warning': 'degraded',
            'degraded': 'degraded',
            'error': 'down',
            'down': 'down',
            'offline': 'down'
        }

        return statusMap[status.toLowerCase()] || 'unknown'
    }

    normalizeServices(services) {
        return services.map(service => ({
            name: service.name || 'Unknown Service',
            status: this.normalizeStatus(service.status),
            responseTime: service.responseTime || 0,
            lastChecked: service.lastChecked || new Date().toISOString()
        }))
    }

    // Fallback data methods
    getFallbackSystemHealth() {
        return {
            status: 'operational',
            uptime: Date.now() - (Math.random() * 30 * 24 * 60 * 60 * 1000), // Random uptime up to 30 days
            version: '1.0.0',
            environment: 'production',
            services: [
                {
                    name: 'Authentication Service',
                    status: 'operational',
                    responseTime: Math.floor(Math.random() * 100 + 50),
                    lastChecked: new Date().toISOString()
                },
                {
                    name: 'User Management',
                    status: 'operational',
                    responseTime: Math.floor(Math.random() * 100 + 50),
                    lastChecked: new Date().toISOString()
                },
                {
                    name: 'Database',
                    status: 'operational',
                    responseTime: Math.floor(Math.random() * 50 + 10),
                    lastChecked: new Date().toISOString()
                }
            ],
            lastChecked: new Date().toISOString()
        }
    }

    getSimulatedMetrics() {
        return {
            cpu: {
                usage: Math.floor(Math.random() * 40 + 20), // 20-60%
                cores: 4
            },
            memory: {
                used: Math.floor(Math.random() * 4000 + 2000), // 2-6GB
                total: 8192, // 8GB
                percentage: Math.floor(Math.random() * 40 + 25) // 25-65%
            },
            network: {
                inbound: Math.floor(Math.random() * 1000 + 500), // KB/s
                outbound: Math.floor(Math.random() * 800 + 200) // KB/s
            },
            requests: {
                total: Math.floor(Math.random() * 100000 + 50000),
                perSecond: Math.floor(Math.random() * 100 + 50),
                averageResponseTime: Math.floor(Math.random() * 200 + 50) // ms
            },
            errors: {
                total: Math.floor(Math.random() * 50 + 10),
                rate: (Math.random() * 2).toFixed(2) // 0-2%
            }
        }
    }

    getFallbackDatabaseHealth() {
        return {
            status: 'operational',
            connections: {
                active: Math.floor(Math.random() * 20 + 5),
                idle: Math.floor(Math.random() * 10 + 2),
                max: 100
            },
            performance: {
                queryTime: Math.floor(Math.random() * 50 + 10), // ms
                slowQueries: Math.floor(Math.random() * 5)
            },
            size: {
                total: 1024, // MB
                used: Math.floor(Math.random() * 500 + 200), // MB
                percentage: Math.floor(Math.random() * 40 + 20) // 20-60%
            }
        }
    }

    estimateActiveSessions() {
        const hour = new Date().getHours()
        let baseSessions = 50
        
        // Peak hours adjustment
        if ((hour >= 7 && hour <= 9) || (hour >= 12 && hour <= 14) || (hour >= 19 && hour <= 22)) {
            baseSessions = Math.floor(baseSessions * 1.5)
        } else if (hour >= 22 || hour <= 6) {
            baseSessions = Math.floor(baseSessions * 0.3)
        }

        const variance = Math.floor(Math.random() * 20 - 10) // +/- 10
        const total = Math.max(baseSessions + variance, 1)

        return {
            total,
            authenticated: Math.floor(total * 0.8),
            anonymous: Math.floor(total * 0.2),
            byRole: {
                'USER': Math.floor(total * 0.7),
                'ADMIN': Math.floor(total * 0.2),
                'SUPERADMIN': Math.floor(total * 0.1)
            },
            recentActivity: [],
            averageDuration: Math.floor(Math.random() * 45 + 15) // 15-60 minutes
        }
    }

    getMockLogs() {
        const levels = ['INFO', 'WARN', 'ERROR', 'DEBUG']
        const messages = [
            'User authentication successful',
            'Database connection established',
            'Cache miss for user data',
            'API rate limit exceeded',
            'Session timeout occurred',
            'New user registration',
            'System backup completed',
            'Memory usage above threshold'
        ]

        return Array.from({ length: 10 }, (_, i) => ({
            id: `log-${Date.now()}-${i}`,
            level: levels[Math.floor(Math.random() * levels.length)],
            message: messages[Math.floor(Math.random() * messages.length)],
            timestamp: new Date(Date.now() - Math.random() * 3600000).toISOString(), // Last hour
            source: 'system'
        }))
    }

    // Cache management
    clearSystemCache() {
        const keys = Object.keys(localStorage).filter(key => key.startsWith('system_cache_'))
        keys.forEach(key => {
            localStorage.removeItem(key)
            // Log cache clear operation
            SingletonLoggerUtil.logCacheOperation(
                'SystemRepository',
                'clear',
                key.replace('system_cache_', ''),
                false,
                0
            )
        })
        console.log('SystemRepository: System cache cleared')
    }
}

export default new SystemRepository()