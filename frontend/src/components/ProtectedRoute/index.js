import React, { useContext } from 'react'
import { Route, Redirect } from 'react-router-dom'
import { GlobalContext } from '../../contexts'

// General protected route (requires authentication)
export const ProtectedRoute = ({ children, ...rest }) => {
  const { user } = useContext(GlobalContext)
  const token = localStorage.getItem('token')

  return (
    <Route
      {...rest}
      render={({ location }) =>
        user && token ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: '/login',
              state: { from: location }
            }}
          />
        )
      }
    />
  )
}

// SuperAdmin only route
export const SuperAdminRoute = ({ children, ...rest }) => {
  const { user } = useContext(GlobalContext)
  const token = localStorage.getItem('token')
  
  // Get user data from localStorage if context is empty
  let userData = user
  if (!userData) {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        userData = JSON.parse(savedUser)
      } catch (e) {
        console.error('Error parsing saved user data:', e)
      }
    }
  }

  return (
    <Route
      {...rest}
      render={({ location }) => {
        if (!token || !userData) {
          return (
            <Redirect
              to={{
                pathname: '/login',
                state: { from: location }
              }}
            />
          )
        }

        if (userData.role !== 'SUPERADMIN') {
          return (
            <Redirect
              to={{
                pathname: '/',
                state: { error: 'Access denied. SuperAdmin privileges required.' }
              }}
            />
          )
        }

        return children
      }}
    />
  )
}

// User only route (both USER and SUPERADMIN can access)
export const UserRoute = ({ children, ...rest }) => {
  const { user } = useContext(GlobalContext)
  const token = localStorage.getItem('token')
  
  // Get user data from localStorage if context is empty
  let userData = user
  if (!userData) {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        userData = JSON.parse(savedUser)
      } catch (e) {
        console.error('Error parsing saved user data:', e)
      }
    }
  }

  return (
    <Route
      {...rest}
      render={({ location }) => {
        if (!token || !userData) {
          return (
            <Redirect
              to={{
                pathname: '/login',
                state: { from: location }
              }}
            />
          )
        }

        // Allow both USER and SUPERADMIN to access user routes
        if (!['USER', 'SUPERADMIN'].includes(userData.role)) {
          return (
            <Redirect
              to={{
                pathname: '/login',
                state: { error: 'Invalid user role.' }
              }}
            />
          )
        }

        return children
      }}
    />
  )
}