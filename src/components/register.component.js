import React, { Component } from "react";


export default class Register extends Component {
    render() {
        return (<div className="auth-wrapper">
            <div className="auth-inner">
                <form>
                    <h3> Registration </h3>

                    <div className="form-group">
                        <label > Username </label>
                        <input name="user"
                            type="text"
                            className="form-control"
                            placeholder="Enter username"
                            value={this.props.user}
                            onChange={this.props.handleUsernameInput}
                        />
                    </div>

                    <div className="form-group" > {this.props.passwordMatch ? <label> Password </label> : <label id="passMatch">Passwords do not match</label >}
                        <input name="pass"
                            type="password"
                            className="form-control"
                            placeholder="Enter password"
                            value={this.props.pass}
                            onChange={this.props.handlePassInput}
                            required
                        />
                    </div>

                    <div className="form-group" >
                        <input name="confirmPass"
                            type="password"
                            className="form-control"
                            placeholder="Confirm password"
                            value={this.props.confirmPass}
                            onChange={this.props.handleConfirmPassInput}
                            required
                        />
                    </div>
                    <button type="submit"
                        className="btn btn-primary btn-block"
                        onClick={this.props.handleRegisterSubmit}
                        disabled={!this.props.passwordMatch}
                    >
                        Register </button>
                    <p className="forgot-password text-right" >
                        Already registered ? < a href="/#/login" > Sign in </a>
                    </p>
                </form>
            </div>
        </div>
        );
    }
}