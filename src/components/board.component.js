import React, { Component,  useState } from "react";
import DotsLoader from './dots.pulse.loader'
import NavTabs from './navtabs.component'
import TablePagination from './tablePagination.component'
import Logger from "../utils/logger"


const logger = new Logger()


export default class Board extends Component {
  // constructor(props) {
  //   super(props)
  //   this.handleLogsPageChange = this.handleLogsPageChange.bind(this)
  //   this.handleErrorsPageChange = this.handleErrorsPageChange.bind(this)
  //   this.handleSearchPageChange = this.handleSearchPageChange.bind(this)
  //   // this.handleTabClick = this.handleTabClick.bind(this)
  //   this.state = { currentLogsPage: 1, currentErrorsPage: 1, currentSearchPage: 1 }
  // }

  // handleLogsPageChange(page) {
  //   logger.console("Setting State Logs Page")
  //   this.setState({ currentLogsPage: page })
  // }

  // handleErrorsPageChange(page) {
  //   logger.console("Setting State Errors Page")
  //   this.setState({ currentErrorsPage: page })
  // }

  // handleSearchPageChange(page) {
  //   logger.console("Setting State Search Page")
  //   this.setState({ currentSearchPage: page })
  // }



  checkAuth = (props) => {

    var element = <DotsLoader color="white" />

    if (this.props.filenames[0] !== undefined) {
      //logger.console("Rendering Filterable Table")
      element = (<div className='board-inner'>
        <FilterableLogTable
          logs={this.props.logs}
          errors={this.props.errors}
          handleSelectSessionChange={this.props.handleSelectSessionChange}
          handleSelectFilenameChange={this.props.handleSelectFilenameChange}
          handleLogsPageChange = {this.props.handleLogsPageChange}
          handleErrorsPageChange = {this.props.handleErrorsPageChange}
          handleSearchPageChange = {this.props.handleSearchPageChange}
          currentLogsPage={this.props.currentLogsPage}
          currentErrorsPage={this.props.currentErrorsPage}
          currentSearchPage={this.props.currentSearchPage}
          sessions={this.props.sessions}
          logsDropdownValue={this.props.logsDropdownValue}
          sessionsDropdownValue={this.props.sessionsDropdownValue}
          filterText={this.props.filterText}
          filterService={this.props.filterService}
          filteredLogs={this.props.filteredLogs}
          filteredErrors={this.props.filteredErrors}
          filteredSearch={this.props.filteredSearch}
          loading={this.props.loading}
          errorsPageCount={this.props.errorsPageCount}
          logsPageCount={this.props.logsPageCount}
          itemsPerPage={this.props.itemsPerPage}
          currentTab={this.props.currentTab}
          handleTabClick={this.props.handleTabClick}
          completeSearchRequest={this.props.completeSearchRequest}
          searchResults={this.props.searchResults}
          searchPageCount={this.props.searchPageCount}
          handleLinkClick={this.props.handleLinkClick}
          logId={this.props.logId}
          searchedPage={this.props.searchedPage}
          filenames={this.props.filenames}
          datePicker={this.props.datePicker}
          dateRange={this.props.dateRange}
          needScrollDown={this.props.needScrollDown}
          currentSlicePage = {this.props.currentSlicePage}
        />
      </div>)
    }

    return (
      <div>
        {element}
      </div>
    )
  }


  render() {
    //logger.console("Rendering Board")
    return (
      this.checkAuth()
    )
  }
}

class FilterableLogTable extends Component {

  constructor(props) {
    super(props);

  }

  tabsArray = ["Logs", "Errors"]
  // firstTab = "Logs"
  // secondTab = "Errors"

 

