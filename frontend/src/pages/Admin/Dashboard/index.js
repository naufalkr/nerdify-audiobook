import React, { useContext, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { GlobalContext } from '../../../contexts'
import AdminRepository from '../../../repositories/AdminRepository'
import './admin.css'

function AdminDashboard() {
    const contextValue = useContext(GlobalContext)
    const { user, logout } = contextValue || {}
    const [stats, setStats] = useState({
        totalAudiobooks: 0,
        totalUsers: 0,
        activeSessions: 0,
        systemStatus: 'operational'
    })
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        document.title = "Admin Dashboard | Nerdify Audiobook"
        loadStats()
    }, [])

    const loadStats = async () => {
        try {
            setLoading(true)
            setError(null)
            console.log('AdminDashboard: Loading dashboard statistics...')
            
            // Use Repository Pattern to fetch real data
            const systemStats = await AdminRepository.getRealTimeStats()
            
            setStats({
                totalAudiobooks: systemStats.totalAudiobooks || 0,
                totalUsers: systemStats.totalUsers || 0,
                activeSessions: systemStats.activeSessions || 0,
                systemStatus: systemStats.systemStatus || 'operational',
                lastUpdated: systemStats.lastUpdated,
                userGrowthRate: systemStats.userGrowthRate || 0,
                newUsersToday: systemStats.newUsersToday || 0
            })

            console.log('AdminDashboard: Statistics loaded successfully')
        } catch (error) {
            console.error('AdminDashboard: Failed to load statistics:', error)
            setError('Failed to load dashboard statistics')
            
            // Fallback to basic stats if repository fails
            setStats({
                totalAudiobooks: 0,
                totalUsers: 0,
                activeSessions: 0,
                systemStatus: 'unknown'
            })
        } finally {
            setLoading(false)
        }
    }

    const handleRefreshStats = () => {
        console.log('AdminDashboard: Manually refreshing statistics...')
        AdminRepository.clearCache()
        loadStats()
    }

    const handleLogout = () => {
        try {
            if (logout && typeof logout === 'function') {
                logout()
            } else {
                // Manual cleanup as fallback
                localStorage.removeItem('token')
                localStorage.removeItem('user')
                localStorage.removeItem('mockUser')
            }
        } catch (error) {
            console.error('Logout error:', error)
            // Force cleanup even if logout fails
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            localStorage.removeItem('mockUser')
        } finally {
            window.location.href = '/login'
        }
    }

    const getGreeting = () => {
        const hour = new Date().getHours()
        if (hour < 12) return 'Good Morning'
        if (hour < 18) return 'Good Afternoon'
        return 'Good Evening'
    }

    const getStatusIcon = (status) => {
        switch (status) {
            case 'operational':
                return '✅'
            case 'degraded':
                return '⚠️'
            case 'down':
                return '❌'
            default:
                return '❓'
        }
    }

    const formatLastUpdated = (timestamp) => {
        if (!timestamp) return 'Never'
        const date = new Date(timestamp)
        return date.toLocaleTimeString()
    }

    return (
        <div className="admin-container">
            <header className="admin-header">
                <div className="admin-header-content">
                    <img 
                        src="/assets/new-logo.svg" 
                        alt="Nerdify Audiobook" 
                        className="admin-logo"
                    />
                    <div className="admin-user-info">
                        <span className="admin-welcome">
                            {getGreeting()}, {user?.full_name || user?.username || 'Admin'}
                        </span>
                        <span className="admin-role-badge">
                            {user?.role || 'SUPERADMIN'}
                        </span>
                        <button onClick={handleLogout} className="admin-logout-btn">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{marginRight: '0.5rem'}}>
                                <path d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.59L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"/>
                            </svg>
                            Logout
                        </button>
                    </div>
                </div>
            </header>

            <main className="admin-main">
                <div className="admin-sidebar">
                    <nav className="admin-nav">
                        <Link to="/admin" className="admin-nav-link active">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>
                            </svg>
                            Dashboard
                        </Link>
                        <Link to="/admin/audiobooks" className="admin-nav-link">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                            </svg>
                            Manage Audiobooks
                        </Link>
                        <Link to="/admin/users" className="admin-nav-link">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M16 7c0-2.21-1.79-4-4-4S8 4.79 8 7s1.79 4 4 4 4-1.79 4-4zM12 13c-2.67 0-8 1.34-8 4v3h16v-3c0-2.66-5.33-4-8-4z"/>
                            </svg>
                            Manage Users
                        </Link>
                        <Link to="/admin/analytics" className="admin-nav-link">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
                            </svg>
                            Analytics
                        </Link>
                    </nav>
                </div>

                <div className="admin-content">
                    {/* Welcome Banner */}
                    <div className="admin-welcome-banner">
                        <h2>Welcome to Admin Dashboard</h2>
                        <p>Manage your audiobook platform efficiently and monitor key metrics</p>
                    </div>

                    <div className="admin-page-header">
                        <h1>Dashboard Overview</h1>
                        <p>Monitor your platform's performance and manage resources</p>
                        <div style={{ marginTop: '1rem' }}>
                            <button 
                                onClick={handleRefreshStats}
                                className="admin-btn admin-btn-secondary"
                                disabled={loading}
                                style={{ marginRight: '1rem' }}
                            >
                                {loading ? 'Loading...' : '🔄 Refresh Data'}
                            </button>
                            {stats.lastUpdated && (
                                <small style={{ color: '#94a3b8' }}>
                                    Last updated: {formatLastUpdated(stats.lastUpdated)}
                                </small>
                            )}
                        </div>
                    </div>

                    {/* Error State */}
                    {error && (
                        <div style={{
                            background: 'rgba(239, 68, 68, 0.1)',
                            border: '1px solid rgba(239, 68, 68, 0.3)',
                            borderRadius: '8px',
                            padding: '1rem',
                            marginBottom: '2rem',
                            color: '#ef4444'
                        }}>
                            <strong>Error:</strong> {error}
                        </div>
                    )}

                    {/* Loading State */}
                    {loading && (
                        <div style={{
                            background: 'rgba(59, 130, 246, 0.1)',
                            border: '1px solid rgba(59, 130, 246, 0.3)',
                            borderRadius: '8px',
                            padding: '1rem',
                            marginBottom: '2rem',
                            color: '#3b82f6'
                        }}>
                            Loading dashboard statistics...
                        </div>
                    )}

                    <div className="admin-quick-actions">
                        <h2>🚀 Quick Actions</h2>
                        <div className="quick-actions-grid">
                            <Link to="/admin/audiobooks" className="quick-action-card">
                                <h3>📚 Manage Audiobooks</h3>
                                <p>Upload new audiobooks, edit existing content, and organize your library</p>
                            </Link>
                            <Link to="/admin/users" className="quick-action-card">
                                <h3>👥 User Management</h3>
                                <p>View user accounts, manage permissions, and handle user support</p>
                            </Link>
                            <Link to="/admin/analytics" className="quick-action-card">
                                <h3>📊 Platform Analytics</h3>
                                <p>View detailed analytics, usage statistics, and performance metrics</p>
                            </Link>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    )
}

export default AdminDashboard