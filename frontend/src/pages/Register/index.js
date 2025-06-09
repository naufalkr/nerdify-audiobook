import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { registerUser } from '../../utils/api'
import '../Auth/auth.css'

function Register({ history }) {
    const [formData, setFormData] = useState({
        user_name: '',
        email: '',
        password: '',
        confirmPassword: '',
        full_name: '',
        alamat: '',
        latitude: -6.2088,
        longitude: 106.8456
    })
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')
    const [success, setSuccess] = useState(false)

    React.useEffect(() => {
        document.title = "Register | Nerdify Audiobook"
    }, [])

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        })
        // Clear error when user starts typing
        if (error) setError('')
    }

    const validateForm = () => {
        if (formData.password !== formData.confirmPassword) {
            setError('Passwords do not match')
            return false
        }

        if (formData.password.length < 6) {
            setError('Password must be at least 6 characters long')
            return false
        }

        if (formData.user_name.length < 3) {
            setError('Username must be at least 3 characters long')
            return false
        }

        if (formData.full_name.length < 2) {
            setError('Full name must be at least 2 characters long')
            return false
        }

        return true
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setLoading(true)
        setError('')

        if (!validateForm()) {
            setLoading(false)
            return
        }

        // Prepare data for API (exclude confirmPassword)
        const { confirmPassword, ...apiData } = formData

        try {
            const response = await registerUser(apiData)
            
            if (response.status === 201 || response.data.success) {
                setSuccess(true)
                setTimeout(() => {
                    history.push('/login')
                }, 3000)
            }
        } catch (err) {
            console.error('Registration error:', err)
            setError(
                err.response?.data?.message || 
                err.response?.data?.error ||
                'Registration failed. Please try again.'
            )
        } finally {
            setLoading(false)
        }
    }

    if (success) {
        return (
            <div className="success-container">
                <div className="success-message">
                    <h3>Registration Successful! 🎉</h3>
                    <p>Your account has been created successfully. Redirecting to login page...</p>
                </div>
            </div>
        )
    }

    return (
        <div className="auth-container">
            <div className="auth-form-wrapper">
                <div className="auth-header">
                    <img 
                        src="/assets/new-logo.svg" 
                        alt="Nerdify Audiobook" 
                        className="auth-logo"
                    />
                    <h1 className="auth-title">Join Nerdify Audiobook</h1>
                    <p className="auth-subtitle">Create your account to start your audiobook journey</p>
                </div>

                <form className="auth-form" onSubmit={handleSubmit}>
                    <div className="form-row">
                        <div className="form-group">
                            <label htmlFor="user_name" className="form-label">Username</label>
                            <input
                                id="user_name"
                                name="user_name"
                                type="text"
                                required
                                value={formData.user_name}
                                onChange={handleChange}
                                className="form-input"
                                placeholder="Choose username"
                                autoComplete="username"
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="full_name" className="form-label">Full Name</label>
                            <input
                                id="full_name"
                                name="full_name"
                                type="text"
                                required
                                value={formData.full_name}
                                onChange={handleChange}
                                className="form-input"
                                placeholder="Your full name"
                                autoComplete="name"
                            />
                        </div>
                    </div>

                    <div className="form-group">
                        <label htmlFor="email" className="form-label">Email Address</label>
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
                        <label htmlFor="alamat" className="form-label">Address</label>
                        <input
                            id="alamat"
                            name="alamat"
                            type="text"
                            required
                            value={formData.alamat}
                            onChange={handleChange}
                            className="form-input"
                            placeholder="Enter your address"
                            autoComplete="address-line1"
                        />
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label htmlFor="password" className="form-label">Password</label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                required
                                value={formData.password}
                                onChange={handleChange}
                                className="form-input"
                                placeholder="Create password"
                                autoComplete="new-password"
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="confirmPassword" className="form-label">Confirm Password</label>
                            <input
                                id="confirmPassword"
                                name="confirmPassword"
                                type="password"
                                required
                                value={formData.confirmPassword}
                                onChange={handleChange}
                                className="form-input"
                                placeholder="Confirm password"
                                autoComplete="new-password"
                            />
                        </div>
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
                                Creating account...
                            </>
                        ) : (
                            'Create Account'
                        )}
                    </button>
                </form>

                <div className="auth-link">
                    <p>
                        Already have an account?{' '}
                        <Link to="/login">Sign in here</Link>
                    </p>
                </div>
            </div>
        </div>
    )
}

export default Register