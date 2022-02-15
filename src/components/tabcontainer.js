<Container>
<Row>
  <Col>
    <Tabs defaultActiveKey="logs" id="uncontrolled-tab-example">
      <Tab eventKey="logs" title="Logs">
        <TablePagination 
          currentPage = {this.props.currentPage}
          totalLenght = {this.props.totalLenght}
          handleBackwardsClick = {this.props.handleBackwardsClick}
          handleFastBackwardsClick = {this.props.handleFastBackwardsClick} 
          handleForwardClick = {this.props.handleForwardClick} 
          handleFastForwardClick = {this.props.handleFastForwardClick}
          handleLastPageClick = {this.props.handleLastPageClick} 
          handleFirstPageClick = {this.props.handleFirstPageClick}
          handlePageClick = {this.props.handlePageClick}
        />
        <LogTable
          names = {NAMES}
          logs = {this.props.logs}
          itemsPerPage = {this.props.itemsPerPage}
          currentPage = {this.props.currentPage}
          filterText = {this.props.filterText}
          filterService = {this.props.filterService}
          filteredLogs = {this.props.filteredLogs}
        />
        <TablePagination 
          currentPage = {this.props.currentPage}
          totalLenght = {this.props.totalLenght}
          handleBackwardsClick = {this.props.handleBackwardsClick}
          handleFastBackwardsClick = {this.props.handleFastBackwardsClick} 
          handleForwardClick = {this.props.handleForwardClick} 
          handleFastForwardClick = {this.props.handleFastForwardClick}
          handleLastPageClick = {this.props.handleLastPageClick} 
          handleFirstPageClick = {this.props.handleFirstPageClick}
          handlePageClick = {this.props.handlePageClick}
        />
      </Tab>
      <Tab eventKey="errors" title="Errors">
        <TablePagination 
          currentPage = {this.props.currentPage}
          totalLenght = {this.props.totalLenght}
          handleBackwardsClick = {this.props.handleBackwardsClick}
          handleFastBackwardsClick = {this.props.handleFastBackwardsClick} 
          handleForwardClick = {this.props.handleForwardClick} 
          handleFastForwardClick = {this.props.handleFastForwardClick}
          handleLastPageClick = {this.props.handleLastPageClick} 
          handleFirstPageClick = {this.props.handleFirstPageClick}
          handlePageClick = {this.props.handlePageClick}
        />
        <LogTable
          names = {NAMES}
          logs = {this.props.logs}
          itemsPerPage = {this.props.itemsPerPage}
          currentPage = {this.props.currentPage}
          filterText = {this.props.filterText}
          filterService = {this.props.filterService}
          filteredLogs = {this.props.filteredLogs}
        />
      <TablePagination 
        currentPage = {this.props.currentPage}
        totalLenght = {this.props.totalLenght}
        handleBackwardsClick = {this.props.handleBackwardsClick}
        handleFastBackwardsClick = {this.props.handleFastBackwardsClick} 
        handleForwardClick = {this.props.handleForwardClick} 
        handleFastForwardClick = {this.props.handleFastForwardClick}
        handleLastPageClick = {this.props.handleLastPageClick} 
        handleFirstPageClick = {this.props.handleFirstPageClick}
        handlePageClick = {this.props.handlePageClick}
      />
      </Tab>
    </Tabs>
  </Col>
  <Col>
    <DropdownSelect
        handleSelectSessionChange = {this.props.handleSelectSessionChange}
        sessions = {this.props.sessions}
        value = {this.props.value}
    />
  </Col>
</Row>
</Container>