class RepositoryLoggerUtil {
    static logCache = new Map()
    static THROTTLE_TIME = 1000 // 1 second throttle

    static log(repository, method, params = null, status = 'loading', duration = null, service = 'Frontend') {
        // Create a unique key for this log entry
        const logKey = `${repository}-${method}-${status}`
        const now = Date.now()
        
        // Check if we've logged this recently
        const lastLogTime = this.logCache.get(logKey)
        if (lastLogTime && (now - lastLogTime) < this.THROTTLE_TIME) {
            return // Skip this log entry
        }
        
        // Update cache
        this.logCache.set(logKey, now)
        
        // Clean old entries from cache
        if (this.logCache.size > 100) {
            const oldestTime = now - (this.THROTTLE_TIME * 2)
            for (const [key, time] of this.logCache.entries()) {
                if (time < oldestTime) {
                    this.logCache.delete(key)
                }
            }
        }

        const event = new CustomEvent('repositoryLog', {
            detail: {
                repository,
                method,
                params,
                status,
                duration,
                service
            }
        })
        window.dispatchEvent(event)
    }

    static logStart(repository, method, params = null, service = 'Frontend') {
        this.log(repository, method, params, 'loading', null, service)
        return Date.now()
    }

    static logSuccess(repository, method, startTime, params = null, service = 'Frontend') {
        const duration = Date.now() - startTime
        this.log(repository, method, params, 'success', duration, service)
    }

    static logError(repository, method, startTime, params = null, service = 'Frontend') {
        const duration = Date.now() - startTime
        this.log(repository, method, params, 'error', duration, service)
    }
}

export default RepositoryLoggerUtil