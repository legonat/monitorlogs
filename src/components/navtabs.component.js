import React, { Component } from "react";
import { Container, Row, Col, Tab, Nav, Dropdown, NavItem, NavLink, Form } from 'react-bootstrap/'
import Datetime from 'react-datetime'
import moment from 'moment'

export default class NavTabs extends Component {



  buildTabsRow = () => {
    var tabsRowElements = []
    this.props.tabsArray.forEach((element, index) => {
      tabsRowElements.push(
        <NavItem onClick={() => this.props.onTabClick(index)} key={index.toString()}>
          <NavLink eventKey={element}>{element}</NavLink>
        </NavItem>
      )
    })
    return tabsRowElements
  }

  buildTabsContent = () => {
    var tabsContentElements = []
    const propContent = this.props.contents

    this.props.tabsArray.forEach((element, index) => {
      tabsContentElements.push(
        <Tab.Pane eventKey={element} key={index.toString()}>
          {propContent[index]}
        </Tab.Pane>
      )
    })
    return tabsContentElements
  }

  buildUnorderedDropdown = (dropdownName, dropdownSelected, selectChange, list, align) => {
    var dropdownValue = dropdownName
    if (dropdownSelected !== "") {
      dropdownValue = dropdownSelected
    }
    if (align === "left") {
      return (
        <NavItem className="ml-auto">
          <DropdownTabUnorderedSelect
            handleSelectChange={selectChange}
            data={list}
            value={dropdownValue}
          />
        </NavItem>
      )
    }
    return (
      <NavItem className="mr-auto">
        <DropdownTabUnorderedSelect
          handleSelectChange={selectChange}
          data={list}
          value={dropdownValue}
        />
      </NavItem>
    )
  }
  buildOrderedDropdown = (dropdownName, dropdownSelected, selectChange, list, align) => {
    var dropdownValue = dropdownName
    if (dropdownSelected !== "") {
      dropdownValue = dropdownSelected
    }
    if (align === "left") {
      return (
        <NavItem className="ml-auto">
          <DropdownTabOrderedSelect
            handleSelectChange={selectChange}
            data={list}
            value={dropdownValue}
          />
        </NavItem>
      )
    }
    return (
      <NavItem className="mr-auto">
        <DropdownTabOrderedSelect
          handleSelectChange={selectChange}
          data={list}
          value={dropdownValue}
        />
      </NavItem>
    )
  }

  buildSessionsDropdown = () => {
    if (this.props.sessions === undefined) {
      return null
    }
    var dropdownValue = "Select session"
    if (this.props.sessionsDropdownValue !== "") {
      dropdownValue = this.props.sessionsDropdownValue
    }
    return (
      <NavItem className="ml-auto">
        <DropdownTabOrderedSelect
          handleSelectChange={this.props.handleSelectSessionChange}
          data={this.props.sessions}
          value={dropdownValue}
        />
      </NavItem>
    )
  }

  buildFilenamesDropdown = () => {
    if (this.props.filenames === undefined) {
      return null
    }
    var dropdownValue = "Select Log Name"
    if (this.props.logsDropdownValue !== "") {
      dropdownValue = this.props.logsDropdownValue
    }
    return (
      <NavItem>
        <DropdownTabUnorderedSelect
          handleSelectChange={this.props.handleSelectFilenameChange}
          data={this.props.filenames}
          value={dropdownValue}
        />
      </NavItem>
    )
  }

  buildDTP = () => {
    return (
      <NavItem>
        <MyDTPicker />
      </NavItem>
    )
  }

  buildDTP2 = () => {
    return (
      <NavItem as={Datetime}>
      </NavItem>
    )
  }

  buildSessionString = () => {
    if (this.props.sessionsDropdownValue !== '') {
      return (
        <p className="session-text">Selected session: {this.props.sessionsDropdownValue}</p>
      )
    }
    if (this.props.dateRange !== undefined && this.props.dateRange[0] !== null && this.props.dateRange[1] !== null) {
  
      return (
        <p className="session-text">Selected session: {moment(this.props.dateRange[0]).format("YYYY-MM-DD")} - {moment(this.props.dateRange[1]).format("YYYY-MM-DD")}</p>
      )
    }
  }

  render() {
    return (
      <div>
        <Tab.Container id="tabs-with-dropdown" defaultActiveKey={this.props.tabsArray[this.props.currentTab]}>
          <Row className="clearfix">
            <Col sm={9}>
              <Nav variant="tabs" defaultActiveKey={this.props.tabsArray[this.props.currentTab]}>
                {this.buildTabsRow()}
                {this.buildFilenamesDropdown()}
                {this.buildSessionString()} 
                {/* {this.buildSessionsDropdown()} */}
              </Nav>
            </Col>
            <Col sm={3}>
              {this.props.datePicker}
            </Col>
          </Row>
          {/* {this.buildSessionString()} */}
          <Row>
            <Col sm={12}>
              <Tab.Content>{/* < animation> */}
                {this.buildTabsContent()}
              </Tab.Content>
            </Col>
          </Row>
        </Tab.Container>
      </div>
    )
  }

}

class DropdownTabOrderedSelect extends Component {
  render() {

    const rows = [];

    this.props.data.forEach((d) => {
      const value = d.value;
      var id = 0
      id = d.id;
      rows.push(
        <Dropdown.Item
          eventKey={id}
          data-key={id}
          key={id.toString()}>
          {id + ".    " + value}
        </Dropdown.Item>
      );
    });

    return (
      <Dropdown as={NavItem} onSelect={this.props.handleSelectChange}>
        <Dropdown.Toggle as={NavLink}>{this.props.value}</Dropdown.Toggle>
        <Dropdown.Menu>
          {rows}
        </Dropdown.Menu>
      </Dropdown>
    );
  }
}

class DropdownTabUnorderedSelect extends Component {
  render() {

    const rows = [];

    this.props.data.forEach((d) => {
      const value = d.value;
      var id = 0
      id = d.id;
      rows.push(
        <Dropdown.Item
          eventKey={id}
          data-key={id}
          key={id.toString()}>
            {value}
        </Dropdown.Item>
      );
    });

    return (
      <Dropdown as={NavItem} onSelect={this.props.handleSelectChange}>
        <Dropdown.Toggle as={NavLink}>{this.props.value}</Dropdown.Toggle>
        <Dropdown.Menu>
          {rows}
        </Dropdown.Menu>
      </Dropdown>
    );
  }
}

class MyDTPicker extends React.Component {
  render() {
    return <Datetime renderInput={this.renderInput} />;
  }
  renderInput(props, openCalendar, closeCalendar) {
    function clear() {
      props.onChange({ target: { value: "" } });
    }
    return (
      <div>
        <input {...props} />
        <button onClick={openCalendar}>open calendar</button>
        <button onClick={closeCalendar}>close calendar</button>
        <button onClick={clear}>clear</button>
      </div>
    );
  }
}
