import React, { useState, useEffect, useContext } from 'react'
import { useHistory } from 'react-router-dom'
import UserRepository from '../../repositories/UserRepository'
import { GlobalContext } from '../../contexts'
import './profile.css'

function Profile() {
    const history = useHistory()
    const { user, setUser } = useContext(GlobalContext)
    const [profile, setProfile] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')
    const [isEditing, setIsEditing] = useState(false)
    const [updating, setUpdating] = useState(false)
    const [editFormData, setEditFormData] = useState({
        userName: '',
        email: '',
        full_name: '',
        alamat: ''
    })

    useEffect(() => {
        document.title = "Profile | The Book Hub"
        
        // Redirect if not logged in
        if (!user) {
            history.push('/login')
            return
        }

        loadProfile()
    }, [user, history])

    const loadProfile = async () => {
        try {
            setLoading(true)
            console.log('Profile: Loading user profile via UserRepository...')
            
            // Use UserRepository instead of direct API call
            const userData = await UserRepository.getCurrentUser()
            
            if (userData) {
                console.log('Profile: UserRepository profile loaded:', userData)
                setProfile(userData)
                // Initialize edit form with current data
                setEditFormData({
                    userName: userData.user_name || '',
                    email: userData.email || '',
                    full_name: userData.full_name || '',
                    alamat: userData.alamat || ''
                })
            } else {
                throw new Error('No profile data received')
            }
        } catch (err) {
            console.error('Profile: Error loading profile via UserRepository:', err)
            setError('Failed to load profile')
        } finally {
            setLoading(false)
        }
    }

    const handleEditToggle = () => {
        if (isEditing) {
            // Cancel editing - reset form data
            setEditFormData({
                userName: profile.user_name || '',
                email: profile.email || '',
                full_name: profile.full_name || '',
                alamat: profile.alamat || ''
            })
        }
        setIsEditing(!isEditing)
        setError('')
    }

    const handleInputChange = (e) => {
        const { name, value } = e.target
        setEditFormData(prev => ({
            ...prev,
            [name]: value
        }))
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        
        // Validate form
        if (!editFormData.userName.trim() || !editFormData.email.trim() || !editFormData.full_name.trim()) {
            setError('Username, email, and full name are required')
            return
        }

        // Email validation
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(editFormData.email)) {
            setError('Please enter a valid email address')
            return
        }

        try {
            setUpdating(true)
            setError('')

            const updatedProfile = await UserRepository.updateUserProfile(editFormData)
            
            // Update local profile state
            const newProfileData = {
                ...profile,
                user_name: editFormData.userName,
                email: editFormData.email,
                full_name: editFormData.full_name,
                alamat: editFormData.alamat
            }
            
            setProfile(newProfileData)
            
            // Update global context and localStorage
            const updatedUser = {
                ...user,
                user_name: editFormData.userName,
                email: editFormData.email,
                full_name: editFormData.full_name,
                alamat: editFormData.alamat
            }
            setUser(updatedUser)
            localStorage.setItem('user', JSON.stringify(updatedUser))
            
            // Exit edit mode
            setIsEditing(false)
            
            // Show success message
            alert('Profile updated successfully!')
            
        } catch (error) {
            console.error('Error updating profile:', error)
            setError(error.message)
        } finally {
            setUpdating(false)
        }
    }

    const getUserInitials = () => {
        if (profile?.full_name) {
            return profile.full_name
                .split(' ')
                .map(name => name.charAt(0))
                .join('')
                .toUpperCase()
                .substring(0, 2)
        }
        if (profile?.user_name) {
            return profile.user_name.substring(0, 2).toUpperCase()
        }
        return "U"
    }

    if (loading) {
        return (
            <div className="profile-container">
                <div className="profile-loading">
                    <div className="loading-spinner"></div>
                    <p>Loading profile via UserRepository...</p>
                </div>
            </div>
        )
    }

    if (error && !profile) {
        return (
            <div className="profile-container">
                <div className="profile-error">
                    <h2>Error Loading Profile</h2>
                    <p>{error || 'Profile not found'}</p>
                    <button onClick={loadProfile} className="retry-btn">
                        Retry with UserRepository
                    </button>
                    <button onClick={() => history.push('/')} className="back-btn">
                        Go Back Home
                    </button>
                </div>
            </div>
        )
    }

    return (
        <div className="profile-container">
            <div className="profile-wrapper">
                <div className="profile-header">
                    <button onClick={() => history.goBack()} className="back-button">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z"/>
                        </svg>
                        Back
                    </button>
                    <h1>Profile Settings</h1>
                    <div className="profile-actions">
                        <button 
                            onClick={handleEditToggle} 
                            className={`edit-profile-btn ${isEditing ? 'cancel' : 'edit'}`}
                            disabled={updating}
                        >
                            {isEditing ? (
                                <>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
                                    </svg>
                                    Cancel
                                </>
                            ) : (
                                <>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
                                    </svg>
                                    Edit Profile
                                </>
                            )}
                        </button>
                    </div>
                    <small style={{ color: '#94a3b8' }}>Powered by UserRepository Pattern</small>
                </div>

                {/* Error Display */}
                {error && (
                    <div className="profile-error-banner">
                        <strong>Error:</strong> {error}
                        <button onClick={() => setError('')} className="error-close">×</button>
                    </div>
                )}

                <div className="profile-content">
                    <div className="profile-avatar-section">
                        <div className="profile-avatar">
                            {getUserInitials()}
                        </div>
                        <div className="profile-basic-info">
                            <h2>{profile.full_name}</h2>
                            <p className="profile-username">@{profile.user_name}</p>
                            <span className={`profile-role ${profile.role?.toLowerCase()}`}>
                                {profile.role}
                            </span>
                            {profile.is_verified && (
                                <span className="verified-badge">
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                                    </svg>
                                    Verified
                                </span>
                            )}
                        </div>
                    </div>

                    {isEditing ? (
                        /* Edit Form */
                        <form onSubmit={handleSubmit} className="profile-edit-form">
                            <div className="profile-section">
                                <h3>Edit Personal Information</h3>
                                <div className="profile-grid">
                                    <div className="profile-field">
                                        <label htmlFor="full_name">Full Name *</label>
                                        <input
                                            type="text"
                                            id="full_name"
                                            name="full_name"
                                            value={editFormData.full_name}
                                            onChange={handleInputChange}
                                            required
                                            className="profile-input"
                                            placeholder="Enter your full name"
                                        />
                                    </div>
                                    <div className="profile-field">
                                        <label htmlFor="userName">Username *</label>
                                        <input
                                            type="text"
                                            id="userName"
                                            name="userName"
                                            value={editFormData.userName}
                                            onChange={handleInputChange}
                                            required
                                            className="profile-input"
                                            placeholder="Enter your username"
                                        />
                                    </div>
                                    <div className="profile-field">
                                        <label htmlFor="email">Email *</label>
                                        <input
                                            type="email"
                                            id="email"
                                            name="email"
                                            value={editFormData.email}
                                            onChange={handleInputChange}
                                            required
                                            className="profile-input"
                                            placeholder="Enter your email"
                                        />
                                    </div>
                                    <div className="profile-field full-width">
                                        <label htmlFor="alamat">Address</label>
                                        <textarea
                                            id="alamat"
                                            name="alamat"
                                            value={editFormData.alamat}
                                            onChange={handleInputChange}
                                            className="profile-textarea"
                                            placeholder="Enter your address"
                                            rows="3"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div className="profile-form-actions">
                                <button 
                                    type="button" 
                                    onClick={handleEditToggle}
                                    className="btn-cancel"
                                    disabled={updating}
                                >
                                    Cancel
                                </button>
                                <button 
                                    type="submit" 
                                    className="btn-save"
                                    disabled={updating}
                                >
                                    {updating ? (
                                        <>
                                            <div className="btn-spinner"></div>
                                            Updating...
                                        </>
                                    ) : (
                                        <>
                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                            </svg>
                                            Save Changes
                                        </>
                                    )}
                                </button>
                            </div>
                        </form>
                    ) : (
                        /* View Mode */
                        <div className="profile-details">
                            <div className="profile-section">
                                <h3>Personal Information</h3>
                                <div className="profile-grid">
                                    <div className="profile-field">
                                        <label>Full Name</label>
                                        <div className="field-value">{profile.full_name}</div>
                                    </div>
                                    <div className="profile-field">
                                        <label>Username</label>
                                        <div className="field-value">{profile.user_name}</div>
                                    </div>
                                    <div className="profile-field">
                                        <label>Email</label>
                                        <div className="field-value">{profile.email}</div>
                                    </div>
                                    <div className="profile-field">
                                        <label>Role</label>
                                        <div className="field-value">{profile.role}</div>
                                    </div>
                                </div>
                            </div>

                            <div className="profile-section">
                                <h3>Location Information</h3>
                                <div className="profile-grid">
                                    <div className="profile-field full-width">
                                        <label>Address</label>
                                        <div className="field-value">{profile.alamat || 'Not specified'}</div>
                                    </div>
                                    <div className="profile-field">
                                        <label>Latitude</label>
                                        <div className="field-value">{profile.latitude || 'Not specified'}</div>
                                    </div>
                                    <div className="profile-field">
                                        <label>Longitude</label>
                                        <div className="field-value">{profile.longitude || 'Not specified'}</div>
                                    </div>
                                </div>
                            </div>

                            <div className="profile-section">
                                <h3>Account Status</h3>
                                <div className="profile-grid">
                                    <div className="profile-field">
                                        <label>Verification Status</label>
                                        <div className={`field-value status ${profile.is_verified ? 'verified' : 'unverified'}`}>
                                            {profile.is_verified ? (
                                                <>
                                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                        <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                                                    </svg>
                                                    Verified
                                                </>
                                            ) : (
                                                <>
                                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                                                    </svg>
                                                    Not Verified
                                                </>
                                            )}
                                        </div>
                                    </div>
                                    <div className="profile-field">
                                        <label>Account ID</label>
                                        <div className="field-value monospace">{profile.id}</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

export default Profile