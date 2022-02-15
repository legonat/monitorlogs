import React, { Component } from "react";
import { Navbar, Nav, Form, FormControl, Button, InputGroup } from "react-bootstrap";
import { faTimesCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { AuthContext } from '../App'

export default class NavbarContainer extends Component {
  render() {
    return (
      <AuthContext.Consumer>
        { value =>
          !value.isAuthenticated
            ? <div></div>
            : <Navbar bg="primary" variant="dark" sticky="top" expand="lg" className="navbar" >
              <Nav className="col mr-4 justify-content-start">
                <Navbar.Brand href="/">HSM Logs</Navbar.Brand>
              </Nav>
              <Navbar.Toggle aria-controls="basic-navbar-nav" />
              <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="col-auto justify-content-space-between">
                  <Form className="search-bar">
                    <InputGroup className="nav-group">
                      <FormControl
                        onChange={this.props.handleSearchInput}
                        value={this.props.searchText}
                        type="text"
                        placeholder="Search All"
                        className="nav-form"
                      />
                      <InputGroup.Append>
                        <Button variant="outline-secondary" onClick={this.props.handleSearchInput} type="button" className="nav-btn-middle">
                          <FontAwesomeIcon icon={faTimesCircle} />
                        </Button>
                        <Button onClick={this.props.handleSearchSubmit} variant="outline-secondary" className="nav-btn" type="submit">
                          Search
                    </Button>
                      </InputGroup.Append>
                    </InputGroup>
                  </Form>
                </Nav>
                {/* <Nav className="col-auto justify-content-space-between">
                  <Form className="nav-filter">
                    <FormControl
                      onChange={this.props.handleFilterTextInput}
                      value={this.props.filterText}
                      type="text"
                      placeholder="Filter..."
                      className="nav-form col-auto"
                    />
                  </Form>
                </Nav> */}
                {/* <Nav className="col-auto justify-content-space-between">
                  <Form className="nav-dropdown">
                    <InputGroup>
                      <DropdownSelect
                        services={this.props.services}
                        value={this.props.value}
                        handleSelectServiceChange={this.props.handleSelectServiceChange}
                        class="dropdown"
                      />
                      <InputGroup.Append>
                        <Button variant="outline-secondary" onClick={this.props.handleSelectServiceChange} className="nav-btn-middle">
                          <FontAwesomeIcon icon={faTimesCircle} />
                        </Button>
                      </InputGroup.Append>
                    </InputGroup>
                  </Form>
                </Nav> */}
                {/* <Nav className="col-auto justify-content-end"> */}
                <Nav className="col ml-4 justify-content-end">
                  <Nav.Link href="/#/logout">Logout</Nav.Link>
                  {/* </Nav> */}
                  {/* <Nav className="col-md-3 ml-auto">
                        <Nav.Link href="/login">Logout</Nav.Link>
                        <Nav.Link href="/register">Register</Nav.Link>
                  </Nav> */}
                </Nav>
              </Navbar.Collapse>
            </Navbar>
        }
      </AuthContext.Consumer>
    )
  }
}


class SelectRow extends Component {
  render() {
    const service = this.props.service
    const serviceInfo = service.service;
    const id = service.id;

    return (
      <option>{serviceInfo}</option>
    );
  }
}

class DropdownSelect extends Component {
  render() {

    const rows = [];

    this.props.services.forEach((service) => {
      rows.push(
        <SelectRow
          service={service}
          id={service.id}
          key={service.id.toString()}
        />
      );
    });

    return (
      <FormControl as="select"
        value={this.props.value}
        onChange={this.props.handleSelectServiceChange}>
        <option value="" selected disabled hidden>Filter services...</option>
        {rows}
      </FormControl>
    );
  }
}
