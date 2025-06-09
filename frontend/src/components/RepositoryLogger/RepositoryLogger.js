import React, { useState, useEffect } from 'react'
import './RepositoryLogger.css'

const RepositoryLogger = () => {
    const [isOpen, setIsOpen] = useState(false)
    const [logs, setLogs] = useState([])
    const [filter, setFilter] = useState('all')

    useEffect(() => {
        // Listen for repository logs
        const handleRepositoryLog = (event) => {
            const newLog = {
                id: Date.now() + Math.random(),
                timestamp: new Date().toLocaleTimeString(),
                repository: event.detail.repository,
                method: event.detail.method,
                params: event.detail.params,
                status: event.detail.status,
                duration: event.detail.duration,
                service: event.detail.service || 'Frontend'
            }
            setLogs(prev => [newLog, ...prev].slice(0, 100)) // Keep last 100 logs
        }

        window.addEventListener('repositoryLog', handleRepositoryLog)
        return () => window.removeEventListener('repositoryLog', handleRepositoryLog)
    }, [])

    const filteredLogs = logs.filter(log => {
        if (filter === 'all') return true
        return log.repository.toLowerCase().includes(filter.toLowerCase())
    })

    const getStatusColor = (status) => {
        switch (status) {
            case 'success': return '#4CAF50'
            case 'error': return '#f44336'
            case 'loading': return '#ff9800'
            default: return '#2196F3'
        }
    }

    const getServiceBadge = (service) => {
        const colors = {
            'Frontend': '#2196F3',
            'User Management': '#9C27B0',
            'Content Management': '#FF5722',
            'Backend': '#607D8B'
        }
        return colors[service] || '#757575'
    }

    return (
        <>
            {/* Floating Button */}
            <div className="repo-logger-button" onClick={() => setIsOpen(!isOpen)}>
                <span className="repo-logger-icon">📊</span>
                {logs.length > 0 && (
                    <span className="repo-logger-badge">{logs.length}</span>
                )}
            </div>

            {/* Floating Panel */}
            {isOpen && (
                <div className="repo-logger-panel">
                    <div className="repo-logger-header">
                        <h3>🏛️ Repository Pattern Monitor</h3>
                        <div className="repo-logger-controls">
                            <select 
                                value={filter} 
                                onChange={(e) => setFilter(e.target.value)}
                                className="repo-logger-filter"
                            >
                                <option value="all">All Repositories</option>
                                <option value="admin">AdminRepository</option>
                                <option value="user">UserRepository</option>
                                <option value="audiobooks">AudiobooksRepository</option>
                                <option value="catalog">CatalogRepository</option>
                                <option value="system">SystemRepository</option>
                            </select>
                            <button 
                                onClick={() => setLogs([])} 
                                className="repo-logger-clear"
                            >
                                Clear
                            </button>
                            <button 
                                onClick={() => setIsOpen(false)} 
                                className="repo-logger-close"
                            >
                                ✕
                            </button>
                        </div>
                    </div>
                    
                    <div className="repo-logger-content">
                        {filteredLogs.length === 0 ? (
                            <div className="repo-logger-empty">
                                <p>🔍 No repository calls yet...</p>
                                <p>Start using the application to see Repository Pattern in action!</p>
                            </div>
                        ) : (
                            filteredLogs.map(log => (
                                <div key={log.id} className="repo-logger-item">
                                    <div className="repo-logger-item-header">
                                        <span 
                                            className="repo-logger-service-badge"
                                            style={{ backgroundColor: getServiceBadge(log.service) }}
                                        >
                                            {log.service}
                                        </span>
                                        <span className="repo-logger-repository">
                                            {log.repository}
                                        </span>
                                        <span className="repo-logger-timestamp">
                                            {log.timestamp}
                                        </span>
                                    </div>
                                    <div className="repo-logger-item-body">
                                        <div className="repo-logger-method">
                                            <strong>Method:</strong> {log.method}
                                        </div>
                                        {log.params && (
                                            <div className="repo-logger-params">
                                                <strong>Params:</strong> 
                                                <code>{JSON.stringify(log.params, null, 2)}</code>
                                            </div>
                                        )}
                                        <div className="repo-logger-status">
                                            <span 
                                                className="repo-logger-status-badge"
                                                style={{ backgroundColor: getStatusColor(log.status) }}
                                            >
                                                {log.status}
                                            </span>
                                            {log.duration && (
                                                <span className="repo-logger-duration">
                                                    {log.duration}ms
                                                </span>
                                            )}
                                        </div>
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

export default RepositoryLogger