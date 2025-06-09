class SingletonLoggerUtil {
    constructor() {
        this.instances = new Map()
        this.cacheHits = 0
        this.cacheMisses = 0
        this.totalMemorySaved = 0
    }

    // Track singleton instance creation
    logInstanceCreation(className, instanceId, memorySize = 0) {
        const timestamp = Date.now()
        
        if (this.instances.has(className)) {
            // Instance already exists - this is the singleton behavior
            this.cacheHits++
            this.totalMemorySaved += memorySize
            
            this.dispatchEvent({
                type: 'instance_reused',
                className,
                instanceId,
                memorySize,
                cacheHits: this.cacheHits,
                memorySaved: this.totalMemorySaved,
                timestamp
            })
        } else {
            // New instance created
            this.instances.set(className, {
                instanceId,
                createdAt: timestamp,
                memorySize,
                accessCount: 1
            })
            this.cacheMisses++
            
            this.dispatchEvent({
                type: 'instance_created',
                className,
                instanceId,
                memorySize,
                totalInstances: this.instances.size,
                timestamp
            })
        }
    }

    // Track cache operations
    logCacheOperation(className, operation, key, hit = false, dataSize = 0) {
        const timestamp = Date.now()
        
        if (hit) {
            this.cacheHits++
            this.totalMemorySaved += dataSize
        } else {
            this.cacheMisses++
        }

        this.dispatchEvent({
            type: 'cache_operation',
            className,
            operation,
            key,
            hit,
            dataSize,
            cacheHitRate: this.getCacheHitRate(),
            memorySaved: this.totalMemorySaved,
            timestamp
        })
    }

    // Track method calls on singleton instances
    logMethodCall(className, methodName, params = null, useCache = false) {
        const startTime = Date.now()
        
        // Update access count
        if (this.instances.has(className)) {
            this.instances.get(className).accessCount++
        }

        this.dispatchEvent({
            type: 'method_start',
            className,
            methodName,
            params,
            useCache,
            startTime,
            timestamp: startTime
        })

        return startTime
    }

    logMethodEnd(className, methodName, startTime, status = 'success', result = null) {
        const endTime = Date.now()
        const duration = endTime - startTime

        this.dispatchEvent({
            type: 'method_end',
            className,
            methodName,
            status,
            duration,
            result: result ? 'data_received' : 'no_data',
            timestamp: endTime
        })
    }

    // Get statistics
    getStats() {
        return {
            totalInstances: this.instances.size,
            cacheHits: this.cacheHits,
            cacheMisses: this.cacheMisses,
            cacheHitRate: this.getCacheHitRate(),
            totalMemorySaved: this.totalMemorySaved,
            instances: Array.from(this.instances.entries()).map(([className, data]) => ({
                className,
                ...data
            }))
        }
    }

    getCacheHitRate() {
        const total = this.cacheHits + this.cacheMisses
        return total > 0 ? ((this.cacheHits / total) * 100).toFixed(1) : 0
    }

    // Dispatch custom events
    dispatchEvent(detail) {
        if (typeof window !== 'undefined') {
            window.dispatchEvent(new CustomEvent('singletonLog', { detail }))
        }
    }

    // Reset statistics
    reset() {
        this.instances.clear()
        this.cacheHits = 0
        this.cacheMisses = 0
        this.totalMemorySaved = 0
    }
}

// Export singleton instance of the logger itself
export default new SingletonLoggerUtil()