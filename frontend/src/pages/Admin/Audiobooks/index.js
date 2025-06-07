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

    // Upload form states
    const [showUploadModal, setShowUploadModal] = useState(false)
    const [uploading, setUploading] = useState(false)
    const [uploadFormData, setUploadFormData] = useState({
        title: '',
        author_id: '',
        description: '',
        image_url: '',
        language: 'English',
        year_of_publishing: new Date().getFullYear(),
        total_duration: '',
        genre_ids: []
    })

    // Edit form states
    const [showEditModal, setShowEditModal] = useState(false)
    const [editing, setEditing] = useState(false)
    const [editingAudiobook, setEditingAudiobook] = useState(null)
    const [editFormData, setEditFormData] = useState({
        title: '',
        author_id: '',
        description: '',
        image_url: '',
        language: 'English',
        year_of_publishing: new Date().getFullYear(),
        total_duration: '',
        genre_ids: []
    })

    // Authors and Genres states
    const [authors, setAuthors] = useState([])
    const [genres, setGenres] = useState([])
    const [filteredAuthors, setFilteredAuthors] = useState([])
    const [filteredGenres, setFilteredGenres] = useState([])
    const [authorSearch, setAuthorSearch] = useState('')
    const [genreSearch, setGenreSearch] = useState('')
    const [loadingAuthors, setLoadingAuthors] = useState(false)
    const [loadingGenres, setLoadingGenres] = useState(false)

    // Edit form separate search states
    const [editAuthorSearch, setEditAuthorSearch] = useState('')
    const [editGenreSearch, setEditGenreSearch] = useState('')
    const [editFilteredAuthors, setEditFilteredAuthors] = useState([])
    const [editFilteredGenres, setEditFilteredGenres] = useState([])

    // useEffect yang benar dengan dependency yang tepat
    useEffect(() => {
        document.title = "Manage Audiobooks | Admin"
        fetchAudiobooks(searchQuery, currentPage)
    }, [currentPage])

    // Load authors and genres when modal opens
    useEffect(() => {
        if (showUploadModal || showEditModal) {
            loadAuthorsAndGenres()
        }
    }, [showUploadModal, showEditModal])

    // Filter authors based on search (Upload)
    useEffect(() => {
        if (authorSearch.trim() === '') {
            setFilteredAuthors(authors)
        } else {
            const filtered = authors.filter(author =>
                author.name.toLowerCase().includes(authorSearch.toLowerCase())
            )
            setFilteredAuthors(filtered)
        }
    }, [authorSearch, authors])

    // Filter genres based on search (Upload)
    useEffect(() => {
        if (genreSearch.trim() === '') {
            setFilteredGenres(genres)
        } else {
            const filtered = genres.filter(genre =>
                genre.name.toLowerCase().includes(genreSearch.toLowerCase())
            )
            setFilteredGenres(filtered)
        }
    }, [genreSearch, genres])

    // Filter authors based on search (Edit)
    useEffect(() => {
        if (editAuthorSearch.trim() === '') {
            setEditFilteredAuthors(authors)
        } else {
            const filtered = authors.filter(author =>
                author.name.toLowerCase().includes(editAuthorSearch.toLowerCase())
            )
            setEditFilteredAuthors(filtered)
        }
    }, [editAuthorSearch, authors])

    // Filter genres based on search (Edit)
    useEffect(() => {
        if (editGenreSearch.trim() === '') {
            setEditFilteredGenres(genres)
        } else {
            const filtered = genres.filter(genre =>
                genre.name.toLowerCase().includes(editGenreSearch.toLowerCase())
            )
            setEditFilteredGenres(filtered)
        }
    }, [editGenreSearch, genres])

    const fetchAudiobooks = async (search = '', page = 1) => {
        try {
            setLoading(true);
            setError('');
            
            console.log(`Fetching audiobooks - Search: "${search}", Page: ${page}`);
            
            const params = {
                page: page,
                limit: 10
            };
            
            if (search && search.trim()) {
                params.search = search;
            }

            const response = await AudiobooksRepository.getAllAudiobooks(params);
            
            if (response.success) {
                setAudiobooks(response.data.items || []);
                
                if (response.data.pagination) {
                    setTotalItems(response.data.pagination.total);
                    setTotalPages(response.data.pagination.total_pages);
                    console.log(`Pagination - Current: ${page}, Total Pages: ${response.data.pagination.total_pages}, Total Items: ${response.data.pagination.total}`);
                } else {
                    setTotalItems(response.data.items?.length || 0);
                    setTotalPages(1);
                }
            } else {
                setError(response.error || 'Failed to fetch audiobooks');
                setAudiobooks([]);
                setTotalItems(0);
                setTotalPages(1);
            }
        } catch (err) {
            console.error('Fetch error:', err);
            setError('Network error occurred');
            setAudiobooks([]);
            setTotalItems(0);
            setTotalPages(1);
        } finally {
            setLoading(false);
        }
    }

    // Load authors and genres for form
    const loadAuthorsAndGenres = async () => {
        try {
            setLoadingAuthors(true)
            setLoadingGenres(true)

            console.log('=== Loading Authors and Genres ===')

            // Load authors
            const authorsResponse = await AudiobooksRepository.getAllAuthors({ limit: 100 })
            console.log('Authors Response:', authorsResponse)
            
            if (authorsResponse.success) {
                const authorsData = authorsResponse.data?.data?.items || 
                              authorsResponse.data?.items || 
                              []
            
                console.log('Authors Data extracted:', authorsData)
                console.log('Authors count:', authorsData.length)
            
                setAuthors(authorsData)
                setFilteredAuthors(authorsData)
                setEditFilteredAuthors(authorsData)
            } else {
                console.error('Authors Response Failed:', authorsResponse.error)
            }

            // Load genres
            const genresResponse = await AudiobooksRepository.getAllGenres({ limit: 100 })
            console.log('Genres Response:', genresResponse)
        
            if (genresResponse.success) {
                const genresData = genresResponse.data?.data?.items || 
                             genresResponse.data?.items || 
                             []
            
                console.log('Genres Data extracted:', genresData)
                console.log('Genres count:', genresData.length)
            
                setGenres(genresData)
                setFilteredGenres(genresData)
                setEditFilteredGenres(genresData)
            } else {
                console.error('Genres Response Failed:', genresResponse.error)
            }
        
            console.log('=== Loading Complete ===')
        } catch (error) {
            console.error('Error loading authors/genres:', error)
        } finally {
            setLoadingAuthors(false)
            setLoadingGenres(false)
        }
    }

    const handleSearch = (e) => {
        e.preventDefault()
        console.log(`Search triggered with query: "${searchQuery}"`);
        setCurrentPage(1)
        fetchAudiobooks(searchQuery, 1)
    }

    const handlePageChange = (newPage) => {
        console.log(`Page change requested: ${currentPage} -> ${newPage}`);
        
        if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
            setCurrentPage(newPage);
        }
    };

    const handleDeleteAudiobook = async (id, title) => {
        if (!window.confirm(`Are you sure you want to delete "${title}"?`)) {
            return
        }

        try {
            const response = await AudiobooksRepository.deleteAudiobook(id)
            
            if (response.success) {
                fetchAudiobooks(searchQuery, currentPage)
                alert('Audiobook deleted successfully!')
            } else {
                alert(`Error deleting audiobook: ${response.error}`)
            }
        } catch (err) {
            alert('Network error occurred while deleting audiobook')
        }
    }

    // Handle upload form
    const handleUploadFormChange = (field, value) => {
        setUploadFormData(prev => ({
            ...prev,
            [field]: value
        }))
    }

    const handleGenreToggle = (genreId) => {
        setUploadFormData(prev => ({
            ...prev,
            genre_ids: prev.genre_ids.includes(genreId)
                ? prev.genre_ids.filter(id => id !== genreId)
                : [...prev.genre_ids, genreId]
        }))
    }

    const handleUploadSubmit = async (e) => {
        e.preventDefault()
        
        if (!uploadFormData.title || !uploadFormData.author_id || uploadFormData.genre_ids.length === 0) {
            alert('Please fill in all required fields')
            return
        }

        try {
            setUploading(true)
            
            const submitData = {
                ...uploadFormData,
                reader_id: 1,
                year_of_publishing: parseInt(uploadFormData.year_of_publishing)
            }

            const response = await AudiobooksRepository.createAudiobook(submitData)
            
            if (response.success) {
                alert('Audiobook uploaded successfully!')
                setShowUploadModal(false)
                resetUploadForm()
                fetchAudiobooks(searchQuery, currentPage)
            } else {
                alert(`Error uploading audiobook: ${response.error}`)
            }
        } catch (error) {
            alert('Network error occurred while uploading audiobook')
        } finally {
            setUploading(false)
        }
    }

    // Handle edit form
    const handleEditAudiobook = async (audiobook) => {
        try {
            setEditing(true)
            
            // Get detailed audiobook data
            const response = await AudiobooksRepository.getAudiobookById(audiobook.id)
            
            if (response.success) {
                const audiobookData = response.data
                setEditingAudiobook(audiobookData)
                
                // Pre-fill form with current data
                setEditFormData({
                    title: audiobookData.title || '',
                    author_id: audiobookData.author?.id || '',
                    description: audiobookData.description || '',
                    image_url: audiobookData.image_url || '',
                    language: audiobookData.language || 'English',
                    year_of_publishing: audiobookData.year_of_publishing || new Date().getFullYear(),
                    total_duration: audiobookData.total_duration || '',
                    genre_ids: audiobookData.genres?.map(genre => genre.id) || []
                })
                
                setShowEditModal(true)
            } else {
                alert(`Error loading audiobook details: ${response.error}`)
            }
        } catch (error) {
            alert('Network error occurred while loading audiobook details')
        } finally {
            setEditing(false)
        }
    }

    const handleEditFormChange = (field, value) => {
        setEditFormData(prev => ({
            ...prev,
            [field]: value
        }))
    }

    const handleEditGenreToggle = (genreId) => {
        setEditFormData(prev => ({
            ...prev,
            genre_ids: prev.genre_ids.includes(genreId)
                ? prev.genre_ids.filter(id => id !== genreId)
                : [...prev.genre_ids, genreId]
        }))
    }

    const handleEditSubmit = async (e) => {
        e.preventDefault()
        
        if (!editFormData.title || !editFormData.author_id || editFormData.genre_ids.length === 0) {
            alert('Please fill in all required fields')
            return
        }

        try {
            setEditing(true)
            
            const submitData = {
                ...editFormData,
                reader_id: 1,
                year_of_publishing: parseInt(editFormData.year_of_publishing)
            }

            const response = await AudiobooksRepository.updateAudiobook(editingAudiobook.id, submitData)
            
            if (response.success) {
                alert('Audiobook updated successfully!')
                setShowEditModal(false)
                resetEditForm()
                fetchAudiobooks(searchQuery, currentPage)
            } else {
                alert(`Error updating audiobook: ${response.error}`)
            }
        } catch (error) {
            alert('Network error occurred while updating audiobook')
        } finally {
            setEditing(false)
        }
    }

    const resetUploadForm = () => {
        setUploadFormData({
            title: '',
            author_id: '',
            description: '',
            image_url: '',
            language: 'English',
            year_of_publishing: new Date().getFullYear(),
            total_duration: '',
            genre_ids: []
        })
        setAuthorSearch('')
        setGenreSearch('')
    }

    const resetEditForm = () => {
        setEditFormData({
            title: '',
            author_id: '',
            description: '',
            image_url: '',
            language: 'English',
            year_of_publishing: new Date().getFullYear(),
            total_duration: '',
            genre_ids: []
        })
        setEditAuthorSearch('')
        setEditGenreSearch('')
        setEditingAudiobook(null)
    }

    const handleCloseUploadModal = () => {
        setShowUploadModal(false)
        resetUploadForm()
    }

    const handleCloseEditModal = () => {
        setShowEditModal(false)
        resetEditForm()
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

    const formatGenres = (genres) => {
        if (!genres || genres.length === 0) return 'No genres'
        return genres.map(genre => genre.name).join(', ')
    }

    const generatePaginationNumbers = () => {
        const delta = 2;
        const range = [];
        const rangeWithDots = [];

        for (let i = Math.max(2, currentPage - delta); 
             i <= Math.min(totalPages - 1, currentPage + delta); 
             i++) {
            range.push(i);
        }

        if (currentPage - delta > 2) {
            rangeWithDots.push(1, '...');
        } else {
            rangeWithDots.push(1);
        }

        rangeWithDots.push(...range);

        if (currentPage + delta < totalPages - 1) {
            rangeWithDots.push('...', totalPages);
        } else {
            if (totalPages > 1) {
                rangeWithDots.push(totalPages);
            }
        }

        return rangeWithDots;
    };

    return (
        <div className="admin-container">
            {/* Header sama seperti sebelumnya */}
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
                {/* Sidebar sama seperti sebelumnya */}
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
                        <h1>Audiobook Management</h1>
                        <p>Upload, edit, and manage your audiobook collection</p>
                    </div>

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
                                Search
                            </button>
                        </form>
                        
                        <button 
                            className="admin-btn admin-btn-primary"
                            onClick={() => setShowUploadModal(true)}
                        >                            
                            Upload New Audiobook
                        </button>
                    </div>

                    <div className="admin-content-section">
                        <h3>Audiobook Library ({totalItems} total)</h3>
                        
                        {loading && (
                            <div style={{ textAlign: 'center', padding: '2rem' }}>
                                <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>⏳</div>
                                <div>Loading audiobooks...</div>
                            </div>
                        )}

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
                                    onClick={() => fetchAudiobooks(searchQuery, currentPage)}
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

                        {!loading && !error && (
                            <div className="admin-table-container">
                                <table className="admin-table">
                                    <thead>
                                        <tr>
                                            <th>Title</th>
                                            <th>Author</th>
                                            <th>Year</th>
                                            <th>Language</th>
                                            <th>Genres</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {audiobooks.length === 0 ? (
                                            <tr>
                                                <td colSpan="5" style={{
                                                    textAlign: 'center', 
                                                    padding: '3rem', 
                                                    color: '#94a3b8',
                                                    fontSize: '1rem'
                                                }}>
                                                    <div style={{marginBottom: '1rem', fontSize: '2.5rem'}}>📚</div>
                                                    <div>No audiobooks found.</div>
                                                    <div style={{fontSize: '0.875rem', marginTop: '0.5rem', opacity: 0.7}}>
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
                                                                    width: '36px',
                                                                    height: '36px',
                                                                    objectFit: 'cover',
                                                                    borderRadius: '4px'
                                                                }}
                                                                onError={(e) => {
                                                                    e.target.src = '/assets/default-book-cover.jpg'
                                                                }}
                                                            />
                                                            <div>
                                                                <div style={{ fontWeight: '500', fontSize: '0.875rem', marginBottom: '2px' }}>
                                                                    {audiobook.title}
                                                                </div>
                                                                <div style={{ fontSize: '0.75rem', color: '#6b7280' }}>
                                                                    ID: {audiobook.id}
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </td>
                                                    <td>{audiobook.author?.name || 'Unknown'}</td>
                                                    <td>{audiobook.year_of_publishing || 'N/A'}</td>
                                                    <td>{audiobook.language || 'N/A'}</td>
                                                    <td>
                                                        <div style={{ fontSize: '0.75rem', maxWidth: '150px' }}>
                                                            {formatGenres(audiobook.genres)}
                                                        </div>
                                                    </td>
                                                    <td>
                                                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                                                            <button 
                                                                className="table-action-btn edit"
                                                                title="Edit audiobook"
                                                                onClick={() => handleEditAudiobook(audiobook)}
                                                                disabled={editing}
                                                            >
                                                                {editing ? (
                                                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                                        <path d="M12,4V2A10,10 0 0,0 2,12H4A8,8 0 0,1 12,4Z" />
                                                                    </svg>
                                                                ) : (
                                                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                                        <path d="M20.71,7.04C21.1,6.65 21.1,6 20.71,5.63L18.37,3.29C18,2.9 17.35,2.9 16.96,3.29L15.12,5.12L18.87,8.87M3,17.25V21H6.75L17.81,9.93L14.06,6.18L3,17.25Z" />
                                                                    </svg>
                                                                )}
                                                                Edit
                                                            </button>
                                                            <button 
                                                                className="table-action-btn delete"
                                                                title="Delete audiobook"
                                                                onClick={() => handleDeleteAudiobook(audiobook.id, audiobook.title)}
                                                            >
                                                                <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                                    <path d="M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z" />
                                                                </svg>
                                                                Delete
                                                            </button>
                                                        </div>
                                                    </td>
                                                </tr>
                                            ))
                                        )}
                                    </tbody>
                                </table>

                                {/* Pagination sama seperti sebelumnya */}
                                {totalPages > 1 && (
                                    <div className="pagination-container">
                                        <div className="pagination-info">
                                            Showing {((currentPage - 1) * 10) + 1} to {Math.min(currentPage * 10, totalItems)} of {totalItems} results
                                        </div>
                                        
                                        <div className="pagination-controls">
                                            <button 
                                                className="pagination-btn"
                                                onClick={() => handlePageChange(currentPage - 1)}
                                                disabled={currentPage === 1}
                                                type="button"
                                            >
                                                ← Previous
                                            </button>
                                            
                                            <div className="pagination-numbers">
                                                {totalPages <= 7 ? (
                                                    Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                                                        <button
                                                            key={page}
                                                            type="button"
                                                            className={`pagination-btn ${currentPage === page ? 'active' : ''}`}
                                                            onClick={() => handlePageChange(page)}
                                                        >
                                                            {page}
                                                        </button>
                                                    ))
                                                ) : (
                                                    generatePaginationNumbers().map((page, index) => (
                                                        page === '...' ? (
                                                            <span key={`dots-${index}`} style={{ padding: '0 0.5rem', color: '#6b7280' }}>
                                                                ...
                                                            </span>
                                                        ) : (
                                                            <button
                                                                key={page}
                                                                type="button"
                                                                className={`pagination-btn ${currentPage === page ? 'active' : ''}`}
                                                                onClick={() => handlePageChange(page)}
                                                            >
                                                                {page}
                                                            </button>
                                                        )
                                                    ))
                                                )}
                                            </div>
                                            
                                            <button 
                                                className="pagination-btn"
                                                onClick={() => handlePageChange(currentPage + 1)}
                                                disabled={currentPage === totalPages}
                                                type="button"
                                            >
                                                Next →
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            </main>

            {/* Upload Modal - sama seperti sebelumnya */}
            {showUploadModal && (
                <div className="modal-overlay" onClick={handleCloseUploadModal}>
                    <div className="modal-content upload-modal" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h3>
                                Upload New Audiobook
                            </h3>
                            <button className="modal-close" onClick={handleCloseUploadModal}>×</button>
                        </div>
                        
                        <form onSubmit={handleUploadSubmit}>
                            <div className="modal-body upload-form">
                                <div className="form-grid">
                                    {/* Basic Information */}
                                    <div className="form-section">
                                        <h4>Basic Information</h4>
                                        
                                        <div className="form-group">
                                            <label htmlFor="title">Title *</label>
                                            <input
                                                type="text"
                                                id="title"
                                                value={uploadFormData.title}
                                                onChange={(e) => handleUploadFormChange('title', e.target.value)}
                                                placeholder="Enter audiobook title"
                                                required
                                            />
                                        </div>

                                        <div className="form-group">
                                            <label htmlFor="description">Description</label>
                                            <textarea
                                                id="description"
                                                value={uploadFormData.description}
                                                onChange={(e) => handleUploadFormChange('description', e.target.value)}
                                                placeholder="Enter book description"
                                                rows="4"
                                            />
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label htmlFor="language">Language</label>
                                                <select
                                                    id="language"
                                                    value={uploadFormData.language}
                                                    onChange={(e) => handleUploadFormChange('language', e.target.value)}
                                                >
                                                    <option value="English">English</option>
                                                    <option value="Indonesian">Indonesian</option>
                                                    <option value="Spanish">Spanish</option>
                                                    <option value="French">French</option>
                                                    <option value="German">German</option>
                                                    <option value="Other">Other</option>
                                                </select>
                                            </div>

                                            <div className="form-group">
                                                <label htmlFor="year">Year</label>
                                                <input
                                                    type="number"
                                                    id="year"
                                                    value={uploadFormData.year_of_publishing}
                                                    onChange={(e) => handleUploadFormChange('year_of_publishing', e.target.value)}
                                                    min="1800"
                                                    max={new Date().getFullYear()}
                                                />
                                            </div>
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label htmlFor="duration">Duration</label>
                                                <input
                                                    type="text"
                                                    id="duration"
                                                    value={uploadFormData.total_duration}
                                                    onChange={(e) => handleUploadFormChange('total_duration', e.target.value)}
                                                    placeholder="e.g., 10:23 or 2hr 30min"
                                                />
                                            </div>

                                            <div className="form-group">
                                                <label htmlFor="image_url">Cover Image URL</label>
                                                <input
                                                    type="url"
                                                    id="image_url"
                                                    value={uploadFormData.image_url}
                                                    onChange={(e) => handleUploadFormChange('image_url', e.target.value)}
                                                    placeholder="https://example.com/cover.jpg"
                                                />
                                            </div>
                                        </div>
                                    </div>

                                    {/* Author Selection */}
                                    <div className="form-section">
                                        <h4>Author *</h4>
                                        
                                        <div className="form-group">
                                            <label htmlFor="author-search">Search Author</label>
                                            <input
                                                type="text"
                                                id="author-search"
                                                value={authorSearch}
                                                onChange={(e) => setAuthorSearch(e.target.value)}
                                                placeholder="Type to search authors..."
                                            />
                                        </div>

                                        <div className="selection-list author-list">
                                            {loadingAuthors ? (
                                                <div className="loading-text">Loading authors...</div>
                                            ) : (
                                                filteredAuthors.map(author => (
                                                    <div 
                                                        key={author.id}
                                                        className={`selection-item ${uploadFormData.author_id === author.id ? 'selected' : ''}`}
                                                        onClick={() => handleUploadFormChange('author_id', author.id)}
                                                    >
                                                        <span>{author.name}</span>
                                                        {uploadFormData.author_id === author.id && (
                                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                                            </svg>
                                                        )}
                                                    </div>
                                                ))
                                            )}
                                        </div>
                                    </div>

                                    {/* Genre Selection */}
                                    <div className="form-section">
                                        <h4>Genres * (Select multiple)</h4>
                                        
                                        <div className="form-group">
                                            <label htmlFor="genre-search">Search Genres</label>
                                            <input
                                                type="text"
                                                id="genre-search"
                                                value={genreSearch}
                                                onChange={(e) => setGenreSearch(e.target.value)}
                                                placeholder="Type to search genres..."
                                            />
                                        </div>

                                        <div className="selected-genres">
                                            {uploadFormData.genre_ids.map(genreId => {
                                                const genre = genres.find(g => g.id === genreId)
                                                return genre ? (
                                                    <span key={genreId} className="genre-tag">
                                                        {genre.name}
                                                        <button 
                                                            type="button"
                                                            onClick={() => handleGenreToggle(genreId)}
                                                        >
                                                            ×
                                                        </button>
                                                    </span>
                                                ) : null
                                            })}
                                        </div>

                                        <div className="selection-list genre-list">
                                            {loadingGenres ? (
                                                <div className="loading-text">Loading genres...</div>
                                            ) : (
                                                filteredGenres.map(genre => (
                                                    <div 
                                                        key={genre.id}
                                                        className={`selection-item ${uploadFormData.genre_ids.includes(genre.id) ? 'selected' : ''}`}
                                                        onClick={() => handleGenreToggle(genre.id)}
                                                    >
                                                        <span>{genre.name}</span>
                                                        {uploadFormData.genre_ids.includes(genre.id) && (
                                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                                            </svg>
                                                        )}
                                                    </div>
                                                ))
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="modal-footer">
                                <button 
                                    type="button" 
                                    className="admin-btn admin-btn-secondary" 
                                    onClick={handleCloseUploadModal}
                                    disabled={uploading}
                                >
                                    Cancel
                                </button>
                                <button 
                                    type="submit" 
                                    className="admin-btn admin-btn-primary"
                                    disabled={uploading || !uploadFormData.title || !uploadFormData.author_id || uploadFormData.genre_ids.length === 0}
                                >
                                    {uploading ? (
                                        <>
                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{marginRight: '0.5rem'}}>
                                                <path d="M12,4V2A10,10 0 0,0 2,12H4A8,8 0 0,1 12,4Z" />
                                            </svg>
                                            Uploading...
                                        </>
                                    ) : (
                                        <>
                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{marginRight: '0.5rem'}}>
                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                            </svg>
                                            Upload Audiobook
                                        </>
                                    )}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Edit Modal */}
            {showEditModal && (
                <div className="modal-overlay" onClick={handleCloseEditModal}>
                    <div className="modal-content upload-modal" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h3>
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor" style={{marginRight: '0.5rem'}}>
                                    <path d="M20.71,7.04C21.1,6.65 21.1,6 20.71,5.63L18.37,3.29C18,2.9 17.35,2.9 16.96,3.29L15.12,5.12L18.87,8.87M3,17.25V21H6.75L17.81,9.93L14.06,6.18L3,17.25Z" />
                                </svg>
                                Edit Audiobook
                            </h3>
                            <button className="modal-close" onClick={handleCloseEditModal}>×</button>
                        </div>
                        
                        <form onSubmit={handleEditSubmit}>
                            <div className="modal-body upload-form">
                                <div className="form-grid">
                                    {/* Basic Information */}
                                    <div className="form-section">
                                        <h4>Basic Information</h4>
                                        
                                        <div className="form-group">
                                            <label htmlFor="edit-title">Title *</label>
                                            <input
                                                type="text"
                                                id="edit-title"
                                                value={editFormData.title}
                                                onChange={(e) => handleEditFormChange('title', e.target.value)}
                                                placeholder="Enter audiobook title"
                                                required
                                            />
                                        </div>

                                        <div className="form-group">
                                            <label htmlFor="edit-description">Description</label>
                                            <textarea
                                                id="edit-description"
                                                value={editFormData.description}
                                                onChange={(e) => handleEditFormChange('description', e.target.value)}
                                                placeholder="Enter book description"
                                                rows="4"
                                            />
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label htmlFor="edit-language">Language</label>
                                                <select
                                                    id="edit-language"
                                                    value={editFormData.language}
                                                    onChange={(e) => handleEditFormChange('language', e.target.value)}
                                                >
                                                    <option value="English">English</option>
                                                    <option value="Indonesian">Indonesian</option>
                                                    <option value="Spanish">Spanish</option>
                                                    <option value="French">French</option>
                                                    <option value="German">German</option>
                                                    <option value="Other">Other</option>
                                                </select>
                                            </div>

                                            <div className="form-group">
                                                <label htmlFor="edit-year">Year</label>
                                                <input
                                                    type="number"
                                                    id="edit-year"
                                                    value={editFormData.year_of_publishing}
                                                    onChange={(e) => handleEditFormChange('year_of_publishing', e.target.value)}
                                                    min="1800"
                                                    max={new Date().getFullYear()}
                                                />
                                            </div>
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label htmlFor="edit-duration">Duration</label>
                                                <input
                                                    type="text"
                                                    id="edit-duration"
                                                    value={editFormData.total_duration}
                                                    onChange={(e) => handleEditFormChange('total_duration', e.target.value)}
                                                    placeholder="e.g., 10:23 or 2hr 30min"
                                                />
                                            </div>

                                            <div className="form-group">
                                                <label htmlFor="edit-image_url">Cover Image URL</label>
                                                <input
                                                    type="url"
                                                    id="edit-image_url"
                                                    value={editFormData.image_url}
                                                    onChange={(e) => handleEditFormChange('image_url', e.target.value)}
                                                    placeholder="https://example.com/cover.jpg"
                                                />
                                            </div>
                                        </div>
                                    </div>

                                    {/* Author Selection */}
                                    <div className="form-section">
                                        <h4>Author *</h4>
                                        
                                        <div className="form-group">
                                            <label htmlFor="edit-author-search">Search Author</label>
                                            <input
                                                type="text"
                                                id="edit-author-search"
                                                value={editAuthorSearch}
                                                onChange={(e) => setEditAuthorSearch(e.target.value)}
                                                placeholder="Type to search authors..."
                                            />
                                        </div>

                                        <div className="selection-list author-list">
                                            {loadingAuthors ? (
                                                <div className="loading-text">Loading authors...</div>
                                            ) : (
                                                editFilteredAuthors.map(author => (
                                                    <div 
                                                        key={author.id}
                                                        className={`selection-item ${editFormData.author_id === author.id ? 'selected' : ''}`}
                                                        onClick={() => handleEditFormChange('author_id', author.id)}
                                                    >
                                                        <span>{author.name}</span>
                                                        {editFormData.author_id === author.id && (
                                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                                            </svg>
                                                        )}
                                                    </div>
                                                ))
                                            )}
                                        </div>
                                    </div>

                                    {/* Genre Selection */}
                                    <div className="form-section">
                                        <h4>Genres * (Select multiple)</h4>
                                        
                                        <div className="form-group">
                                            <label htmlFor="edit-genre-search">Search Genres</label>
                                            <input
                                                type="text"
                                                id="edit-genre-search"
                                                value={editGenreSearch}
                                                onChange={(e) => setEditGenreSearch(e.target.value)}
                                                placeholder="Type to search genres..."
                                            />
                                        </div>

                                        <div className="selected-genres">
                                            {editFormData.genre_ids.map(genreId => {
                                                const genre = genres.find(g => g.id === genreId)
                                                return genre ? (
                                                    <span key={genreId} className="genre-tag">
                                                        {genre.name}
                                                        <button 
                                                            type="button"
                                                            onClick={() => handleEditGenreToggle(genreId)}
                                                        >
                                                            ×
                                                        </button>
                                                    </span>
                                                ) : null
                                            })}
                                        </div>

                                        <div className="selection-list genre-list">
                                            {loadingGenres ? (
                                                <div className="loading-text">Loading genres...</div>
                                            ) : (
                                                editFilteredGenres.map(genre => (
                                                    <div 
                                                        key={genre.id}
                                                        className={`selection-item ${editFormData.genre_ids.includes(genre.id) ? 'selected' : ''}`}
                                                        onClick={() => handleEditGenreToggle(genre.id)}
                                                    >
                                                        <span>{genre.name}</span>
                                                        {editFormData.genre_ids.includes(genre.id) && (
                                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                                            </svg>
                                                        )}
                                                    </div>
                                                ))
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="modal-footer">
                                <button 
                                    type="button" 
                                    className="admin-btn admin-btn-secondary" 
                                    onClick={handleCloseEditModal}
                                    disabled={editing}
                                >
                                    Cancel
                                </button>
                                <button 
                                    type="submit" 
                                    className="admin-btn admin-btn-primary"
                                    disabled={editing || !editFormData.title || !editFormData.author_id || editFormData.genre_ids.length === 0}
                                >
                                    {editing ? (
                                        <>
                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{marginRight: '0.5rem'}}>
                                                <path d="M12,4V2A10,10 0 0,0 2,12H4A8,8 0 0,1 12,4Z" />
                                            </svg>
                                            Updating...
                                        </>
                                    ) : (
                                        <>
                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{marginRight: '0.5rem'}}>
                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                            </svg>
                                            Update Audiobook
                                        </>
                                    )}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    )
}

export default AdminAudiobooks