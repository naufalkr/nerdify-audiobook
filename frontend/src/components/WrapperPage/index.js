import React from 'react'

function WrapperPage({ children, style = {} }) {
    return (
        <div 
            className="rest-page" 
            style={{
                // Ensure no overflow properties are added
                ...style
            }}
        >
            {children}
        </div>
    )
}

export default WrapperPage