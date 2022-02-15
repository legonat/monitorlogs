import React, { Component } from "react";
import { Redirect } from 'react-router-dom'
import DotsLoader from './dots.pulse.loader';


export default class Logout extends Component {

handleRender = () => {
    if (!this.props.isAuthenticated) {
        return <Redirect to="/" />
    } else {
        this.props.handleLogoutSubmit()
        return <DotsLoader color="white" />
    }
}

// componentDidMount(){
// }


render() {
    return (
        <div>
            {this.handleRender()}
            {/* {this.props.handleLogoutSubmit()} */}
        </div>
    )
}
}
{/* <div className="auth-wrapper">
<div className="auth-inner">
<h3>Logout Success</h3>
</div>
</div> */}