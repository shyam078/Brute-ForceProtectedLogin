import { useState } from 'react'
import LoginForm from './components/LoginForm'
import './App.css'

function App() {
  return (
    <div className="app">
      <div className="container">
        <h1 className="title">Brute-Force Protected Login</h1>
        <p className="subtitle">Secure authentication with advanced protection</p>
        <LoginForm />
      </div>
    </div>
  )
}

export default App

