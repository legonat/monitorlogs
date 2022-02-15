import React, { Component } from "react";
import DotsLoader from './dots.pulse.loader';

export default class Login extends Component {
    render() {
        return (
            
                
               
            <div>{this.props.submitted 
                 ? <div><DotsLoader color = "white"/></div>
                 :<div className="auth-wrapper">
                 <div className="auth-inner">
                 <form>
                    <h3>HSM Wallet Logs Access</h3>

                    <div className="form-group">
                        <label>Username</label>
                        <input name="user" type="text" className="form-control" placeholder="Enter username"
                        value={this.props.user} onChange={this.props.handleUsernameInput}/>
                    </div>

                    <div className="form-group">
                        <label>Password</label>
                        <input name="pass" type="password" className="form-control" placeholder="Enter password"
                        value={this.props.pass} onChange={this.props.handlePassInput}/>
                    </div>

                    <div className="form-group">
                        <div className="custom-control custom-checkbox">
                            <input type="checkbox" className="custom-control-input" id="customCheck1" defaultChecked={this.props.canRememberMe} onChange={this.props.handleRememberMeCheckbox} />
                            <label className="custom-control-label" htmlFor="customCheck1">Remember me</label>
                        </div>
                    </div>


                    <button type="submit" className="btn btn-primary btn-block" onClick={this.props.handleLoginSubmit}>Submit</button>
                    <p className="forgot-password text-right">
                        Not registered yet? <a href="/#/registration">Sign up</a>
                    </p>
                </form>
                </div>
                </div>
                }</div>
                
                
            
        );
    }
}

