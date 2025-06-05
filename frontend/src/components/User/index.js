import React, { useContext, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { GlobalContext } from '../../contexts'
import { logoutUser } from '../../utils/api'
import UserRepository from '../../repositories/UserRepository'
import { handlePromise } from '../../utils/promises'
import './style.css'

function User(){
    const { user, setUser } = useContext(GlobalContext)
    const [userProfile, setUserProfile] = useState(null)
    const [showDropdown, setShowDropdown] = useState(false)
    const [loading, setLoading] = useState(false)

    const handleLogout = async () => {
        setLoading(true)
        try {
            // Call logout endpoint
            await logoutUser()
            console.log('Logout successful')
        } catch (err) {
            console.error('Logout error:', err)
            // Continue with logout even if API call fails
        } finally {
            // Clear local storage and state
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            localStorage.removeItem('mockUser')
            setUser("")
            setUserProfile(null)
            setLoading(false)
            window.location.href = "/login"
        }
    }

    useEffect(() => {
        const func = async() => {
            // Check if token exists
            const token = localStorage.getItem('token')
            if (!token) {
                console.log('No token found, user not logged in')
                return
            }

            console.log('Token found, validating user session...')

            // Check if user is already set from localStorage
            const savedUser = localStorage.getItem('user')
            if (savedUser && !user) {
                try {
                    const userData = JSON.parse(savedUser)
                    console.log('Restoring user from localStorage:', userData)
                    setUser(userData) // Store full user object
                    setUserProfile(userData)
                    return
                } catch (e) {
                    console.error('Error parsing saved user data:', e)
                }
            }

            // Get fresh user profile using UserRepository
            console.log('UserComponent: Fetching user profile via UserRepository...')
            const [response, err] = await handlePromise(UserRepository.getCurrentUser())
            
            if(err){
                console.error('UserRepository profile fetch failed:', err)
                // Token might be expired, remove it
                localStorage.removeItem('token')
                localStorage.removeItem('user')
                localStorage.removeItem('mockUser')
                setUser("")
                setUserProfile(null)
                return
            }

            console.log('UserRepository profile fetched:', response)
            
            // Handle UserRepository response
            let userData = response
            
            // Update user state and save to localStorage
            setUser(userData) // Store full user object
            setUserProfile(userData)
            localStorage.setItem('user', JSON.stringify(userData))
        }
        func()
    }, [setUser, user])

    // Get user display name
    const getDisplayName = () => {
        if (userProfile?.full_name) {
            return userProfile.full_name
        }
        if (userProfile?.user_name) {
            return userProfile.user_name
        }
        // Handle case where user is an object with email
        if (user && typeof user === 'object' && user.email) {
            return user.email.substring(0, user.email.lastIndexOf("@"))
        }
        // Handle case where user is a string (email)
        if (user && typeof user === 'string') {
            return user.substring(0, user.lastIndexOf("@"))
        }
        return "User"
    }

    // Get user initials for avatar
    const getUserInitials = () => {
        if (userProfile?.full_name) {
            return userProfile.full_name
                .split(' ')
                .map(name => name.charAt(0))
                .join('')
                .toUpperCase()
                .substring(0, 2)
        }
        if (userProfile?.user_name) {
            return userProfile.user_name.substring(0, 2).toUpperCase()
        }
        return "U"
    }

    // Check if user is logged in
    const isLoggedIn = () => {
        return user && (typeof user === 'string' || typeof user === 'object')
    }

    return (
        <div className="user">
            {
                !isLoggedIn()
                &&
                <div 
                    className="login-btn user-btn"
                    onClick={() => window.location.href = "/login"}
                >
                    Login
                </div>
            }
            {
                isLoggedIn()
                &&
                <div className="user-profile-wrapper">
                    <div 
                        className="user-profile-trigger"
                        onClick={() => setShowDropdown(!showDropdown)}
                    >
                        <div className="user-avatar">
                            {getUserInitials()}
                        </div>
                        <span className="user-name">{getDisplayName()}</span>
                        <svg className="dropdown-icon" width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
                            <path d="M6 8.5L2 4.5h8L6 8.5z"/>
                        </svg>
                    </div>
                    
                    {showDropdown && (
                        <div className="user-dropdown">
                            <div className="dropdown-header">
                                <div className="user-avatar large">
                                    {getUserInitials()}
                                </div>
                                <div className="user-info">
                                    <div className="user-full-name">{userProfile?.full_name || userProfile?.user_name}</div>
                                    <div className="user-email">{userProfile?.email}</div>
                                    <div className="user-role">{userProfile?.role}</div>
                                </div>
                            </div>
                            
                            <div className="dropdown-divider"></div>
                            
                            <div className="dropdown-menu">
                                <Link to="/profile" className="dropdown-item">
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
                                    </svg>
                                    Profile Settings
                                </Link>
                                
                                <button 
                                    className="dropdown-item logout-item"
                                    onClick={handleLogout}
                                    disabled={loading}
                                >
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.59L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"/>
                                    </svg>
                                    {loading ? 'Logging out...' : 'Logout'}
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            }
            <br/>
            <br/>
        </div>
    )
}

export default User