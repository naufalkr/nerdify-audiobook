import React, { useContext, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { GlobalContext } from '../../../contexts'
import { AudiobooksRepository } from '../../../repositories'
import '../Dashboard/admin.css'

function AdminAudiobooks() {
    const { user, logout } = useContext(GlobalContext)
    
    // State management
    const [audiobooks, setAudiobooks] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [searchQuery, setSearchQuery] = useState('')
    const [currentPage, setCurrentPage] = useState(1)
    const [totalPages, setTotalPages] = useState(1)
    const [totalItems, setTotalItems] = useState(0)

    useEffect(() => {
        document.title = "Manage Audiobooks | Admin"
        fetchAudiobooks()
    }, [currentPage])

    // Fetch audiobooks from API
    const fetchAudiobooks = async (search = '') => {
        setLoading(true)
        setError(null)
        
        try {
            const params = {
                page: currentPage,
                limit: 10
            }
            
            if (search) {
                params.search = search
            }

            const response = await AudiobooksRepository.getAllAudiobooks(params)
            
            if (response.success) {
                setAudiobooks(response.data.items || [])
                setTotalItems(response.data.total || 0)
                setTotalPages(Math.ceil((response.data.total || 0) / 10))
            } else {
                setError(response.error || 'Failed to fetch audiobooks')
                setAudiobooks([])
            }
        } catch (err) {
            setError('Network error occurred')
            setAudiobooks([])
        } finally {
            setLoading(false)
        }
    }

    // Handle search
    const handleSearch = (e) => {
        e.preventDefault()
        setCurrentPage(1)
        fetchAudiobooks(searchQuery)
    }

    // Handle page change
    const handlePageChange = (newPage) => {
        setCurrentPage(newPage)
    }

    // Handle delete audiobook
    const handleDeleteAudiobook = async (id, title) => {
        if (!window.confirm(`Are you sure you want to delete "${title}"?`)) {
            return
        }

        try {
            const response = await AudiobooksRepository.deleteAudiobook(id)
            
            if (response.success) {
                // Refresh the list
                fetchAudiobooks(searchQuery)
                alert('Audiobook deleted successfully!')
            } else {
                alert(`Error deleting audiobook: ${response.error}`)
            }
        } catch (err) {
            alert('Network error occurred while deleting audiobook')
        }
    }

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

    // Format genres display
    const formatGenres = (genres) => {
        if (!genres || genres.length === 0) return 'No genres'
        return genres.map(genre => genre.name).join(', ')
    }

    // Format date
    const formatDate = (dateString) => {
        if (!dateString) return 'N/A'
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        })
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

                    {/* Search and Actions Bar */}
                    <div className="admin-actions-bar">
                        <form onSubmit={handleSearch} style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                            <input
                                type="text"
                                placeholder="Search audiobooks..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                style={{
                                    padding: '0.5rem 1rem',
                                    border: '1px solid #d1d5db',
                                    borderRadius: '0.375rem',
                                    fontSize: '0.875rem',
                                    minWidth: '300px'
                                }}
                            />
                            <button type="submit" className="admin-btn admin-btn-secondary">
                                🔍 Search
                            </button>
                        </form>
                        
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
                    </div>

                    {/* Content Section */}
                    <div className="admin-content-section">
                        <h3>Audiobook Library</h3>
                        
                        {/* Loading State */}
                        {loading && (
                            <div style={{ textAlign: 'center', padding: '2rem' }}>
                                <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>⏳</div>
                                <div>Loading audiobooks...</div>
                            </div>
                        )}

                        {/* Error State */}
                        {error && (
                            <div style={{ 
                                textAlign: 'center', 
                                padding: '2rem', 
                                color: '#dc2626',
                                backgroundColor: '#fef2f2',
                                border: '1px solid #fecaca',
                                borderRadius: '0.5rem',
                                margin: '1rem 0'
                            }}>
                                <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>❌</div>
                                <div>Error: {error}</div>
                                <button 
                                    onClick={() => fetchAudiobooks(searchQuery)}
                                    style={{
                                        marginTop: '1rem',
                                        padding: '0.5rem 1rem',
                                        backgroundColor: '#dc2626',
                                        color: 'white',
                                        border: 'none',
                                        borderRadius: '0.375rem',
                                        cursor: 'pointer'
                                    }}
                                >
                                    Retry
                                </button>
                            </div>
                        )}

                        {/* Table */}
                        {!loading && !error && (
                            <div className="admin-table-container">
                                <table className="admin-table">
                                    <thead>
                                        <tr>
                                            <th>Title</th>
                                            <th>Author</th>
                                            <th>Reader</th>
                                            <th>Genres</th>
                                            <th>Language</th>
                                            <th>Duration</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {audiobooks.length === 0 ? (
                                            <tr>
                                                <td colSpan="8" style={{
                                                    textAlign: 'center', 
                                                    padding: '3rem', 
                                                    color: '#94a3b8',
                                                    fontSize: '1.1rem'
                                                }}>
                                                    <div style={{marginBottom: '1rem', fontSize: '3rem'}}>📚</div>
                                                    <div>No audiobooks found.</div>
                                                    <div style={{fontSize: '0.9rem', marginTop: '0.5rem'}}>
                                                        {searchQuery ? 'Try adjusting your search terms.' : 'Upload your first audiobook to get started.'}
                                                    </div>
                                                </td>
                                            </tr>
                                        ) : (
                                            audiobooks.map((audiobook) => (
                                                <tr key={audiobook.id}>
                                                    <td>
                                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                                            <img 
                                                                src={audiobook.image_url || '/assets/default-book-cover.jpg'}
                                                                alt={audiobook.title}
                                                                style={{
                                                                    width: '40px',
                                                                    height: '40px',
                                                                    objectFit: 'cover',
                                                                    borderRadius: '0.25rem'
                                                                }}
                                                                onError={(e) => {
                                                                    e.target.src = '/assets/default-book-cover.jpg'
                                                                }}
                                                            />
                                                            <div>
                                                                <div style={{ fontWeight: '600', fontSize: '0.875rem' }}>
                                                                    {audiobook.title}
                                                                </div>
                                                                <div style={{ fontSize: '0.75rem', color: '#6b7280' }}>
                                                                    ID: {audiobook.id}
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </td>
                                                    <td>{audiobook.author?.name || 'Unknown'}</td>
                                                    <td>{audiobook.reader?.name || 'Unknown'}</td>
                                                    <td>
                                                        <div style={{ fontSize: '0.75rem', maxWidth: '150px' }}>
                                                            {formatGenres(audiobook.genres)}
                                                        </div>
                                                    </td>
                                                    <td>{audiobook.language || 'N/A'}</td>
                                                    <td>{audiobook.total_duration || 'N/A'}</td>                                                    
                                                    <td>
                                                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                                                            <button 
                                                                className="admin-btn-small admin-btn-outline"
                                                                title="Edit"
                                                                style={{
                                                                    padding: '0.25rem 0.5rem',
                                                                    fontSize: '0.75rem',
                                                                    border: '1px solid #d1d5db',
                                                                    backgroundColor: 'white',
                                                                    color: '#374151',
                                                                    borderRadius: '0.25rem'
                                                                }}
                                                            >
                                                                ✏️
                                                            </button>
                                                            <button 
                                                                className="admin-btn-small admin-btn-danger"
                                                                title="Delete"
                                                                onClick={() => handleDeleteAudiobook(audiobook.id, audiobook.title)}
                                                                style={{
                                                                    padding: '0.25rem 0.5rem',
                                                                    fontSize: '0.75rem',
                                                                    border: '1px solid #dc2626',
                                                                    backgroundColor: '#dc2626',
                                                                    color: 'white',
                                                                    borderRadius: '0.25rem'
                                                                }}
                                                            >
                                                                🗑️
                                                            </button>
                                                        </div>
                                                    </td>
                                                </tr>
                                            ))
                                        )}
                                    </tbody>
                                </table>
                            </div>
                        )}

                        {/* Pagination */}
                        {!loading && !error && totalPages > 1 && (
                            <div style={{ 
                                display: 'flex', 
                                justifyContent: 'center', 
                                alignItems: 'center', 
                                gap: '1rem', 
                                marginTop: '2rem' 
                            }}>
                                <button 
                                    onClick={() => handlePageChange(currentPage - 1)}
                                    disabled={currentPage === 1}
                                    style={{
                                        padding: '0.5rem 1rem',
                                        border: '1px solid #d1d5db',
                                        backgroundColor: currentPage === 1 ? '#f9fafb' : 'white',
                                        color: currentPage === 1 ? '#9ca3af' : '#374151',
                                        borderRadius: '0.375rem',
                                        cursor: currentPage === 1 ? 'not-allowed' : 'pointer'
                                    }}
                                >
                                    ← Previous
                                </button>
                                
                                <span style={{ fontSize: '0.875rem', color: '#6b7280' }}>
                                    Page {currentPage} of {totalPages} ({totalItems} total items)
                                </span>
                                
                                <button 
                                    onClick={() => handlePageChange(currentPage + 1)}
                                    disabled={currentPage === totalPages}
                                    style={{
                                        padding: '0.5rem 1rem',
                                        border: '1px solid #d1d5db',
                                        backgroundColor: currentPage === totalPages ? '#f9fafb' : 'white',
                                        color: currentPage === totalPages ? '#9ca3af' : '#374151',
                                        borderRadius: '0.375rem',
                                        cursor: currentPage === totalPages ? 'not-allowed' : 'pointer'
                                    }}
                                >
                                    Next →
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </main>
        </div>
    )
}

export default AdminAudiobooks