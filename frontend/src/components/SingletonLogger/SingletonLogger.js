import React, { useState, useEffect } from 'react'
import SingletonLoggerUtil from '../../utils/singletonLogger'
import './SingletonLogger.css'

const SingletonLogger = () => {
    const [isOpen, setIsOpen] = useState(false)
    const [logs, setLogs] = useState([])
    const [stats, setStats] = useState(SingletonLoggerUtil.getStats())
    const [filter, setFilter] = useState('all')

    useEffect(() => {
        // Listen for singleton logs
        const handleSingletonLog = (event) => {
            const newLog = {
                id: Date.now() + Math.random(),
                timestamp: new Date().toLocaleTimeString(),
                ...event.detail
            }
            setLogs(prev => [newLog, ...prev].slice(0, 100)) // Keep last 100 logs
            setStats(SingletonLoggerUtil.getStats())
        }

        window.addEventListener('singletonLog', handleSingletonLog)
        return () => window.removeEventListener('singletonLog', handleSingletonLog)
    }, [])

    const filteredLogs = logs.filter(log => {
        if (filter === 'all') return true
        if (filter === 'instances') return log.type.includes('instance')
        if (filter === 'cache') return log.type.includes('cache')
        if (filter === 'methods') return log.type.includes('method')
        return log.className?.toLowerCase().includes(filter.toLowerCase())
    })

    const getTypeColor = (type) => {
        switch (type) {
            case 'instance_created': return '#4CAF50'
            case 'instance_reused': return '#2196F3'
            case 'cache_operation': return '#FF9800'
            case 'method_start': return '#9C27B0'
            case 'method_end': return '#607D8B'
            default: return '#757575'
        }
    }

    const getTypeIcon = (type) => {
        switch (type) {
            case 'instance_created': return '🆕'
            case 'instance_reused': return '♻️'
            case 'cache_operation': return '💾'
            case 'method_start': return '▶️'
            case 'method_end': return '✅'
            default: return '📋'
        }
    }

    const formatMemorySize = (bytes) => {
        if (bytes < 1024) return `${bytes}B`
        if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}KB`
        return `${(bytes / (1024 * 1024)).toFixed(1)}MB`
    }

    return (
        <>
            {/* Floating Button */}
            <div className="singleton-logger-button" onClick={() => setIsOpen(!isOpen)}>
                <span className="singleton-logger-icon">🔄</span>
                {logs.length > 0 && (
                    <span className="singleton-logger-badge">{logs.length}</span>
                )}
            </div>

            {/* Floating Panel */}
            {isOpen && (
                <div className="singleton-logger-panel">
                    <div className="singleton-logger-header">
                        <h3>🔄 Singleton Pattern Monitor</h3>
                        <div className="singleton-logger-controls">
                            <select 
                                value={filter} 
                                onChange={(e) => setFilter(e.target.value)}
                                className="singleton-logger-filter"
                            >
                                <option value="all">All Events</option>
                                <option value="instances">Instance Events</option>
                                <option value="cache">Cache Operations</option>
                                <option value="methods">Method Calls</option>
                                <option value="admin">AdminRepository</option>
                            </select>
                            <button 
                                onClick={() => {
                                    setLogs([])
                                    SingletonLoggerUtil.reset()
                                    setStats(SingletonLoggerUtil.getStats())
                                }} 
                                className="singleton-logger-clear"
                            >
                                Clear
                            </button>
                            <button 
                                onClick={() => setIsOpen(false)} 
                                className="singleton-logger-close"
                            >
                                ✕
                            </button>
                        </div>
                    </div>
                    
                    {/* Statistics Panel */}
                    <div className="singleton-logger-stats">
                        <div className="singleton-stat-item">
                            <span className="stat-label">Instances:</span>
                            <span className="stat-value">{stats.totalInstances}</span>
                        </div>
                        <div className="singleton-stat-item">
                            <span className="stat-label">Cache Hit Rate:</span>
                            <span className="stat-value">{stats.cacheHitRate}%</span>
                        </div>
                        <div className="singleton-stat-item">
                            <span className="stat-label">Memory Saved:</span>
                            <span className="stat-value">{formatMemorySize(stats.totalMemorySaved)}</span>
                        </div>
                        <div className="singleton-stat-item">
                            <span className="stat-label">Cache Hits:</span>
                            <span className="stat-value success">{stats.cacheHits}</span>
                        </div>
                    </div>
                    
                    <div className="singleton-logger-content">
                        {filteredLogs.length === 0 ? (
                            <div className="singleton-logger-empty">
                                <p>🔍 No singleton events yet...</p>
                                <p>Use admin features to see Singleton Pattern in action!</p>
                            </div>
                        ) : (
                            filteredLogs.map(log => (
                                <div key={log.id} className="singleton-logger-item">
                                    <div className="singleton-logger-item-header">
                                        <span className="singleton-logger-type-icon">
                                            {getTypeIcon(log.type)}
                                        </span>
                                        <span className="singleton-logger-class">
                                            {log.className || 'System'}
                                        </span>
                                        <span className="singleton-logger-timestamp">
                                            {log.timestamp}
                                        </span>
                                    </div>
                                    <div className="singleton-logger-item-body">
                                        <div className="singleton-logger-type">
                                            <span 
                                                className="singleton-logger-type-badge"
                                                style={{ backgroundColor: getTypeColor(log.type) }}
                                            >
                                                {log.type.replace('_', ' ').toUpperCase()}
                                            </span>
                                        </div>
                                        
                                        {log.methodName && (
                                            <div className="singleton-logger-method">
                                                <strong>Method:</strong> {log.methodName}
                                            </div>
                                        )}
                                        
                                        {log.operation && (
                                            <div className="singleton-logger-operation">
                                                <strong>Operation:</strong> {log.operation} 
                                                {log.key && `(${log.key})`}
                                                {log.hit !== undefined && (
                                                    <span className={`cache-result ${log.hit ? 'hit' : 'miss'}`}>
                                                        {log.hit ? 'HIT' : 'MISS'}
                                                    </span>
                                                )}
                                            </div>
                                        )}
                                        
                                        {log.cacheHitRate && (
                                            <div className="singleton-logger-metrics">
                                                <span>Hit Rate: {log.cacheHitRate}%</span>
                                                {log.memorySaved > 0 && (
                                                    <span>Memory Saved: {formatMemorySize(log.memorySaved)}</span>
                                                )}
                                            </div>
                                        )}
                                        
                                        {log.duration && (
                                            <div className="singleton-logger-duration">
                                                Duration: {log.duration}ms
                                            </div>
                                        )}
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>
            )}
        </>
    )
}

export default SingletonLogger