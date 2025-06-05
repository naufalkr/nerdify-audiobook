import React, { useState, useEffect, useContext, useCallback } from 'react'
import { Link, useHistory } from 'react-router-dom'
import { GlobalContext } from '../../../contexts'
import UserRepository from '../../../repositories/UserRepository'
import './adminUsers.css'

function AdminUsers() {
    const history = useHistory()
    const { user } = useContext(GlobalContext)
    const [users, setUsers] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')
    const [searchTerm, setSearchTerm] = useState('')
    const [currentPage, setCurrentPage] = useState(1)
    const [totalPages, setTotalPages] = useState(1)
    const [totalUsers, setTotalUsers] = useState(0)
    const [selectedRole, setSelectedRole] = useState('')
    const [selectedStatus, setSelectedStatus] = useState('')
    const [debugInfo, setDebugInfo] = useState('')

    // Modal states
    const [showEditModal, setShowEditModal] = useState(false)
    const [showDeleteModal, setShowDeleteModal] = useState(false)
    const [selectedUser, setSelectedUser] = useState(null)
    const [editFormData, setEditFormData] = useState({
        userName: '',
        email: '',
        full_name: '',
        alamat: ''
    })
    const [processing, setProcessing] = useState(false)

    const itemsPerPage = 20

    const loadUsers = useCallback(async () => {
        try {
            setLoading(true)
            setError('')
            setDebugInfo('Starting to load users...')
            
            console.log('=== AdminUsers: loadUsers START ===')
            console.log('Current user:', user)
            console.log('Token in localStorage:', localStorage.getItem('token') ? 'EXISTS' : 'NOT FOUND')
            
            const result = await UserRepository.getAllUsersForAdmin(
                currentPage, 
                itemsPerPage, 
                searchTerm
            )
            
            console.log('AdminUsers: UserRepository result:', result)
            setDebugInfo(`Received ${result.users.length} users from API (Total: ${result.total})`)
            
            // Filter by role and status if selected
            let filteredUsers = result.users
            if (selectedRole) {
                filteredUsers = filteredUsers.filter(u => u.role === selectedRole)
                setDebugInfo(prev => prev + ` | Filtered by role: ${filteredUsers.length} users`)
            }
            if (selectedStatus) {
                filteredUsers = filteredUsers.filter(u => u.status === selectedStatus)
                setDebugInfo(prev => prev + ` | Filtered by status: ${filteredUsers.length} users`)
            }

            setUsers(filteredUsers)
            setTotalUsers(result.total)
            setTotalPages(result.totalPages)
            
            console.log(`AdminUsers: Final filtered users count: ${filteredUsers.length}`)
            console.log('AdminUsers: Users data:', filteredUsers)
            console.log('=== AdminUsers: loadUsers END ===')
            
            if (result.users.length === 0 && result.total === 0) {
                setDebugInfo('No users found in API response. Check backend logs.')
            }
        } catch (error) {
            console.error('=== AdminUsers: loadUsers ERROR ===')
            console.error('Error:', error)
            console.error('Error message:', error.message)
            console.error('Error response:', error.response)
            
            setError(`Failed to load users: ${error.message}`)
            setDebugInfo(`Error: ${error.message}`)
        } finally {
            setLoading(false)
        }
    }, [currentPage, searchTerm, selectedRole, selectedStatus, user])

    useEffect(() => {
        document.title = "User Management | Admin Dashboard"
        
        console.log('=== AdminUsers: Component mounted ===')
        console.log('Current user:', user)
        console.log('User role:', user?.role)
        
        // Check if user is admin
        if (!user || (user.role !== 'ADMIN' && user.role !== 'SUPERADMIN')) {
            console.log('AdminUsers: Access denied, redirecting to /admin')
            setError('Access denied. Only ADMIN and SUPERADMIN can access this page.')
            history.push('/admin')
            return
        }

        loadUsers()
    }, [user, history, loadUsers])

    // Modal handlers
    const handleEditUser = (userItem) => {
        setSelectedUser(userItem)
        setEditFormData({
            userName: userItem.user_name || '',
            email: userItem.email || '',
            full_name: userItem.full_name || '',
            alamat: userItem.alamat || ''
        })
        setShowEditModal(true)
    }

    const handleDeleteUser = (userItem) => {
        setSelectedUser(userItem)
        setShowDeleteModal(true)
    }

    const handleEditSubmit = async (e) => {
        e.preventDefault()
        if (!selectedUser) return

        try {
            setProcessing(true)
            setError('')

            await UserRepository.updateUser(selectedUser.id, editFormData)
            
            // Refresh users list
            await loadUsers()
            
            // Close modal
            setShowEditModal(false)
            setSelectedUser(null)
            
            // Show success message
            alert('User updated successfully!')
        } catch (error) {
            console.error('Error updating user:', error)
            setError(error.message)
        } finally {
            setProcessing(false)
        }
    }

    const handleDeleteConfirm = async () => {
        if (!selectedUser) return

        try {
            setProcessing(true)
            setError('')

            await UserRepository.deleteUser(selectedUser.id)
            
            // Refresh users list
            await loadUsers()
            
            // Close modal
            setShowDeleteModal(false)
            setSelectedUser(null)
            
            // Show success message
            alert('User deleted successfully!')
        } catch (error) {
            console.error('Error deleting user:', error)
            setError(error.message)
        } finally {
            setProcessing(false)
        }
    }

    const handleCloseModal = () => {
        setShowEditModal(false)
        setShowDeleteModal(false)
        setSelectedUser(null)
        setEditFormData({
            userName: '',
            email: '',
            full_name: '',
            alamat: ''
        })
    }

    // Event handlers
    const handleSearch = (e) => {
        setSearchTerm(e.target.value)
        setCurrentPage(1) // Reset to first page
    }

    const handleRoleFilter = (e) => {
        setSelectedRole(e.target.value)
        setCurrentPage(1)
    }

    const handleStatusFilter = (e) => {
        setSelectedStatus(e.target.value)
        setCurrentPage(1)
    }

    const handlePageChange = (newPage) => {
        setCurrentPage(newPage)
    }

    const getUserInitials = (user) => {
        if (user.full_name) {
            return user.full_name
                .split(' ')
                .map(name => name.charAt(0))
                .join('')
                .toUpperCase()
                .substring(0, 2)
        }
        return user.user_name?.substring(0, 2)?.toUpperCase() || 'U'
    }

    const formatDate = (dateString) => {
        if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'Never'
        return new Date(dateString).toLocaleDateString()
    }

    return (
        <div className="admin-container">
            <header className="admin-header">
                <div className="admin-header-content">
                    <Link to="/admin" className="admin-logo-link">
                        <img 
                            src="/assets/new-logo.svg" 
                            alt="The Book Hub" 
                            className="admin-logo"
                        />
                    </Link>
                    <div className="admin-user-info">
                        <span className="admin-welcome">
                            {user?.full_name || user?.username || 'Admin'}
                        </span>
                        <span className="admin-role-badge">
                            {user?.role || 'ADMIN'}
                        </span>
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
                        <Link to="/admin/audiobooks" className="admin-nav-link">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                            </svg>
                            Manage Audiobooks
                        </Link>
                        <Link to="/admin/users" className="admin-nav-link active">
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
                        <h1>User Management</h1>
                        <p>Manage platform users, roles, and permissions</p>
                    </div>

                    {/* Error Display */}
                    {error && (
                        <div className="error-banner">
                            <strong>Error:</strong> {error}
                            <button onClick={() => setError('')} className="error-close">×</button>
                        </div>
                    )}

                    {/* Filters and Search */}
                    <div className="users-controls">
                        <div className="search-box">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                            </svg>
                            <input
                                type="text"
                                placeholder="Search users by name or email..."
                                value={searchTerm}
                                onChange={handleSearch}
                            />
                        </div>

                        <div className="filters">
                            <select value={selectedRole} onChange={handleRoleFilter}>
                                <option value="">All Roles</option>
                                <option value="SUPERADMIN">Super Admin</option>
                                <option value="ADMIN">Admin</option>
                                <option value="USER">User</option>
                            </select>

                            <select value={selectedStatus} onChange={handleStatusFilter}>
                                <option value="">All Status</option>
                                <option value="active">Active</option>
                                <option value="pending">Pending</option>
                                <option value="inactive">Inactive</option>
                            </select>

                            <button onClick={loadUsers} className="refresh-btn" disabled={loading}>
                                {loading ? 'Loading...' : 'Refresh'}
                            </button>
                        </div>
                    </div>

                    {/* Users Statistics */}
                    <div className="users-stats">
                        <div className="stat-item">
                            <span className="stat-number">{totalUsers}</span>
                            <span className="stat-label">Total Users</span>
                        </div>
                        <div className="stat-item">
                            <span className="stat-number">{users.filter(u => u.status === 'active').length}</span>
                            <span className="stat-label">Active</span>
                        </div>
                        <div className="stat-item">
                            <span className="stat-number">{users.filter(u => u.status === 'pending').length}</span>
                            <span className="stat-label">Pending</span>
                        </div>
                        <div className="stat-item">
                            <span className="stat-number">{users.filter(u => u.role === 'ADMIN').length}</span>
                            <span className="stat-label">Admins</span>
                        </div>
                    </div>

                    {/* Loading State */}
                    {loading && (
                        <div className="loading-state">
                            <div className="loading-spinner"></div>
                            <p>Loading users via UserRepository...</p>
                        </div>
                    )}

                    {/* Users Table */}
                    {!loading && (
                        <div className="users-table-container">
                            {users.length > 0 ? (
                                <table className="users-table">
                                    <thead>
                                        <tr>
                                            <th>User</th>
                                            <th>Email</th>
                                            <th>Role</th>
                                            <th>Status</th>
                                            <th>Verified</th>
                                            <th>Created</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {users.map(userItem => (
                                            <tr key={userItem.id}>
                                                <td>
                                                    <div className="user-info">
                                                        <div className="user-avatar">
                                                            {getUserInitials(userItem)}
                                                        </div>
                                                        <div className="user-details">
                                                            <div className="user-name">{userItem.full_name}</div>
                                                            <div className="user-username">@{userItem.user_name}</div>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td>{userItem.email}</td>
                                                <td>
                                                    <span style={{
                                                        color: userItem.role === 'SUPERADMIN' ? '#dc2626' : userItem.role === 'ADMIN' ? '#2563eb' : '#059669',
                                                        backgroundColor: userItem.role === 'SUPERADMIN' ? '#fef2f2' : userItem.role === 'ADMIN' ? '#eff6ff' : '#ecfdf5',
                                                        padding: '4px 8px',
                                                        borderRadius: '4px',
                                                        fontSize: '12px',
                                                        fontWeight: '500'
                                                    }}>
                                                        {userItem.role}
                                                    </span>
                                                </td>
                                                <td>
                                                    <span style={{
                                                        color: userItem.status === 'active' ? '#10b981' : userItem.status === 'pending' ? '#f59e0b' : '#ef4444',
                                                        backgroundColor: userItem.status === 'active' ? '#ecfdf5' : userItem.status === 'pending' ? '#fffbeb' : '#fef2f2',
                                                        padding: '4px 8px',
                                                        borderRadius: '4px',
                                                        fontSize: '12px',
                                                        fontWeight: '500'
                                                    }}>
                                                        {userItem.status || 'Unknown'}
                                                    </span>
                                                </td>
                                                <td>
                                                    <span className={`verification-badge ${userItem.is_verified ? 'verified' : 'unverified'}`}>
                                                        {userItem.is_verified ? 'Verified' : 'Unverified'}
                                                    </span>
                                                </td>
                                                <td>{formatDate(userItem.created_at)}</td>
                                                <td>
                                                    <div className="action-buttons">
                                                        <button 
                                                            className="action-btn edit-btn"
                                                            onClick={() => handleEditUser(userItem)}
                                                            disabled={processing}
                                                            title="Edit User"
                                                        >
                                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                                <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
                                                            </svg>
                                                        </button>
                                                        <button 
                                                            className="action-btn delete-btn"
                                                            onClick={() => handleDeleteUser(userItem)}
                                                            disabled={processing}
                                                            title="Delete User"
                                                        >
                                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                                <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                                                            </svg>
                                                        </button>
                                                    </div>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            ) : (
                                <div className="no-users">
                                    <h3>No users found</h3>
                                    <p>CORS issue resolved. Click "Test API Direct" to manually load users.</p>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Pagination */}
                    {totalPages > 1 && (
                        <div className="pagination">
                            <button
                                onClick={() => handlePageChange(currentPage - 1)}
                                disabled={currentPage === 1}
                                className="pagination-btn"
                            >
                                Previous
                            </button>
                            
                            <span className="pagination-info">
                                Page {currentPage} of {totalPages}
                            </span>
                            
                            <button
                                onClick={() => handlePageChange(currentPage + 1)}
                                disabled={currentPage === totalPages}
                                className="pagination-btn"
                            >
                                Next
                            </button>
                        </div>
                    )}
                </div>
            </main>

            {/* Edit User Modal */}
            {showEditModal && (
                <div className="modal-overlay" onClick={handleCloseModal}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h3>Edit User</h3>
                            <button className="modal-close" onClick={handleCloseModal}>×</button>
                        </div>
                        <form onSubmit={handleEditSubmit}>
                            <div className="modal-body">
                                <div className="form-group">
                                    <label htmlFor="userName">Username</label>
                                    <input
                                        type="text"
                                        id="userName"
                                        value={editFormData.userName}
                                        onChange={(e) => setEditFormData(prev => ({...prev, userName: e.target.value}))}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="email">Email</label>
                                    <input
                                        type="email"
                                        id="email"
                                        value={editFormData.email}
                                        onChange={(e) => setEditFormData(prev => ({...prev, email: e.target.value}))}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="full_name">Full Name</label>
                                    <input
                                        type="text"
                                        id="full_name"
                                        value={editFormData.full_name}
                                        onChange={(e) => setEditFormData(prev => ({...prev, full_name: e.target.value}))}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="alamat">Address</label>
                                    <textarea
                                        id="alamat"
                                        value={editFormData.alamat}
                                        onChange={(e) => setEditFormData(prev => ({...prev, alamat: e.target.value}))}
                                        rows="3"
                                    />
                                </div>
                            </div>
                            <div className="modal-footer">
                                <button type="button" className="btn-cancel" onClick={handleCloseModal}>
                                    Cancel
                                </button>
                                <button type="submit" className="btn-save" disabled={processing}>
                                    {processing ? 'Saving...' : 'Save Changes'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Delete User Modal */}
            {showDeleteModal && (
                <div className="modal-overlay" onClick={handleCloseModal}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h3>Delete User</h3>
                            <button className="modal-close" onClick={handleCloseModal}>×</button>
                        </div>
                        <div className="modal-body">
                            <p>Are you sure you want to delete this user?</p>
                            <div className="user-info-preview">
                                <strong>{selectedUser?.full_name}</strong><br/>
                                <span>{selectedUser?.email}</span><br/>
                                <span>@{selectedUser?.user_name}</span>
                            </div>
                            <p className="warning-text">This action cannot be undone.</p>
                        </div>
                        <div className="modal-footer">
                            <button type="button" className="btn-cancel" onClick={handleCloseModal}>
                                Cancel
                            </button>
                            <button 
                                type="button" 
                                className="btn-delete" 
                                onClick={handleDeleteConfirm}
                                disabled={processing}
                            >
                                {processing ? 'Deleting...' : 'Delete User'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default AdminUsers