import React, { useContext } from 'react'
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
import RepositoryLogger from './components/RepositoryLogger/RepositoryLogger'
import SingletonLogger from './components/SingletonLogger/SingletonLogger'

// Import route protections
import { ProtectedRoute, SuperAdminRoute, UserRoute } from './components/ProtectedRoute'

// Admin components
import AdminDashboard from './pages/Admin/Dashboard'
import AdminAudiobooks from './pages/Admin/Audiobooks'
import AdminUsers from './pages/Admin/Users'

// Import context provider (not just the context)
import { GlobalProvider, GlobalContext } from './contexts'

import './App.css'

function AppContent() {
  const { currentAudio } = useContext(GlobalContext)
  
  return (
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
        <SuperAdminRoute exact path="/admin/users">
          <AdminUsers />
        </SuperAdminRoute>

        {/* User routes with sidebar layout */}
        <UserRoute path="/">
          <div 
            className="main-container" 
            data-show-player={currentAudio ? "true" : "false"}
          >
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
  )
}

function App() {
  return (
    <GlobalProvider>
      <AppContent />
      
      {/* Repository Logger - only show in development or demo mode */}
      {(process.env.NODE_ENV === 'development' || process.env.REACT_APP_SHOW_REPO_LOGGER === 'true') && (
        <RepositoryLogger />
      )}
      
      {/* Singleton Logger - only show in development or demo mode */}
      {(process.env.NODE_ENV === 'development' || process.env.REACT_APP_SHOW_SINGLETON_LOGGER === 'true') && (
        <SingletonLogger />
      )}
    </GlobalProvider>
  )
}

export default App