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
                            Welcome, {user?.full_name || user?.username}
                        </span>
                        <span className="admin-role-badge">SUPERADMIN</span>
                        <button onClick={handleLogout} className="admin-logout-btn">
                            Logout
                        </button>
                    </div>
                </div>
            </header>

            <main className="admin-main">
                <div className="admin-sidebar">
                    <nav className="admin-nav">
                        <Link to="/admin" className="admin-nav-link">
                            Dashboard
                        </Link>
                        <Link to="/admin/audiobooks" className="admin-nav-link active">
                            Manage Audiobooks
                        </Link>
                        <Link to="/admin/users" className="admin-nav-link">
                            Manage Users
                        </Link>
                        <Link to="/admin/analytics" className="admin-nav-link">
                            Analytics
                        </Link>
                    </nav>
                </div>

                <div className="admin-content">
                    <div className="admin-page-header">
                        <h1>Manage Audiobooks</h1>
                        <p>Upload, edit, and manage audiobook content</p>
                    </div>

                    <div className="admin-actions-bar">
                        <button className="admin-btn admin-btn-primary">
                            📚 Upload New Audiobook
                        </button>
                        <button className="admin-btn admin-btn-secondary">
                            📊 View Statistics
                        </button>
                    </div>

                    <div className="admin-content-section">
                        <h3>Recent Audiobooks</h3>
                        <div className="admin-table-container">
                            <table className="admin-table">
                                <thead>
                                    <tr>
                                        <th>Title</th>
                                        <th>Author</th>
                                        <th>Upload Date</th>
                                        <th>Status</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td colSpan="5" style={{textAlign: 'center', padding: '2rem', color: '#6c757d'}}>
                                            No audiobooks found. Upload your first audiobook to get started.
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