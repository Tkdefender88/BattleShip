import React, { useState } from 'react';

const NavBar = (props) => {


    return (
        <nav className="navbar">
            <ul className="navbar-nav">{props.children}</ul>
        </nav>
    )
}

const NavItem = (props) => {

    const [open, setOpen] = useState(false);
    return (
        <li className="nav-item">
            <a href="#" className="icon-button" onClick={() => setOpen(!open)}>
                {props.icon}
            </a>

            {open && props.children}
        </li>
    );
}

export {
    NavBar,
    NavItem
};