  toggleAnimation = () => {
    // var elements = []
    // this.tabsArray.forEach(() => {
    //   elements.push(<DotsLoader color="blue" />)
    // });
    // logger.console("building blue loader")
    // //var elements = [<DotsLoader color = "blue"/>]

    if (this.props.loading === true) {
      logger.console("LOADING Board Content")
      return (<DotsLoader color="blue" />)
    }

    var elements = []
    logger.console("passing content into tabs")
    elements.push(
      <div>
        <TablePagination
          currentPage={this.props.currentLogsPage}
          pageCount={this.props.logsPageCount}
          handlePageChange={this.props.handleLogsPageChange}
        />
        <LogTable
          names={NAMES}
          logs={this.props.logs}
          itemsPerPage={this.props.itemsPerPage}
          currentPage={this.props.currentSlicePage}
          filterText={this.props.filterText}
          filterService={this.props.filterService}
          filteredLogs={this.props.filteredLogs}
          handleLinkClick={this.props.handleLinkClick}
          logId={this.props.logId}
          needScrollDown={this.props.needScrollDown}
        />
        <TablePagination
          currentPage={this.props.currentLogsPage}
          pageCount={this.props.logsPageCount}
          handlePageChange={this.props.handleLogsPageChange}
        />
      </div>
    )

    elements.push(
      <div>
        <TablePagination
          currentPage={this.props.currentErrorsPage}
          pageCount={this.props.errorsPageCount}
          handlePageChange={this.props.handleErrorsPageChange}
        />
        <LogTable
          names={NAMES}
          logs={this.props.logs}
          itemsPerPage={this.props.itemsPerPage}
          currentPage={this.props.currentSlicePage}
          filterText={this.props.filterText}
          filterService={this.props.filterService}
          filteredLogs={this.props.filteredErrors}
          handleLinkClick={this.props.handleLinkClick}
          logId={this.props.logId}
          needScrollDown={this.props.needScrollDown}
        />
        <TablePagination
          currentPage={this.props.currentErrorsPage}
          pageCount={this.props.errorsPageCount}
          handlePageChange={this.props.handleErrorsPageChange}
        />

      </div>
    )


   // if (this.props.searchResults !== undefined && (this.props.searchResults).length !== 0) {
     if (this.props.completeSearchRequest !== '') {
      logger.console("Searching...")
      logger.console(this.props.searchResults)
      return (
        <div>
          <TablePagination
            currentPage={this.props.currentSearchPage}
            pageCount={this.props.searchPageCount}
            handlePageChange={this.props.handleSearchPageChange}
          />
          <LogTable
            names={NAMES}
            logs={this.props.logs}
            itemsPerPage={this.props.itemsPerPage}
            currentPage={this.props.currentSlicePage}
            filterText={this.props.filterText}
            filterService={this.props.filterService}
            filteredSearch={this.props.filteredSearch}
            handleLinkClick={this.props.handleLinkClick}
            logId={this.props.logId}
            needScrollDown={this.props.needScrollDown}
          />
          <TablePagination
            currentPage={this.props.currentSearchPage}
            pageCount={this.props.searchPageCount}
            handlePageChange={this.props.handleSearchPageChange}
          />
        </div>
      )
    } else {
      logger.console("building NAVTABS")
      return (

        <NavTabs
          tabsArray={this.tabsArray}
          handleSelectSessionChange={this.props.handleSelectSessionChange}
          handleSelectFilenameChange={this.props.handleSelectFilenameChange}
          sessions={this.props.sessions}
          filenames={this.props.filenames}
          logsDropdownValue={this.props.logsDropdownValue}
          sessionsDropdownValue={this.props.sessionsDropdownValue}
          contents={elements}
          onTabClick={this.props.handleTabClick}
          currentTab={this.props.currentTab}
          datePicker={this.props.datePicker}
          dateRange={this.props.dateRange}
        />

      )
    }
  }


  render() {
    return (
      this.toggleAnimation()
    );
  }
}

class LogTable extends Component {


