import React, { createContext, useState, useEffect } from 'react'

export const GlobalContext = createContext()

export const GlobalProvider = ({ children }) => {
  const [currentAudio, setCurrentAudio] = useState({})
  const [user, setUser] = useState(null)

  // Check for existing user session on app load
  useEffect(() => {
    const checkAuthStatus = () => {
      try {
        const token = localStorage.getItem('token')
        const userData = localStorage.getItem('user')
        
        if (token && userData) {
          const parsedUser = JSON.parse(userData)
          console.log('🔄 Restoring user session:', parsedUser)
          setUser(parsedUser)
        }
      } catch (error) {
        console.error('❌ Error checking auth status:', error)
        localStorage.removeItem('token')
        localStorage.removeItem('user')
      }
    }

    checkAuthStatus()
  }, [])

  // Helper function to check if user is SUPERADMIN
  const isSuperAdmin = () => {
    return user && user.role === 'SUPERADMIN'
  }

  // Helper function to check if user is authenticated
  const isAuthenticated = () => {
    return user && localStorage.getItem('token')
  }

  // Helper function to get user role
  const getUserRole = () => {
    return user?.role || null
  }

  // Logout function
  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('mockUser')
    setUser(null)
    setCurrentAudio({})
  }

  const contextValue = {
    currentAudio,
    setCurrentAudio,
    user,
    setUser,
    logout,
    isSuperAdmin,
    isAuthenticated,
    getUserRole
  }

  return (
    <GlobalContext.Provider value={contextValue}>
      {children}
    </GlobalContext.Provider>
  )
}