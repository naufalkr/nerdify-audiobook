// filepath: d:\Kampus\SEMESTER 6\PPL\nerdify-audiobook\monolith-base\frontend\src\pages\Login\index.js
import React, { useState, useContext } from 'react'
import { Link, useHistory } from 'react-router-dom'
import { loginUser } from '../../utils/api'
import { GlobalContext } from '../../contexts'
import '../Auth/auth.css'

function Login() {
    const history = useHistory()
    const { setUser } = useContext(GlobalContext)
    const [formData, setFormData] = useState({
        email: '',
        password: ''
    })
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    React.useEffect(() => {
        document.title = "Login | The Book Hub"
    }, [])

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        })
        if (error) setError('')
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setLoading(true)
        setError('')

        try {
            console.log('🔄 Attempting login...')
            const response = await loginUser(formData)
            
            console.log('✅ Login response:', response)
            console.log('📊 Response data structure:', response.data)
            
            // BE-LecSens response structure based on your example:
            // response.data.data contains the user data and tokens
            let token = null
            let userData = null
            
            if (response.data && response.data.data) {
                const data = response.data.data
                
                // Extract token from response.data.data.access_token
                token = data.access_token
                
                // Extract user data
                userData = {
                    id: data.id,
                    email: data.email,
                    username: data.username,
                    user_name: data.username, // BE-LecSens uses 'username'
                    full_name: data.full_name,
                    role: data.role,
                    role_id: data.role_id,
                    is_verified: data.is_verified,
                    alamat: data.alamat,
                    latitude: data.latitude,
                    longitude: data.longitude
                }
            }
            
            console.log('🔑 Extracted token:', token ? 'Found' : 'Not found')
            console.log('👤 Extracted user data:', userData)
            
            if (token && userData) {
                // Store authentication data
                localStorage.setItem('token', token)
                localStorage.setItem('user', JSON.stringify(userData))
                
                // Update global context
                setUser(userData.email)
                
                console.log('✅ Login successful, redirecting to home...')
                
                // Redirect to home
                history.push('/')
            } else {
                console.error('❌ No token or user data found in response')
                console.log('🔍 Full response structure:', JSON.stringify(response.data, null, 2))
                setError('Login successful but authentication data is missing')
            }
        } catch (err) {
            console.error('❌ Login error:', err)
            console.error('❌ Error response:', err.response?.data)
            
            setError(
                err.response?.data?.message || 
                err.response?.data?.error ||
                err.message ||
                'Login failed. Please check your credentials and try again.'
            )
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="auth-container">
            <div className="auth-form-wrapper">
                <div className="auth-header">
                    <img 
                        src="/assets/new-logo.svg" 
                        alt="The Book Hub" 
                        className="auth-logo"
                    />
                    <h1 className="auth-title">Welcome Back</h1>
                    <p className="auth-subtitle">Sign in to your account to continue listening</p>
                </div>

                <form className="auth-form" onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="email" className="form-label">
                            Email Address
                        </label>
                        <input
                            id="email"
                            name="email"
                            type="email"
                            required
                            value={formData.email}
                            onChange={handleChange}
                            className="form-input"
                            placeholder="Enter your email address"
                            autoComplete="email"
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password" className="form-label">
                            Password
                        </label>
                        <input
                            id="password"
                            name="password"
                            type="password"
                            required
                            value={formData.password}
                            onChange={handleChange}
                            className="form-input"
                            placeholder="Enter your password"
                            autoComplete="current-password"
                        />
                    </div>

                    {error && (
                        <div className="error-message">{error}</div>
                    )}

                    <button
                        type="submit"
                        disabled={loading}
                        className="auth-button"
                    >
                        {loading ? (
                            <>
                                <div className="loading-spinner"></div>
                                Signing in...
                            </>
                        ) : (
                            'Sign In'
                        )}
                    </button>
                </form>

                <div className="auth-link">
                    <p>
                        Don't have an account?{' '}
                        <Link to="/register">Create one here</Link>
                    </p>
                </div>
            </div>
        </div>
    )
}

export default Login