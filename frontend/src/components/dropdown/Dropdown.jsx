import React, { useState } from 'react';

import {ReactComponent as CogIcon} from '../../icons/cog.svg';
import {ReactComponent as ArrowIcon} from '../../icons/arrow.svg';
import {ReactComponent as ChevronIcon} from '../../icons/chevron.svg';
import {ReactComponent as PlusIcon} from '../../icons/plus.svg';

import { CSSTransition } from 'react-transition-group';


const DropdownMenu = () => {

    const [activeMenu, setActiveMenu] = useState('main');
    const [menuHeight, setMenuHeight] = useState(null);

    function calcHeight(el) {
        const height = el.offsetHeight;
        setMenuHeight(height);
    }

    const DropdownItem = (props) => {
        return (
            <a href="#" className="menu-item" onClick={() => props.goToMenu && setActiveMenu(props.goToMenu)}>
                <span className="icon-button">{props.leftIcon}</span>
                {props.children}
                <span className="icon-right">{props.rightIcon}</span>
            </a>
        )
    }

    return (
        <div className="dropdown" style={{height : menuHeight}} >
            <CSSTransition 
                in={activeMenu === 'main'}
                unmountOnExit 
                timeout={500}
                onEnter={calcHeight}
                classNames="menu-primary"
                >
                    <div className="menu">
						<DropdownItem 
						leftIcon={<PlusIcon/>}
						rightIcon={<ChevronIcon/>}
						goToMenu="game">
							Game
						</DropdownItem>
                        <DropdownItem
                            leftIcon={<CogIcon/>}
                            rightIcon={<ChevronIcon/>}
                            goToMenu="info">
								Info
                        </DropdownItem>
                    </div>
            </CSSTransition>

            <CSSTransition
                in={activeMenu === 'info'}
                unmountOnExit timeout={500}
                classNames="menu-secondary"
                onEnter={calcHeight}
                >
                    <div className="menu">
                        <DropdownItem
                            leftIcon={<ArrowIcon/>}
                            goToMenu="main">
                            Back
                        </DropdownItem>
                        <DropdownItem>
                            Help 
                        </DropdownItem>
                        <DropdownItem>
                            About 
                        </DropdownItem>
                    </div>
            </CSSTransition>

			<CSSTransition
                in={activeMenu === 'game'}
                unmountOnExit timeout={500}
                classNames="menu-secondary"
                onEnter={calcHeight}
                >
                    <div className="menu">
                        <DropdownItem
                            leftIcon={<ArrowIcon/>}
                            goToMenu="main">
                            Back
                        </DropdownItem>
                        <DropdownItem>
                            New Game 
                        </DropdownItem>
                        <DropdownItem>
                            Load Game 
                        </DropdownItem>                       
                        <DropdownItem>
                            Save Game 
                        </DropdownItem>
                        <DropdownItem>
                            Exit Game 
                        </DropdownItem>
                    </div>
            </CSSTransition>
        </div>
    )
}

export default DropdownMenu