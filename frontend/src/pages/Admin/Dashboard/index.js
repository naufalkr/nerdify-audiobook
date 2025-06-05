import React, { useContext, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { GlobalContext } from '../../../contexts'
import './admin.css'

function AdminDashboard() {
    const { user, logout } = useContext(GlobalContext)

    useEffect(() => {
        document.title = "Admin Dashboard | The Book Hub"
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
                        <Link to="/admin" className="admin-nav-link active">
                            Dashboard
                        </Link>
                        <Link to="/admin/audiobooks" className="admin-nav-link">
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
                        <h1>Dashboard</h1>
                        <p>Manage your audiobook platform</p>
                    </div>

                    <div className="admin-stats-grid">
                        <div className="admin-stat-card">
                            <h3>Total Audiobooks</h3>
                            <div className="stat-number">-</div>
                            <p>Manage audiobook collection</p>
                        </div>
                        <div className="admin-stat-card">
                            <h3>Total Users</h3>
                            <div className="stat-number">-</div>
                            <p>Registered users</p>
                        </div>
                        <div className="admin-stat-card">
                            <h3>Active Sessions</h3>
                            <div className="stat-number">-</div>
                            <p>Currently listening</p>
                        </div>
                        <div className="admin-stat-card">
                            <h3>System Status</h3>
                            <div className="stat-number">✅</div>
                            <p>All systems operational</p>
                        </div>
                    </div>

                    <div className="admin-quick-actions">
                        <h2>Quick Actions</h2>
                        <div className="quick-actions-grid">
                            <Link to="/admin/audiobooks" className="quick-action-card">
                                <h3>📚 Add New Audiobook</h3>
                                <p>Upload and manage audiobook content</p>
                            </Link>
                            <Link to="/admin/users" className="quick-action-card">
                                <h3>👥 Manage Users</h3>
                                <p>View and manage user accounts</p>
                            </Link>
                            <Link to="/admin/analytics" className="quick-action-card">
                                <h3>📊 View Analytics</h3>
                                <p>Check platform usage statistics</p>
                            </Link>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    )
}

export default AdminDashboard