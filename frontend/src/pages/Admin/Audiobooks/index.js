import React, { useContext, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { GlobalContext } from '../../../contexts'
import '../Dashboard/admin.css'

function AdminAudiobooks() {
    const { user, logout } = useContext(GlobalContext)

    useEffect(() => {
        document.title = "Manage Audiobooks | Admin"
    }, [])

    const handleLogout = () => {
        logout()
        window.location.href = '/login'
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
                        <Link to="/admin" className="admin-nav-link">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>
                            </svg>
                            Dashboard
                        </Link>
                        <Link to="/admin/audiobooks" className="admin-nav-link active">
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
                    <div className="admin-page-header">
                        <h1>📚 Audiobook Management</h1>
                        <p>Upload, edit, and manage your audiobook collection</p>
                    </div>

                    <div className="admin-actions-bar">
                        <button className="admin-btn admin-btn-primary">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6z"/>
                                <polyline points="14,2 14,8 20,8"/>
                                <line x1="16" y1="13" x2="8" y2="13"/>
                                <line x1="16" y1="17" x2="8" y2="17"/>
                                <polyline points="10,9 9,9 8,9"/>
                            </svg>
                            Upload New Audiobook
                        </button>
                        <button className="admin-btn admin-btn-secondary">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
                            </svg>
                            View Statistics
                        </button>
                    </div>

                    <div className="admin-content-section">
                        <h3>📖 Audiobook Library</h3>
                        <div className="admin-table-container">
                            <table className="admin-table">
                                <thead>
                                    <tr>
                                        <th>📚 Title</th>
                                        <th>✍️ Author</th>
                                        <th>📅 Upload Date</th>
                                        <th>🔄 Status</th>
                                        <th>⚙️ Actions</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td colSpan="5" style={{
                                            textAlign: 'center', 
                                            padding: '3rem', 
                                            color: '#94a3b8',
                                            fontSize: '1.1rem'
                                        }}>
                                            <div style={{marginBottom: '1rem', fontSize: '3rem'}}>📚</div>
                                            <div>No audiobooks found.</div>
                                            <div style={{fontSize: '0.9rem', marginTop: '0.5rem'}}>
                                                Upload your first audiobook to get started.
                                            </div>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    )
}

export default AdminAudiobooks