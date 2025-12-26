import { useState } from 'react'
import axios from 'axios'
import './LoginForm.css'

const API_URL = import.meta.env.VITE_API_URL || 'http://${process.env.HOST}:8080/api'

function LoginForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')
  const [messageType, setMessageType] = useState('') // 'success' or 'error'

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setMessage('')
    setMessageType('')

    try {
      const response = await axios.post(`${API_URL}/login`, {
        email,
        password,
      })

      if (response.data.success) {
        setMessageType('success')
        setMessage(response.data.message || 'Login successful!')
        // Store token if needed
        if (response.data.token) {
          localStorage.setItem('token', response.data.token)
        }
        // Reset form
        setEmail('')
        setPassword('')
      } else {
        setMessageType('error')
        setMessage(response.data.message || 'Login failed')
      }
    } catch (error) {
      setMessageType('error')
      if (error.response && error.response.data && error.response.data.message) {
        setMessage(error.response.data.message)
      } else if (error.response && error.response.data && error.response.data.error) {
        setMessage(error.response.data.error)
      } else {
        setMessage('An error occurred. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="login-form">
      <div className="form-group">
        <label htmlFor="email">Email</label>
        <input
          type="email"
          id="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="alice@example.com"
          required
          disabled={loading}
          autoComplete="email"
        />
      </div>

      <div className="form-group">
        <label htmlFor="password">Password</label>
        <input
          type="password"
          id="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Enter your password"
          required
          disabled={loading}
          autoComplete="current-password"
        />
      </div>

      {message && (
        <div className={`message ${messageType}`}>
          {message}
        </div>
      )}

      <button
        type="submit"
        className="submit-button"
        disabled={loading}
      >
        {loading ? 'Logging in...' : 'Login'}
      </button>

      <div className="info-box">
        <p className="info-title">Test Credentials:</p>
        <p className="info-text">Email: alice@example.com</p>
        <p className="info-text">Password: password123</p>
        <div className="info-divider"></div>
        <p className="info-note">
          <strong>Protection Features:</strong>
        </p>
        <ul className="info-list">
          <li>5 failed attempts → Account suspended (15 min)</li>
          <li>100 failed attempts from IP → IP blocked</li>
        </ul>
      </div>
    </form>
  )
}

export default LoginForm

