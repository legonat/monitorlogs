package models

type LogStruct struct {
	Id			int		`json:"id"`
	SessionId	int		`json:"sessionId"`
	Date        int64	`json:"date"`
	DateUtc		string  `json:"dateUtc"`
	ServiceInfo string	`json:"service"`
	Description string	`json:"description"`
}

type RequestLogStruct struct {
	Id			int		`json:"id"`
	SessionId	int		`json:"sessionId"`
	Date        string	`json:"date"`
	ServiceInfo string	`json:"service"`
	Description string	`json:"description"`
}

type LogFileStruct struct {
	Id              int
	LogfileName     string
	FileLength      int
	LastSessionDate int64
	PreviousDate	int64
	SessionCount    int
}

type LogFilenameStruct struct {
	Id          int		`json:"id"`
	LogfileName string	`json:"value"`
}

type ErrorStruct struct {
	LogStruct
	LogId		int		`json:"logId"`
}

type LogSessionStruct struct {
	Id			int	`json:"id"`
	Dates       string	`json:"value"`
}

type GetLogsBySessionStruct struct {
	SessionId	int 	`json:"sessionId"`
	Filename	string 	`json:"filename"`
}

type GetLogsBySessionWithLimitStruct struct {
	GetLogsBySessionStruct
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

	type GetLogsByDateStruct struct {
	StartDate	string 	`json:"startDate"`
	EndDate		string 	`json:"endDate"`
	Filename	string 	`json:"filename"`
}

type GetLogsByDateWithLimitStruct struct {
	GetLogsByDateStruct
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type FindLogsStruct struct {
	SearchText string `json:"text"`
	Filename   string `json:"filename"`
}


type FindLogsStructWithLimit struct {
	FindLogsStruct
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}


type LineStruct struct {
	Number int
	Length int
}