  render() {


    const filterText = this.props.filterText;
    const filterService = this.props.filterService

    var i = 0
    const rows = [];
    var itemsPerPage = this.props.itemsPerPage
    var currentPage = this.props.currentPage

    let lastItem = (currentPage + 1) * itemsPerPage
    let firstItem = currentPage * itemsPerPage


    if (filterText === "" && filterService === "") {
      let logs = this.props.logs
      let slice = logs.slice(firstItem, lastItem)
      slice.forEach((log) => {
        rows.push(
          <LogRow
            handleLinkClick={this.props.handleLinkClick}
            logId={this.props.logId}
            needScrollDown={this.props.needScrollDown}
            log={log}
            key={log.id}
          />
        );
      });
    }

    if (this.props.filteredLogs !== undefined) {
      var filteredLogs = this.props.filteredLogs
      let slice = filteredLogs.slice(firstItem, lastItem)
      slice.forEach((log) => {
        rows.push(
          <LogRow
            logId={this.props.logId}
            handleLinkClick={this.props.handleLinkClick}
            needScrollDown={this.props.needScrollDown}
            log={log}
            key={log.id}

          />
        );
      });
    }

    if (this.props.filteredSearch !== undefined) {
      var filteredSearch = this.props.filteredSearch
      let slice = filteredSearch.slice(firstItem, lastItem)
      slice.forEach((log) => {
        rows.push(
          <LogRow
            logId={this.props.logId}
            handleLinkClick={this.props.handleLinkClick}
            needScrollDown={this.props.needScrollDown}
            log={log}
            key={log.id}
          />
        );
      });
    }

    return (

      <table>
        <thead>
          <tr>
            <td id="datecol">Date</td>
            <td id="servcol">Service</td>
            <td id="desccol">Description</td>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </table>
    );
  }
}


// class LogTable extends Component {


//   render() {


//     const filterText = this.props.filterText;
//     const filterService = this.props.filterService

//     var i = 0
//     const rows = [];
//     var itemsPerPage = this.props.itemsPerPage
//     var currentPage = this.props.currentPage

//     let lastItem = currentPage * itemsPerPage
//     let firstItem = (currentPage - 1) * itemsPerPage


//     if (filterText === "" && filterService === "") {
//       let logs = this.props.logs
//       let slice = logs.slice(firstItem, lastItem)
//       slice.forEach((log) => {
//         rows.push(
//           <LogRow
//             handleLinkClick={this.props.handleLinkClick}
//             logId={this.props.logId}
//             needScrollDown={this.props.needScrollDown}
//             log={log}
//             key={log.id}
//           />
//         );
//       });
//     }

//     if (this.props.filteredLogs !== undefined) {
//       var filteredLogs = this.props.filteredLogs
//       let slice = filteredLogs.slice(firstItem, lastItem)
//       slice.forEach((log) => {
//         rows.push(
//           <LogRow
//             logId={this.props.logId}
//             handleLinkClick={this.props.handleLinkClick}
//             needScrollDown={this.props.needScrollDown}
//             log={log}
//             key={log.id}

//           />
//         );
//       });
//     }

//     if (this.props.filteredSearch !== undefined) {
//       var filteredSearch = this.props.filteredSearch
//       let slice = filteredSearch.slice(firstItem, lastItem)
//       slice.forEach((log) => {
//         rows.push(
//           <LogRow
//             logId={this.props.logId}
//             handleLinkClick={this.props.handleLinkClick}
//             needScrollDown={this.props.needScrollDown}
//             log={log}
//             key={log.id}
//           />
//         );
//       });
//     }

//     return (

//       <table>
//         <thead>
//           <tr>
//             <td id="datecol">Date</td>
//             <td id="servcol">Service</td>
//             <td id="desccol">Description</td>
//           </tr>
//         </thead>
//         <tbody>{rows}</tbody>
//       </table>
//     );
//   }
// }


class LogRow extends Component {

  constructor(props) {
    super(props)
    this.scrollRef = React.createRef()
  }

