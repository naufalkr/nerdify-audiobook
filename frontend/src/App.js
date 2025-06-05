import React, { useState } from 'react'
import { BrowserRouter, Route, Switch } from 'react-router-dom'

import SideMenu from './components/SideMenu'
import Player from './components/Player'

import GenreAudiobooks from './pages/GenreAudiobooks'
import ListingPage from './pages/ListingPage'
import Audiobook from './pages/Audiobook'
import Search from './pages/Search'
import Login from './pages/Login'
import Register from './pages/Register'
import Profile from './pages/Profile'

import { GlobalContext } from './contexts'

import './App.css';

function App() {

  const [currentAudio, setCurrentAudio] = useState({})
  const [user, setUser] = useState("")

  const contextValue = {
    currentAudio, setCurrentAudio,
    user, setUser
  }

  return (
    <GlobalContext.Provider value={contextValue}>
      <BrowserRouter>
        <Switch>
          {/* Auth routes - full screen dengan context */}
          <Route exact path="/login" component={Login} />
          <Route exact path="/register" component={Register} />
          <Route exact path="/profile" component={Profile} />
          
          {/* Main app routes - dengan sidebar layout */}
          <Route path="/" render={(props) => (
            <div className="main-container" data-show-player={Boolean(currentAudio?.chapter?.Link)}>
              <SideMenu {...props} />
              <Switch>
                <Route exact path="/genre/:genre" component={GenreAudiobooks} />
                <Route exact path="/search" component={Search} />
                <Route exact path="/audiobook/" component={Audiobook} />
                <Route exact path="/audiobook/:id" component={Audiobook} />
                <Route exact path="/" component={ListingPage} />
              </Switch>
              <Player />
            </div>
          )} />
        </Switch>
      </BrowserRouter>
    </GlobalContext.Provider>
  );
}

export default App;
