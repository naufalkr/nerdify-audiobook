import React, { useState, useEffect, useContext } from 'react'
import { useHistory } from 'react-router-dom'
import UserRepository from '../../repositories/UserRepository'
import { GlobalContext } from '../../contexts'
import './profile.css'

function Profile() {
    const history = useHistory()
    const { user } = useContext(GlobalContext)
    const [profile, setProfile] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

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

    if (error || !profile) {
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
                    <small style={{ color: '#94a3b8' }}>Powered by UserRepository Pattern</small>
                </div>

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
                                    <div className="field-value">{profile.alamat}</div>
                                </div>
                                <div className="profile-field">
                                    <label>Latitude</label>
                                    <div className="field-value">{profile.latitude}</div>
                                </div>
                                <div className="profile-field">
                                    <label>Longitude</label>
                                    <div className="field-value">{profile.longitude}</div>
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
                </div>
            </div>
        </div>
    )
}

export default Profile