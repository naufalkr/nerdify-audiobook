import React, { useContext, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { GlobalContext } from '../../../contexts'
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

    useEffect(() => {
        document.title = "Admin Dashboard | The Book Hub"
        
        // Simulate loading stats (replace with actual API calls)
        const loadStats = () => {
            // This would be replaced with actual API calls
            setStats({
                totalAudiobooks: 127,
                totalUsers: 1453,
                activeSessions: 89,
                systemStatus: 'operational'
            })
        }
        
        loadStats()
    }, [])

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

    return (
        <div className="admin-container">
            <header className="admin-header">
                <div className="admin-header-content">
                    <img 
                        src="/assets/new-logo.svg" 
                        alt="The Book Hub" 
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
                        <h1>📊 Dashboard Overview</h1>
                        <p>Monitor your platform's performance and manage resources</p>
                    </div>

                    <div className="admin-stats-grid">
                        <div className="admin-stat-card">
                            <h3>📚 Total Audiobooks</h3>
                            <div className="stat-number">{stats.totalAudiobooks}</div>
                            <p>Published audiobook titles</p>
                        </div>
                        <div className="admin-stat-card">
                            <h3>👥 Total Users</h3>
                            <div className="stat-number">{stats.totalUsers}</div>
                            <p>Registered platform users</p>
                        </div>
                        <div className="admin-stat-card">
                            <h3>🎧 Active Sessions</h3>
                            <div className="stat-number">{stats.activeSessions}</div>
                            <p>Currently listening users</p>
                        </div>
                        <div className="admin-stat-card">
                            <h3>⚡ System Status</h3>
                            <div className="stat-number">
                                {stats.systemStatus === 'operational' ? '✅' : '⚠️'}
                            </div>
                            <p>All systems operational</p>
                        </div>
                    </div>

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