  render() {

    const scrollTo = (ref) => {
      if (ref && this.props.needScrollDown) {
        ref.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }
    }

    //<div ref={scrollTo}>Item</div>
    const log = this.props.log
    const id = log.id
    const date = log.dateUtc;
    const service = log.service;
    const description = log.description;
    const link = "http://localhost:5000/v2/private/getLogById?id=" + id

    if (log.id === this.props.logId) {
      logger.console("Found log")
      return (
        <tr className="selectedRow" ref={scrollTo}>
          <td id="datecol"><a href={link} onClick={this.props.handleLinkClick} data-log-id={id}>{date}</a></td>
          <td id="servcol">{service}</td>
          <td id="desccol">{description}</td>
        </tr>
      );
    }

    return (
      <tr>
        <td id="datecol"><a href={link} onClick={this.props.handleLinkClick} data-log-id={id}  className="table-link">{date}</a></td>
        <td id="servcol">{service}</td>
        <td id="desccol">{description}</td>
      </tr>
    );
  }
}

class SelectRow extends Component {
  render() {
    const session = this.props.session
    const dates = session.dates;
    var id = 0
    id = session.id;

    return (
      <option data-key={id}>{id + ".    " + dates}</option>
    );
  }
}

class DropdownSelect extends Component {
  render() {

    const rows = [];

    this.props.sessions.forEach((session) => {
      rows.push(
        <SelectRow
          session={session}
          key={session.id}
          id={session.id}
        />
      );
    });

    return (
      <div className="filter-component">
        <select
          value={this.props.value}
          onChange={this.props.handleSelectSessionChange}>
          <option value="" selected disabled hidden>Select session</option>
          {rows}
        </select>
      </div>
    );
  }
}




{/* <NavDropdown eventKey="3" title={this.props.value} id="nav-dropdown-within-tab" onSelect={this.props.handleSelectSessionChange}>
          {rows}
        </NavDropdown>  */}



const NAMES = { id: 0, date: "Date", service: "Source", description: "Description" }
  // class SearchBar extends Component {
  //   constructor(props) {
  //     super(props);
  //     this.handleFilterTextChange = this.handleFilterTextChange.bind(this);
  //   }

  //   handleFilterTextChange(e) {
  //     this.props.handleFilterTextChange(e.target.value);
  //   }


  //   render() {
  //     return (
  //       <div className="filter-component">
  //       <DropdownSelect
  //           handleSelectSessionChange = {this.props.handleSelectSessionChange}
  //           sessions = {this.props.sessions}
  //           value = {this.props.value}
  //         />
  //       <form id="board-filter"> 
  //         <input
  //           type="text"
  //           placeholder="Filter"
  //           value={this.props.filterText}
  //           onChange={this.handleFilterTextChange}
  //         />
  //       </form>
  //       </div>
  //     );
  //   }
  // }

  // class FilterableLogTable extends Component {
  //   constructor(props) {
  //     super(props);
  //     this.state = {
  //       filterText: ''
  //     };


  //     this.handleFilterTextChange = this.handleFilterTextChange.bind(this);
  //   }

  //   handleFilterTextChange(filterText) {
  //     this.setState({
  //       filterText: filterText
  //     });
  //   }

  //   render() {
  //     return (
  //       <div>
  //         <SearchBar
  //         handleSelectSessionChange = {this.props.handleSelectSessionChange}
  //         sessions = {this.props.sessions}
  //         value = {this.props.value}
  //         filterText={this.state.filterText}
  //         handleFilterTextChange={this.handleFilterTextChange}
  //         />
  //         <LogTable
  //           names = {NAMES}
  //           logs={this.props.logs}
  //           filterText={this.state.filterText}
  //         />
  //       </div>
  //     );
  //   }
  // }




//   ReactDOM.render(
//     <FilterableProductTable products={PRODUCTS} />,
//     document.getElementById('container')
//   );