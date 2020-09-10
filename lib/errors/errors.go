package errors

import (
	"errors"
	"log"
	"strconv"
)

/*** Match Enum squence with xxxList ***/
const (
	// code for MiddleWare
	cello uint64 = (1 + iota) * 10000
	clarinet
	flute
	harp
	oboe
	piano
	piccolo
	viola
	violin
	violinNoVNC
	violinScheduler
)

var middleWareList = [...]string{"", "Cello", "Clarinet", "Flute", "Harp", "Oboe", "Piano", "Piccolo", "Viola", "Violin", "NoVNC", "Scheduler"}

const (
	internal uint64 = (1 + iota) * 1000 // lib
	driver                              // driver
	graphql                             // action
	grpc
	sql
	rabbitmq
)

var functionList = [...]string{"", "Internal", "Driver", "GraphQL", "Grpc", "SQL", "RabbitMQ"}

const (
	// Use Generally
	initFail uint64 = 1 + iota
	connectionFail
	undefinedError
	argumentError
	jsonMarshalError
	jsonUnmarshalError
	requestError  // send Request fail
	responseError // get Response fail or has error
	sendError     // send error to client
	receiveError  // get error as result from server
	parsingError
	tokenExpired
	operationFail
	noResult
	timestampConversionError
	UUIDGenerationError

	// clarinet specific

	// piccolo specific
	prepareError
	executeError
	tokenGenerationError
	loginFailed

	// cello specific

	// violin-scheduler specific

	// flute specific
	ipmiError

	// viola specific

	// piano specific

	// harp specific
	interfaceAddrLookupError
	pfError
	dhcpdError
	fileError
	ifconfigError
	IPAddressError
	subnetInUseError
	subnetNotAllocatedError
	adaptiveIPAllocatedError

	// violin-novnc specific

	// violin specific
	createServerFailed
	createServerRoutineError
	getAvailableNodesError
	getNodesError
	serverNodePresentError
)

var actionList = [...]string{
	"",
	"Initialize fail -> ",
	"Connection fail -> ",
	"Undefined error -> ",
	"Argumnet error -> ",
	"JSON marshal fail -> ",
	"JSON unmarshal fail -> ",
	"Request error -> ",
	"Response error -> ",
	"Send error -> ",
	"Receive error -> ",
	"Parsing error -> ",
	"Token Expired -> ",
	"DB operationfail -> ",
	"DB no result",
	"timestamp conversion error -> ",
	"UUID generation error -> ",

	// clarinet specific

	// piccolo specific
	"Prepare error -> ",
	"Execute error -> ",
	"Token Generation Error -> ",
	"Login failed -> ",

	// cello specific

	// violin-scheduler specific

	// flute specific

	// viola specific

	// piano specific

	// harp specific
	"interface address lookup error -> ",
	"PF error -> ",
	"DHCPD error -> ",
	"file error -> ",
	"ifconfig error -> ",
	"IP address error -> ",
	"Subnet In Use error -> ",
	"Subnet not allocated error -> ",
	"AdaptiveIP allocated error -> ",

	// violin-novnc specific

	// violin specific
	"Create Server failed -> ",
	"Create Server routine error ->",
	"Get available nodes error ->",
	"Get nodes error ->",
	"ServerNode present error ->",
}

var errlogger *log.Logger

func SetErrLogger(l *log.Logger) {
	errlogger = l
}

/*    HCCERROR    */

type HccError struct {
	ErrCode uint64 `json:"errcode"` // decimal error code
	ErrText string `json:"errtext"` // error string
}

func NewHccError(errorCode uint64, errorText string) *HccError {
	return &HccError{
		ErrText: errorText,
		ErrCode: errorCode,
	}
}

func (e HccError) New() error {
	return errors.New(e.ToString())
}

func (e HccError) Error() string {
	return e.ToString()
}

func (e HccError) Code() uint64 {
	return e.ErrCode
}

func (e HccError) Text() string {
	return e.ErrText
}

func (e HccError) ToString() string {
	m := e.ErrCode / 10000
	f := e.ErrCode % 10000 / 1000
	a := e.ErrCode % 1000

	return "[" + middleWareList[m] + "] Code :" + strconv.FormatUint(e.ErrCode, 10) + " (" + functionList[f] + ") " + actionList[a] + " " + e.ErrText
}

func (e HccError) Println() {
	errlogger.Println(e.ToString())
}

func (e HccError) Fatal() {
	errlogger.Fatal(e.ToString())
}

/*    HCCERRORSTACK    */

type HccErrorStack []HccError

func NewHccErrorStack(errList ...*HccError) *HccErrorStack {
	es := HccErrorStack{HccError{ErrCode: 0, ErrText: ""}}

	for _, err := range errList {
		es.Push(err)
	}
	return &es
}

func (es *HccErrorStack) Len() int {
	return es.len() - 1
}

func (es *HccErrorStack) len() int {
	return len(*es)
}

func (es *HccErrorStack) Pop() *HccError {
	l := es.len()
	if l > 1 {
		err := (*es)[l-1]
		*es = (*es)[:l-1]
		return &err
	}
	return nil
}

func (es *HccErrorStack) Push(err *HccError) {
	*es = append(*es, *err)
}

// Dump() will clean stack
func (es *HccErrorStack) Dump() *HccError {
	var firstErr *HccError = nil
	if es.Len() == 0 {
		return nil
	}

	if (*es)[0].ErrCode == 0 {
		errlogger.Fatal("Error Stack is already converted to report form. Cannot dump.\n")
	}

	errlogger.Printf("------ [Dump Error Stack] ------\n")
	errlogger.Printf("Stack Size : %v\n", es.Len())
	firstErr = es.Pop()
	firstErr.Println()
	for err := es.Pop(); err != nil; err = es.Pop() {
		err.Println()
	}
	errlogger.Println("--------- [ End Here ] ---------")
	return firstErr
}

func (es *HccErrorStack) ConvertReportForm() *HccErrorStack {
	if es.Len() > 0 {
		*es = (*es)[1:]
		for idx := 0; idx < es.len(); idx++ {
			(*es)[idx].ErrText = "#" + strconv.Itoa(idx) + " " + (*es)[idx].ToString()
		}
	} else {
		*es = nil
	}
	return es
}
