import React from 'react';
import PropTypes from 'prop-types';
import { Component, useState } from 'react';


class ToggleSwitch extends Component {

    state = {
        checked: this.props.defaultChecked
    }

    onChange = (e) => {
        this.setState({
            checked: e.target.checked
        });
        if (typeof this.props.onChange === "function") this.props.onChange();
    }

    render() {
        const [checked, setChecked] = useState(false); 

        return (
            <label class="switch">
                <input type="checkbox" className="toggle" name="toggleSwitch">
                </input>
                <span className="slider">gray</span>
            </label>
        )
    }
}


ToggleSwitch.propTypes = {

}

export default ToggleSwitch;