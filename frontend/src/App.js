import React, { useState } from 'react'
import { BrowserRouter, Switch, Route } from 'react-router-dom'

// Import components
import Login from './pages/Login'
import Register from './pages/Register'
import Profile from './pages/Profile'
import GenreAudiobooks from './pages/GenreAudiobooks'
import Search from './pages/Search'
import Audiobook from './pages/Audiobook'
import ListingPage from './pages/ListingPage'
import SideMenu from './components/SideMenu'
import Player from './components/Player'

// Import route protections
import { ProtectedRoute, SuperAdminRoute, UserRoute } from './components/ProtectedRoute'

// Admin components
import AdminDashboard from './pages/Admin/Dashboard'
import AdminAudiobooks from './pages/Admin/Audiobooks'

// Import context
import { GlobalContext } from './contexts'

import './App.css'

function App() {
  const [currentAudio, setCurrentAudio] = useState({})
  const [user, setUser] = useState("")

  const contextValue = {
    currentAudio,
    setCurrentAudio,
    user,
    setUser
  }

  return (
    <GlobalContext.Provider value={contextValue}>
      <BrowserRouter>
        <Switch>
          {/* Public Auth routes */}
          <Route exact path="/login" component={Login} />
          <Route exact path="/register" component={Register} />
          
          {/* Protected Profile route */}
          <ProtectedRoute exact path="/profile">
            <Profile />
          </ProtectedRoute>

          {/* SuperAdmin only routes */}
          <SuperAdminRoute exact path="/admin">
            <AdminDashboard />
          </SuperAdminRoute>
          <SuperAdminRoute exact path="/admin/audiobooks">
            <AdminAudiobooks />
          </SuperAdminRoute>

          {/* User routes with sidebar layout */}
          <UserRoute path="/">
            <div className="main-container" data-show-player={Boolean(currentAudio?.chapter?.Link)}>
              <SideMenu />
              <Switch>
                <Route exact path="/genre/:genre" component={GenreAudiobooks} />
                <Route exact path="/search" component={Search} />
                <Route exact path="/audiobook/:id" component={Audiobook} />
                <Route exact path="/audiobook/" component={Audiobook} />
                <Route exact path="/" component={ListingPage} />
              </Switch>
              <Player />
            </div>
          </UserRoute>
        </Switch>
      </BrowserRouter>
    </GlobalContext.Provider>
  )
}

export default App