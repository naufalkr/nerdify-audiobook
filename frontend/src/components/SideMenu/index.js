import React from 'react'
import { useLocation, useHistory } from 'react-router-dom'
import { BsSearch, BsFillHouseDoorFill } from 'react-icons/bs';

import './style.css'

function SideMenu() {
    const location = useLocation()
    const history = useHistory()

    const isSelected = (path) => location.pathname === path

    return (
        <div className="side-menu">
            <div className="logo-and-name">
                <img 
                    src="/assets/new-logo.svg"
                    className="logo"
                    alt="Nerdify Audiobook"
                />
                <h3 className="org-name">Nerdify Audiobook</h3>
            </div>
            <div className="actions">
                <div 
                    className="action" 
                    onClick={() => history.push("/")}
                    data-is-selected={isSelected("/")}
                >
                       <BsFillHouseDoorFill size={"20px"}/> 
                       <span className="text">
                            Home
                       </span>
                    </div>
                <div 
                    className="action"
                    onClick={() => history.push("/search")}
                    data-is-selected={isSelected("/search")}
                >
                    <BsSearch size={"20px"}/> 
                    <span className="text">
                        Search
                    </span>
                </div>
                {/* <div className="action">Your Library</div> */}
            </div>
        </div>
    )
}

export default SideMenu